package sender

//Imports
import (
	"io"
	"log"
	"net"
	"time"

	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"
)

// reader function
func reader(r io.Reader) {
	// Read 1024 bit
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		logger.Notice("Received response from server")
		logger.Info("Server: " + string(buf[0:n]))
	}
}

// Send JSON to Server
func Send(requestJson []byte, socketPath string) {
	// Open socket connection
	logger.Notice("Opening socket connection to " + socketPath)
	con, err := net.Dial("unix", socketPath)
	// Check if err exists
	if err != nil {
		panic(err)
	}
	defer con.Close()

	go reader(con)
	logger.Notice("Sending message to Server")
	_, er := con.Write([]byte(requestJson))
	if er != nil {
		log.Fatal("Write error: ", er)
		return
	}
	time.Sleep(1e9)
}
