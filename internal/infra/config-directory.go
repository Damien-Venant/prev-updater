package infra

import (
	"errors"
	"io"
	"os"
	"path"
)

var (
	logFileName string = ""
	file        *os.File
)

func ConfigDirectory() (string, error) {
	dirName, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	prevUpdaterPath := path.Join(dirName, "prev-udpater")
	if _, err = os.Stat(prevUpdaterPath); os.IsNotExist(err) {
		err = os.MkdirAll(prevUpdaterPath, 0750)
		return "", err
	}
	logFileName = path.Join(prevUpdaterPath)
	return prevUpdaterPath, nil
}

func OpenLogFile() (io.Writer, error) {
	var err error
	file, err = os.OpenFile(logFileName, os.O_CREATE, 0750)
	if err != nil {
		return nil, err
	}
	return file, err
}

func CloseLogFile() error {
	if file == nil {
		return errors.New("log file not opened")
	}
	return file.Close()
}
