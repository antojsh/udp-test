package main

import (
    "fmt"
    "net"
    "os"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/gorilla/websocket"
    //"github.com/gorilla/mux"
    "time"
    "net/http"
    "log"
    //"code.google.com/p/go.net/websocket"
)
type datagrams struct {
    ID        bson.ObjectId `bson:"_id,omitempty"`
    Datagram   string        `json:"datagram" bson:"datagram"`
    Timestamp time.Time
}



func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}


func main() {
    go runserver()
    c := getSession().DB("udpdatagram").C("datagrams")
    ServerAddr,err := net.ResolveUDPAddr("udp",":33333")
    ServerConn, err := net.ListenUDP("udp", ServerAddr)
    CheckError(err)
    defer ServerConn.Close()

    buf := make([]byte, 1024)

    for {
        n,_,err := ServerConn.ReadFromUDP(buf)
        //fmt.Println("Received ",string(buf[0:n]), " from ",addr)

        if err != nil {
            fmt.Println("Error: ",err)
        }
        err = c.Insert(&datagrams{Datagram:string(buf[0:n]),Timestamp: time.Now()})

        if err != nil {
            panic(err)
        } 
    }
    
}
func runserver(){
      http.HandleFunc("/echo", echoHandler)
      //http.Handle("/", http.FileServer(http.Dir(".")))
      err := http.ListenAndServe(":3000", nil)
      if err != nil {
        panic("Error: " + err.Error())
      }
}
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}
 
func print_binary(s []byte) {
  fmt.Printf("Received b:");
  for n := 0;n < len(s);n++ {
    fmt.Printf("%d,",s[n]);
  }
  fmt.Printf("\n");
}
 
func echoHandler(w http.ResponseWriter, r *http.Request) {

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
 
    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            return
        }
 
        print_binary(p)
        //fmt.Printf(string(messageType));
        err := conn.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
        if err != nil {
            log.Fatal(err)
        }
    }
}

func getSession() *mgo.Session {  
    s, err := mgo.Dial("mongodb://localhost/")
    if err != nil {
        panic(err)
    }else{
        fmt.Println("CONECTADO")
    }
    return s
}