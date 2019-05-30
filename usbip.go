package usbip

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type UsbCom struct {
	Flag1  byte
	Flag2  byte
	Length uint16
	Com1   byte
	Com2   byte
}

type UsbComEx struct {
	UsbCom
	Port uint32
}

type caInfo struct {
	Port    byte
	State   byte
	HasRead byte
	Name    [130]byte
	TaxCode [20]byte
}

type caInfos struct {
	Count uint8
	Info  [100]caInfo
}

type UsbIP struct {
	IPAddr string
}

func NewUsbIP(ipAddr string) *UsbIP { // 返回值指向UsbIP结构体的指针
	return &UsbIP{
		IPAddr: ipAddr,
	}
}

func cUSBCom(Com1 byte, Com2 byte) UsbCom {
	return UsbCom{0x01, 0x10, 6, Com1, Com2}
}

func dUSBCom(Com1 byte, Com2 byte, Port uint32) UsbComEx {
	return UsbComEx{UsbCom{0x01, 0x10, 10, Com1, Com2}, Port}
}

func (u UsbIP) Reboot() string {
	usbCom := cUSBCom(0x03, 0x08)
	sendBuf := new(bytes.Buffer)
	err1 := binary.Write(sendBuf, binary.BigEndian, &usbCom)
	if err1 != nil {
	}

	conn, err := net.Dial("tcp", u.IPAddr)
	if err != nil {
		fmt.Println(err)
		return "300U100129"
	}
	defer conn.Close()
	checkErr(err)
	n, err := conn.Write(sendBuf.Bytes())
	checkErr(err)
	fmt.Println("Write to server ", u.IPAddr, DecimalByteSlice2HexString(sendBuf.Bytes()))

	var readBuf [512]byte
	n, err = conn.Read(readBuf[0:])
	checkErr(err)
	fmt.Println("Reply from server ", u.IPAddr, DecimalByteSlice2HexString(readBuf[0:n]))
	return "0"
}

func (u UsbIP) Info() (map[string]string, error) {
	result := make(map[string]string)
	conn, err := net.Dial("tcp", u.IPAddr)
	if err != nil {
		log.Println(err)
		return result, nil
	}

	usbCom := cUSBCom(0x02, 0x10)
	WriteBuf := new(bytes.Buffer)
	err1 := binary.Write(WriteBuf, binary.BigEndian, &usbCom)
	if err1 != nil {
	}

	tn, err := conn.Write(WriteBuf.Bytes())
	checkErr(err)
	fmt.Println("Question to server ", u.IPAddr, tn, DecimalByteSlice2HexString(WriteBuf.Bytes()))
	///////////////////////////////////////////////////////////////////
	// Read Command
	var readBuf [16]byte
	nn, err := conn.Read(readBuf[0:6])
	fmt.Println("Reply from server ", nn, DecimalByteSlice2HexString(readBuf[0:6]))
	buf := bytes.NewBuffer(readBuf[0:nn])
	var command UsbCom
	binary.Read(buf, binary.BigEndian, &command)
	//////////////////////////////////////////////////////////////////
	// Read Data
	infoBuf := make([]byte, command.Length)
	var HadRead uint16
	HadRead = 0
	tn, err = conn.Read(infoBuf[HadRead:])

	for ii := int8(0); ii < int8(infoBuf[0]); ii++ {
		result[strconv.Itoa(int(ii)+1)] = strconv.Itoa(int(infoBuf[ii+1]))
	}
	return result, nil
}

func (u UsbIP) Close(UPort uint8) int {
	usbComEx := dUSBCom(0x02, 0x10, uint32(UPort))
	sendBuf := new(bytes.Buffer)
	err1 := binary.Write(sendBuf, binary.BigEndian, &usbComEx)
	if err1 != nil {
	}
	fmt.Println("Write to server ", DecimalByteSlice2HexString(sendBuf.Bytes()))
	return 0
}

func (u UsbIP) CaInfo() caInfos {
	var infos caInfos
	var readBuf [12]byte
	conn, err := net.Dial("tcp", u.IPAddr)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	usbCom := cUSBCom(0x02, 0x13)
	WriteBuf := new(bytes.Buffer)
	err1 := binary.Write(WriteBuf, binary.BigEndian, &usbCom)
	if err1 != nil {
	}
	n, err := conn.Write(WriteBuf.Bytes())
	checkErr(err)
	fmt.Println("question to server ", u.IPAddr, DecimalByteSlice2HexString(WriteBuf.Bytes()))
	///////////////////////////////////////////////////////////////////
	// Read Command
	nn, err := conn.Read(readBuf[0:6])
	fmt.Println("Reply from server ", nn)
	buf := bytes.NewBuffer(readBuf[0:nn])
	var command UsbCom
	err1 = binary.Read(buf, binary.BigEndian, &command)
	if err1 != nil {

	}
	//////////////////////////////////////////////////////////////////
	// Read Data
	infoBuf := make([]byte, command.Length)
	var HadRead uint16
	HadRead = 0
	for {
		n, err = conn.Read(infoBuf[HadRead:])
		if err != nil {
			fmt.Println(err)
			break
		}
		HadRead += uint16(n)
		if HadRead >= uint16(command.Length-6) {
			break
		}
	}

	fmt.Println("Reply from server ", HadRead)
	caBuf := bytes.NewBuffer(infoBuf[0:HadRead])
	err = binary.Read(caBuf, binary.LittleEndian, &infos)
	checkErr(err)
	return infos
}

func DecimalByteSlice2HexString(DecimalSlice []byte) string {
	var sa = make([]string, 0)
	for _, v := range DecimalSlice {
		sa = append(sa, fmt.Sprintf("%02X", v))
	}
	ss := strings.Join(sa, " ")
	return ss
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		//os.Exit(1)
	}
}
