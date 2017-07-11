package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"crypto/md5"
	"io"
	"math"
	"math/rand"
	"strconv"
)

func GetPathByMd5(imageId string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", conf.Storage, imageId[0:8], imageId[8:16], imageId[16:24], imageId[24:32])
}

func FileExist(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		return false
	} else {
		return true
	}
}

func MkdirByMd5(imageId string) error {
	return os.MkdirAll(fmt.Sprintf("%s/%s/%s/%s", conf.Storage, imageId[0:8], imageId[8:16], imageId[16:24]), fileAuth)
}

func getTempFilePath(name string) string {
	random := strconv.Itoa(rand.Int())
	if name == "" {
		name = "base64.jpg"
	}
	path := conf.Storage + "/temp/" + random + "/"
	os.MkdirAll(path, fileAuth)
	return path + name

}

const filechunk = 8192 // we settle for 8KB

func getFileMd5(path string) string {

	file, err := os.Open(path)

	if err != nil {
		panic(err.Error())
	}

	defer file.Close()

	// calculate the file size
	info, _ := file.Stat()

	filesize := info.Size()

	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))

	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)

		file.Read(buf)
		io.WriteString(hash, string(buf)) // append into the hash
	}

	//fmt.Printf("%s checksum is %x\n", file.Name(), hash.Sum(nil))
	return hex.EncodeToString(hash.Sum(nil))

}
