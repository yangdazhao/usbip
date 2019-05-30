package usbip_test

import (
	"fmt"
	"testing"
	"usbip"
)

func TestHelloWorld(t *testing.T) {
	t.Log("hello world")
}

func TestCaInfo(t *testing.T) {
	usb := usbip.NewUsbIP("182.18.75.50:10007")
	infos := usb.CaInfo()
	fmt.Println(infos.Count)
	fmt.Println(infos.Info)
}

func TestCloseCaInfo(t *testing.T) {
	usb := usbip.NewUsbIP("182.18.75.50:10007")
	fmt.Println(usb.Close(1))
}
