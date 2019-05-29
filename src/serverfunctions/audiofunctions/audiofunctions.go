package audiofunctions

// Imports
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

// PlayMusic function
func PlayMusic(data structs.Data, serverData *structs.ServerData) string {
	logger.Info("Executing: Play Music")
	logger.Info("Path given " + data.Path)
	var songs []string
	// Get Songs
	if len(data.Values) == 0 {
		songs = ParseSongs([]string{data.Path}, data.Depth, serverData)
	} else {
		songs = ParseSongs(data.Values, data.Depth, serverData)
	} // End of else
	if len(songs) != 0 {
		serverData.SongQueue = serverData.SongQueue[:0]
		serverData.CurrentSong = 0
		// Append songs to queue
		for _, song := range songs {
			serverData.SongQueue = append(serverData.SongQueue, song)
		} // End of for

		// Check if a song is currently playing
		PlayCurrentSong(serverData)
		return "Now playing " + util.PrintMp3Infos(serverData.SongQueue[serverData.CurrentSong])
	} else {
		if len(serverData.SongQueue) != 0 {
			// Check if a song is currently playing
			PlayCurrentSong(serverData)
			return "Now playing " + util.PrintMp3Infos(serverData.SongQueue[serverData.CurrentSong])
		} else {
			logger.Error("No supported input file and no file in queue")
			return ("No supported input file and no file in queue")
		} // End foe else
	} // End of if
	return "should never be shown"
} // End of playMusic

// PlayCurrentSong function
func PlayCurrentSong(serverData *structs.ServerData) {
	// Check if status is play or pause
	if portaudiofunctions.GetStatus() == "play" || portaudiofunctions.GetStatus() == "pause" {
		logger.Info("A song is currently playing")
		_ = StopMusic()
	} // End of if
	logger.Info(serverData.SongQueue[serverData.CurrentSong])
	go portaudiofunctions.PlayAudio(serverData.SongQueue[serverData.CurrentSong])
} // End of playCurrentSong

// StopMusic function
func StopMusic() string {
	logger.Info("Execution: Stop Music")
	portaudiofunctions.StopAudio()

	wg.Add(1)
	go CheckIfStatusStop()
	wg.Wait()
	return "Stopped music"
}

// PlayPreviousSong function
func PlayPreviousSong(serverData *structs.ServerData) string {
	// Check if currentSong is > 0 in queue
	if serverData.CurrentSong > 0 {
		// Decrease currentSong and play
		serverData.CurrentSong--
		PlayCurrentSong(serverData)
		return "Playing now " + util.PrintMp3Infos(serverData.SongQueue[serverData.CurrentSong])
	} else {
		// Check if loop is enabled
		if serverData.SaveLoop {
			if len(serverData.SongQueue) > 0 {
				serverData.CurrentSong = len(serverData.SongQueue) - 1
				PlayCurrentSong(serverData)
				return "Playing now " + util.PrintMp3Infos(serverData.SongQueue[serverData.CurrentSong])
			} else {
				return "Error: The queue is empty. You could't go a song back"
			} // End of else
		} else {
			return "You are currently playing the first song"
		} // End of else
	} // End of else
	return "should never be shown"
} // End of playLastSong

// PauseMusic function
func PauseMusic(data structs.Data, serverData *structs.ServerData) string {
	logger.Info("Executing: Pause Music")
	go portaudiofunctions.PauseAudio()
	return "Music paused"
} // End of pauseMusic

// ResumeMusic function
func ResumeMusic() string {
	logger.Info("Execution: Resume Music")
	go portaudiofunctions.ResumeAudio()
	return "Resuming music"
}

// RepeatSong function
func RepeatSong(serverData *structs.ServerData) string {
	// Check sonqQueue length
	if len(serverData.SongQueue) > 0 {
		PlayCurrentSong(serverData)
		return "Playing now " + util.PrintMp3Infos(serverData.SongQueue[serverData.CurrentSong])
	} else {
		return "There is no current song"
	} // End of else
} // End of repeatSong

