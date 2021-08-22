package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/beevik/etree"
	"github.com/use-go/onvif"
	onvif_media "github.com/use-go/onvif/media"
	onvif_xsd_onvif "github.com/use-go/onvif/xsd/onvif"
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

//StreamURI ...
type StreamURI struct {
	networkInterface string
}

func (s *StreamURI) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	devices := onvif.GetAvailableDevicesAtSpecificEthernetInterface(s.networkInterface)
	if len(devices) == 0 {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(s.errorResponse("Camera not found."))
		return
	}
	device := devices[0]

	profileToken, err := s.getProfileToken(device)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(s.errorResponse(err.Error()))
		return
	}

	streamURI, err := s.getStreamURI(device, profileToken)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(s.errorResponse(err.Error()))
		return
	}

	body := map[string]interface{}{
		"result": "OK",
		"uri":    streamURI,
	}
	bytes, _ := json.Marshal(body)
	writer.Write(bytes)
}

func (s *StreamURI) read(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	return buf.String()
}

func (s *StreamURI) getProfileToken(device onvif.Device) (string, error) {
	res, err := device.CallMethod(onvif_media.GetProfiles{})
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("Failed to get profiles")
	}
	xml := s.read(res.Body)
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		return "", err

	}
	profileElements := doc.Root().FindElements("./Body/GetProfilesResponse/Profiles")
	for _, attr := range profileElements[0].Attr {
		if attr.Key == "token" {
			return attr.Value, nil
		}
	}
	return "", fmt.Errorf("Failed to read profile token")
}

func (s *StreamURI) getStreamURI(device onvif.Device, profileToken string) (string, error) {
	res, err := device.CallMethod(s.createGetStreamURIMessage(profileToken))
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("Failed to get streamURI")
	}
	xml := s.read(res.Body)
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		return "", err

	}
	element := doc.Root().FindElement("./Body/GetStreamUriResponse/MediaUri/Uri")
	if element == nil {
		return "", fmt.Errorf("Failed to read streamURI")

	}
	return element.Text(), nil

}

func (s *StreamURI) createGetStreamURIMessage(profileToken string) onvif_media.GetStreamUri {
	return onvif_media.GetStreamUri{
		StreamSetup: onvif_xsd_onvif.StreamSetup{
			Stream: onvif_xsd_onvif.StreamType("RTP-Unicast"),
			Transport: onvif_xsd_onvif.Transport{
				Protocol: "RTSP",
			},
		},
		ProfileToken: onvif_xsd_onvif.ReferenceToken(profileToken),
	}

}

func (s *StreamURI) errorResponse(message string) []byte {
	body := map[string]interface{}{
		"result":  "NG",
		"message": message,
	}
	bytes, _ := json.Marshal(body)
	return bytes
}

func main() {
	networkInterface, ipNet, err := findNetwork()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Network interface: %s\n", networkInterface.Name)
	fmt.Printf("Ip address       : %s\n", ipNet.IP)

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.Handle("/streamUri", &StreamURI{
		networkInterface: networkInterface.Name,
	})
	http.ListenAndServe("0.0.0.0:3333", nil)
}
