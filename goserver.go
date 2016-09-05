package main

import (
    "fmt"
    "net"
    "os"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/gorilla/websocket"
    "github.com/gorilla/mux"
    "time"
    "net/http"
)
type datagrams struct {
    ID        bson.ObjectId `bson:"_id,omitempty"`
    Datagram   string        `json:"datagram" bson:"datagram"`
    Timestamp time.Time
}
type Client struct{
    Id int
    Websocket *websocket.Conn
}
type responseJSON struct{
    Data string
}

type ClientTcp struct{
    Id int
    TelnetConn net.Conn
}

var Clients = make(map[int] Client)
var ClientsTcp = make(map[int] ClientTcp)
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}


func main() {
    go runServerSocket()// Funcion que levanta el servidor websocket concurrente
    go runServerTcp()
    //Configuracion del servidor UDP
    c := getSession().DB("udpdatagram").C("datagrams")
    ServerAddr,err := net.ResolveUDPAddr("udp",":33333")
    ServerConn, err := net.ListenUDP("udp", ServerAddr)
    CheckError(err)
    defer ServerConn.Close()
    buf := make([]byte, 1024)

    
   
    for {
        n,addr,err := ServerConn.ReadFromUDP(buf)
        fmt.Println("Received ",string(buf[0:n]), " from ",addr)

        if err != nil {
            fmt.Println("Error: ",err)
        }
       
        err = c.Insert(&datagrams{Datagram:string(buf[0:n]),Timestamp: time.Now()}) //Inserta en Mongo DB

        if err != nil {
            fmt.Println("Error al insertar en Mongo")
        }

        //Envia el datagram a todas las websockets conectadas
        for _, client := range Clients {
           if err := client.Websocket.WriteJSON(responseJSON{Data:string(buf[0:n])}); err != nil{
            return
           }
        }
        for _, clientNew := range ClientsTcp{
            if _,err := clientNew.TelnetConn.Write([]byte(buf[0:n])); err != nil{
                fmt.Println("Error al enviar peticion tcp")
                delete(ClientsTcp, clientNew.Id)
            }
        }
        

    }
    
}
//Configuracion de la websocket
var wsupgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func runServerSocket(){
    mux := mux.NewRouter()
    mux.HandleFunc("/ws",newSocketConn).Methods("GET")
    http.Handle("/",mux)
    http.ListenAndServe(":3000",nil)
}

func newSocketConn(w http.ResponseWriter, r *http.Request) {
    ws, err := websocket.Upgrade(w,r,nil,1024,1024)
    if err != nil{
        return
    }
    count := len(Clients)
    new_client := Client{count,ws}
    Clients[count] = new_client

    for{
        _, _, err :=ws.ReadMessage()
        if err !=nil{
            delete(Clients, new_client.Id)
            return
        }

    }

}

func runServerTcp(){
    // Listen for incoming connections.
    l, err := net.Listen("tcp", "192.168.129.178"+":"+"3333")
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Println("Listening on localhost 3333")

    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()

        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        count := len(ClientsTcp)
        fmt.Println(count)
        new_clientTcp := ClientTcp{count,conn}
        ClientsTcp[count] = new_clientTcp
    }
}

func getSession() *mgo.Session {  
    s, err := mgo.Dial("mongodb://localhost/")
    if err != nil {
        fmt.Println("No de puedo conectar a Mongo DB")
    }else{
        fmt.Println("CONECTADO")
    }
    return s
}