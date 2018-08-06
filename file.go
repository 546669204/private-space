package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	bufferSize        = 1024 * 32
	encryptBufferSize = bufferSize + aes.BlockSize
)

type MyFileInfo struct {
	Size    int64
	MD5     string
	ModTime time.Time
}

func encryptFile(srcPath, dstPath, password string, progress handlerProgressModel) (MyFileInfo, error) {
	srcFile, err := os.OpenFile(srcPath, os.O_RDWR, 0666)
	if err != nil {
		return MyFileInfo{}, err
	}
	defer srcFile.Close()
	FileInfo, _ := srcFile.Stat()
	h := md5.New()
	srcFileInfo := MyFileInfo{
		Size:    FileInfo.Size(),
		ModTime: FileInfo.ModTime(),
	}
	dstFile, err := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err, 42)
		return MyFileInfo{}, err
	}
	defer dstFile.Close()

	key := []byte(password)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	ciphertext := make([]byte, encryptBufferSize)
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)

	tempByte := make([]byte, bufferSize)
	dstFile.Write(iv)
	tempPullSize := int64(0)
	for {
		n, err := srcFile.Read(tempByte)
		if err != nil {
			break
		}
		h.Write(tempByte[:n])
		tempPullSize += int64(n)
		log.Println(tempPullSize, srcFileInfo.Size)
		progress.Percentage = (tempPullSize * 100 / srcFileInfo.Size)
		progress.Conn.WriteJSON(gin.H{"id": progress.ID, "method": "progress", "data": progress})

		stream.XORKeyStream(ciphertext[aes.BlockSize:], tempByte[:n])

		if n < bufferSize {
			dstFile.Write(ciphertext[aes.BlockSize : aes.BlockSize+n])
			continue
		}
		dstFile.Write(ciphertext[aes.BlockSize:])
	}
	srcFileInfo.MD5 = hex.EncodeToString(h.Sum(nil))

	return srcFileInfo, nil
}
func decryptFile(srcPath, dstPath, password string, progress handlerProgressModel) error {
	srcFile, err := os.OpenFile(srcPath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	FileInfo, _ := srcFile.Stat()

	dstFile, err := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	key := []byte(password)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, bufferSize)

	tempByte := make([]byte, bufferSize)

	iv := make([]byte, aes.BlockSize)
	n, err := srcFile.Read(iv)
	if err != nil {
		return errors.New("iv不存在")
	}
	if n != aes.BlockSize {
		return errors.New("iv不存在")
	}
	stream := cipher.NewCFBDecrypter(block, iv)
	tempPullSize := int64(aes.BlockSize)
	for {
		n, err := srcFile.Read(tempByte)
		if err != nil {
			break
		}
		tempPullSize += int64(n)
		stream.XORKeyStream(ciphertext, tempByte[:n])
		dstFile.Write(ciphertext[:n])
		progress.Percentage = (tempPullSize * 100 / FileInfo.Size())
		progress.Conn.WriteJSON(gin.H{"id": progress.ID, "method": "progress", "data": progress})

	}
	return nil
}

func randomPath() string {
	asf := make([]byte, 64)
	rand.Read(asf)
	strp := hex.EncodeToString(asf)
	strr := make([]string, 0)
	for index := 0; index+32 < len(strp); index += 32 {
		strr = append(strr, strp[index:index+32])
	}
	return strings.Join(strr, "/")
}
func randomFileName() string {
	asf := make([]byte, 16)
	rand.Read(asf)
	strp := hex.EncodeToString(asf)
	return strp
}
func randomPassword() string {
	asf := make([]byte, 16)
	rand.Read(asf)
	h := md5.New()
	h.Write([]byte("hcaiyue"))
	h.Write(asf)
	return hex.EncodeToString(h.Sum(nil))
}
func getFileName(FilePath string) string {
	FilePath = strings.Replace(FilePath, "\\", "/", -1)
	return FilePath[strings.LastIndex(FilePath, "/")+1:]
}
func getFileExt(FilePath string) string {
	return FilePath[strings.LastIndex(FilePath, "."):]
}
