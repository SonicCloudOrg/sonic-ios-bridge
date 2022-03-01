package util

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const DownLoadTimeOut = 30 * time.Second

var versionMap = map[string]string{
	"12.5": "12.4",
	"15.3": "15.2",
}

var urlList = [...]string{"https://tool.appetizer.io", "https://github.com"}

func downloadZip(url, version string) error {
	if versionMap[version] != "" {
		version = versionMap[version]
	}
	_, errT := os.Stat(fmt.Sprintf(".sib/%s.zip", version))
	if errT != nil {
		_, err := os.Stat(".sib")
		if err != nil {
			os.MkdirAll(".sib", os.ModePerm)
		}
		client := http.Client{
			Timeout: DownLoadTimeOut,
		}
		res, err := client.Get(fmt.Sprintf("%s/JinjunHan/iOSDeviceSupport/raw/master/DeviceSupport/%s.zip", url, version))
		if err != nil {
			return err
		}
		defer res.Body.Close()
		r := bufio.NewReaderSize(res.Body, 32*1024)
		newFile, err := os.Create(fmt.Sprintf(".sib/%s.zip", version))
		w := bufio.NewWriter(newFile)
		io.Copy(w, r)
		abs, _ := filepath.Abs(newFile.Name())
		errZip := unzip(abs, ".sib")
		if errZip != nil {
			fmt.Println(errZip)
		}
	}
	return nil
}

func unzip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func LoadDevelopImage(version string) bool {
	var done = false
	for _, s := range urlList {
		err1 := downloadZip(s, version)
		if err1 == nil {
			done = true
			break
		}
	}
	return done
}
