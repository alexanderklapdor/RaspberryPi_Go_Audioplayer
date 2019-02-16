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
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/structs"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/util"
	"github.com/tkanos/gonfig"
)

// Global var definition
var wg = &sync.WaitGroup{}
var configuration = structs.ServerConfiguration{}
var supportedFormats []string
var songQueue []string
var currentSong int = 0
var saveLoop bool

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
		message = increaseVolume()
	case "next":
		message = nextMusic(data)
	case "pause":
		message = pauseMusic(data)
	case "play":
		message = playMusic(data)
	case "quieter", "setVolumeDown":
		message = decreaseVolume()
	case "repeat":
		message = repeatSong()
	case "remove", "delete", "removeAt", "deleteAt":
		message = removeSong(data)
	case "resume":
		message = resumeMusic()
	case "setUp":
		message = setUpMusicPlayer(data)
	case "setVolume":
		message = setVolume(data)
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

func setUpMusicPlayer(data structs.Data) string {
	portaudiofunctions.SetVolume(strconv.Itoa(data.Volume))
	addToQueue(data)
	setLoop(data)
	if data.Shuffle {
		shuffleQueue()
	} // end of if
	return "Set up Music Player" + printInfo()
}

func setLoop(data structs.Data) string {
	if len(data.Values) > 0 {
		value_string := data.Values[0]
		if strings.Contains(value_string, "on") || strings.Contains(value_string, "true") {
			saveLoop = true
		} else if strings.Contains(value_string, "off") || strings.Contains(value_string, "false") {
			saveLoop = false
		} // end of else
	} else {
		saveLoop = data.Loop
	} // end of else
	return "Set loop to " + strconv.FormatBool(saveLoop)
} // end of setLoop

func closeConnection(c net.Conn) {
	socketPath := configuration.Socket_Path
	logger.Warning("Connection  will be closed")
	defer c.Close()
	err := syscall.Unlink(socketPath)
	if err != nil {
		logger.Error("Error during unlink process of the socket: " + err.Error())
		logger.Info("Pls run manually unlink 'unlink" + socketPath + "'")
		os.Exit(69)
	}
	os.Exit(0)
} // end of closeConnection

func playPreviousSong() string {
	if currentSong > 0 {
		currentSong--
		playCurrentSong()
		return "Playing now " + songQueue[currentSong]
	} else {
		if saveLoop {
			if len(songQueue) > 0 {
				currentSong = len(songQueue) - 1
				playCurrentSong()
				return "Playing now " + songQueue[currentSong]
			} else {
				return "Error: The queue is empty. You could't go a song back"
			} // end of else
		} else {
			return "You are currently playing the first song"
		} // end of else
	} // end of else
	return "should never be shown"
} // end of playLastSong

func repeatSong() string {
	if len(songQueue) > 0 {
		playCurrentSong()
		return "Playing now " + songQueue[currentSong]
	} else {
		return "There is no current song"
	} // end of else
} // end of repeatSong

