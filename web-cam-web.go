// Uses small scraps from stack overflow.
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

var listenFlag = flag.String("listen", ":60606", "listen host:port")
var pwFlag = flag.String("pw", "", "basic auth username:password")
var dirFlag = flag.String("dir", "/tmp/xyzzy", "host directory to serve")
var uploadFlag = flag.String("upload", "/tmp/xyzzy/image.jpg", "host filename to upload")

func authWrapperFunc(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if *pwFlag == "" {
			h(w, r)
			return
		}
		user, pass, _ := r.BasicAuth()
		if user+":"+pass != *pwFlag {
			w.Header().Set("WWW-Authenticate", `Basic realm="REALM"`)
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
			return
		}
		h(w, r)
	}
}

func authWrapper(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 * 1000 * 1000)
	file, handler, err := r.FormFile("image.jpg")
	_ = handler
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	f, err := os.Create(*uploadFlag)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
}

func main() {
	flag.Parse()
	http.HandleFunc("/xyzzy.upload", authWrapperFunc(upload))
	http.Handle("/", authWrapper(http.FileServer(http.Dir(*dirFlag))))
	http.ListenAndServe(*listenFlag, nil)
}
