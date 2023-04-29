package services

type Downloader interface {
	Download(fileURL, filename string) (<-chan FileInfo, error)
}
