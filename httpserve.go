package offthegrid

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"io"
	"net/http"
	"time"
)

func (h *OTGServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Server", fmt.Sprintf("offthegrid/%d", VERSION))

	// Only accept GET or HEAD
	if req.Method != "GET" && req.Method != "HEAD" {
		http.Error(w, "Method not allowed.", 405)
		return
	}

	db := h.DB()
	defer db.Close()

	filename := req.URL.Path[1:]
	file, err := db.DB("").GridFS(h.Config.GridFSPrefix).Open(filename)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, "File not found.", 404)
		} else {
			http.Error(w, "Internal server error.", 500)
		}
		return
	}

	defer file.Close()

	if h.Config.CORSHeader != "" {
		w.Header().Set("Access-Control-Allow-Origin", h.Config.CORSHeader)
	}

	// Set cache-control and expiry
	w.Header().Set("Cache-Control", fmt.Sprintf("maxage=%d", h.Config.MaxAge))
	w.Header().Set("Expires", time.Now().Add(time.Duration(h.Config.MaxAge)*time.Second).Format(time.RFC1123))

	ctype := file.ContentType()
	if ctype == "" {
		ctype = "application/octet-stream"
	}
	w.Header().Set("Content-Type", ctype)

	w.Header().Set("ETag", file.MD5())

	// And send the file.
	if req.Method == "GET" {
		io.Copy(w, file)
	}
}
