package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

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
		paths, err := filepath.Glob(inputPath)
		if err != nil {
			return err
		}
		input = append(input, paths...)
	}

	Zip(p.Output, input)
	return nil
}

func Zip(dst string, src []string) {
	abs, err := filepath.Abs(dst)
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Printf("dst path: %v \n", abs)
	fw, err := os.Create(dst)
	defer func(fw *os.File) {
		err := fw.Close()
		if err != nil {
			logrus.Fatal(err)
		}
	}(fw)
	if err != nil {
		logrus.Fatal(err)
	}

	zw := zip.NewWriter(fw)
	defer func() {
		// 检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			logrus.Fatal(err)
		}
	}()

	for _, path := range src {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			logrus.Errorf("file: %v does not exist",path)
		}

		file, err := os.Lstat(path)

		if err != nil {
			logrus.Errorf("get %s fileinfo error: %v", path, err)
		}

		fh, err := zip.FileInfoHeader(file)
		if err != nil {
			logrus.Errorf("get zip FileInfoHeader err: %v", err)
		}

		fh.Name = file.Name()

		if file.IsDir() {
			fh.Name += "/"
		}
		w, err := zw.CreateHeader(fh)
		if err != nil {
			logrus.Fatal(err)
		}

		if !fh.Mode().IsRegular() {
			continue
		}

		n, err := CopyFileToWriter(path, w)
		if err != nil {
			logrus.Fatal(err)
		}

		fmt.Printf("Successful file compression: %s, A total of %.2f KB bytes of data was written\n", path, float64(n)/float64(1024))

	}

}

func CopyFileToWriter(path string, writer io.Writer) (int64, error) {
	fr, err := os.Open(path)

	defer func(fr *os.File) {
		err := fr.Close()
		if err != nil {
			logrus.Fatal(err)
		}
	}(fr)

	if err != nil {
		logrus.Fatal(err)
	}

	number, err := io.Copy(writer, fr)
	if err != nil {
		return 0, err
	}

	return number, err
}
