package mw

import (
	"bufio"
	"compress/flate"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

type writeResetter interface {
	io.Writer
	Reset(io.Writer)
}

type compressResponseWriter struct {
	writeResetter
	http.ResponseWriter
}

func (w *compressResponseWriter) WriteHeader(c int) {
	w.ResponseWriter.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(c)
}

func (w *compressResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *compressResponseWriter) Write(b []byte) (int, error) {
	h := w.ResponseWriter.Header()
	if h.Get("Content-Type") == "" {
		h.Set("Content-Type", http.DetectContentType(b))
	}
	h.Del("Content-Length")

	return w.writeResetter.Write(b)
}

func (w *compressResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("compress response does not implement http.Hijacker")
	}

	conn, bufrw, err := h.Hijack()
	if err == nil {
		w.writeResetter.Reset(ioutil.Discard)
	}

	return conn, bufrw, err
}

// CompressHandler gzip compresses HTTP responses for clients that support it
// via the 'Accept-Encoding' header.
func Compress(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	L:
		for _, enc := range strings.Split(r.Header.Get("Accept-Encoding"), ",") {
			switch strings.TrimSpace(enc) {
			case "gzip":
				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Add("Vary", "Accept-Encoding")

				gw := gzip.NewWriter(w)
				defer gw.Close()

				w = &compressResponseWriter{
					writeResetter:  gw,
					ResponseWriter: w,
				}

				break L
			case "deflate":
				fw, err := flate.NewWriter(w, flate.DefaultCompression)
				if err != nil {
					log.Println(err)
					next.ServeHTTP(w, r)
				}
				w.Header().Set("Content-Encoding", "deflate")
				w.Header().Add("Vary", "Accept-Encoding")
				defer fw.Close()

				w = &compressResponseWriter{
					writeResetter:  fw,
					ResponseWriter: w,
				}

				break L
			}
		}

		next.ServeHTTP(w, r)
	}
}
