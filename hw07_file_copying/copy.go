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
)

type Params struct {
	inF, outF, tmpF   *os.File
	buf               []byte
	ch                chan int64
	size, pSize       int64
	wg                sync.WaitGroup
	isSame, fIn, fOut bool
	toPathTmp         string
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	params := &Params{
		ch: make(chan int64, 1),
	}

	if limit == 0 {
		limit = math.MaxInt64
	}

	if limit < int64(N) {
		params.buf = make([]byte, limit)
	} else {
		params.buf = make([]byte, N)
	}

	defer closeFiles(params)

	err := check(fromPath, toPath, offset, params)
	if err != nil {
		return err
	}

	params.pSize = params.size
	params.pSize -= offset
	if params.pSize > limit {
		params.pSize = limit
	}
	bar := pb.StartNew(int(params.pSize)).SetUnits(pb.U_BYTES)
	params.wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for d := range params.ch {
			bar.Add(int(d))
		}
	}(&params.wg)

	wOffset := int64(0)
	var fErr error
	var n int
	for wOffset < limit && !errors.Is(fErr, io.EOF) {
		n, fErr = params.inF.ReadAt(params.buf, offset)
		if fErr != nil && !errors.Is(fErr, io.EOF) {
			if fErr = params.outF.Close(); fErr != nil {
				return fmt.Errorf(ErrCloseFile.Error()+" "+fErr.Error()+" %s\n", toPath)
			}
			if fErr = removeFile(params.outF); fErr != nil {
				return fmt.Errorf(ErrReadFile.Error()+" "+fErr.Error()+" %s\n", fromPath)
			}
		}
		tmp := wOffset + int64(n)
		if tmp > limit {
			params.buf = params.buf[0 : limit-wOffset]
		} else if errors.Is(fErr, io.EOF) {
			params.buf = params.buf[0:n]
		}

		n, err = params.outF.WriteAt(params.buf, wOffset)
		if err != nil {
			return fmt.Errorf(ErrWriteFile.Error()+" "+err.Error()+" %s\n", toPath)
		}
		wOffset += int64(n)
		offset += int64(n)
		params.ch <- int64(n)
	}
	close(params.ch)
	if params.isSame {
		if err = params.inF.Close(); err != nil {
			return fmt.Errorf(ErrCloseFile.Error()+" "+err.Error()+" %s\n", fromPath)
		}
		params.fIn = true
		if err = os.Rename(params.toPathTmp, toPath); err != nil {
			return fmt.Errorf(ErrCreateFile.Error()+" "+err.Error()+" %s\n", toPath)
		}
	}
	params.wg.Wait()
	bar.Finish()

	return nil
}

func removeFile(f *os.File) (er error) {
	fName := f.Name()
	if er = os.Remove(fName); er != nil {
		return fmt.Errorf(ErrDeleteFile.Error()+" "+er.Error()+" %s\n", fName)
	}
	return nil
}

func check(fromPath, toPath string, offset int64, params *Params) error {
	var err error
	params.inF, err = os.Open(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(ErrFileNotFound.Error()+" "+err.Error()+" %s\n", fromPath)
		}
		return fmt.Errorf(ErrOpenFile.Error()+" "+err.Error()+" %s\n", fromPath)
	}

	inFInf, er := params.inF.Stat()
	if er != nil {
		return fmt.Errorf(ErrFileStat.Error()+" "+er.Error()+" %s\n", fromPath)
	}
	params.size = inFInf.Size()
	if params.size == 0 {
		return ErrUnsupportedFile
	}

	if offset >= inFInf.Size() {
		return ErrOffsetExceedsFileSize
	}

	params.isSame, er = isTheSame(toPath, &inFInf, params)
	if er != nil {
		return er
	}
	if params.isSame {
		params.toPathTmp = toPath + "_" + strconv.FormatInt(time.Now().UnixMilli(), 10)
		if params.outF, err = os.Create(params.toPathTmp); err != nil {
			return fmt.Errorf(ErrCreateFile.Error()+" "+err.Error()+" %s\n", params.toPathTmp)
		}
	} else {
		if params.outF, err = os.Create(toPath); err != nil {
			return fmt.Errorf(ErrCreateFile.Error()+" "+err.Error()+" %s\n", toPath)
		}
	}

	return nil
}

func isTheSame(toPath string, inFInf *os.FileInfo, params *Params) (bool, error) {
	var tmpInf os.FileInfo
	var er error

	if params.tmpF, er = os.Open(toPath); er != nil {
		if os.IsNotExist(er) {
			return false, nil
		}
		return false, fmt.Errorf(ErrOpenFile.Error()+" "+er.Error()+" %s\n", toPath)
	}
	defer params.tmpF.Close()
	if tmpInf, er = params.tmpF.Stat(); er != nil {
		return false, fmt.Errorf(ErrFileStat.Error()+" "+er.Error()+" %s\n", toPath)
	}
	return os.SameFile(*inFInf, tmpInf), nil
}

func closeFiles(params *Params) {
	if !params.fIn {
		params.inF.Close()
	}
	if !params.fOut {
		params.outF.Close()
	}
}
