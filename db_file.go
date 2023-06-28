package minidb

import "os"

const FileName = "miniDB.data"
const MergeFilename = "miniDB.data.merge"

type DBFile struct {
	File   *os.File
	Offset int64
}
