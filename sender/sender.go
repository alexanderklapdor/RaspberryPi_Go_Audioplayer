package sender

//Imports
import (
	"io"
	"log"
	"net"
	"time"

	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"
)

// reader
func reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		logger.Log.Notice("Received response from server")
		logger.Log.Info("Server: '" + string(buf[0:n]) + "'")
	}
}

// Send JSON to Server
func Send(requestJson []byte) {
	//todo: socket path from external config file

	// Open socket connection
	socketPath := "/tmp/mp.sock"
	logger.Log.Notice("Opening socket connection to " + socketPath)
	con, err := net.Dial("unix", socketPath)
	if err != nil {
		panic(err)
	}
	defer con.Close()

	go reader(con)
	logger.Log.Notice("Sending message to Server")
	_, er := con.Write([]byte(requestJson))
	if er != nil {
		log.Fatal("Write error: ", er)
		return
	}
	time.Sleep(1e9)
}
