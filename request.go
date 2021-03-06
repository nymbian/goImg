package main

import (
	"bytes"
	"crypto/md5"

	"encoding/hex"
	"fmt"

	"io"
	"io/ioutil"
	"log"

	"net/http"

	"os"

	"strings"
)

func getImgFromOtherServer(imageId string) bool {
	result := false
	servers := conf.Servers
	if len(servers) == 0 {
		return result
	}
	fmt.Println(servers)
	for _, v := range servers {
		if v == conf.ListenAddr {
			continue
		} else {
			if getImg(imageId, v) {
				result = true
				break
			}
		}
	}
	return result
}

func getImg(imageId string, server string) bool {

	result := false
	url := "http://" + server + "/" + imageId + "_sync"
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return result
	}

	defer resp.Body.Close()
	fmt.Println(resp)

	if resp.StatusCode == 200 {

		urlPath := strings.Split(url, "/")
		var name string
		if len(urlPath) > 1 {
			name = urlPath[len(urlPath)-1]
		}
		name = getTempFilePath(name)
		out, err := os.Create(name)
		if err != nil {
			log.Println(err)
			return result
		}
		defer out.Close()

		pix, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return result
		}
		_, err = io.Copy(out, bytes.NewReader(pix))
		if err != nil {
			log.Println(err)
			return result
		}
		file, err := os.Open(name)
		if err != nil {
			log.Println(err)
			return result
		}
		defer file.Close()
		//file type
		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			log.Println(err)
			return result
		}
		fileType := http.DetectContentType(buff)
		fmt.Println(fileType)

		if checkFileType(fileType) == false {

			return result
		}

		if _, err = file.Seek(0, 0); err != nil {
			log.Println(err)
		}

		md5h := md5.New()
		io.Copy(md5h, file)

		imageIdNew := hex.EncodeToString(md5h.Sum(nil))

		if imageIdNew != imageId {

			return false
		}

		fmt.Println(imageId)
		//imageId = "86427d1debefe65f0da3a7affdc204f2"

		path := GetPathByMd5(imageId)
		//path := "E:/1037u/1.gif"

		err = MkdirByMd5(imageId)
		if err != nil {
			log.Println(err)
		}

		if FileExist(path) {

			return true
		}

		if _, err = file.Seek(0, 0); err != nil {
			log.Println(err)
		}

		fmt.Println(path)
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, fileAuth)
		if err != nil {
			log.Println(err)

			return result
		}
		defer f.Close()
		bytesWritten, err := io.Copy(f, file)
		if err != nil {
			log.Println(err)

			return result
		}
		fmt.Println(bytesWritten)
		result = true
	}
	return result

}
