package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

type History struct {
	imageDir string
	camera   *Camera
	interval time.Duration
	ttl      time.Duration
}

func (h *History) Start() {
	ticker := time.NewTicker(h.interval)
	go func() {
		h.saveImage()
		for {
			select {
			case <-ticker.C:
				h.saveImage()
				h.cleanupDir()
			}
		}
	}()
}

func (h *History) saveImage() {

	uri, err := h.camera.GetSnapshotUri()
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Getting image from ", uri)
	resp, err := http.Get(uri)
	if err != nil {
		log.Println("Failed to get image:", err.Error())
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read body:", err.Error())
		return
	}
	timestamp := time.Now().Unix()
	path := h.imageDir + "/" + fmt.Sprint(timestamp) + ".jpg"
	log.Println("Writing to file:", path, ", Size:", len(buf), "bytes")
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		log.Println("Failed to create file:", err.Error())
		return
	}
	fw := bufio.NewWriter(f)
	_, err = fw.Write(buf)
	if err != nil {
		log.Println("Failed to write file:", err.Error())
		return
	}
	err = fw.Flush()
	if err != nil {
		log.Println("Failed to flush file:", err.Error())
		return
	}
}

func (h *History) cleanupDir() error {
	imageinfos, err := getImageList(h.imageDir)
	if err != nil {
		return err
	}
	for _, info := range imageinfos {
		now := time.Now().Unix()
		if info.timestamp < now-int64(h.ttl.Seconds()) {
			path := h.imageDir + "/" + info.name
			log.Println("Removing", path)
			err = os.Remove(path)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (h *History) FileHandler(mountPath string) http.Handler {
	fileServer := http.FileServer(http.Dir(h.imageDir))
	return http.StripPrefix(mountPath, &fileHandler{fileServer: fileServer})
}

type fileHandler struct {
	fileServer http.Handler
}

func (h *fileHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("Cache-Control", "public, max-age=262800, immutable")
	h.fileServer.ServeHTTP(writer, req)
}

func (h *History) ListHandler(urlFilesPath string) http.Handler {
	return &listHandler{
		imageDir:  h.imageDir,
		filesPath: urlFilesPath,
	}
}

type listHandler struct {
	imageDir  string
	filesPath string
	Handler
}

func (h *listHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	imageInfos, err := getImageList((h.imageDir))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(h.errorMessage(err.Error()))
		return
	}

	respItems := [](interface{}){}
	for _, info := range imageInfos {
		item := map[string]interface{}{
			"timestamp": info.timestamp,
			"path":      h.filesPath + "/" + info.name,
		}
		respItems = append(respItems, item)
	}
	body := map[string]interface{}{
		"files": respItems,
	}
	h.responseJson(writer, body)
}

type imageInfo struct {
	timestamp int64
	name      string
}

func getImageList(dir string) ([]imageInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filenameMatcher := regexp.MustCompile(`^(\d+)\.jpg$`)
	list := []imageInfo{}
	for _, file := range files {
		groups := filenameMatcher.FindStringSubmatch(file.Name())
		if file.IsDir() || groups == nil {
			continue
		}
		timestamp, _ := strconv.ParseInt(groups[1], 10, 64)
		item := imageInfo{
			timestamp: timestamp,
			name:      file.Name(),
		}
		list = append(list, item)
	}
	return list, nil
}
