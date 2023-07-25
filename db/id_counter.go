package db

import (
	"encoding/binary"
	"os"
	"sync"
)

const fileName = "indexcounter"

type IdCounter struct {
	sync.Mutex
	count uint64
	file  *os.File
}

func NewIdCounter(filesPath string) (*IdCounter, error) {
	file, err := os.OpenFile(filesPath+"/"+fileName, os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	var counter uint64
	if stat.Size() > 0 {
		if err = binary.Read(file, binary.LittleEndian, &counter); err != nil {
			return nil, err
		}
	}

	return &IdCounter{
		count: counter,
		file:  file,
	}, nil
}

func (c *IdCounter) FetchAndInc() (uint64, error) {
	c.Lock()
	defer c.Unlock()

	prev := c.count
	c.count++

	_, err := c.file.Seek(0, 0)
	if err != nil {
		return 0, err
	}

	if err = binary.Write(c.file, binary.LittleEndian, &c.count); err != nil {
		return 0, err
	}

	_, err = c.file.Seek(0, 0)
	if err != nil {
		return 0, err
	}

	return prev, nil
}
