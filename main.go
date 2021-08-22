package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/use-go/onvif"
	"github.com/use-go/onvif/media"
	xsd_onvif "github.com/use-go/onvif/xsd/onvif"
)

func findNetwork() (*net.Interface, *net.IPNet, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}

	for _, networkInterface := range interfaces {
		addrs, _ := networkInterface.Addrs()
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.IsGlobalUnicast() && ipNet.IP.To4() != nil {
				return &networkInterface, ipNet, nil
			}
		}
	}

	return nil, nil, fmt.Errorf("No network interface is available")
}

// StreamUriHandler
type StreamUriHandler struct {
	networkInterface string
}

func (s *StreamUriHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	fmt.Printf("Searching camera on network interface %s\n", s.networkInterface)
	devs := onvif.GetAvailableDevicesAtSpecificEthernetInterface(s.networkInterface)
	if len(devs) == 0 {
		fmt.Println("Camera not found")
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(s.errorResponse("Camera not found"))
		return
	}
	dev := devs[0]
	fmt.Println("Camera found")

	profileToken, err := s.getProfileToken(dev)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(s.errorResponse(err.Error()))
		return
	}
	fmt.Println("ProfileToken:", profileToken)

	streamUri, err := s.getStreamUri(dev, profileToken)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(s.errorResponse(err.Error()))
		return
	}
	fmt.Println("StreamUri:", streamUri)

	body := map[string]interface{}{
		"result": "OK",
		"uri":    streamUri,
	}
	bytes, _ := json.Marshal(body)
	writer.Write(bytes)
}

func (s *StreamUriHandler) read(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	return buf.String()
}

type GetProfilesResponseWrapper struct {
	XMLName             xml.Name                  `xml:"Envelope"`
	GetProfilesResponse media.GetProfilesResponse `xml:"Body>GetProfilesResponse"`
}

func (s *StreamUriHandler) getProfileToken(device onvif.Device) (string, error) {
	fmt.Println("Getting Profiles")
	res, err := device.CallMethod(media.GetProfiles{})
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("Status Code: %d", res.StatusCode)
	}
	str := s.read(res.Body)
	wrapper := GetProfilesResponseWrapper{}
	err = xml.Unmarshal([]byte(str), &wrapper)
	if err != nil {
		return "", err
	}
	profiles := wrapper.GetProfilesResponse.Profiles
	if len(profiles) == 0 {
		return "", fmt.Errorf("no profile found")
	}
	return string(profiles[0].Token), nil
}

type GetStreamUriResponseWrapper struct {
	XMLName              xml.Name                   `xml:"Envelope"`
	GetStreamUriResponse media.GetStreamUriResponse `xml:"Body>GetStreamUriResponse"`
}

func (s *StreamUriHandler) getStreamUri(device onvif.Device, profileToken string) (string, error) {
	fmt.Println("Getting Stream URI")
	getStreamUri := media.GetStreamUri{
		StreamSetup: xsd_onvif.StreamSetup{
			Stream: xsd_onvif.StreamType("RTP-Unicast"),
			Transport: xsd_onvif.Transport{
				Protocol: "RTSP",
			},
		},
		ProfileToken: xsd_onvif.ReferenceToken(profileToken),
	}
	res, err := device.CallMethod(getStreamUri)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("Status Code: %d", res.StatusCode)
	}
	str := s.read(res.Body)
	wrapper := GetStreamUriResponseWrapper{}
	err = xml.Unmarshal([]byte(str), &wrapper)
	if err != nil {
		return "", err
	}
	uri := string(wrapper.GetStreamUriResponse.MediaUri.Uri)
	return uri, nil
}

func (s *StreamUriHandler) errorResponse(message string) []byte {
	body := map[string]interface{}{
		"result":  "NG",
		"message": message,
	}
	bytes, _ := json.Marshal(body)
	return bytes
}

func main() {
	const port = 3333
	networkInterface, ipNet, err := findNetwork()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Server Info:")
	fmt.Printf("  Network interface: %s\n", networkInterface.Name)
	fmt.Printf("  Ip address       : %s\n", ipNet.IP)
	fmt.Printf("  Port             : %d\n", port)

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.Handle("/streamUri", &StreamUriHandler{
		networkInterface: networkInterface.Name,
	})
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
}
