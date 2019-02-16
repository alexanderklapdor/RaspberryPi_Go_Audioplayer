package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/portaudiofunctions"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/screener"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/serverfunctions/volumefunctions"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/structs"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/util"
	"github.com/tkanos/gonfig"
)

// Global var definition
var wg = &sync.WaitGroup{}
var configuration = structs.ServerConfiguration{}
var serverData = structs.ServerData{}

//receiveCommand function
func receiveCommand(c net.Conn) {
	// read message
	buf := make([]byte, 512)
	nr, err := c.Read(buf)
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
		message = addToQueue(data)
	case "back", "previous":
		message = playPreviousSong()
	case "exit":
		closeConnection(c)
	case "info", "default":
		message = printInfo()
	case "loop", "setLoop":
		message = setLoop(data)
	case "louder", "setVolumeUp":
		message = volumefunctions.IncreaseVolume()
	case "next":
		message = nextMusic(data)
	case "pause":
		message = pauseMusic(data)
	case "play":
		message = playMusic(data)
	case "quieter", "setVolumeDown":
		message = volumefunctions.DecreaseVolume()
	case "repeat":
		message = repeatSong()
	case "remove", "delete", "removeAt", "deleteAt":
		message = removeSong(data)
	case "resume":
		message = resumeMusic()
	case "setup":
		message = setupMusicPlayer(data)
	case "setVolume":
		message = volumefunctions.SetVolume(data)
	case "shuffle", "setShuffle":
		message = shuffleQueue()
	case "stop":
		message = stopMusic()
	default:
		logger.Error("Unknown command received")
	}

	// write to client
	logger.Notice("Send a message back to the client")
	_, err = c.Write([]byte(message))
	if err != nil {
		log.Fatal("Write: ", err)
	}
} // end of receiveCommand

