package main

import (
	"io"
	"log"
	"net/http"
)

func UrlServer(w http.ResponseWriter, req *http.Request) {
	urlPath := req.URL.Path
	io.WriteString(w, "Url: "+urlPath)
	io.WriteString(w, "\n")
	io.WriteString(w, "IP: "+req.RemoteAddr)
	io.WriteString(w, "\n")
	io.WriteString(w, "User agent: ")
	io.WriteString(w, req.Header.Get("User-Agent"))
}

func UrlRedirect(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, string("http://hackurls.com"), 301)
}

func main() {
	http.HandleFunc("/", UrlServer)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
