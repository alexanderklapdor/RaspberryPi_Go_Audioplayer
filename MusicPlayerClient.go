package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/screener"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/sender"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/util"
	id3 "github.com/mikkyang/id3-go"
	// "github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/util"
)

type Request struct {
	Command string
	Data    Data
}

// Data struct
type Data struct {
	Depth   int
	FadeIn  int
	FadeOut int
	Path    string
	Shuffle bool
	Loop    bool
	Values  []string
	Volume  int
}

func main() {
	// Set up Logger
	//todo: logging path in configuration file
	logger.Setup("logs/client.log", true)
	socket_path := "/tmp/mp.sock" // todo: config file socket path

	// Start Screen
	screener.StartScreen()

	// check if server is running
	if checkServerStatus() {
		logger.Info("Server is running")
	} else {
		logger.Info("Server is not running")
		startServer()
		ind := 0
		for {
			if _, err := os.Stat(socket_path); err == nil || ind > 10 {
				break
			} // end of if
			logger.Info("Waiting for server")
			time.Sleep(1 * time.Second)
		} // end of for
		if _, err2 := os.Stat(socket_path); err2 == nil {
			logger.Info("Server started succesfully")
		} else if os.IsNotExist(err2) {
			logger.Info("Server not started succesfully")
			os.Exit(304)
		} else {
			logger.Info("Something unexpected happened")
			os.Exit(777)
		}
	} // end of else

	// check if no argument is given
	if len(os.Args) < 2 {
		logger.Error("Missing required argument")
		return
	}

	// define flags
	command := flag.String("c", "default", "command for the server")
	input := flag.String("i", "", "input music file/folder")
	volume := flag.Int("v", 50, "music volume in percent (default 50)")
	depth := flag.Int("d", 2, "audio file searching depth (default/recommended 2)")
	shuffle := flag.Bool("s", false, "shuffle (default false)")
	loop := flag.Bool("l", false, "loop (default false)")
	fadeIn := flag.Int("fi", 0, "fadein in milliseconds (default 0)")
	fadeOut := flag.Int("fo", 0, "fadeout in milliseconds (default 0)")

	// parsing flags
	logger.Notice("Start Parsing cli parameters")
	flag.Parse()

	var values []string
	// if argument without flagname is given parse it as command
	if flag.NArg() > 1 {
		*command = flag.Arg(0)
		for id, arg := range flag.Args() {
			if id != 0 {
				values = append(values, arg)
			} //end of if
		} // end of for
	} else {
		if flag.NArg() == 1 && *command == "default" {
			*command = flag.Arg(0)
		} // end of if
	} // end of else

	// check received arguments
	logger.Notice("Check received arguments")
	if *volume < 0 || *depth < 0 || *fadeIn < 0 || *fadeOut < 0 {
		logger.Error("no negative values allowed")
		return
	}
	if *volume > 100 {
		logger.Info("No volume above 100 allowed")
		*volume = 100
	}

	// print received argument
	logger.Notice("Given arguments:")
	logger.Info("Command   " + *command)
	logger.Info("Input:    " + *input)
	logger.Info("Volume:   " + strconv.Itoa(*volume))
	logger.Info("Depth:    " + strconv.Itoa(*depth))
	logger.Info("Shuffle:  " + strconv.FormatBool(*shuffle))
	logger.Info("Loop:     " + strconv.FormatBool(*loop))
	logger.Info("Fade in:  " + strconv.Itoa(*fadeIn))
	logger.Info("Fade out: " + strconv.Itoa(*fadeOut))
	//logger.Info("Tail:     " + flag.Args())

	// parsings songs

	// parsing to json
	logger.Notice("Parsing argument to json")

	dataInfo := &Data{
		Depth:   *depth,
		FadeIn:  *fadeIn,
		FadeOut: *fadeOut,
		Shuffle: *shuffle,
		Loop:    *loop,
		Path:    *input,
		Values:  values,
		Volume:  *volume}
	requestInfo := &Request{
		Command: string(*command),
		Data:    *dataInfo}
	requestJson, _ := json.Marshal(requestInfo)
	logger.Info("JSON String : " + string(requestJson))

	sender.Send(requestJson) //todo: socket should be given to the sender

}

func checkServerStatus() bool {
	socket_path := "/tmp/mp.sock"
	if _, err := os.Stat(socket_path); err != nil {
		return false //unix socket does not exists
	} else {
		// check if process exists
		cmd := "ps -ef | grep MusicPlayerServer"
		output, err := exec.Command("bash", "-c", cmd).Output()
		util.Check(err)
		for _, pi := range strings.Split(string(output), "\n") {
			if strings.Contains(pi, "go run") {
				return true
			} // end of if
		} // end of for
		return false
	} // end of else
} // end of checkServerStatus

func startServer() {
	logger.Info("Starting Server process")
	var attr = os.ProcAttr{
		Dir: ".",
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			nil,
			nil,
		},
	}
	process, err := os.StartProcess("/usr/local/go/bin/go", []string{"go", "run", "MusicPlayerServer.go"}, &attr)
	util.Check(err)
	logger.Info("Detaching process")
	err = process.Release()
	util.Check(err)
}

func printMp3Infos(filePath string) {
	//Check if Path exists
	if _, err := os.Stat(filePath); err == nil {
		//open file for id3 tags
		mp3File, err := id3.Open(filePath)
		util.Check(err)
		//close file at the end
		defer mp3File.Close()
		//get Tag Infos
		title := mp3File.Title()
		artist := mp3File.Artist()
		album := mp3File.Album()
		//get Audio length
		blength, lengtherr := exec.Command("mp3info", "-p", "%S", filePath).Output()
		util.Check(lengtherr)

		//check if one information is empty
		if title == "" || artist == "" || album == "" || string(blength[:]) == "" {
			fmt.Println(filePath)
		} else {
			//print Infos
			length, err := strconv.Atoi(string(blength[:]))
			util.Check(err)
			fmt.Println("Title: " + title + "\t\t\t\tArtist: " + artist + "\t\t\t\tAlbum: " + album + "\t\t\t\tLength: " + secondsToMinutes(length))
		}
	}
}

//Get Minute and Secons from Seconds
func secondsToMinutes(inSeconds int) string {
	minutes := inSeconds / 60
	seconds := inSeconds % 60
	str := fmt.Sprintf("%dmin %dsec", minutes, seconds)
	return str
}
