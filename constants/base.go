package constants

import (
	"runtime"
	"time"
)

const (
	DefaultLanguage        = "en-US"
	DefaultDir             = "./downloaded"
	DefaultLogFormat       = "[%lvl%]: %time% - %msg% \n"
	DefaultTimestampFormat = time.DateTime

	FolderName          = "[%d] %s"
	FilenameClassData   = "class_data.json"
	FilenameVideoData   = "%03d_%s_data.json"
	FilenameVideo       = "%03d_%s%s"
	FilenameSubtitle    = "%03d_%s%s"
	ProgressBarTemplate = `{{counters .}} - {{ bar . "[" "=" (cycle . ">" ) "-" "]"}} {{percent .}} {{speed .}}`

	// Credentials SKillshare
	PolicyKey = "BCpkADawqM2OOcM6njnM7hf9EaK6lIFlqiXB0iWjqGWUQjU7R8965xUvIQNqdQbnDTLz0IAO7E6Ir2rIbXJtFdzrGtitoee0n1XXRliD-RH9A-svuvNW9qgo3Bh34HEZjXjG4Nml4iyz3KqF"
)

var (
	MaxWorker     int = runtime.NumCPU()
	DefaultWorker int = MaxWorker

	// Credentials SKillshare
	BrightcoveAccountId int64 = 3695997568001
)

// Splash Screen
const SplashScreen = `
    +=----------------------------------------------------------------------------------------------------------=+
    
    .d88b. 8    w 8 8      8                        888b.                        8               8            
    YPwww. 8.dP w 8 8 d88b 8d8b. .d88 8d8b .d88b    8   8 .d8b. Yb  db  dP 8d8b. 8 .d8b. .d88 .d88 .d88b 8d8b 
        d8 88b  8 8 8 ` + "`" + `Yb. 8P Y8 8  8 8P   8.dP'    8   8 8' .8  YbdPYbdP  8P Y8 8 8' .8 8  8 8  8 8.dP' 8P   
    ` + "`" + `Y88P' 8 Yb 8 8 8 Y88P 8   8 ` + "`" + `Y88 8    ` + "`" + `Y88P    888P' ` + "`" + `Y8P'   YP  YP   8   8 8 ` + "`" + `Y8P' ` + "`" + `Y88 ` + "`" + `Y88 ` + "`" + `Y88P 8
    
    +=-------------------------------------------- By Rizal Arfiyan --------------------------------------------=+
`
