package minidb

import "sync"

type MiniDB struct {
	indexes map[string]int64
	dbFile  *DBFile
	dirPath string
	mu      sync.Mutex
}
