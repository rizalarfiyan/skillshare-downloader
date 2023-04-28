package constants

import "time"

const (
	DefaultLanguage        = "en-US"
	DefaultDir             = "./downloaded"
	DefaultLogFormat       = "[%lvl%]: %time% - %msg% \n"
	DefaultTimestampFormat = time.DateTime
	DefaultWorker          = 8

	MaxWorker = 32

	FolderName        = "[%d] %s"
	FilenameClassData = "class_data.json"
	FilenameVideoData = "%03d_%s_data.json"

	// Credentials SKillshare
	PolicyKey           = "BCpkADawqM2OOcM6njnM7hf9EaK6lIFlqiXB0iWjqGWUQjU7R8965xUvIQNqdQbnDTLz0IAO7E6Ir2rIbXJtFdzrGtitoee0n1XXRliD-RH9A-svuvNW9qgo3Bh34HEZjXjG4Nml4iyz3KqF"
	BrightcoveAccountId = 3695997568001
)
