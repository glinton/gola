package api

import (
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/glinton/gola/models/core"
)

// hopefully these are correct - from https://www.webucator.com/blog/2010/03/saying-hello-world-in-your-language-using-javascript/
var (
	greetings = core.Phrases{
		"你好世界",             // Chinese
		"Hallo wereld",     // Dutch
		"Hello world",      // English
		"Bonjour monde",    // French
		"Hallo Welt",       // German
		"γειά σου κόσμος",  // Greek
		"Ciao mondo",       // Italian
		"こんにちは世界",          // Japanese
		"여보세요 세계",          // Korean
		"Olá mundo",        // Portuguese
		"Здравствулте мир", // Russian
		"Hola mundo",       // Spanish
	}

	requests uint64 = 0
)

// hello's handler (GET only)
func helloHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		hello(rw, req)
		atomic.AddUint64(&requests, 1)
	})
}

// hello accepts data and will store it in a database
func hello(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusOK)
	rw.Write(append([]byte(greetings[requests%12]), '\n'))
}

// requests' handler (GET only)
func requestsHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(http.StatusOK)
		rw.Write(append([]byte(strconv.Itoa(int(requests))), '\n'))
	})
}
