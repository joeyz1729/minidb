package main

import (
	"fmt"
	"strconv"
	"zouyi/minidb"
)

func main() {
	db, err := minidb.Open("./tmp/minidb")
	if err != nil {
		panic(err)
	}

	var (
		key   = []byte("dbname")
		value = []byte("minidb")
	)
	err = db.Put(key, value)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 1000; i++ {

		err = db.Put([]byte("key"+strconv.Itoa(i)), []byte("value"+strconv.Itoa(i)))
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("1. put kv successfully, key: %s, value: %s.\n", string(key), string(value))

	cur, err := db.Get([]byte("key50"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("2. get value of key %s, the value is %s. \n", string(key), string(cur))

	err = db.Del(key)
	if err != nil {
		panic(err)
	}

	fmt.Printf("3. delete key %s. \n", string(key))

	db.Merge()
	fmt.Println("4. compact data to new db file.")

	db.Close()
	fmt.Println("5. close mini db. ")

}
