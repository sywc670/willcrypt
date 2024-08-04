package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sywc670/willcrypt/internal/utils"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("RemoteAddr: %s, URI: %s\n", r.RemoteAddr, r.RequestURI)
	key := r.PostFormValue("key")
	id := r.PostFormValue("id")

	if r.Method != "POST" {
		reject(w, r, "HTTP method is not POST, got "+r.Method)
		return
	}

	if id == "" {
		reject(w, r, "id parameter not set or empty")
		return
	}

	if key == "" {
		reject(w, r, "key parameter is not set or empty")
		return
	}

	for _, pair := range Pairs {
		if pair.id == id {
			reject(w, r, "key already exists")
			return
		}
	}

	pair := Pair{id, key}
	Pairs = append(Pairs, pair)

	// store pair into file
	file, err := os.OpenFile(StoreFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Println("open store file err:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, http.StatusText(http.StatusInternalServerError))
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprint(pair.id, "\t", utils.EncodeBase64(pair.key), "\n"))
	if err != nil {
		log.Println("write store file err:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, http.StatusText(http.StatusInternalServerError))
		return
	}
}

func retrieveHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("RemoteAddr: %s, URI: %s\n", r.RemoteAddr, r.RequestURI)
	id := r.PostFormValue("id")

	if r.Method != "POST" {
		reject(w, r, "HTTP method is not POST, got "+r.Method)
		return
	}

	if id == "" {
		reject(w, r, "id parameter is not set")
		return
	}

	for _, pair := range Pairs {
		if pair.id == id {
			fmt.Fprint(w, pair.key)
			return
		}
	}

	reject(w, r, "no key found for id "+id)
}

func reject(w http.ResponseWriter, r *http.Request, reason string) {
	log.Println("Rejecting ", r.RemoteAddr+":", reason)
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, http.StatusText(http.StatusNotFound))
}
