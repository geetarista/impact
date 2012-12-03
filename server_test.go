package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testItem struct {
	url  string
	code int
	body string
}

func runTests(tests []testItem, server *httptest.Server, t *testing.T) {
	for _, test := range tests {
		res, err := http.Get(server.URL + test.url)
		if err != nil {
			t.Fatal(err)
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := res.StatusCode, test.code; got != want {
			t.Errorf("%s: code = %d, want %d", test.url, got, want)
		}
		if got, want := string(body), test.body; got != want {
			t.Errorf("%s: body = %q, want %q", test.url, got, want)
		}
	}
}

func TestFileHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(fileHandler))
	defer server.Close()

	bi, _ := ioutil.ReadFile("./index.html")
	bw, _ := ioutil.ReadFile("./weltmeister.html")

	var tests = []testItem{
		{"/", 200, string(bi)},
		{"/wm", 200, string(bw)},
		{"/test/test.js", 200, "fdsa"},
		{"/test/scripts/foo.js", 200, ""},
		{"/asdf", 404, "404 page not found\n"},
	}

	runTests(tests, server, t)
}

func TestGlobHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(globHandler))
	defer server.Close()

	var tests = []testItem{
		{
			"?glob[]=test/scripts/*.js",
			200,
			`["test/scripts/bar.js","test/scripts/baz.js","test/scripts/foo.js"]`,
		},
		{
			"?glob[]=test/images/*.png",
			200,
			`["test/images/bar.png","test/images/baz.png","test/images/foo.png"]`,
		},
		{
			"?glob[]=test/images/*.js",
			200,
			`[]`,
		},
	}

	runTests(tests, server, t)
}

func TestBrowseHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(browseHandler))
	defer server.Close()

	var tests = []testItem{
		{
			"?dir=test/scripts&type=script",
			200,
			`{"dirs":[],"files":["test/scripts/bar.js","test/scripts/baz.js","test/scripts/foo.js"],"parent":"test"}`,
		},
		{
			"?dir=test/scripts&type=",
			200,
			`{"dirs":[],"files":["test/scripts/bar.js","test/scripts/baz.js","test/scripts/foo.js"],"parent":"test"}`,
		},
		{
			"?dir=test/scripts&type=images",
			200,
			`{"dirs":[],"files":[],"parent":"test"}`,
		},
		{
			"?dir=test&type=",
			200,
			`{"dirs":["test/images","test/mix","test/scripts"],"files":["test/test.js"],"parent":"."}`,
		},
		{
			"?dir=&type=",
			200,
			`{"dirs":["test"],"files":["LICENSE","README.md","impact","index.html","server.go","server_test.go","version.go","weltmeister.html"],"parent":false}`,
		},
	}

	runTests(tests, server, t)
}

func TestSaveHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(saveHandler))
	defer server.Close()

	var tests = []testItem{
		{"?path=test/test.js&data=fdsa", 200, `{"message":"","error":0}`},
		{"?path=&data=fdsa", 200, `{"message":"No Data or Path specified","error":1}`},
		{"?path=asdf&data=", 200, `{"message":"No Data or Path specified","error":1}`},
		{"?path=derp.js&data=fdsa", 200, `{"message":"Couldn't write to file: open derp.js: no such file or directory","error":2}`},
		{"?path=asdf&data=fdsa", 200, `{"message":"File must have a .js suffix","error":3}`},
	}

	runTests(tests, server, t)
}
