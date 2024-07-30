package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileNotFound          = errors.New("file not found")
	ErrFileStat              = errors.New("file stat error")
	ErrOpenFile              = errors.New("file open error")
	ErrReadFile              = errors.New("file read error")
	ErrWriteFile             = errors.New("file write error")
	ErrCloseFile             = errors.New("file close error")
	ErrDeleteFile            = errors.New("file delete error")
	ErrCreateFile            = errors.New("file create error")
	N                        = 1024
	inF, outF, tmpF          *os.File
	buf                      []byte
	ch                       chan int64
	size, pSize              int64
	wg                       sync.WaitGroup
	isSame, fIn, fOut        bool
	toPathTmp                string
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

	defer closeFiles()

	err := check(fromPath, toPath, offset, &size)
	if err != nil {
		return err
	}

	pSize = size
	pSize -= offset
	if pSize > limit {
		pSize = limit
	}
	bar := pb.StartNew(int(pSize)).SetUnits(pb.U_BYTES)
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for d := range ch {
			bar.Add(int(d))
		}
	}(&wg)

	wOffset := int64(0)
	var fErr error
	var n int
	for wOffset < limit && !errors.Is(fErr, io.EOF) {
		n, fErr = inF.ReadAt(buf, offset)
		if fErr != nil && !errors.Is(fErr, io.EOF) {
			if fErr = removeToFileIfError(outF); fErr != nil {
				return fmt.Errorf(ErrReadFile.Error()+" "+fErr.Error()+" %s\n", fromPath)
			}
		}
		tmp := wOffset + int64(n)
		if tmp > limit {
			buf = buf[0 : limit-wOffset]
		} else if errors.Is(fErr, io.EOF) {
			buf = buf[0:n]
		}

		n, err = outF.WriteAt(buf, wOffset)
		if err != nil {
			return fmt.Errorf(ErrWriteFile.Error()+" "+err.Error()+" %s\n", toPath)
		}
		wOffset += int64(n)
		offset += int64(n)
		ch <- int64(n)
	}
	close(ch)
	if isSame {
		if err = inF.Close(); err != nil {
			return fmt.Errorf(ErrCloseFile.Error()+" "+err.Error()+" %s\n", fromPath)
		}
		fIn = true
		if err = os.Rename(toPathTmp, toPath); err != nil {
			return fmt.Errorf(ErrCreateFile.Error()+" "+err.Error()+" %s\n", toPath)
		}
	}
	wg.Wait()
	bar.Finish()

	return nil
}

func removeToFileIfError(f *os.File) (er error) {
	fName := f.Name()
	if er = f.Close(); er != nil {
		return fmt.Errorf(ErrCloseFile.Error()+" "+er.Error()+" %s\n", fName)
	}
	if er = os.Remove(fName); er != nil {
		return fmt.Errorf(ErrDeleteFile.Error()+" "+er.Error()+" %s\n", fName)
	}
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

	isSame, er = isTheSame(toPath, &inFInf)
	if er != nil {
		return er
	}
	if isSame {
		toPathTmp = toPath + "_" + strconv.FormatInt(time.Now().UnixMilli(), 10)
		if outF, err = os.Create(toPathTmp); err != nil {
			return fmt.Errorf(ErrCreateFile.Error()+" "+err.Error()+" %s\n", toPathTmp)
		}
	} else {
		if outF, err = os.Create(toPath); err != nil {
			return fmt.Errorf(ErrCreateFile.Error()+" "+err.Error()+" %s\n", toPath)
		}
	}

	return nil
}

func isTheSame(toPath string, inFInf *os.FileInfo) (bool, error) {
	var tmpInf os.FileInfo
	var er error

	if tmpF, er = os.Open(toPath); er != nil {
		if os.IsNotExist(er) {
			return false, nil
		}
		return false, fmt.Errorf(ErrOpenFile.Error()+" "+er.Error()+" %s\n", toPath)
	}
	defer tmpF.Close()
	if tmpInf, er = tmpF.Stat(); er != nil {
		return false, fmt.Errorf(ErrFileStat.Error()+" "+er.Error()+" %s\n", toPath)
	}
	return os.SameFile(*inFInf, tmpInf), nil
}

func closeFiles() {
	if !fIn {
		inF.Close()
	}
	if !fOut {
		outF.Close()
	}
}
