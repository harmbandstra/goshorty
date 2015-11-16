package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pilu/go-base62"
	"gopkg.in/yaml.v2"
	// "io"
	"io/ioutil"
	"log"
	"net/http"
	// "strconv"
	"strings"
)

var DB *sql.DB

const (
	YamlConfig = "goshorty.yml"
)

type (
	DatabaseConfig struct {
		Host     string
		Database string
		Port     string
		User     string
		Password string
	}
	UrlRow struct {
		Url string
	}
	GS struct {
	}
)

func UrlServer(w http.ResponseWriter, req *http.Request) {

	// lookup redirct Url, redirect or payst 404 if not found

	// log request details

	urlPath := req.URL.Path

	// io.WriteString(w, "Url: "+urlPath)
	// io.WriteString(w, "\n")
	// io.WriteString(w, "IP: "+req.RemoteAddr)
	// io.WriteString(w, "\n")
	// io.WriteString(w, "User agent: ")
	// io.WriteString(w, req.Header.Get("User-Agent"))
	// io.WriteString(w, "\n")
	// io.WriteString(w, "Referrer: "+req.Header.Get("Referer"))
	// io.WriteString(w, "\n")

	path := strings.TrimLeft(urlPath, "/")
	pathId := base62.Decode(path)
	// io.WriteString(w, strconv.Itoa(pathId))

	row := DB.QueryRow("select url from urls where id=?", pathId)
	ur := new(UrlRow)
	err := row.Scan(&ur.Url)
	if err != nil {
		log.Fatalf("DB Error: %v", err)
	}

	log.Println(ur.Url)
	http.Redirect(w, req, ur.Url, 301)
}

func UrlRedirect(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, string("http://hackurls.com"), 301)
}

func connectDatabase() {
	data, fileError := ioutil.ReadFile(YamlConfig)
	if fileError != nil {
		log.Fatalf("Unable to read config file %s.", YamlConfig)
	}

	config := DatabaseConfig{}
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Unable to parse config file: %v", err)
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.User, config.Password, config.Host, config.Port, config.Database)
	log.Println(connectionString)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("Unable to connect to db: %v", err)
	}
	DB = db

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Unable to connect to db: %v", err)
	}
}

func main() {

	// ➜  goshorty git:(develop) ✗ export PATH=$PATH:/usr/local/opt/go/libexec/bin
	// ➜  goshorty git:(develop) ✗ export GOPATH=$HOME/Source/golang/own
	connectDatabase()

	// defer db.Close()

	http.HandleFunc("/", UrlServer)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
