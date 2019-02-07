package main

import (
        "log"
        "net"
        "encoding/json"
        "os"
        "syscall"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"
)

type Request struct {
    Command string
    Data Data
}

type Data struct {
    Depth   int
    FadeIn  int
    FadeOut int
    Path    string
    Shuffle bool
    Volume  int
}

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
    case "pause":
        pauseMusic(data)
    case "setVolume":
        setVolume(data)
    case "addToQueue":
        addToQueue(data)
    case "quieter":
        increaseVolume()
    case "louder":
        decreaseVolume()
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
    err :=  syscall.Unlink(socketPath)
    if err != nil {
        logger.Log.Error("Error during unlink process of the socket: " + err.Error())
        logger.Log.Info("Pls run manually unlink 'unlink" + socketPath + "'")
    }
    os.Exit(0)
} // end of closeConnection

func playMusic(data Data) {
    logger.Log.Info("Executing: Play Music")
}


func pauseMusic(data Data) {
    logger.Log.Info("Executing: Pause Music")
}

func setVolume(data Data) {
    logger.Log.Info("Executing: Set Volume")
}

func addToQueue(data Data) {
    logger.Log.Info("Executing: Add to queue")
}

func increaseVolume() {
    logger.Log.Info("Executing: Increase volume")
}

func decreaseVolume() {
    logger.Log.Info("Executing: Decrease volume")
}

func printInfo() {
    logger.Log.Info("Executing: Print info ")
}

func main() {
    unixSocket := "/tmp/mp.sock"
    // create server socket mp.sock
    logger.Log.Notice("Creating unixSocket.")
    logger.Log.Info("Listening on " + unixSocket)
    ln, err := net.Listen("unix", unixSocket)
    if err != nil {
        log.Fatal("listen error", err)
    }

    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Fatal("accept error: ", err)
        }
        go receiveCommand(conn)
    }
} // end of main

func parseSongs(path string, depth int) []string {
    // check if given file/folder exists
    logger.Log.Notice("Check if folder/file exists", path)
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
        if util.StringInArray(extension, supportedFormats) {
            logger.Log.Notice("Extension supported")
            return path
        } else {
            logger.Log.Warning("Extension not supported")
            return []
        }
    default:
        logger.Log.Error("Path is not a file or a folder")
        return []
    } // end of switch
}

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
				}
			case mode.IsRegular():
				var extension = filepath.Ext(filename)
				if util.StringInArray(extension, supportedExtensions) {
					fileList = append(fileList, filename)
				}
			}
		}
	} else {
		//fmt.Println("Max depth reached")
	}
	return fileList
}

func joinPath(source, target string) string {
	if path.IsAbs(target) {
		return target
	}
	return path.Join(path.Dir(source), target)
}

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
		}
	}

	util.Check(scanner.Err())
	return supportedFormats

}
