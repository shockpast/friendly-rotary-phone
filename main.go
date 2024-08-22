package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		return
	}

	key := r.Header.Get("X-Key")
	if key != "shh!" {
		return
	}

	start := time.Now()
	file, header, err := r.FormFile("file")
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return
	}

	fileName := RandomString(6)
	fileExt := filepath.Ext(filepath.Clean(header.Filename))[1:]

	out, err := os.Create(fmt.Sprintf("./tmp/%s.%s", fileName, fileExt))
	if err != nil {
		return
	}

	out.Write(buf.Bytes())
	out.Close()

	log.Printf("%s uploaded %s.%s in %s", r.RemoteAddr, fileName, fileExt, time.Since(start))
	fmt.Fprintf(w, "%s/ss/%s.%s", r.Host, fileName, fileExt)
}

func main() {
	http.Handle("/ss/", http.StripPrefix("/ss/", http.FileServer(http.Dir("./tmp"))))
	http.HandleFunc("/upload", upload)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
