package audiofunctions

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/portaudiofunctions"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/screener"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/structs"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/util"
)

var wg = &sync.WaitGroup{}

//*******Music Functions********

// playMusic function
func PlayMusic(data structs.Data, serverData *structs.ServerData) string {
	logger.Info("Executing: Play Music")
	logger.Info("Path given " + data.Path)
	var songs []string
	//get Songs
	if len(data.Values) == 0 {
		songs = ParseSongs([]string{data.Path}, data.Depth, serverData)
	} else {
		songs = ParseSongs(data.Values, data.Depth, serverData)
	} // end of else
	if len(songs) != 0 {
		serverData.SongQueue = serverData.SongQueue[:0]
		serverData.CurrentSong = 0
		// Append songs to queue
		for _, song := range songs {
			serverData.SongQueue = append(serverData.SongQueue, song)
		} // end of for

		//Check if a song is currently playing
		PlayCurrentSong(serverData)
		return "Playing " + serverData.SongQueue[serverData.CurrentSong]
	} else {
		if len(serverData.SongQueue) != 0 {
			//Check if a song is currently playing
			PlayCurrentSong(serverData)
			return "Playing " + serverData.SongQueue[serverData.CurrentSong]
		} else {
			logger.Error("No supported input file and no file in queue")
			return ("No supported input file and no file in queue")
		} // end foe else
	} // end of if
	return "should never be shown"
} // end of playMusic

// playCurrentSong function
func PlayCurrentSong(serverData *structs.ServerData) {
	// Check if status is play or pause
	if portaudiofunctions.GetStatus() == "play" || portaudiofunctions.GetStatus() == "pause" {
		logger.Info("A song is currently playing")
		_ = StopMusic()
	} // end of if
	logger.Info(serverData.SongQueue[serverData.CurrentSong])
	go portaudiofunctions.PlayAudio(serverData.SongQueue[serverData.CurrentSong])
} // end of playCurrentSong

// stopMusic function
func StopMusic() string {
	logger.Info("Execution: Stop Music")
	portaudiofunctions.StopAudio()

	wg.Add(1)
	go CheckIfStatusStop()
	wg.Wait()
	return "Stopped music"
}

// playPreviousSong function
func PlayPreviousSong(serverData *structs.ServerData) string {
	//check if currentSong is > 0 in queue
	if serverData.CurrentSong > 0 {
		// decrease currentSong and play
		serverData.CurrentSong--
		PlayCurrentSong(serverData)
		return "Playing now " + serverData.SongQueue[serverData.CurrentSong]
	} else {
		// check if loop is enabled
		if serverData.SaveLoop {
			if len(serverData.SongQueue) > 0 {
				serverData.CurrentSong = len(serverData.SongQueue) - 1
				PlayCurrentSong(serverData)
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

// pauseMusic function
func PauseMusic(data structs.Data, serverData *structs.ServerData) string {
	logger.Info("Executing: Pause Music")
	go portaudiofunctions.PauseAudio()
	return "Music paused"
} // end of pauseMusic

// resumeMusic function
func ResumeMusic() string {
	logger.Info("Execution: Resume Music")
	go portaudiofunctions.ResumeAudio()
	return "Resuming music"
}

// repeatSong function
func RepeatSong(serverData *structs.ServerData) string {
	// check sonqQueue length
	if len(serverData.SongQueue) > 0 {
		PlayCurrentSong(serverData)
		return "Playing now " + serverData.SongQueue[serverData.CurrentSong]
	} else {
		return "There is no current song"
	} // end of else
} // end of repeatSong

// nextMusic function
func NextMusic(data structs.Data, serverData *structs.ServerData) string {
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
		logger.Info("Now playing" + serverData.SongQueue[serverData.CurrentSong])
		PlayCurrentSong(serverData)
		return "Now playing" + serverData.SongQueue[serverData.CurrentSong]
	}
	return "Should never be shown "
} // end of nextMusic

//*******Help Functions********

// checkIfStatusStop function
func CheckIfStatusStop() {
	defer wg.Done()
	for {
		if portaudiofunctions.GetStatus() == "stop" {
			return
		}
	}
}

// parseSongs
func ParseSongs(paths []string, depth int, serverData *structs.ServerData) []string {
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
		if err != nil {
			logger.Error("Path is not a file or a folder")
			continue
		}
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

//*******Queue Functions********

//addToQueue function
func AddToQueue(data structs.Data, serverData *structs.ServerData) string {
	logger.Info("Executing: Add to queue")
	var songs []string
	// get songs
	if len(data.Values) == 0 {
		songs = ParseSongs([]string{data.Path}, data.Depth, serverData)
	} else {
		songs = ParseSongs(data.Values, data.Depth, serverData)
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

// removeSong function
func RemoveSong(data structs.Data, serverData *structs.ServerData) string {
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
