package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// StreamUriHandler
type StreamUriHandler struct {
	camera *Camera
	Handler
}

func (s *StreamUriHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	log.Println("Request:", req.Method, req.URL.Path)
	streamUri, err := s.camera.GetStreamUri()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(s.errorMessage(err.Error()))
		return
	}

	body := map[string]interface{}{
		"result": "OK",
		"uri":    streamUri,
	}
	s.responseJson(writer, body)
}

func main() {
	const port = 3333
	const imageDir = "/tmp/ipcamera-images"

	err := os.MkdirAll(imageDir, 0755)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	camera := Camera{}
	camera.Init()

	history := History{
		imageDir: imageDir,
		camera:   &camera,
		interval: 10 * time.Minute,
		ttl:      3 * 24 * time.Hour,
	}
	history.Start()

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.Handle("/api/streamUri", &StreamUriHandler{camera: &camera})
	http.Handle("/api/history", history.ListHandler("/historyimages"))
	http.Handle("/historyimages/", history.FileHandler("/historyimages"))
	log.Println("Starting server, port:", port)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
}
