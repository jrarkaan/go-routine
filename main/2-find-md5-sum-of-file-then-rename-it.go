package main

import (
	"crypto/md5"
	"fmt"
	"go-routine/main/dto"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	totalCpu := runtime.NumCPU()
	fmt.Println("Total CPU: ", totalCpu)

	totalThread := runtime.GOMAXPROCS(-1)
	fmt.Println("Total Threads: ", totalThread)
	log.Println("start")
	log.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~ start program without concurrent/ worker threads ~~~~~~~~~~~~~~~~~~~~~~~~~~")
	start := time.Now()

	proceed()

	duration := time.Since(start)
	log.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~ done in", duration.Seconds(), "seconds ~~~~~~~~~~~~~~~~~~~~~~~~~~")
}

func proceed() {
	counterTotal := 0
	counterRenamed := 0
	err := filepath.Walk(dto.TempPath, func(path string, info os.FileInfo, err error) error {
		// if there is an error, return immediatelly
		if err != nil {
			return err
		}
		// if it is a sub directory, return immediatelly
		if info.IsDir() {
			return nil
		}
		counterTotal++
		// read file
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		// sum it
		sum := fmt.Sprintf("%x", md5.Sum(buf))
		// rename file
		destinationPath := filepath.Join(dto.TempPath, fmt.Sprintf("file-%s.txt", sum))
		err = os.Rename(path, destinationPath)
		if err != nil {
			return err
		}
		counterRenamed++
		return nil
	})
	if err != nil {
		log.Println("ERROR:", err.Error())
	}
	log.Printf("%d/%d files renamed", counterRenamed, counterTotal)

}
