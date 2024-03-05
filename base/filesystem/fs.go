package filesystem

import (
	"fmt"
	"github.com/illuminatingKong/kongming-kit/base/logx/logrusx"
	"os"
	"path/filepath"
	"strings"
)

var formatter logrusx.JsonFormatter
var log = logrusx.New(logrusx.WithFormatter(formatter))

func WorkDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	wd := filepath.Dir(execPath)
	if filepath.Base(wd) == "bin" {
		wd = filepath.Dir(wd)
	}

	return wd, nil
}

func FileExist(file string) bool {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	if os.IsPermission(err) {
		return false
	}

	return true
}

func CreateDirIfNotExists(path ...string) {
	for _, value := range path {
		if FileExist(value) {
			continue
		}
		err := os.Mkdir(value, 0755)
		if err != nil {
			log.Fatal(fmt.Sprintf("创建目录失败:%s", err.Error()))
		}
	}
}

func CreateAllDirIfNotExists(filePath ...string) {
	for _, value := range filePath {
		if FileExist(value) {
			continue
		}
		err := os.MkdirAll(value, 0755)
		if err != nil {
			log.Fatal(fmt.Sprintf("创建目录失败:%s", err.Error()))
		}
	}
}

func AbPath(inPath string) string {
	log.Info("absolute path: ", inPath)

	if inPath == "$HOME" || strings.HasPrefix(inPath, "$HOME"+string(os.PathSeparator)) {
		inPath = userHomeDir() + inPath[5:]
	}
	inPath = os.ExpandEnv(inPath)
	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}
	p, err := filepath.Abs(inPath)
	if err == nil {
		return filepath.Clean(p)
	}

	log.Fatalf(fmt.Errorf("could not discover absolute path: %w", err).Error())

	return ""
}

func userHomeDir() string {
	return os.Getenv("HOME")
}
