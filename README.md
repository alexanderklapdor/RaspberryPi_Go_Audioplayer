# GO_Audioplayer

This project is a Raspberry API written in GO, which allows to play music from the console on the Raspberry Pi. It is implemented as part of the student research project at the DHBW Stuttgart.

### Console arguments

```
Usage of ./main:
  -d int
    	audio file searching depth (default/recommended 2) (default 2)
  -fi int
    	fadein in milliseconds (default 0)
  -fo int
    	fadeout in milliseconds (default 0)
  -i string
    	input music file/folder
  -s	shuffle (default false)
  -v int
   	music volume in percent (default 50) (default 50)
```

### Logger functions
```
logger.Log.Info("info")
logger.Log.Notice("notice")
logger.Log.Warning("warning")
logger.Log.Error("err")
logger.Log.Critical("crit")
```

### Logger Level
**Info** - Generally useful information to log (service start/stop, configuration assumptions, etc). Info I want to always have available but usually don't care about under normal circumstances. This is my out-of-the-box config level.

**Notice** - Simply a statement that is non-actionable, use these to alert the user of something smaller and passive that you want the use to notice, such as an event that has happened like successful submit.

**Warn** - Anything that can potentially cause application oddities, but for which I am automatically recovering. (Such as switching from a primary to backup server, retrying an operation, missing secondary data, etc.)

**Error** - Any error which is fatal to the operation, but not the service or application (can't open a required file, missing data, etc.). These errors will force user (administrator, or direct user) intervention. These are usually reserved (in my apps) for incorrect connection strings, missing services, etc.

**Critical** - Any error that is forcing a shutdown of the service or application to prevent data loss (or further data loss). I reserve these only for the most heinous errors and situations where there is guaranteed to have been data corruption or loss.

### Portaudio

[Portaudio Repository](https://github.com/gordonklaus/portaudio)
Wird benötigt:
https://github.com/gordonklaus/portaudio
https://github.com/bobertlo/go-mpg123

### Installed packages

```
sudo apt-get install libmpg123-dev
sudo apt-get install libasound-dev
sudo apt-get install portaudio19-dev
sudo apt-get install pkg-config
sudo apt-get install xauth
sudo apt-get install jackd2
```

### Start portaudio without X11

```
 if test -z "$DBUS_SESSION_BUS_ADDRESS" ; then
        eval `dbus-launch --sh-syntax`
        echo "D-Bus per-session daemon address is:"
        echo "$DBUS_SESSION_BUS_ADDRESS"
    fi
```
[Link to solution post](https://bugs.debian.org/cgi-bin/bugreport.cgi?bug=690530)

### Play Music works perfect
alsamixer
portaudio
pulseaudio
Raspberry Pi Update!

```
(Raspberry auf Analog Output stellen)
alsamixer (nur zum überprüfen)
pulseaudio -D
go run mp3.go *MP3 Datei*
    -> go run ./main.go -i="/home/pi/Music/testmusic/JNSPTRS - Chasing.mp3" -v=100

```

### Für Fade-IN und Fade-OUT und Volume generell
```
pactl set-sink-mute 0 toggle  # toggle mute
pactl set-sink-volume 0 0     # mute (force)
pactl set-sink-volume 0 100%  # max
pactl set-sink-volume 0 +5%   # +5% (up)
pactl set-sink-volume 0 -5%   # -5% (down)
```
