package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sync"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileNotFound          = errors.New("file not found")
	ErrFileStat              = errors.New("file stat error")
	ErrOpenFile              = errors.New("file open error")
	ErrReadFile              = errors.New("file read error")
	ErrWriteFile             = errors.New("file write error")
	N                        = 1024
	inF                      *os.File
	outF                     *os.File
	buf                      []byte
	ch                       chan int64
	wg                       sync.WaitGroup
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	ch = make(chan int64, 1)
	if limit == 0 {
		limit = math.MaxInt64
	}

	if limit < int64(N) {
		buf = make([]byte, limit)
	} else {
		buf = make([]byte, N)
	}

	var size int64

	defer inF.Close()
	defer outF.Close()

	err := check(fromPath, toPath, offset, &size)
	if err != nil {
		return err
	}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for d := range ch {
			fmt.Printf(" %d"+"%%", 100*d/size)
		}
		fmt.Print(" 100%")
	}(&wg)

	wOffset := int64(0)
	var fErr error
	var n int
	for wOffset < limit && fErr != io.EOF {
		n, fErr = inF.ReadAt(buf, offset)
		if fErr != nil && fErr != io.EOF {
			return fmt.Errorf(ErrReadFile.Error()+" "+fErr.Error()+" %s\n", fromPath)
		}
		tmp := wOffset + int64(n)
		if tmp > limit {
			buf = buf[0 : limit-wOffset]
		} else if fErr == io.EOF {
			buf = buf[0:n]
		}

		n, err = outF.WriteAt(buf, wOffset)
		if err != nil {
			return fmt.Errorf(ErrWriteFile.Error()+" "+err.Error()+" %s\n", toPath)
		}
		wOffset += int64(n)
		offset += int64(n)
		ch <- wOffset
	}
	close(ch)
	wg.Wait()

	return nil
}

func check(fromPath, toPath string, offset int64, size *int64) error {
	var err error
	inF, err = os.Open(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(ErrFileNotFound.Error()+" "+err.Error()+" %s\n", fromPath)
		}
		return fmt.Errorf(ErrOpenFile.Error()+" "+err.Error()+" %s\n", fromPath)
	}

	inFInf, er := inF.Stat()
	if er != nil {
		return fmt.Errorf(ErrFileStat.Error()+" "+er.Error()+" %s\n", fromPath)
	}
	*size = inFInf.Size()
	if *size == 0 {
		return ErrUnsupportedFile
	}

	if offset >= inFInf.Size() {
		return ErrOffsetExceedsFileSize
	}

	outF, err = os.Create(toPath)
	if err != nil {
		return fmt.Errorf(ErrOpenFile.Error()+" "+err.Error()+" %s\n", toPath)
	}

	return nil
}
