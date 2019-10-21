package test_test

import (
	"fmt"
	"github.com/yangdazhao/usbip"
	"testing"
)

func TestCaInfo(t *testing.T) {
	usb := usbip.NewUsbIP("182.18.75.50:10007")
	infos := usb.CaInfo()
	fmt.Println(infos.Count)
	fmt.Println(infos.Info)
}

func TestInfo(t *testing.T) {
	usb := usbip.NewUsbIP("182.18.75.50:10007")
	infos, _ := usb.Info()
	fmt.Println(infos)
}

func TestCloseCaInfo(t *testing.T) {
	usb := usbip.NewUsbIP("182.18.75.50:10003")
	fmt.Println(usb.Close(27))
}
