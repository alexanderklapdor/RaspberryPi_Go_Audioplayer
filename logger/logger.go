package logger

// Imports
import (
	"fmt"
	"os"
	"time"
)

// Var Definition
var filename = "Log.txt"
var debug = false

// Setup Logger
func Setup(filepath string, debugMode bool) {
	filename = filepath
	debug = debugMode
}

// Debug Option
func Debug(message string) {
	//append to file
	appendString("DEBU", message)
	//print to stdoutput if debug is true
	if debug == true {
		fmt.Println("DEBUG -> " + message)
	}
}

// Info Option
func Info(message string) {
	//append to file
	appendString("INFO", message)
	//print to stdoutput if debug is true
	if debug == true {
		fmt.Println("INFO -> " + message)
	}
}

// Notice Option
func Notice(message string) {
	//append to file
	appendString("NOTI", message)
	//print to stdoutput if debug is true
	if debug == true {
		fmt.Println("NOTICE -> " + message)
	}
}

// Warning Option
func Warning(message string) {
	//append to file
	appendString("WARN", message)
	//print to stdoutput if debug is true
	if debug == true {
		fmt.Println("WARNING -> " + message)
	}
}

// Error Option
func Error(message string) {
	//append to file
	appendString("ERRO", message)
	//log to stdoutput
	fmt.Println("ERROR -> " + message)

}

// Critical Option
func Critical(message string) {
	//append to file
	appendString("CRIT", message)
	//log to stdoutput
	fmt.Println("CRITICAL -> " + message)
	panic(message)
}



func appendString(mType string, message string) {
	//append to file
        file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
        if err != nil {
                panic(err)
        }
	if _, err = file.WriteString(time.Now().Format("2006.01.02 15:04:05") + " " + mType + " -> " + message + "\n"); err != nil {
		panic(err)
	}
}

