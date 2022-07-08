package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrMessageOpenFile       = errors.New("open file error")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileFrom, err := os.OpenFile(fromPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return ErrMessageOpenFile
	}
	defer fileFrom.Close()

	statFileFrom, errStat := fileFrom.Stat()

	switch {
	case errStat != nil:
		return err
	case offset > statFileFrom.Size():
		return ErrOffsetExceedsFileSize
	case statFileFrom.Size() == 0:
		return ErrUnsupportedFile
	}

	fileFrom.Seek(offset, 0)
	diff := statFileFrom.Size() - offset
	var bytesLen int64
	if diff > limit {
		bytesLen = limit
	} else {
		bytesLen = diff
	}
	bar := pb.New(int(bytesLen)).SetUnits(pb.U_BYTES)
	bar.Start()
	fileTo, _ := os.Create(toPath)

	writer := io.MultiWriter(fileTo, bar)
	io.CopyN(writer, fileFrom, bytesLen)
	bar.Finish()

	return nil
}
