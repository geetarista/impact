package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	Types = map[string][]string{
		"images": {".png", ".gif", ".jpg", ".jpeg", ".bmp"},
		"script": {".js"}}
)

type BrowseResponse struct {
	Dirs   []string    `json:"dirs"`
	Files  []string    `json:"files"`
	Parent interface{} `json:"parent"`
}

type SaveResponse struct {
	Msg   string `json:"message"`
	Error int    `json:"error"`
}

func printError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func encodeJson(resp interface{}) string {
	json, err := json.Marshal(resp)
	printError(err)
	return string(json)
}

func writeJson(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, encodeJson(resp))
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/" {
		path = "/index.html"
	} else if path == "/wm" {
		path = "/weltmeister.html"
	}

	b, err := ioutil.ReadFile("." + path)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(path)))
	io.WriteString(w, string(b))
}

func globHandler(w http.ResponseWriter, r *http.Request) {
	glob, err := filepath.Glob(r.FormValue("glob[]"))
	printError(err)
	if glob == nil {
		glob = make([]string, 0)
	}

	writeJson(w, glob)
}

func browseHandler(w http.ResponseWriter, r *http.Request) {
	dir := ""
	var parent interface{}
	if r.FormValue("dir") != "" {
		dir = filepath.Dir(r.FormValue("dir") + "/")
		dirs := strings.Split(dir, "/")
		if dir != "" && len(dirs) > 1 {
			parent = filepath.Base(strings.Join(dirs[0:len(dirs)-1], "/"))
		} else {
			parent = false
		}
	}

	paths, err := filepath.Glob("./" + dir + "/*")
	printError(err)
	dirs := []string{}
	files := []string{}
	for _, v := range paths {
		stat, err := os.Stat(v)
		printError(err)
		if stat.IsDir() == true {
			dirs = append(dirs[:], v)
		} else {
			kind, ok := Types[r.FormValue("type")]
			if ok && len(kind) >= 0 {
				ext := path.Ext(v)
				for _, e := range kind {
					if e == ext {
						files = append(files[:], v)
					}
				}
			} else {
				files = append(files[:], v)
			}
		}
	}

	writeJson(w, BrowseResponse{Parent: parent, Files: files, Dirs: dirs})
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	pth := r.FormValue("path")
	data := r.FormValue("data")
	resp := SaveResponse{}

	if pth != "" && data != "" {
		dir := filepath.Join(".", pth)

		if path.Ext(dir) == ".js" {
			f, err := os.OpenFile(dir, os.O_WRONLY, 0666)

			if err != nil {
				writeJson(w, SaveResponse{Error: 2, Msg: "Couldn't write to file: " + err.Error()})
				return
			}

			defer f.Close()

			f.WriteString(data)
		} else {
			resp = SaveResponse{Error: 3, Msg: "File must have a .js suffix"}
		}
	} else {
		resp = SaveResponse{Error: 1, Msg: "No Data or Path specified"}
	}

	writeJson(w, resp)
}

func main() {
	http.HandleFunc("/", fileHandler)
	http.HandleFunc("/wm", fileHandler)
	http.HandleFunc("/lib/weltmeister/api/glob.php", globHandler)
	http.HandleFunc("/lib/weltmeister/api/browse.php", browseHandler)
	http.HandleFunc("/lib/weltmeister/api/save.php", saveHandler)
	http.HandleFunc("*", fileHandler)
	fmt.Println("Started impact server at localhost:8080. Visit /wm for the weltmeister editor")
	http.ListenAndServe(":8080", nil)
}
