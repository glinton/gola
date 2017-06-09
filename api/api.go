// Package api defines the routes and how the application will serve them
package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	file404 []byte
)

type (
	apiMsg struct {
		MsgString string `json:"msg"`
		Success   bool   `json:"success"`
	}
)

// Start starts the server
func Start(listenAddr string) error {
	var err error

	// preload the 404
	// todo: readfile will allocate memory the size of the file reading, don't use for very large 404 files
	file404, err = ioutil.ReadFile("./app/404.html")
	if err != nil {
		return fmt.Errorf("Failed reading cusom 404 doc - %s", err)
	}

	// serve the 'site'
	http.Handle("/requests", logHandler(requestsHandler()))          // returns request counter
	http.Handle("/greet-me", logHandler(helloHandler()))             // show greeting
	http.Handle("/README.md", logHandler(http.NotFoundHandler()))    // don't serve README
	http.Handle("/", logHandler(http.FileServer(http.Dir("./app")))) // serve static files

	// start server
	fmt.Printf("Starting server on %s...\n", listenAddr)
	return http.ListenAndServe(listenAddr, nil)
}

// parseBody parses the request into v
func parseBody(req *http.Request, v interface{}) error {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	fmt.Printf("Parsed body - %s\n", b)

	if err := json.Unmarshal(b, v); err != nil {
		return err
	}

	return nil
}

// cors middleware
func corsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "foxtrotguns.com")
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// run the handler if not OPTIONS
		if req.Method != "OPTIONS" {
			h.ServeHTTP(rw, req)
		}
	})
}

// logging middleware
func logHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()

		writer := &LogRW{rw, 0, 0}
		h.ServeHTTP(writer, req)
		// if writer.status != http.StatusOK {
		if writer.status == http.StatusNotFound {
			writer.CustomErr(rw, req, writer.status)
		}

		remoteAddr := req.RemoteAddr
		if fwdFor := req.Header.Get("X-Forwarded-For"); len(fwdFor) > 0 {
			remoteAddr = fwdFor
		}

		fmt.Printf("%s %s %s%s %s %d(%d) - %s [User-Agent: %s] (%s)\n",
			time.Now().Format(time.RFC3339), req.Method, req.Host, req.RequestURI,
			req.Proto, writer.status, writer.wrote, remoteAddr,
			req.Header.Get("User-Agent"), time.Since(start))
	})
}

// LogRW is provides the logging functionality i've always wanted, giving access
// to the number bytes written, as well as the status. (I try to always writeheader
// prior to write, so status works fine for me)
type LogRW struct {
	http.ResponseWriter
	status int
	wrote  int
}

// WriteHeader matches the response writer interface, and stores the status
func (n *LogRW) WriteHeader(status int) {
	n.status = status

	// http.FileServer and its (http.)Error() function will write text/plain headers
	// which cause the browser to not render the html from our custom error page.
	// write 404 page to current url rather than redirect so refreshing the page will
	// work properly (if the page becomes available later)
	if status != 404 {
		n.ResponseWriter.WriteHeader(status)
	}
}

// Write matches the response writer interface, and stores the number of bytes written
func (n *LogRW) Write(p []byte) (int, error) {
	if n.status == http.StatusNotFound {
		n.wrote = len(p)
		return n.wrote, nil
	}
	wrote, err := n.ResponseWriter.Write(p)
	n.wrote = wrote
	return wrote, err
}

// CustomErr allows us to write a custom error file to the user. It is part of
// LogRW so we can track the bytes written.
func (n *LogRW) CustomErr(w http.ResponseWriter, r *http.Request, status int) {
	if status == http.StatusNotFound {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusNotFound)

		n.wrote, _ = w.Write(file404)
	}
}
