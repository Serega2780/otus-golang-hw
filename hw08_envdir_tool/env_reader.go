package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	WhiteSpace = " "
	Tab        = "\t"
	NullSymbol = "\x00"
	NewLine    = "\n"
	Slash      = "/"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrIsDir = errors.New("unable to read %s: is a directory")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Place your code here
	var f *os.File
	var fErr error
	var fInf os.FileInfo
	var env Environment
	var buf []byte

	env = Environment{}

	dirs, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, fi := range dirs {
		fileName := dir + Slash + fi.Name()
		if strings.Contains(fileName, "=") {
			continue
		}
		if fi.IsDir() {
			return nil, fmt.Errorf(ErrIsDir.Error(), fileName)
		}
		if fInf, err = fi.Info(); err != nil {
			return nil, err
		}
		size := fInf.Size()
		if size == 0 {
			env[fi.Name()] = EnvValue{Value: "", NeedRemove: true}
			continue
		}

		if f, err = os.Open(fileName); err != nil {
			return nil, err
		}
		envValue := make([]byte, 0)
		buf = make([]byte, 1)
		for !errors.Is(fErr, io.EOF) {
			if _, fErr = f.Read(buf); fErr != nil && !errors.Is(fErr, io.EOF) {
				return nil, fErr
			}
			if buf[0] == 10 || errors.Is(fErr, io.EOF) {
				break
			}
			envValue = append(envValue, buf[0])
		}
		f.Close()
		fErr = nil
		if len(envValue) == 0 {
			env[fi.Name()] = EnvValue{Value: "", NeedRemove: false}
			continue
		}
		envValue = bytes.ReplaceAll(envValue, []byte(NullSymbol), []byte(NewLine))
		str := strings.TrimRight(strings.TrimRight(string(envValue), WhiteSpace), Tab)
		env[fi.Name()] = EnvValue{Value: str, NeedRemove: false}
	}
	return env, nil
}
