package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	var err error
	config, err := LoadConfig()
	if config == nil {
		os.Exit(1)
	}

	err = os.MkdirAll(config.History.ImageDir, 0755)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	camera := Camera{}
	camera.Init()

	history := History{
		imageDir: config.History.ImageDir,
		camera:   &camera,
		interval: config.History.Interval.Duration,
		ttl:      config.History.Ttl.Duration,
	}
	history.Start()

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.Handle("/api/streamUri", &StreamUriHandler{camera: &camera})
	http.Handle("/api/history", history.ListHandler("/historyimages"))
	http.Handle("/historyimages/", history.FileHandler("/historyimages"))
	log.Println("Starting server, port:", config.Port)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.Port), nil)
}