func removeSong(data structs.Data) string {
	if len(data.Values) != 0 {
		number, err := strconv.Atoi(data.Values[0])
		// todo: remove multiple values (problem with changing position)
		if err == nil {
			if number > 0 && number < (len(songQueue)-currentSong) {
				number = number + currentSong
				song_name := songQueue[number]
				songQueue = append(songQueue[:number], songQueue[number+1:]...)
				return "Removed song" + song_name
			} else if number >= len(songQueue)-currentSong && number < len(songQueue) && saveLoop {
				number = currentSong - len(songQueue) + number
				song_name := songQueue[number]
				songQueue = append(songQueue[:number], songQueue[number+1:]...)
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

func stopMusic() string {
	logger.Info("Execution: Stop Music")
	portaudiofunctions.StopAudio()

	wg.Add(1)
	go checkIfStatusStop()
	wg.Wait()
	return "Stopped music"
}

func checkIfStatusStop() {
	defer wg.Done()
	for {
		if portaudiofunctions.GetStatus() == "stop" {
			return
		}
	}
}

func playMusic(data structs.Data) string {
	logger.Info("Executing: Play Music")
	logger.Info("Path given " + data.Path)
	var songs []string
	if len(data.Values) == 0 {
		songs = parseSongs([]string{data.Path}, supportedFormats, data.Depth)
	} else {
		songs = parseSongs(data.Values, supportedFormats, data.Depth)
	} // end of else
	if len(songs) != 0 {
		songQueue = songQueue[:0]
		currentSong = 0
		for _, song := range songs {
			songQueue = append(songQueue, song)
		} // end of for

		//Check if a song is currently playing
		playCurrentSong()
		return "Playing " + songQueue[currentSong]
	} else {
		if len(songQueue) != 0 {
			//Check if a song is currently playing
			playCurrentSong()
			return "Playing " + songQueue[currentSong]
		} else {
			logger.Error("No input file and no Song in Queue")
			return ("No input file and no Song in Queue")
		} // end foe else
	} // end of if
	return "should never be shown"
} // end of playMusic

func playCurrentSong() {
	if portaudiofunctions.GetStatus() == "play" || portaudiofunctions.GetStatus() == "pause" {
		logger.Info("A song is currently playing")
		_ = stopMusic()
	} // end of if
	logger.Info(songQueue[currentSong])
	go portaudiofunctions.PlayAudio(songQueue[currentSong])
} // end of playCurrentSong

func pauseMusic(data structs.Data) string {
	logger.Info("Executing: Pause Music")
	go portaudiofunctions.PauseAudio()
	return "Music paused"
} // end of pauseMusic

func resumeMusic() string {
	logger.Info("Execution: Resume Music")
	go portaudiofunctions.ResumeAudio()
	return "Resuming music"
}

func nextMusic(data structs.Data) string {
	// check if loop was set by "playMusic" - if yes..than change data.loop to true
	if saveLoop == true { //comment: why here
		data.Loop = true
	}
	if currentSong < (len(songQueue) - 1) {
		currentSong += 1
	} else {
		currentSong = 0
	}
	if data.Loop == false && currentSong == 0 {
		logger.Info("Loop is not active and queue has ended -> Music stopped")
		return "Loop is not active and queue has ended -> Music stopped"
	} else {
		logger.Info(songQueue[currentSong])
		playCurrentSong()
		return "Now playing" + songQueue[currentSong]
	}
	return "Should never be shown "
} // end of nextMusic

func setVolume(data structs.Data) string {
	var volume string
	if len(data.Values) != 0 {
		volume = data.Values[0]
	} else {
		volume = strconv.Itoa(data.Volume)
	} // end of else
	portaudiofunctions.SetVolume(volume)
	logger.Info("Executing: Set volume to " + volume)
	return "Set volume to " + volume
} // end of setVolume

func addToQueue(data structs.Data) string {
	logger.Info("Executing: Add to queue")
	var songs []string
	if len(data.Values) == 0 {
		songs = parseSongs([]string{data.Path}, supportedFormats, data.Depth)
	} else {
		songs = parseSongs(data.Values, supportedFormats, data.Depth)
	} // end of else
	if len(songs) != 0 {
		for _, song := range songs {
			songQueue = append(songQueue, song)
		} // end of for
	} // end of if
	message := "Added " + string(len(songs)) + " songs to queue"
	return message
} // end of addToQueue

func increaseVolume() string {
	logger.Info("Executing: Increase volume")
	portaudiofunctions.SetVolumeUp("10")
	return "Increased volume by 10 \n" + printVolume()
} // end of increaseVolume

func decreaseVolume() string {
	logger.Info("Executing: Decrease volume")
	portaudiofunctions.SetVolumeDown("10")
	return "Decreased volume by 10 \n" + printVolume()
} // end of decreaseVolume

func getVolume() (string, string) {
	left, right := portaudiofunctions.GetVolume()
	return left, right
} // end of getVolume()

func printVolume() string {
	left, right := getVolume()
	return "Current Volume:  Left(" + left + ")  Right(" + right + ")"
}

func printInfo() string {
	logger.Info("Executing: Print info ")
	message := "\n"
	if len(songQueue) != 0 {
		message = message + ("Current Song: " + songQueue[currentSong] + "\n")
		if (len(songQueue) - 1 - currentSong) != 0 {
			message = message + ("Song Queue: \n")
			//songs from current to end
			for index, song := range songQueue[currentSong+1:] {
				message = message + (strconv.Itoa(index+1) + ". " + song + "\n")
			} // enf of for
			// songs from beginning to current
			if saveLoop {
				for index, song := range songQueue[:currentSong] {
					message = message + (strconv.Itoa(len(songQueue)+index-currentSong) + ". " + song + "\n")
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
	message = message + printVolume()
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
	supportedFormats = getSupportedFormats()
	// print supported formats
	printSupportedFormats(supportedFormats)

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

func parseSongs(paths []string, supportedFormats []string, depth int) []string {
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
			fileList := util.GetFilesInFolder(path, supportedFormats, depth)
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
			if util.StringInArray(extension, supportedFormats) {
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
		//fmt.Println(line)
		if !strings.ContainsAny(line, "#") {
			supportedFormats = append(supportedFormats, line)
			//fmt.Println("format", line)
		} //end of if
	} // end of for
	util.Check(scanner.Err())
	return supportedFormats
} // End of getSupportedFormats

func printSupportedFormats(supportedFormats []string) {
	formatString := ""
	for _, format := range supportedFormats {
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
	if len(songQueue) > 0 {
		songQueue = util.Shuffle(songQueue)
		return "Queue has been shuffled"
	} else {
		return "Queue is not filled - shuffle failed"
	}
}
