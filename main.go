package main

import (
	"github.com/gorilla/sessions"

  	"html/template"
  	"io"
  	"log"
  	"net"
  	"net/http"
  	"os"
  	"time"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	t := template.New("index")
    t, err := template.ParseFiles("templates\\example.html")
    if(err != nil) {
    	panic(err)
    }

    cookie, err := r.Cookie("id")
    if err != nil {
    	log.Printf(err.Error())
    } else {
    	log.Printf(cookie.Value)
	}

  	name := "Sahil"
  	err = t.Execute(w, name) 
  	if(err != nil) {
  		panic(err)
  	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	//POST takes the uploaded file(s) and saves it to disk.
	case "POST":
		//get the multipart reader for the request.
		reader, err := r.MultipartReader()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//copy each part to destination.
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}

			//if part.FileName() is empty, skip this iteration.
			if part.FileName() == "" {
				continue
			}
			dst, err := os.Create("mine.png")
			defer dst.Close()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			if _, err := io.Copy(dst, part); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func download_s3(w http.ResponseWriter, r *http.Request) {
	url := "http://somefilehere.com/asdfasdlf"

	timeout := time.Duration(5) * time.Second
	transport := &http.Transport{
		ResponseHeaderTimeout: timeout,
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		},
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: transport,
	}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//copy the relevant headers. If you want to preserve the downloaded file name, extract it with go's url parser.
	w.Header().Set("Content-Disposition", "attachment; filename=Wiki.png")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", r.Header.Get("Content-Length"))

	//stream the body to the client without fully loading it into memory
	io.Copy(w, resp.Body)
}

func download(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "mine.png")
}

func main() {
	fs := http.FileServer(http.Dir("public"))
  	http.Handle("/public/", http.StripPrefix("/public/", fs))
  	http.HandleFunc("/upload", upload)
  	http.HandleFunc("/download", download)
  	http.HandleFunc("/", serveTemplate)

  	log.Println("Listening...")
  	http.ListenAndServe(":3000", nil)
}