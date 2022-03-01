package util

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadZip(url, version string) error {
	_, err := os.Stat(".sib")
	if err != nil {
		os.MkdirAll(".sib", os.ModePerm)
	}
	res, err := http.Get(fmt.Sprintf("%s/JinjunHan/iOSDeviceSupport/raw/master/DeviceSupport/%s.zip", url, version))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	r := bufio.NewReaderSize(res.Body, 32*1024)
	newFile, err := os.Create(fmt.Sprintf(".sib/%s.zip", version))
	w := bufio.NewWriter(newFile)
	io.Copy(w, r)
	return nil
}
