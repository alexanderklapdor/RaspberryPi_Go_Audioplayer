package connectionfunctions

import "net"

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

func Read(buf []byte) (int, error) {
	return connection.Read(buf)
}

func Write(message []byte) (int, error) {
	return connection.Write(message)
}

func Close() {
	connection.Close()
}
