package logger

// Imports
import (
	"fmt"
	"os"
	"time"
)

// Var Definition
var file *os.File
var filename = "./Log.txt"
var debug = false

// Setup Logger
func Setup(filepath string, debugMode bool) {
	filename = filepath
	debug = debugMode
}

// Debug Option
func Debug(message string) {
	//check file
	checkFile()
	//open file
	openFile()
	//close file after writing
	defer file.Close()
	//append to file
	appendString("DEBU", message)
	//print to stdoutput if debug is true
	if debug == true {
		fmt.Println("DEBUG -> " + message)
	}
}

// Info Option
func Info(message string) {
	//check file
	checkFile()
	//open file
	openFile()
	//close file after writing
	defer file.Close()
	//append to file
	appendString("INFO", message)
	//print to stdoutput if debug is true
	if debug == true {
		fmt.Println("INFO -> " + message)
	}
}

// Notice Option
func Notice(message string) {
	//check file
	checkFile()
	//open file
	openFile()
	//close file after writing
	defer file.Close()
	//append to file
	appendString("NOTI", message)
	//print to stdoutput if debug is true
	if debug == true {
		fmt.Println("NOTICE -> " + message)
	}
}

// Warning Option
func Warning(message string) {
	//check file
	checkFile()
	//open file
	openFile()
	//close file after writing
	defer file.Close()
	//append to file
	appendString("WARN", message)
	//print to stdoutput if debug is true
	if debug == true {
		fmt.Println("WARNING -> " + message)
	}
}

// Error Option
func Error(message string) {
	//check file
	checkFile()
	//open file
	openFile()
	//close file after writing
	defer file.Close()
	//append to file
	appendString("ERRO", message)
	//log to stdoutput
	fmt.Println("ERROR -> " + message)

}

// Critical Option
func Critical(message string) {
	//check file
	checkFile()
	//open file
	openFile()
	//close file after writing
	defer file.Close()
	//append to file
	appendString("CRIT", message)
	//log to stdoutput
	fmt.Println("CRITICAL -> " + message)
	panic(message)
}

func checkFile() {
	if _, err := os.Stat("./Log.txt"); err == nil {
		//file exists
	} else if os.IsNotExist(err) {
		// file not exists -> create file
		file, err = os.Create(filename)
		check(err)
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		fmt.Println(err)
		panic(err)
	}
}

func openFile() {
	//open file
	var err error
	file, err = os.OpenFile(filename, os.O_APPEND, 0600)
	check(err)
}

func appendString(mType string, message string) {
	//append to file
	var err error
	if _, err = file.WriteString(time.Now().Format("2006.01.02 15:04:05") + " " + mType + " -> " + message + "\n"); err != nil {
		panic(err)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
