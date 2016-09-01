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
    "io/ioutil"
   // "log"
    //"sync"
    //"gopkg.in/igm/sockjs-go.v2/sockjs"
    //"golang.org/x/net/websocket"
//"io"
    //"github.com/gin-gonic/gin"
    //"code.google.com/p/go.net/websocket"
)
type datagrams struct {
    ID        bson.ObjectId `bson:"_id,omitempty"`
    Datagram   string        `json:"datagram" bson:"datagram"`
    Timestamp time.Time
}

// var usersConects []*websocket.Conn
type msg struct {
    Num string
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

        // for _, user := range usersConects {
        //     fmt.Println("USER")
        //     fmt.Println(len(user))
        //     websocket.JSON.Send(user, "xdjlhkl")
        // }
        
        

    }
    
}
var wsupgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
    conn, err := wsupgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("Failed to set websocket upgrade: %+v", err)
        return
    }

    for {
        t, msg, err := conn.ReadMessage()
        if err != nil {
            break
        }
        conn.WriteMessage(t, msg)
    }
}
func runserver(){
    http.HandleFunc("/ws", wsHandler)
    http.HandleFunc("/", rootHandler)

    panic(http.ListenAndServe(":8080", nil))
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
    content, err := ioutil.ReadFile("index.html")
    if err != nil {
        fmt.Println("Could not open file.", err)
    }
    fmt.Fprintf(w, "%s", content)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("Origin") != "http://"+r.Host {
        http.Error(w, "Origin not allowed", 403)
        return
    }
    conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
    if err != nil {
        http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
    }

    go echo(conn)
}

func echo(conn *websocket.Conn) {
    for {
        
        m := msg{}
        err := conn.ReadJSON(&m)
        if err != nil {
            fmt.Println("Error reading json.", err)
        }

        fmt.Printf("Got message: %#v\n", m)

        if err = conn.WriteJSON(m); err != nil {
            fmt.Println(err)
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