package main

import(
        "io"
        "log"
        "net"
        "time"
        "flag"
        "fmt"
        "encoding/json"
)

type Request struct {
    Command string
    Data Data
}

type Data struct {
    Volume int
    Songs []string
}


func reader(r io.Reader) {
    buf := make([]byte, 1024)
    for {
        n, err := r.Read(buf[:])
        if err != nil {
            return
        }
        fmt.Println("Client received: ", string(buf[0:n]))
    }
}

func main () {
    command := flag.String("c", "command", "command for the server")
    flag.Parse()
    fmt.Println("Received command argument: ", *command)

    dataInfo := &Data{
        Volume:   50,
        Songs:    []string{"Song1", "Song2"}}
    requestInfo := &Request{
        Command:  string(*command),
        Data:     *dataInfo}
    requestJson, _ := json.Marshal(requestInfo)

    fmt.Println(string(requestJson))



    con, err := net.Dial("unix", "/tmp/mp.sock")
    if err != nil {
        panic(err)
    }
    defer con.Close()

    go reader(con)
    _, er := con.Write([]byte(requestJson) )
    if er != nil {
        log.Fatal("Write error: ", er)
        return
    }
    time.Sleep(1e9)
}
