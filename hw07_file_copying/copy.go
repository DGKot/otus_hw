package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

type PbWritter struct {
	current int64
	total   int64
}

func (pb *PbWritter) Write(p []byte) (n int, err error) {
	n = len(p)
	pb.current += int64(n)
	fmt.Printf("Progress of copying file: %.2f%%\n", float64(pb.current)*100/float64(pb.total))
	return n, nil
}

func CopyCore(to io.Writer, from io.Reader, limit int64) error {
	progresWritter := &PbWritter{total: limit}
	_, err := io.CopyN(io.MultiWriter(to, progresWritter), from, limit)
	if err != nil {
		return err
	}
	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	file, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}
	if fileInfo.Size() < offset {
		return ErrOffsetExceedsFileSize
	}
	resFile, err := os.OpenFile(toPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer resFile.Close()

	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	if limit == 0 {
		limit = fileInfo.Size()
	} else {
		limit = min(limit, fileInfo.Size()-offset)
	}

	return CopyCore(resFile, file, limit)
}
