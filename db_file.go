package minidb

import (
	"os"
	"path/filepath"
)

const Filename = "miniDB.data"
const MergeFilename = "miniDB.data.merge"

type DBFile struct {
	File   *os.File
	Offset int64
}

func newInternal(filename string) (*DBFile, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	stat, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	return &DBFile{File: file, Offset: stat.Size()}, nil
}

// NewDBFile create a new db file
func NewDBFile(path string) (*DBFile, error) {
	filename := filepath.Join(path, Filename)
	return newInternal(filename)
}

// NewMergeDBFile create a merge db file
func NewMergeDBFile(path string) (*DBFile, error) {
	filename := filepath.Join(path, MergeFilename)
	return newInternal(filename)
}

// Read from offset
func (df *DBFile) Read(offset int64) (e *Entry, err error) {

	return
}
