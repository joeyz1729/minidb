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

// Read db file from offset, return entry record
func (df *DBFile) Read(offset int64) (e *Entry, err error) {
	buf := make([]byte, entryHeaderSize)
	if _, err = df.File.ReadAt(buf, offset); err != nil {
		return
	}
	if e, err = Decode(buf); err != nil {
		return
	}

	offset += entryHeaderSize
	if e.KeySize > 0 {
		key := make([]byte, e.KeySize)
		if _, err = df.File.ReadAt(buf, offset); err != nil {
			return
		}
		e.Key = key
	}

	offset += int64(e.KeySize)
	if e.ValueSize > 0 {
		value := make([]byte, e.ValueSize)
		if _, err = df.File.ReadAt(buf, offset); err != nil {
			return
		}
		e.Value = value
	}

	return
}

// Write entry
func (df *DBFile) Write(e *Entry) (err error) {
	enc, err := e.Encode()
	if err != nil {
		return err
	}
	_, err = df.File.WriteAt(enc, df.Offset)
	df.Offset += e.GetSize()
	return
}
