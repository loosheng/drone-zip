package main

import (
	"archive/zip"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

type Plugin struct {
	Input []string
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
		paths, err := filepath.Glob(filepath.Join(inputPath))
		if err != nil {
			return err
		}
		input = append(input, paths...)
	}

	if err := Zip(p.Output, input); err != nil {
		return fmt.Errorf("packagingFailure : %v", err)
	}
	return nil
}

func Zip(dst string, src []string) (err error) {
	fw, err := os.Create(dst)
	defer func(fw *os.File) {
		err := fw.Close()
		if err != nil {
			logrus.Fatal(err)
		}
	}(fw)
	if err != nil {
		return err
	}

	zw := zip.NewWriter(fw)
	defer func() {
		// 检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			logrus.Fatal(err)
		}
	}()

	for _, path := range src {
		file, err := os.Lstat(path)
		if err != nil {
			return fmt.Errorf("get %s fileinfo error: %v \n", path, err)
		}

		fh, err := zip.FileInfoHeader(file)
		if err != nil {
			return fmt.Errorf("get zip FileInfoHeader err: %v", err)
		}

		fh.Name = file.Name()

		if file.IsDir() {
			fh.Name += "/"
		}
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}

		if !fh.Mode().IsRegular() {
			return nil
		}

		n, err := CopyFileToWriter(path, w)
		if err != nil {
			return err
		}

		fmt.Printf("Successful file compression: %s, A total of %.2f KB bytes of data was written\n", path, float64(n)/float64(1024))

	}
	return nil
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
