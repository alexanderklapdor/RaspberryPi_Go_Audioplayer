# GO Audioplayer

This project is a Raspberry API written in GO, which allows to play music from the console on the Raspberry Pi. It is implemented as part of the student research project at the DHBW Stuttgart.


### Commands of the Client
| Command | Description |
| ---: | :--- |
| exit | Terminates the MusicPlayerServer|
| play | Plays the music. |
| stop | Stops the music. |
| pause | Pauses the music |
| resume | Resumes the music |
| next | Skips current song and plays the next one |
| back, previous | Skips current song and plays the previous one |
| add, addToQueue | Adds a song or a directory to the song queue |
| remove, delete, removeAt, deleteAt | Remove a song from song queue |
| setVolume | Sets the volume (value between 0 and 100) |
| quieter, setVolumeDown | Decreases the volume by 10 |
| louder, setVolumeUp | Increases the volume by 10 |
| shuffle, setShuffle | Shuffles the song queue |
| loop, setLoop | Set loop to true or false |
| repeat | Repeat the current song |
| info | Prints information like the current song or the song queue|

### Console arguments

```
Usage of ./MusicPlayerClient:
  -c string
        command for the server (default:"info")
  -d int
    	audio file searching depth (default/recommended 2) (default 2)
  -fi int
    	fadein in milliseconds (default 0)
  -fo int
    	fadeout in milliseconds (default 0)
  -i string
    	input music file/folder
  -l    loop (default false)
  -s	shuffle (default false) -not working-
  -v int
   	music volume in percent (default 50) (default 50)
```

### Logger functions
```
logger.Info("info")
logger.Notice("notice")
logger.Warning("warning")
logger.Error("err")
logger.Critical("crit") -> throw panic
```

### Logger Level
**Info** - Generally useful information to log (service start/stop, configuration assumptions, etc). Info I want to always have available but usually don't care about under normal circumstances. This is my out-of-the-box config level.

**Notice** - Simply a statement that is non-actionable, use these to alert the user of something smaller and passive that you want the use to notice, such as an event that has happened like successful submit.

**Warn** - Anything that can potentially cause application oddities, but for which I am automatically recovering. (Such as switching from a primary to backup server, retrying an operation, missing secondary data, etc.)

**Error** - Any error which is fatal to the operation, but not the service or application (can't open a required file, missing data, etc.). These errors will force user (administrator, or direct user) intervention. These are usually reserved (in my apps) for incorrect connection strings, missing services, etc.

**Critical** - Any error that is forcing a shutdown of the service or application to prevent data loss (or further data loss). I reserve these only for the most heinous errors and situations where there is guaranteed to have been data corruption or loss.

### Used GO-Packages

[Portaudio](https://github.com/gordonklaus/portaudio)
[MP3 Decoder](https://github.com/bobertlo/go-mpg123)
[ID3 Decoder](https://github.com/mikkyang/id3-go)
[Config Parser](https://github.com/tkanos/gonfig)

### Used packages on Raspberry PI
portaudio19-dev
libmpg123-dev
mp3info





# Installation
This chapter contains detailed instructions for installing the Go Musicplayer for Raspberry PI

### Update Raspberry PI
```
sudo apt-get update 
sudo apt-get upgrade 
sudo apt-get install rpi-update 
sudo rpi-update 
```

### Install Dependencis
```
sudo apt-get install portaudio19-dev
sudo apt-get install libmpg123-dev
sudo apt-get install mp3info 
```



## For development purposes

### Installing GO
* Download
>Link: https://golang.org/dl/ 
Kind: Archive 
OS: Linux 
Arch: ARMv6 
Befehl: wget [LINK] 
* Extract Archive
>tar -C /usr/local -xzf [Filename]
* Set export Path
>export PATH=$PATH:/usr/local/go/bin
* Check if go is correct installed
>go version
* Go-Ordnerstruktur erstellen: 
>**bin** will contain all Go executable's you have installed using go install command. 
**pkg** will contain all compiled packages that can be imported into your projects. 
**src** will contain all your source files, either your own or sources downloaded from external repositories. 


>Aufbau: 
>* /home 
>    * /pi 
>       * /go 
>           * /src 
>           * /pkg 
>           * /bin 

>Befehle: 
mkdir go 
cd go 
mkdir src 
mkdir pkg 
mkdir bin 

* Clone Repository
>Starting from /home/pi/go/src/ 
mkdir github.com 
cd github.com 
mkdir alexanderklapdor 
cd alexanderklapdor 
git clone https://github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer.git 

* Install Go Dependencis
> Starting from ../RaspberryPi_Go_Audioplayer/
go get ./...

* Edit Alsa Lib files
>Folgende Datei muss editiert werden damit die Fehlermeldungen der ALSA Lib nicht mehr mit ausgegeben werden. Diese weisen nur darauf hin, dass die Folgenden Anschlüsse an dem Raspberry Pi nicht vorhanden sind. 

>Datei: sudo nano /usr/share/alsa/alsa.conf 

>Muss aus der Datei gelöscht werden: 
pcm.rear cards.pcm.rear 
pcm.center_lfe cards.pcm.center_lfe 
pcm.side cards.pcm.side 
pcm.surround21 cards.pcm.surround21 
pcm.surround40 cards.pcm.surround40 
pcm.surround41 cards.pcm.surround41 
pcm.surround50 cards.pcm.surround50 
pcm.surround51 cards.pcm.surround51 
pcm.surround71 cards.pcm.surround71 
pcm.iec958 cards.pcm.iec958 
pcm.spdif iec958 
pcm.hdmi cards.pcm.hdmi 
pcm.dmix cards.pcm.dmix 
pcm.dsnoop cards.pcm.dsnoop 
pcm.modem cards.pcm.modem 
pcm.phoneline cards.pcm.phoneline 

* You can now execute the MusicPlayer
> go run MusicPlayerClient.go


## For use purposes

* Download latest [Release](https://github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/releases)

> wget [Link]

* Extract Archive
> tar -xvf RaspberryPi_Go_Audioplayer_v***.tar

* Start -> Finish
> ./MusicPlayerClient