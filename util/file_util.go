package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

func RemoveDirContents(dirNames ...string) error {
	for _, dirName := range dirNames {
		dir, err := ioutil.ReadDir(dirName)
		for _, d := range dir {
			os.RemoveAll(path.Join([]string{dirName, d.Name()}...))
		}
		if err != nil {
			return err
		}

	}
	return nil
}

// roatates the provided file by adding the current timestamp to its name
func RotateFile(rotationCount int, baseFileName string, baseFileExtension string) error {

	var baseFile = fmt.Sprintf("%s.%s", baseFileName, baseFileExtension)
	var rotatedFile = fmt.Sprintf("%s-%s.%s", baseFileName, time.Now().Format("2006-01-02-150405"), baseFileExtension)
	if _, err := os.Stat(baseFile); os.IsNotExist(err) {
		return nil
	}

	os.Rename(baseFile, rotatedFile)

	globPattern := fmt.Sprintf("%s-%s.%s", baseFileName, "*", baseFileExtension)
	matches, err := filepath.Glob(globPattern)
	if err != nil {
		return err
	}

	var toUnlink []string

	if rotationCount > 0 {
		// Only delete if we have more than rotationCount
		if rotationCount >= len(matches) {
			return nil
		}

		toUnlink = matches[:len(matches)-rotationCount]
	}

	if len(toUnlink) <= 0 {
		return nil
	}

	for _, path := range toUnlink {
		os.Remove(path)
	}
	return nil
}

func GetExecutableDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

func CheckFileExists(path string) bool {
	fileInfo, err := os.Stat(path)
	return err == nil && !fileInfo.IsDir()
}
func CheckDirectoryExists(path string) bool {
	fileInfo, err := os.Stat(path)
	return err == nil && fileInfo.IsDir()
}
func CheckFilesExists(paths ...string) bool {
	for _, path := range paths {
		if !CheckFileExists(path) {
			return false
		}
	}
	return true
}
func CheckDirectoriesExists(paths ...string) bool {
	for _, path := range paths {
		if !CheckDirectoryExists(path) {
			return false
		}
	}
	return true
}