// setupMusicPlayer function
func setupMusicPlayer(data structs.Data) string {
	// SetVolume
	portaudiofunctions.SetVolume(strconv.Itoa(data.Volume))
	// Add files to queue
	addToQueue(data)
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

// closeConnection function
func closeConnection(c net.Conn) {
	// get socket path
	socketPath := configuration.Socket_Path
	logger.Warning("Connection  will be closed")
	defer c.Close()
	//unlink Socket
	err := syscall.Unlink(socketPath)
	if err != nil {
		logger.Error("Error during unlink process of the socket: " + err.Error())
		logger.Info("Pls run manually unlink 'unlink" + socketPath + "'")
		os.Exit(69)
	}
	os.Exit(0)
} // end of closeConnection

// playPreviousSong function
func playPreviousSong() string {
	//check if currentSong is > 0 in queue
	if serverData.CurrentSong > 0 {
		// decrease currentSong and play
		serverData.CurrentSong--
		playCurrentSong()
		return "Playing now " + serverData.SongQueue[serverData.CurrentSong]
	} else {
		// check if loop is enabled
		if serverData.SaveLoop {
			if len(serverData.SongQueue) > 0 {
				serverData.CurrentSong = len(serverData.SongQueue) - 1
				playCurrentSong()
				return "Playing now " + serverData.SongQueue[serverData.CurrentSong]
			} else {
				return "Error: The queue is empty. You could't go a song back"
			} // end of else
		} else {
			return "You are currently playing the first song"
		} // end of else
	} // end of else
	return "should never be shown"
} // end of playLastSong

// repeatSong function
func repeatSong() string {
	// check sonqQueue length
	if len(serverData.SongQueue) > 0 {
		playCurrentSong()
		return "Playing now " + serverData.SongQueue[serverData.CurrentSong]
	} else {
		return "There is no current song"
	} // end of else
} // end of repeatSong

// removeSong function
func removeSong(data structs.Data) string {
	if len(data.Values) != 0 {
		number, err := strconv.Atoi(data.Values[0])
		// todo: remove multiple values (problem with changing position)
		if err == nil {
			//Check loop is off and song in queue
			if number > 0 && number < (len(serverData.SongQueue)-serverData.CurrentSong) {
				number = number + serverData.CurrentSong
				song_name := serverData.SongQueue[number]
				serverData.SongQueue = append(serverData.SongQueue[:number], serverData.SongQueue[number+1:]...)
				return "Removed song" + song_name
				//check loop is on and song is queue
			} else if number >= len(serverData.SongQueue)-serverData.CurrentSong && number < len(serverData.SongQueue) && serverData.SaveLoop {
				number = serverData.CurrentSong - len(serverData.SongQueue) + number
				song_name := serverData.SongQueue[number]
				serverData.SongQueue = append(serverData.SongQueue[:number], serverData.SongQueue[number+1:]...)
				return "Removed song" + song_name
			} else {
				return "There is no song with the given number (" + strconv.Itoa(number) + ")"
			} // end of else
		} else {
			return "Remove is only allowed with the number of the song in the queue. Pls use 'info' to see the queue"
		} // else
	} else {
		return "No argument given"
	}
	return "should never be shown"
} // end of removeSong

// stopMusic function
func stopMusic() string {
	logger.Info("Execution: Stop Music")
	portaudiofunctions.StopAudio()

	wg.Add(1)
	go checkIfStatusStop()
	wg.Wait()
	return "Stopped music"
}

// checkIfStatusStop function
func checkIfStatusStop() {
	defer wg.Done()
	for {
		if portaudiofunctions.GetStatus() == "stop" {
			return
		}
	}
}

// playMusic function
func playMusic(data structs.Data) string {
	logger.Info("Executing: Play Music")
	logger.Info("Path given " + data.Path)
	var songs []string
	//get Songs
	if len(data.Values) == 0 {
		songs = parseSongs([]string{data.Path}, data.Depth)
	} else {
		songs = parseSongs(data.Values, data.Depth)
	} // end of else
	if len(songs) != 0 {
		serverData.SongQueue = serverData.SongQueue[:0]
		serverData.CurrentSong = 0
		// Append songs to queue
		for _, song := range songs {
			serverData.SongQueue = append(serverData.SongQueue, song)
		} // end of for

		//Check if a song is currently playing
		playCurrentSong()
		return "Playing " + serverData.SongQueue[serverData.CurrentSong]
	} else {
		if len(serverData.SongQueue) != 0 {
			//Check if a song is currently playing
			playCurrentSong()
			return "Playing " + serverData.SongQueue[serverData.CurrentSong]
		} else {
			logger.Error("No input file and no Song in Queue")
			return ("No input file and no Song in Queue")
		} // end foe else
	} // end of if
	return "should never be shown"
} // end of playMusic

// playCurrentSong function
func playCurrentSong() {
	// Check if status is play or pause
	if portaudiofunctions.GetStatus() == "play" || portaudiofunctions.GetStatus() == "pause" {
		logger.Info("A song is currently playing")
		_ = stopMusic()
	} // end of if
	logger.Info(serverData.SongQueue[serverData.CurrentSong])
	go portaudiofunctions.PlayAudio(serverData.SongQueue[serverData.CurrentSong])
} // end of playCurrentSong

// pauseMusic function
func pauseMusic(data structs.Data) string {
	logger.Info("Executing: Pause Music")
	go portaudiofunctions.PauseAudio()
	return "Music paused"
} // end of pauseMusic

// resumeMusic function
func resumeMusic() string {
	logger.Info("Execution: Resume Music")
	go portaudiofunctions.ResumeAudio()
	return "Resuming music"
}

// nextMusic function
func nextMusic(data structs.Data) string {
	// check if loop was set by "playMusic" - if yes..than change data.loop to true
	if serverData.SaveLoop == true { //comment: why here
		data.Loop = true
	}
	//check if nextsong can be played
	if serverData.CurrentSong < (len(serverData.SongQueue) - 1) {
		serverData.CurrentSong += 1
	} else {
		serverData.CurrentSong = 0
	}
	// check if loop is enabled
	if data.Loop == false && serverData.CurrentSong == 0 {
		logger.Info("Loop is not active and queue has ended -> Music stopped")
		return "Loop is not active and queue has ended -> Music stopped"
	} else {
		logger.Info(serverData.SongQueue[serverData.CurrentSong])
		playCurrentSong()
		return "Now playing" + serverData.SongQueue[serverData.CurrentSong]
	}
	return "Should never be shown "
} // end of nextMusic

//addToQueue function
func addToQueue(data structs.Data) string {
	logger.Info("Executing: Add to queue")
	var songs []string
	// get songs
	if len(data.Values) == 0 {
		songs = parseSongs([]string{data.Path}, data.Depth)
	} else {
		songs = parseSongs(data.Values, data.Depth)
	} // end of else
	if len(songs) != 0 {
		//append songs to queue
		for _, song := range songs {
			serverData.SongQueue = append(serverData.SongQueue, song)
		} // end of for
	} // end of if
	message := "Added " + string(len(songs)) + " songs to queue"
	return message
} // end of addToQueue

// printInfo function
func printInfo() string {
	logger.Info("Executing: Print info ")
	message := "\n"
	if len(serverData.SongQueue) != 0 {
		message = message + ("Current Song: " + serverData.SongQueue[serverData.CurrentSong] + "\n")
		if (len(serverData.SongQueue) - 1 - serverData.CurrentSong) != 0 {
			message = message + ("Song Queue: \n")
			//songs from current to end
			for index, song := range serverData.SongQueue[serverData.CurrentSong+1:] {
				message = message + (strconv.Itoa(index+1) + ". " + song + "\n")
			} // enf of for
			// songs from beginning to current
			if serverData.SaveLoop {
				for index, song := range serverData.SongQueue[:serverData.CurrentSong] {
					message = message + (strconv.Itoa(len(serverData.SongQueue)+index-serverData.CurrentSong) + ". " + song + "\n")
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

func main() {
	// set up configuration
	err := gonfig.GetConf("config.json", &configuration)
	// set up logger
	logger.Setup(util.JoinPath(configuration.Log_Dir, configuration.Server_Log), false)
	// create server socket mp.sock
	unixSocket := configuration.Socket_Path
	logger.Notice("Creating unixSocket.")
	logger.Info("Listening on " + unixSocket)
	ln, err := net.Listen("unix", unixSocket)
	if err != nil {
		log.Fatal("listen error", err)
	}
	// check supported formats
	logger.Notice("Parsing supported formats")
	serverData.SupportedFormats = getSupportedFormats()
	// print supported formats
	printSupportedFormats()
	// start pulseAudio
	logger.Notice("Start Pulseaudio")
	portaudiofunctions.StartPulseaudio()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("accept error: ", err)
		}
		go receiveCommand(conn)
	}
} // end of main

// parseSongs
func parseSongs(paths []string, depth int) []string {
	var songs []string
	for _, path := range paths {
		// check if given file/folder exists
		logger.Notice("Check if folder/file exists: " + path)
		// check if path is empty
		if len(path) == 0 {
			logger.Error("Path is not a file or a folder")
			continue
		}
		fi, err := os.Stat(path)
		util.Check(err)
		switch mode := fi.Mode(); {
		case mode.IsDir():
			// directory given
			logger.Info("Directory found")
			logger.Notice("Getting files inside of the folder")
			fileList := util.GetFilesInFolder(path, serverData.SupportedFormats, depth)
			//Print Supported Filelist
			screener.PrintFiles(fileList, false)
			for _, song := range fileList {
				songs = append(songs, song)
			}
		case mode.IsRegular():
			// file given
			logger.Notice("File found")
			var extension = filepath.Ext(path)
			logger.Info("Extension: " + extension)
			if util.StringInArray(extension, serverData.SupportedFormats) {
				logger.Notice("Extension supported")
				songs = append(songs, path)
			} else {
				logger.Warning("Extension not supported")
			}
		default:
			logger.Error("Path is not a file or a folder")
		} // end of switch
	} // end of for
	return songs
} // end of parseSongs

// getSupportedFormats function
func getSupportedFormats() []string {
	// get supported audio formats of 'supportedFormats.cfg' file
	supportedFormats := make([]string, 0)

	// Opening file
	file, err := os.Open("supportedFormats.cfg")
	util.Check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.ContainsAny(line, "#") {
			supportedFormats = append(supportedFormats, line)
		} //end of if
	} // end of for
	util.Check(scanner.Err())
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
