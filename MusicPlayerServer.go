package main

import (
	"bufio"
	"encoding/json"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/portaudiofunctions"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/sender"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/serverfunctions/audiofunctions"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/serverfunctions/connectionfunctions"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/serverfunctions/volumefunctions"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/structs"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/util"
	"github.com/tkanos/gonfig"
)

// Global var definition
var configuration = structs.ServerConfiguration{}
var serverData = structs.ServerData{}

// main
func main() {
	// set up configuration
	err := gonfig.GetConf("config.json", &configuration)
	// set up logger
	logger.Setup(util.JoinPath(configuration.Log_Dir, configuration.Server_Log), configuration.Debug_Infos)
	logger.Notice("Starting MusicPlayerServer...")
	// create server socket mp.sock
	unixSocket := configuration.Socket_Path
	logger.Notice("Creating unixSocket.")
	logger.Info("Listening on " + unixSocket)
	ln, err := net.Listen("unix", unixSocket)
	util.Check(err, "Server")
	// set socketPath
	sender.SetSocketPath(unixSocket)
	connectionfunctions.SetSocketPath(unixSocket)
	// check supported formats
	logger.Notice("Parsing supported formats")
	serverData.SupportedFormats = getSupportedFormats()
	// print supported formats
	printSupportedFormats()
	// start pulseAudio
	logger.Notice("Start Pulseaudio")
	portaudiofunctions.StartPulseaudio()
	for {
		connection, err := ln.Accept()
		connectionfunctions.SetConnection(connection)
		util.Check(err, "Server")
		go receiveCommand()
	}
} // end of main

//receiveCommand function
func receiveCommand() {
	// read message
	buf := make([]byte, 512)
	nr, err := connectionfunctions.Read(buf)
	if err != nil {
		return
	}
	receivedBytes := buf[0:nr]
	logger.Info("Server received message: " + string(receivedBytes))

	// convert message back to a request-object
	logger.Notice("Converting message back to a Request-Object")
	received := structs.Request{}
	json.Unmarshal(receivedBytes, &received)
	command := received.Command
	data := received.Data
	logger.Notice("Command: " + command)
	//logger.Notice("Data   : " + string(data))
	// todo check if values are different from default and different from current values
	// if yes -> change them on every command

	message := "Default-message"
	// switch case commands
	switch command {
	case "addToQueue", "add":
		message = audiofunctions.AddToQueue(data, &serverData)
	case "back", "previous":
		message = audiofunctions.PlayPreviousSong(&serverData)
	case "exit":
		audiofunctions.StopMusic()
		portaudiofunctions.StopPulseaudio()
		connectionfunctions.Close()
	case "info", "default":
		message = printInfo()
	case "loop", "setLoop":
		message = setLoop(data)
	case "louder", "setVolumeUp":
		message = volumefunctions.IncreaseVolume()
	case "next":
		message = audiofunctions.NextMusic(data, &serverData)
	case "pause":
		message = audiofunctions.PauseMusic(data, &serverData)
	case "play":
		message = audiofunctions.PlayMusic(data, &serverData)
	case "quieter", "setVolumeDown":
		message = volumefunctions.DecreaseVolume()
	case "repeat":
		message = audiofunctions.RepeatSong(&serverData)
	case "remove", "delete", "removeAt", "deleteAt":
		message = audiofunctions.RemoveSong(data, &serverData)
	case "resume":
		message = audiofunctions.ResumeMusic()
	case "setup":
		message = setupMusicPlayer(data)
	case "setVolume":
		message = volumefunctions.SetVolume(data)
	case "shuffle", "setShuffle":
		message = shuffleQueue()
	case "stop":
		message = audiofunctions.StopMusic()
	default:
		message = "Unknown command received"
		logger.Error("Unknown command received")
	}

	// write to client
	logger.Notice("Send a message back to the client")
	_, err = connectionfunctions.Write([]byte(message))
	util.Check(err, "Server")
} // end of receiveCommand

// setupMusicPlayer function
func setupMusicPlayer(data structs.Data) string {
	// SetVolume
	portaudiofunctions.SetVolume(strconv.Itoa(data.Volume))
	// Add files to queue
	audiofunctions.AddToQueue(data, &serverData)
	// Set Liio
	setLoop(data)
	// Shuffle Songs
	if data.Shuffle {
		shuffleQueue()
	} // end of if
	return "Set up Music Player" + printInfo()
}

// SetLoop function
func setLoop(data structs.Data) string {
	//check if data.Values length > 0
	if len(data.Values) > 0 {
		value_string := data.Values[0]
		// Check if loop is set to on or off
		if strings.Contains(value_string, "on") || strings.Contains(value_string, "true") {
			serverData.SaveLoop = true
		} else if strings.Contains(value_string, "off") || strings.Contains(value_string, "false") {
			serverData.SaveLoop = false
		} // end of else
	} else {
		serverData.SaveLoop = data.Loop
	} // end of else
	return "Set loop to " + strconv.FormatBool(serverData.SaveLoop)
} // end of setLoop

// printInfo function
func printInfo() string {
	logger.Info("Executing: Print info ")
	message := "\n"
	if len(serverData.SongQueue) != 0 {
		message = message + ("Current Song: " + util.PrintMp3Infos(serverData.SongQueue[serverData.CurrentSong]) + "\n")
		if (len(serverData.SongQueue) - 1 - serverData.CurrentSong) != 0 {
			message = message + ("Song Queue: \n")
			//songs from current to end
			for index, song := range serverData.SongQueue[serverData.CurrentSong+1:] {
				message = message + (strconv.Itoa(index+1) + ". " + util.PrintMp3Infos(song) + "\n")
			} // enf of for
			// songs from beginning to current
			if serverData.SaveLoop {
				for index, song := range serverData.SongQueue[:serverData.CurrentSong] {
					message = message + (strconv.Itoa(len(serverData.SongQueue)+index-serverData.CurrentSong) + ". " + util.PrintMp3Infos(song) + "\n")
				} //end of for

			} // end of if
			for _, line := range strings.Split(message, "\n") {
				logger.Info(line)
			} // end of for
		} else {
			message = message + "The Song Queue is empty. \n"
		} // end of else
	} else {
		message = message + ("Currently there is no song playing \n")
		message = message + "The Song Queue is empty. \n"
	} // end of else
	message = message + volumefunctions.PrintVolume()
	return message
} // end of printInfo

// getSupportedFormats function
func getSupportedFormats() []string {
	// get supported audio formats of 'supportedFormats.cfg' file
	supportedFormats := make([]string, 0)

	// Opening file
	file, err := os.Open("supportedFormats.cfg")
	util.Check(err, "Server")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.ContainsAny(line, "#") {
			supportedFormats = append(supportedFormats, line)
		} //end of if
	} // end of for
	util.Check(scanner.Err(), "Server")
	return supportedFormats
} // End of getSupportedFormats

// printSupportedFormats function
func printSupportedFormats() {
	formatString := ""
	for _, format := range serverData.SupportedFormats {
		if formatString != "" {
			formatString = formatString + ", "
		} // end of if
		formatString = formatString + format
	} // end of for
	logger.Info("Supported formats: " + formatString)
} // end of printSupportedFormats

// Shuffle Queue Function
func shuffleQueue() string {
	// Check if Queue is filled
	if len(serverData.SongQueue) > 0 {
		serverData.SongQueue = util.Shuffle(serverData.SongQueue)
		return "Queue has been shuffled"
	} else {
		return "Queue is not filled - shuffle failed"
	}
}
