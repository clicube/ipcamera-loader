package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"

	"github.com/use-go/onvif"
	"github.com/use-go/onvif/media"
	xsd_onvif "github.com/use-go/onvif/xsd/onvif"
)

type Camera struct {
	device   *onvif.Device
	probeReq chan interface{}
	probeRes chan error
}

func (c *Camera) Init() error {
	c.probeReq = make(chan interface{}, 8)
	c.probeRes = make(chan error)
	go c.probeLoop()
	return nil
}

func (c *Camera) probeLoop() {
	for {
		<-c.probeReq
		device, err := c.doProbe()
		c.device = device
		c.probeRes <- err
		// consume all request
		for len(c.probeReq) > 0 {
			<-c.probeReq
			c.probeRes <- err
		}
	}
}

func (c *Camera) doProbe() (*onvif.Device, error) {
	log.Println("Probing")
	netif, _, err := c.findNetwork()
	if err != nil {
		return nil, err
	}

	log.Println("Searching camera on network interface", netif.Name)
	devs := onvif.GetAvailableDevicesAtSpecificEthernetInterface(netif.Name)
	if len(devs) == 0 {
		return nil, fmt.Errorf("Camera not found")
	}
	dev := devs[0]
	log.Println("Camera found")
	return &dev, nil
}

func (c *Camera) probe() error {
	c.probeReq <- 0
	return <-c.probeRes
}

func (c *Camera) GetSnapshotUri() (string, error) {
	streamUri, err := c.GetStreamUri()
	if err != nil {
		return "", err
	}
	u, err := url.Parse(streamUri)
	if err != nil {
		return "", err
	}
	snapshotUri := "http://" + u.Host + "/snapshot"
	return snapshotUri, nil
}

func (c *Camera) GetStreamUri() (string, error) {
	if c.device == nil {
		err := c.probe()
		if err != nil {
			return "", err
		}
	}
	profileToken, err := c.getProfileToken(c.device)
	if err != nil {
		log.Println(err.Error())
		log.Println("Re-probing")
		err = c.probe()
		if err != nil {
			return "", err
		}
		profileToken, err = c.getProfileToken(c.device)
		if err != nil {
			return "", err
		}
	}
	log.Println("ProfileToken:", profileToken)

	streamUri, err := c.getStreamUri(c.device, profileToken)
	if err != nil {
		return "", nil
	}
	log.Println("StreamUri:", streamUri)

	return streamUri, nil
}

type GetProfilesResponseWrapper struct {
	XMLName             xml.Name                  `xml:"Envelope"`
	GetProfilesResponse media.GetProfilesResponse `xml:"Body>GetProfilesResponse"`
}

func (c *Camera) getProfileToken(device *onvif.Device) (string, error) {
	log.Println("Getting Profiles")
	res, err := device.CallMethod(media.GetProfiles{})
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("Status Code: %d", res.StatusCode)
	}
	str := c.read(res.Body)
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

func (c *Camera) getStreamUri(device *onvif.Device, profileToken string) (string, error) {
	log.Println("Getting Stream URI")
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
	str := c.read(res.Body)
	wrapper := GetStreamUriResponseWrapper{}
	err = xml.Unmarshal([]byte(str), &wrapper)
	if err != nil {
		return "", err
	}
	uri := string(wrapper.GetStreamUriResponse.MediaUri.Uri)
	return uri, nil
}

func (c *Camera) read(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	return buf.String()
}

func (c *Camera) findNetwork() (*net.Interface, *net.IPNet, error) {
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
