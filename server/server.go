package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/sywc670/willcrypt/internal/utils"
)

const (
	UploadPath = "/upload"

	RetrievePath = "/retrieve"

	Address = ":8080"

	StoreFile = "pairs.txt"
)

type Pair struct {
	id  string
	key string
}

// could use a map
var Pairs []Pair

func init() {
	readStoreFile()
}

func main() {
	http.HandleFunc(UploadPath, uploadHandler)
	http.HandleFunc(RetrievePath, retrieveHandler)

	log.Fatalln(http.ListenAndServe(Address, nil))
}

func readStoreFile() {
	_, err := os.Stat(StoreFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(StoreFile, "not exists")
			return
		} else {
			panic(err)
		}
	}

	bs, err := os.ReadFile(StoreFile)
	if err != nil {
		panic(err)
	}

	scan := bufio.NewScanner(bytes.NewReader(bs))

	for scan.Scan() {
		elem := strings.Fields(scan.Text())

		if len(elem) != 2 {
			panic("store file has wrong format:" + scan.Text())
		}

		id, key := elem[0], elem[1]

		key = utils.DecodeBase64(key)

		Pairs = append(Pairs, Pair{id, key})
	}
}
