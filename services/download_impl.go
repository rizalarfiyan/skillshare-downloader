/*
   Original Code: https://github.com/gojunkie/goget
   Modified by: Rizal Arfiyan
   Description: Download file from url with multiple connection
*/

package services

import (
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/rizalarfiyan/skillshare-downloader/utils"
)

type FileInfo struct {
	Err       error
	ChunkSize int64
	TotalSize int64
}

func NewDownloader(maxConn int) Downloader {
	return &downloader{maxConn: maxConn}
}

type downloader struct {
	maxConn int
}

type chunkInfo struct {
	totalSize int64
	ranges    []int64
	index     int
	fileURL   string
	filename  string
}

func (s *downloader) Download(fileURL, filename string) (<-chan FileInfo, error) {
	totalSize, err := s.getFileSize(fileURL)
	if err != nil {
		return nil, err
	}

	ch := make(chan FileInfo)
	out := make(chan FileInfo)
	ranges := s.calcRanges(totalSize)
	w := utils.NewWorkerPool(s.maxConn)

	var wg sync.WaitGroup
	wg.Add(len(ranges))

	for i := 0; i < len(ranges); i++ {
		info := chunkInfo{
			totalSize: totalSize,
			ranges:    ranges[i],
			index:     i + 1,
			fileURL:   fileURL,
			filename:  filename,
		}
		go func() {
			w.RunFunc(func() {
				s.download(info, ch)
			})
			wg.Done()
		}()
	}

	go func() {
		for p := range ch {
			if p.Err != nil {
				close(ch)
				close(out)
				w.Stop()
				return
			}
			out <- p
		}
	}()

	go func() {
		wg.Wait()
		w.Stop()

		if err := s.combineFiles(filename, len(ranges)); err != nil {
			out <- FileInfo{Err: err}
		}

		close(ch)
		close(out)
	}()

	return out, nil
}

func (s *downloader) combineFiles(filename string, n int) error {
	os.Remove(filename)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for i := 0; i < n; i++ {
		chunkName := fmt.Sprintf("%s.part%.3d", filename, i+1)
		cf, err := os.Open(chunkName)
		if err != nil {
			return err
		}

		buf, err := io.ReadAll(cf)
		if err != nil {
			return err
		}

		if _, err := f.Write(buf); err != nil {
			return err
		}

		cf.Close()
		os.Remove(chunkName)
	}

	return nil
}

func (s *downloader) download(info chunkInfo, ch chan FileInfo) {
	req, err := http.NewRequest("GET", info.fileURL, nil)

	if err != nil {
		ch <- FileInfo{Err: err}
		return
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", info.ranges[0], info.ranges[1]))

	c := &http.Client{Timeout: 60 * time.Second}
	res, err := c.Do(req)
	if err != nil {
		ch <- FileInfo{Err: err}
		return
	}
	defer res.Body.Close()

	dir := filepath.Dir(info.filename)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		ch <- FileInfo{Err: err}
		return
	}

	path := fmt.Sprintf("%s.part%.3d", info.filename, info.index)
	file, err := os.Create(path)
	if err != nil {
		ch <- FileInfo{Err: err}
		return
	}
	defer file.Close()

	buf := make([]byte, 1024*1024)
	if _, err := io.CopyBuffer(file, res.Body, buf); err != nil {
		os.Remove(path)
		ch <- FileInfo{Err: err}
		return
	}

	ssize := res.Header.Get("Content-Length")
	size, err := strconv.ParseInt(ssize, 10, 64)
	if err != nil {
		os.Remove(path)
		ch <- FileInfo{Err: err}
		return
	}

	ch <- FileInfo{
		ChunkSize: size,
		TotalSize: info.totalSize,
	}
}

func (s *downloader) calcRanges(totalSize int64) [][]int64 {
	var bounds [][]int64
	maxSize := int64(1024 * 512)

	n := int(math.Ceil(float64(totalSize) / float64(maxSize)))
	for i := 0; i < n; i++ {
		startOffset := int64(i) * maxSize
		endOffset := startOffset + maxSize - 1
		if endOffset > totalSize-1 {
			endOffset = totalSize - 1
		}
		bounds = append(bounds, []int64{startOffset, endOffset})
	}

	return bounds
}

func (s *downloader) getFileSize(fileURL string) (int64, error) {
	c := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("HEAD", fileURL, nil)
	if err != nil {
		return 0, err
	}

	res, err := c.Do(req)
	if err != nil {
		return 0, err
	}

	ssize := res.Header.Get("Content-Length")
	size, err := strconv.ParseInt(ssize, 10, 64)
	if err != nil {
		return 0, errors.New("chunk download not support")
	}

	return size, nil
}
