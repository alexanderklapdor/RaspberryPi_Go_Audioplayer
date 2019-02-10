package logger

//Imports
import "os"
import "github.com/op/go-logging"
import "log"
import "io"

// Var Definition
var Log = logging.MustGetLogger("Test")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func SetUpLogger(file_path string) {
        // create logFile
        file, err := os.OpenFile(file_path, os.O_CREATE|os.O_APPEND, 0644)
        if err != nil {
                log.Fatal(err)
        }
        defer file.Close()
        writer := io.MultiWriter(file, os.Stdout)
	// Create backend for os.Stderr.
	backend1 := logging.NewLogBackend(writer, "", 0)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend1Formatter)
}
