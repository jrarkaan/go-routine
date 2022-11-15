package main

import (
	"crypto/md5"
	"fmt"
	"go-routine/main/dto"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type FileInfo struct {
	FilePath  string
	Content   []byte
	Sum       string // md5 sum of content
	IsRenamed bool   // indicator whether the particular file is renamed already or not
}

func main() {
	totalCpu := runtime.NumCPU()
	fmt.Println("Total CPU: ", totalCpu)

	totalThread := runtime.GOMAXPROCS(-1)
	fmt.Println("Total Threads: ", totalThread)
	log.Println("start")
	log.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~ start program with concurrent/ worker threads ~~~~~~~~~~~~~~~~~~~~~~~~~~")
	start := time.Now()

	// pipeline 1: loop all files and read it
	chanFileContent := readFiles()
	// pipleline 2: calculate md5sum
	chanFilesum1 := getSum(chanFileContent)
	chanFilesum2 := getSum(chanFileContent)
	chanFilesum3 := getSum(chanFileContent)
	chanFileSum := mergeChanFileInfo(chanFilesum1, chanFilesum2, chanFilesum3)
	// pipleline 3: rename files
	chanRename1 := rename(chanFileSum)
	chanRename2 := rename(chanFileSum)
	chanRename3 := rename(chanFileSum)
	chanRename4 := rename(chanFileSum)
	chanRename := mergeChanFileInfo(chanRename1, chanRename2, chanRename3, chanRename4)
	// print output
	counterRenamed := 0
	counterTotal := 0
	for fileInfo := range chanRename {
		if fileInfo.IsRenamed {
			counterRenamed++
		}
		counterTotal++
	}
	log.Printf("%d/%d files renamed", counterRenamed, counterTotal)

	duration := time.Since(start)
	log.Println(" ~~~~~~~~~~~~~~~~~~~~~~~~~~ done in", duration.Seconds(), "seconds ~~~~~~~~~~~~~~~~~~~~~~~~~~")

}

func readFiles() <-chan FileInfo {
	chanOut := make(chan FileInfo)
	fmt.Println("READ FILES - 1: ", chanOut)
	go func() {
		err := filepath.Walk(dto.TempPath, func(path string, info os.FileInfo, err error) error {
			// if there is an error, return immediatelly
			if err != nil {
				return err
			}
			// if it is a sub directory, return immediatelly
			if info.IsDir() {
				return nil
			}
			buf, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			chanOut <- FileInfo{
				FilePath: path,
				Content:  buf,
			}
			return nil
		})
		if err != nil {
			log.Println("ERROR: ", err.Error())
		}
		close(chanOut)
	}()
	return chanOut
}

func getSum(chanIn <-chan FileInfo) <-chan FileInfo {
	chanOut := make(chan FileInfo)
	go func() {
		for fileInfo := range chanIn {
			fileInfo.Sum = fmt.Sprintf("%x", md5.Sum(fileInfo.Content))
			chanOut <- fileInfo
		}
		close(chanOut)
	}()
	return chanOut
}

func mergeChanFileInfo(chanInMany ...<-chan FileInfo) <-chan FileInfo {
	wg := new(sync.WaitGroup)
	chanOut := make(chan FileInfo)

	wg.Add(len(chanInMany))
	for _, eachChan := range chanInMany {
		go func(eachChan <-chan FileInfo) {
			for eachChanData := range eachChan {
				chanOut <- eachChanData
			}
			wg.Done()
		}(eachChan)
	}

	go func() {
		wg.Wait()
		close(chanOut)
	}()

	return chanOut
}

func rename(chanIn <-chan FileInfo) <-chan FileInfo {
	chanOut := make(chan FileInfo)

	go func() {
		for fileInfo := range chanIn {
			newPath := filepath.Join(dto.TempPath, fmt.Sprintf("file-%s.txt", fileInfo.Sum))
			err := os.Rename(fileInfo.FilePath, newPath)
			fileInfo.IsRenamed = err == nil
			chanOut <- fileInfo
		}

		close(chanOut)
	}()

	return chanOut
}