// NextMusic function
func NextMusic(data structs.Data, serverData *structs.ServerData) string {
	// Check if loop was set by "playMusic" - if yes..than change data.loop to true
	if serverData.SaveLoop == true { //comment: why here
		data.Loop = true
	}
	// Check if nextsong can be played
	if serverData.CurrentSong < (len(serverData.SongQueue) - 1) {
		serverData.CurrentSong += 1
	} else {
		serverData.CurrentSong = 0
	}
	// Check if loop is enabled
	if data.Loop == false && serverData.CurrentSong == 0 {
		logger.Info("Loop is not active and queue has ended -> Music stopped")
		return "Loop is not active and queue has ended -> Music stopped"
	} else {
		logger.Info("Playing now" + serverData.SongQueue[serverData.CurrentSong])
		PlayCurrentSong(serverData)
		return "Playing now " + util.PrintMp3Infos(serverData.SongQueue[serverData.CurrentSong])
	}
	return "Should never be shown "
} // End of nextMusic

//*******Help Functions********

// CheckIfStatusStop function
func CheckIfStatusStop() {
	defer wg.Done()
	for {
		if portaudiofunctions.GetStatus() == "stop" {
			return
		}
	}
}

// ParseSongs
func ParseSongs(paths []string, depth int, serverData *structs.ServerData) []string {
	var songs []string
	for _, path := range paths {
		// Check if given file/folder exists
		logger.Notice("Check if folder/file exists: " + path)
		// Check if path is empty
		if len(path) == 0 {
			logger.Error("Path is not a file or a folder")
			continue
		}
		fi, err := os.Stat(path)
		if err != nil {
			logger.Error("Path is not a file or a folder")
			continue
		}
		util.Check(err, "Server")
		switch mode := fi.Mode(); {
		case mode.IsDir():
			// Directory given
			logger.Info("Directory found")
			logger.Notice("Getting files inside of the folder")
			fileList := util.GetFilesInFolder(path, serverData.SupportedFormats, depth)
			// Print Supported Filelist
			screener.PrintFiles(fileList, false)
			for _, song := range fileList {
				songs = append(songs, song)
			}
		case mode.IsRegular():
			// File given
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
		} // End of switch
	} // End of for
	return songs
} // End of parseSongs

//*******Queue Functions********

// AddToQueue function
func AddToQueue(data structs.Data, serverData *structs.ServerData) string {
	logger.Info("Executing: Add to queue")
	var songs []string
	// Get songs
	if len(data.Values) == 0 {
		songs = ParseSongs([]string{data.Path}, data.Depth, serverData)
	} else {
		songs = ParseSongs(data.Values, data.Depth, serverData)
	} // End of else
	if len(songs) != 0 {
		// Append songs to queue
		for _, song := range songs {
			serverData.SongQueue = append(serverData.SongQueue, song)
		} // End of for
	} // End of if
	message := "Added " + string(len(songs)) + " songs to queue"
	return message
} // End of addToQueue

// RemoveSong function
func RemoveSong(data structs.Data, serverData *structs.ServerData) string {
	if len(data.Values) != 0 {
		number, err := strconv.Atoi(data.Values[0])
		if err == nil {
			// Check loop is off and song in queue
			if number > 0 && number < (len(serverData.SongQueue)-serverData.CurrentSong) {
				number = number + serverData.CurrentSong
				song_name := serverData.SongQueue[number]
				serverData.SongQueue = append(serverData.SongQueue[:number], serverData.SongQueue[number+1:]...)
				return "Removed song" + util.PrintMp3Infos(song_name)
				// Check loop is on and song is queue
			} else if number >= len(serverData.SongQueue)-serverData.CurrentSong && number < len(serverData.SongQueue) && serverData.SaveLoop {
				number = serverData.CurrentSong - len(serverData.SongQueue) + number
				song_name := serverData.SongQueue[number]
				serverData.SongQueue = append(serverData.SongQueue[:number], serverData.SongQueue[number+1:]...)
				return "Removed song" + util.PrintMp3Infos(song_name)
			} else {
				return "There is no song with the given number (" + strconv.Itoa(number) + ")"
			} // End of else
		} else {
			return "Remove is only allowed with the number of the song in the queue. Pls use 'info' to see the queue"
		} // Else
	} else {
		return "No argument given"
	}
	return "should never be shown"
} // End of removeSong
