package connectionfunctions

// Imports
import (
	"net"
	"os"
	"syscall"

	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"
)

// Variable definition
var connection net.Conn
var socketPath string

// SetSocketPath function
func SetSocketPath(tempSocketPath string) {
	socketPath = tempSocketPath
}

// SetConnection function
func SetConnection(tempConnection net.Conn) {
	connection = tempConnection
}

// Read fuction
func Read(buf []byte) (int, error) {
	return connection.Read(buf)
}

// Write funtion
func Write(message []byte) (int, error) {
	return connection.Write(message)
}

// Close function
func Close() {
	// Get socket path
	logger.Warning("Connection will be closed")
	defer connection.Close()
	// Unlink Socket
	err := syscall.Unlink(socketPath)
	if err != nil {
		logger.Error("Error during unlink process of the socket: " + err.Error())
		logger.Info("Pls run manually unlink 'unlink" + socketPath + "'")
		os.Exit(69)
	}
	logger.Info("Closing MusicPlayerServer...\n")
	os.Exit(0)
}
