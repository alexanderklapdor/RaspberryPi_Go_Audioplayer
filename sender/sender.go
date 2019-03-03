package sender

//Imports
import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"
)

//global var definition
var socketPath string

// reader function
func reader(r io.Reader) {
	// Read 2048 bit
	buf := make([]byte, 2048)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		logger.Notice("Received response from server")
		logger.Info("Server: " + string(buf[0:n]))
		fmt.Println("Received response from Server...")
		fmt.Println("Server: " + string(buf[0:n]))
	}
}

// Send JSON to Server
func Send(requestJson []byte) {
	// Open socket connection
	logger.Notice("Opening socket connection to " + socketPath)
	con, err := net.Dial("unix", socketPath)
	// Check if err exists
	if err != nil {
		panic(err)
	}
	defer con.Close()

	go reader(con)
	fmt.Println("Sending Command to Server...")
	logger.Notice("Sending Command to Server...")
	_, er := con.Write([]byte(requestJson))
	if er != nil {
		log.Fatal("Write error: ", er)
		return
	}
	time.Sleep(1e9)
}

func SetSocketPath(tempSocketPath string) {
	socketPath = tempSocketPath
}
