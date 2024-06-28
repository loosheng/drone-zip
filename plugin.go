package main

import (
	"io"
	"math"
	"os"
	"runtime"
	"sync"

	glob "github.com/bmatcuk/doublestar/v4"
	"github.com/klauspost/compress/zip"
	"github.com/sirupsen/logrus"
)

type Plugin struct {
	Input  []string
	Output string
}

var mutex = &sync.Mutex{}

func (p Plugin) Exec() error {

	if len(p.Input) == 0 {
		logrus.Fatalf("please enter the file or directory to be packed")
	}

	if p.Output == "" {
		logrus.Fatalf("please enter the zip output path")
	}

	var (
		input []string
	)

	for _, inputPath := range p.Input {
		filePath := getFilePaths(inputPath)
		input = append(input, filePath...)
	}
	logrus.Infof("glob path: %s", input)

	Zip(p.Output, input)
	return nil
}

func Zip(fileName string, inputList []string) {
	fw, err := os.Create(fileName)
	if err != nil {
		logrus.Fatalf("create %s error: %v", fileName, err)
	}
	defer func(fw *os.File) {
		err := fw.Close()
		if err != nil {
			logrus.Errorf("close %s error: %v", fileName, err)
		}
	}(fw)

	w := zip.NewWriter(fw)
	defer func(w *zip.Writer) {
		err := w.Close()
		if err != nil {
			logrus.Errorf("close %s error: %v", fileName, err)
		}
	}(w)

	// set maxGoroutines
	cpu80Percent := int(math.Ceil(float64(runtime.NumCPU()) * 0.8))
	fileCount := len(inputList)
	maxGoroutines := getMin(cpu80Percent, fileCount)

	sem := make(chan struct{}, maxGoroutines)
	errCh := make(chan error, fileCount)
	quit := make(chan struct{})
	wg := &sync.WaitGroup{}

	for _, input := range inputList {
		select {
		case <-quit:
			return
		default:
			wg.Add(1)
			sem <- struct{}{}

			go func(filePath string) {
				defer wg.Done()
				defer func() { <-sem }() // Release semaphore.

				if err := addFileToZip(w, filePath); err != nil {
					select {
					case errCh <- err:
					case <-quit: // If other goroutines have already detected an error, do not send.
					}
				}
			}(input)
		}
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			close(quit) // When an error is detected, notify all goroutines to stop.
			logrus.Fatalf("%v", err)
			return
		}
	}
}

func addFileToZip(w *zip.Writer, filePath string) error {
	mutex.Lock()
	defer mutex.Unlock()

	targetFile, err := w.Create(filePath)
	if err != nil {
		return err
	}

	sourceFile, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {
		}
	}(sourceFile)

	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		return err
	}

	logrus.Infof("compression %s file success", filePath)
	return nil
}

func getMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsDir()
}

func getFilePaths(path string) []string {
	var paths, resultPaths []string
	var patternPath string

	if IsDir(path) {
		patternPath = path + "/**/*"
	} else {
		patternPath = path
	}

	globedPaths, err := glob.FilepathGlob(patternPath)
	if err != nil {
		logrus.Fatalf("glob error: %v", err)
	}

	paths = append(paths, globedPaths...)

	// remove directory
	for _, path := range paths {
		if !IsDir(path) {
			resultPaths = append(resultPaths, path)
		}
	}

	return resultPaths
}
