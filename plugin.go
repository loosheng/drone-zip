package main

import (
	"fmt"
	"io"
	"os"

	glob "github.com/bmatcuk/doublestar/v4"
	"github.com/klauspost/compress/zip"
	"github.com/sirupsen/logrus"
)

type Plugin struct {
	Input  []string
	Output string
}


func (p Plugin) Exec() error {

	if len(p.Input) == 0 {
		return fmt.Errorf("please enter the file or directory to be packed")
	}

	if p.Output == "" {
		return fmt.Errorf("please enter the zip output path")
	}

	var (
		input []string
	)

	for _, inputPath := range p.Input {
		filePath, err := getFilePaths(inputPath)
		if err != nil {
			return err
		}
		input = append(input, filePath...)
	}

	for _, inputPath := range input {
		logrus.Infof("match file: %s", inputPath)
	}

	return Zip(p.Output, input)
}

func Zip(fileName string, inputList []string) error {
	fw, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("create %s error: %v", fileName, err)
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

	for _, filePath := range inputList {
		if err := addFileToZip(w, filePath); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(w *zip.Writer, filePath string) error {

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


func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsDir()
}

func getFilePaths(path string) ([]string, error) {
	var paths, resultPaths []string
	var patternPath string

	if IsDir(path) {
		patternPath = path + "/**/*"
	} else {
		patternPath = path
	}

	globedPaths, err := glob.FilepathGlob(patternPath)
	if err != nil {
		return nil, fmt.Errorf("glob error: %v", err)
	}

	paths = append(paths, globedPaths...)

	// remove directory
	for _, path := range paths {
		if !IsDir(path) {
			resultPaths = append(resultPaths, path)
		}
	}

	return resultPaths, nil
}
