package main

import (
        "log"
        "net"
        "fmt"
        "encoding/json"
        "os"
        "syscall"
)

type Request struct {
    Command string
    Data Data
}

type Data struct {
    Volume int
    Songs []string
}

func receiveCommand(c net.Conn) {
    // read message
    buf := make([]byte, 512)
    nr, err := c.Read(buf)
    if err != nil {
        return
    }
    receivedBytes := buf[0:nr]
    println("Server received command: ", string(receivedBytes))
    // convert message back to json
    received := Request{}
    json.Unmarshal(receivedBytes, &received)
    command := received.Command
    data := received.Data

    // switch case commands
    switch command {
    case "exit":
        closeConnection(c)
    case "play":
        playMusic(data)
    case "pause":
        pauseMusic(data)
    case "setVolume":
        setVolume(data)
    case "addToQueue":
        addToQueue(data)
    case "quieter":
        increaseVolume()
    case "louder":
        decreaseVolume()
    case "info":
        printInfo()
    default:
        println("Unknown command")
    }


    if command == "exit" {
       closeConnection(c)
    }


    // write to client
    _, err = c.Write([]byte("Ok perfect"))
    if err != nil {
        log.Fatal("Write: ", err)
    }
} // end of receiveCommand

func closeConnection(c net.Conn) {
    println("Connection will be closed")
    defer c.Close()
    err :=  syscall.Unlink("/tmp/mp2.sock")
    if err != nil {
        println("Unlink()", err)
    }
    os.Exit(0)
} // end of closeConnection

func playMusic(data Data) {
    println("Play Music")
}

func pauseMusic(data Data) {
    println("Pause Music")
}

func setVolume(data Data) {
    println("Set Volume")
}

func addToQueue(data Data) {
    println("Add to queue")
}

func increaseVolume() {
    println("Increase volume")
}

func decreaseVolume() {
    println("Decrease volume")
}

func printInfo() {
    println("Info ")
}

func main() {
    unixSocket := "/tmp/mp.sock"
    // create server socket mp.sock
    fmt.Println("Starting Unix socket on '" + unixSocket + "'")
    ln, err := net.Listen("unix", unixSocket)
    if err != nil {
        log.Fatal("listen error", err)
    }

    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Fatal("accept error: ", err)
        }
        go receiveCommand(conn)
    }
} // end of main
