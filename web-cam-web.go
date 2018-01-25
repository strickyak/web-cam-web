// Uses small scraps from stack overflow.
package main

import (
	"github.com/strickyak/web-cam-web/imagediff"

	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var listenFlag = flag.String("listen", ":60606", "listen host:port")
var pwFlag = flag.String("pw", "", "basic auth username:password")
var dirFlag = flag.String("dir", "/tmp/xyzzy", "host directory to serve")

const FilenameFormat = "cam_20060102_15_04.jpg"

func main() {
	flag.Parse()
	http.Handle("/-upload", authWrapper(Upload))
	http.Handle("/", authWrapper(http.FileServer(http.Dir(*dirFlag))))
	http.ListenAndServe(*listenFlag, nil)
}

var previousFilename string

type UploadType int

var Upload UploadType

func (UploadType) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("UploadType::ServeHTTP: %v", r)
	r.ParseMultipartForm(16 * 1000 * 1000)
	rfile, _, err := r.FormFile("image.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rfile.Close()
	filename := *dirFlag + "/" + time.Now().Format(FilenameFormat)
	wfile, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	func() {
		defer wfile.Close()
		io.Copy(wfile, rfile)
	}()

	if previousFilename != "" {
		imagediff.DiffFilenames(previousFilename, filename, filename+".diff.png")
	}
	previousFilename = filename
}

func authWrapper(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("ServeHTTP: %v", r)
		if *pwFlag == "" {
			h.ServeHTTP(w, r)
			return
		}
		user, pass, _ := r.BasicAuth()
		if user+":"+pass != *pwFlag {
			w.Header().Set("WWW-Authenticate", `Basic realm="REALM"`)
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
			return
		}
		h.ServeHTTP(w, r)
		log.Printf("ServeHTTP: OK")
	}
}
