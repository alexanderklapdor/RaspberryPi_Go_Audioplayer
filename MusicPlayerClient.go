package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/screener"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/sender"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/structs"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/util"
	id3 "github.com/mikkyang/id3-go"
	"github.com/tkanos/gonfig"
)

// Global var declaration
var configuration = structs.ClientConfiguration{}

// Main function
func main() {
	// Set up configuration
	err := gonfig.GetConf("config.json", &configuration)
	util.Check(err)

	// Set up Logger
	logger.Setup(path.Join(configuration.Log_Dir, configuration.Client_Log), true)

	// Start Screen
	screener.StartScreen()
	logger.Notice("Starting MusicPlayerClient...")

	socket_path := configuration.Socket_Path

	// check if server is running
	if checkServerStatus() {
		logger.Info("Server is running")
	} else {
		logger.Info("Server is not running")
		// Start Server
		startServer()
		// Wait for server has been started
		ind := 0
		for {
			if _, err := os.Stat(socket_path); err == nil ||
				ind >= configuration.Server_Connection_Attempts {
				break
			} // end of if
			logger.Info("Waiting for server")
			time.Sleep(1 * time.Second)
			ind++
		} // end of for

		//check Server Stat
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
	command := flag.String("c", configuration.Default_Command, "command for the server (default "+
		configuration.Default_Command+")")
	input := flag.String("i", configuration.Default_Input, "input music file/folder (default "+
		configuration.Default_Input+")")
	volume := flag.Int("v", configuration.Default_Volume, "music volume in percent (default "+
		strconv.Itoa(configuration.Default_Volume)+")")
	depth := flag.Int("d", configuration.Default_Depth, "audio file searching depth (default/recommended "+
		strconv.Itoa(configuration.Default_Depth)+")")
	shuffle := flag.Bool("s", configuration.Default_Shuffle, "shuffle (default "+
		strconv.FormatBool(configuration.Default_Shuffle)+")")
	loop := flag.Bool("l", configuration.Default_Loop, "loop (default "+
		strconv.FormatBool(configuration.Default_Loop)+")")
	fadeIn := flag.Int("fi", configuration.Default_FadeIn, "fadein in milliseconds (default "+
		strconv.Itoa(configuration.Default_FadeIn)+")")
	fadeOut := flag.Int("fo", configuration.Default_FadeOut, "fadeout in milliseconds (default "+
		strconv.Itoa(configuration.Default_FadeOut)+")")

	// parsing flags
	logger.Notice("Start Parsing cli parameters")
	flag.Parse()

	var values []string
	// if argument without flagname is given parse it as command
	if flag.NArg() > 1 {
		// command argument
		*command = flag.Arg(0)
		// value arguments
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

	// check volume
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

	// parsing to json
	logger.Notice("Parsing argument to json")
	dataInfo := &structs.Data{
		Depth:   *depth,
		FadeIn:  *fadeIn,
		FadeOut: *fadeOut,
		Shuffle: *shuffle,
		Loop:    *loop,
		Path:    *input,
		Values:  values,
		Volume:  *volume}
	requestInfo := &structs.Request{
		Command: string(*command),
		Data:    *dataInfo}
	requestJson, _ := json.Marshal(requestInfo)
	logger.Info("JSON String : " + string(requestJson))

	// Send command
	sender.SetSocketPath(configuration.Socket_Path)
	sender.Send(requestJson)

	// Closing Client
	logger.Info("Closing MusicPlayerClient...\n")

}

// checkServerStatus function
func checkServerStatus() bool {
	// get socket_path
	socket_path := configuration.Socket_Path
	// check if socket exists
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

// startServer function
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

	//Start process
	process, err := os.StartProcess(util.GetGoExPath(), []string{"go", "run", "MusicPlayerServer.go"}, &attr)
	util.Check(err)
	logger.Info("Detaching process")
	err = process.Release()
	util.Check(err)
}

// printMp3Infos function
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
			fmt.Println("Title: " + title + "\t\t\t\tArtist: " + artist + "\t\t\t\tAlbum: " + album + "\t\t\t\tLength: " + util.SecondsToMinutes(length))
		}
	}
}
