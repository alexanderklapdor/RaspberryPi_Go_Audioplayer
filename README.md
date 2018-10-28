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


