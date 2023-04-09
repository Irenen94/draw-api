package util

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type idutil struct{}

var IDUtil idutil

func (i *idutil) GetImageId(base64str, source, sourceApp string) string {
	strbuilder := strings.Builder{}
	strbuilder.WriteString(base64str)
	strbuilder.WriteString(source)
	strbuilder.WriteString(sourceApp)

	has := md5.Sum([]byte(strbuilder.String()))
	return fmt.Sprintf("%x", has)
}

func (i *idutil) GetImageId2(base64str string) string {
	strbuilder := strings.Builder{}
	strbuilder.WriteString(base64str)
	has := md5.Sum([]byte(strbuilder.String()))
	return fmt.Sprintf("%x", has)
}

func (i *idutil) GetPersonId(account_id, source, sourceApp string) string {
	strbuilder := strings.Builder{}
	strbuilder.WriteString(account_id)
	strbuilder.WriteString(source)
	strbuilder.WriteString(sourceApp)
	strbuilder.WriteString(strconv.FormatInt(time.Now().UnixNano()/1e6, 10))

	has := md5.Sum([]byte(strbuilder.String()))
	return fmt.Sprintf("%x", has)
}


func RandStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}
