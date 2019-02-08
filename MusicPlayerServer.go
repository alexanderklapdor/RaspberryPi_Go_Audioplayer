package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/audiofunctions"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/screener"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/util"
)

var wg = &sync.WaitGroup{}

// Request struct
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
	Volume  int
}

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
	logger.Log.Info("Server received message: " + string(receivedBytes))

	// convert message back to a request-object
	logger.Log.Notice("Converting message back to a Request-Object")
	received := Request{}
	json.Unmarshal(receivedBytes, &received)
	command := received.Command
	data := received.Data
	logger.Log.Notice("Command: " + command)
	//logger.Log.Notice("Data   : " + string(data))

	// switch case commands
	switch command {
	case "exit":
		closeConnection(c)
	case "play":
		playMusic(data)
	case "stop":
		stopMusic()
	case "pause":
		pauseMusic(data)
	case "resume":
		resumeMusic()
	case "next":
		nextMusic(data)
	case "setVolume":
		setVolume(data)
	case "addToQueue":
		addToQueue(data)
	case "quieter":
		decreaseVolume()
	case "louder":
		increaseVolume()
	case "info":
		printInfo()
	default:
		logger.Log.Error("Unknown command received")
	}

	// write to client
	logger.Log.Notice("Send a message back to the client")
	message := "Default-message"
	_, err = c.Write([]byte(message))
	if err != nil {
		log.Fatal("Write: ", err)
	}
} // end of receiveCommand

func closeConnection(c net.Conn) {
	socketPath := "/tmp/mp.sock" // todo: should be passed as an argument or be written out of a config file
	logger.Log.Warning("Connection  will be closed")
	defer c.Close()
	err := syscall.Unlink(socketPath)
	if err != nil {
		logger.Log.Error("Error during unlink process of the socket: " + err.Error())
		logger.Log.Info("Pls run manually unlink 'unlink" + socketPath + "'")
	}
	os.Exit(0)
} // end of closeConnection

func stopMusic() {
	logger.Log.Info("Execution: Stop Music")
	audiofunctions.StopAudio()

	wg.Add(1)
	go checkIfStatusStop()
	wg.Wait()
}

func checkIfStatusStop() {
	defer wg.Done()
	for {
		if audiofunctions.GetStatus() == "stop" {
			return
		}
	}
}

func playMusic(data Data) {
	logger.Log.Info("Executing: Play Music")
	logger.Log.Info("Path given " + data.Path)
	songs := parseSongs(data.Path, supportedFormats, data.Depth)
	if len(songs) != 0 {
		songQueue = songQueue[:0]
		currentSong = 0
		for _, song := range songs {
			songQueue = append(songQueue, song)
		} // end of for

		//Check if a song is currently playing
		if audiofunctions.GetStatus() == "play" || audiofunctions.GetStatus() == "pause" {
			logger.Log.Info("A song is currently playing")
			stopMusic()
		}
		logger.Log.Info(songQueue[currentSong])
		go audiofunctions.PlayAudio(songQueue[currentSong])
		// set loop variable
		saveLoop = data.Loop
	} else {
		if len(songQueue) != 0 {

			//Check if a song is currently playing
			if audiofunctions.GetStatus() == "play" || audiofunctions.GetStatus() == "pause" {
				logger.Log.Info("A song is currently playing")
				stopMusic()
			}

			logger.Log.Info(songQueue[currentSong])
			go audiofunctions.PlayAudio(songQueue[currentSong])
		} else {
			logger.Log.Error("No input file and no Song in Queue")
		}
	} // end of if
} // end of playMusic

func pauseMusic(data Data) {
	logger.Log.Info("Executing: Pause Music")
	go audiofunctions.PauseAudio()
} // end of pauseMusic

func resumeMusic() {
	logger.Log.Info("Execution: Resume Music")
	go audiofunctions.ResumeAudio()
}

func nextMusic(data Data) {
	//check if loop was set by "playMusic" - if yes..than change data.loop to true
	if saveLoop == true {
		data.Loop = true
	}
	if currentSong < (len(songQueue) - 1) {
		currentSong += 1
	} else {
		currentSong = 0
	}
	//Check if a song is currently playing
	if audiofunctions.GetStatus() == "play" || audiofunctions.GetStatus() == "pause" {
		logger.Log.Info("A song is currently playing")
		stopMusic()
	}
	if data.Loop == false && currentSong == 0 {
		logger.Log.Info("Loop is not activate and Queue has ended -> Music stops")
	} else {
		logger.Log.Info(songQueue[currentSong])
		go audiofunctions.PlayAudio(songQueue[currentSong])
	}

}

