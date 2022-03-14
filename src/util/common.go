package util

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const DownLoadTimeOut = 30 * time.Second

var versionMap = map[string]string{
	"12.5": "12.4",
	"15.3": "15.2",
}

var urlList = [...]string{"https://tool.appetizer.io/JinjunHan", "https://code.aliyun.com/hanjinjun", "https://github.com/JinjunHan"}

func downloadZip(url, version string) (string, error) {
	if versionMap[version] != "" {
		version = versionMap[version]
	}
	f, err := os.Stat(".sib")
	if err != nil {
		os.MkdirAll(".sib", os.ModePerm)
		f, err = os.Stat(".sib")
	}
	localAbs, _ := filepath.Abs(f.Name())
	_, errT := os.Stat(fmt.Sprintf(".sib/%s.zip", version))
	if errT != nil {
		client := http.Client{
			Timeout: DownLoadTimeOut,
		}
		res, err := client.Get(fmt.Sprintf("%s/iOSDeviceSupport/raw/master/DeviceSupport/%s.zip", url, version))
		if err != nil {
			return "", err
		}
		defer res.Body.Close()
		r := bufio.NewReaderSize(res.Body, 32*1024)
		newFile, err := os.Create(fmt.Sprintf(".sib/%s.zip", version))
		w := bufio.NewWriter(newFile)
		io.Copy(w, r)
		abs, _ := filepath.Abs(newFile.Name())
		errZip := unzip(abs, ".sib", version)
		if errZip != nil {
			os.Remove(newFile.Name())
			return "", errZip
		}
	}
	return localAbs, nil
}

func unzip(zipFile, destDir, version string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		var fpath string
		if strings.HasPrefix(f.Name, version) && f.FileInfo().IsDir() {
			fpath = filepath.Join(destDir, version)
		} else {
			fpath = filepath.Join(destDir, version+"/"+path.Base(f.Name))
		}
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

func LoadDevelopImage(version string) (string, bool) {
	var done = false
	var path = ""
	for _, s := range urlList {
		p, err1 := downloadZip(s, version)
		if err1 == nil {
			path = p
			done = true
			break
		}
	}
	return path, done
}
