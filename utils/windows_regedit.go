package utils

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"strconv"
)

func SetProxy(proxyIp, pacUrl string) {
	/*
		- 0F	全部开启(ALL)
		- 01	全部禁用(Off)
		- 03	使用代理服务器(ProxyOnly)
		- 05	使用自动脚本(PacOnly)；
		- 07	使用脚本和代理(ProxyAndPac)
		- 09	打开自动检测设置(D)；
		- 0B	打开自动检测并使用代理(DIP)
		- 0D	打开自动检测并使用脚本(DS)；
	*/

	data := "460000003A160000" + "05" +
		"000000" + fmt.Sprintf("%02x", len(proxyIp)) +
		"000000" + hex.EncodeToString([]byte(proxyIp)) + "070000003c6c6f63616c3e2b" +
		"000000" + hex.EncodeToString([]byte(pacUrl)) + "00000000000000000000000000000000000000000000000000000000000000"

	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings\Connections`, registry.ALL_ACCESS)
	defer key.Close()
	if err != nil {
		panic(err)
	}

	//把16进制字符串转为byte切片
	bytedata := []byte{}
	for i := 0; i < len(data)-2; i = i + 2 {
		t := data[i : i+2]
		n, err := strconv.ParseUint(t, 16, 32)
		if err != nil {
			panic(err)
		}
		n2 := byte(n)
		bytedata = append(bytedata, n2)
	}

	err = key.SetBinaryValue("DefaultConnectionSettings", bytedata)
	if err != nil {
		panic(err)
	}
	fmt.Println("代理设置成功")
}