func setVolume(data Data) {
	logger.Log.Info("Executing: Set Volume")
} // end of setVolume

func addToQueue(data Data) {
	logger.Log.Info("Executing: Add to queue")
	songs := parseSongs(data.Path, supportedFormats, data.Depth)
	if len(songs) != 0 {
		for _, song := range songs {
			songQueue = append(songQueue, song)
		} // end of for
	} // end of if
} // end of addToQueue

func increaseVolume() {
	logger.Log.Info("Executing: Increase volume")
	go audiofunctions.SetVolumeUp()
} // end of increaseVolume

func decreaseVolume() {
	logger.Log.Info("Executing: Decrease volume")
	go audiofunctions.SetVolumeDown()
} // end of decreaseVolume

func printInfo() {
	logger.Log.Info("Executing: Print info ")
} // end of printInfo

func main() {
	// create server socket mp.sock
	unixSocket := "/tmp/mp.sock"
	logger.Log.Notice("Creating unixSocket.")
	logger.Log.Info("Listening on " + unixSocket)
	ln, err := net.Listen("unix", unixSocket)
	if err != nil {
		log.Fatal("listen error", err)
	}

	// check supported formats
	logger.Log.Notice("Parsing supported formats")
	supportedFormats = getSupportedFormats()
	// print supported formats
	printSupportedFormats(supportedFormats)

	// start pulseAudio
	logger.Log.Notice("start PulseAudio")
	audiofunctions.StartPulseaudio()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("accept error: ", err)
		}
		go receiveCommand(conn)
	}
} // end of main

func parseSongs(path string, supportedFormats []string, depth int) []string {
	// check if given file/folder exists
	logger.Log.Notice("Check if folder/file exists: ", path)

	// check if path is empty
	if len(path) == 0 {
		logger.Log.Error("Path is not a file or a folder")
		var songs []string
		return songs
	}

	fi, err := os.Stat(path)
	util.Check(err)

	switch mode := fi.Mode(); {
	case mode.IsDir():
		// directory given
		logger.Log.Info("Directory found")
		logger.Log.Notice("Getting files inside of the folder")
		fileList := getFilesInFolder(path, supportedFormats, depth)
		//Print Supported Filelist
		screener.PrintFiles(fileList, false)
		return fileList
	case mode.IsRegular():
		// file given
		logger.Log.Notice("File found")
		var extension = filepath.Ext(path)
		logger.Log.Info("Extension: " + extension)
		if util.StringInArray(extension, supportedFormats) {
			logger.Log.Notice("Extension supported")
			var songs []string
			return append(songs, path)
		} else {
			logger.Log.Warning("Extension not supported")
			var songs []string
			return songs
		}
	default:
		logger.Log.Error("Path is not a file or a folder")
		var songs []string
		return songs
	} // end of switch
} // end of parseSongs

func getFilesInFolder(folder string, supportedExtensions []string, depth int) []string {
	// fmt.Println("get files in ", folder)
	fileList := make([]string, 0)
	if depth > 0 {
		files, err := ioutil.ReadDir(folder)
		util.Check(err)
		for _, file := range files {
			filename := joinPath(folder, file.Name())

			fi, err := os.Stat(filename)
			util.Check(err)

			switch mode := fi.Mode(); {
			case mode.IsDir():
				newFolder := filename + "/"
				newFiles := getFilesInFolder(newFolder, supportedExtensions, depth-1)
				for _, newFile := range newFiles {
					fileList = append(fileList, newFile)
				} // end of for
			case mode.IsRegular():
				var extension = filepath.Ext(filename)
				if util.StringInArray(extension, supportedExtensions) {
					fileList = append(fileList, filename)
				} // end of if
			} // end of switch
		} // end of for
	} else {
		logger.Log.Info("Max depth reached")
	}
	return fileList
} // end of getFilesInFolder

func joinPath(source, target string) string {
	if path.IsAbs(target) {
		return target
	} // end of if
	return path.Join(path.Dir(source), target)
} // end of JoinPath

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
	logger.Log.Info("Supported formats: " + formatString)
} // end of printSupportedFormats
