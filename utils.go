package sptty

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path"

	"github.com/rs/xid"
)

type RequestError struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func NewRequestError(code string, msg string) []byte {
	b, _ := json.Marshal(RequestError{
		Code: code,
		Msg:  msg,
	})

	return b
}

func GenerateUID() string {
	return xid.New().String()
}

func Sha1(data string) string {
	s := sha1.New()
	s.Write([]byte(data))
	return hex.EncodeToString(s.Sum([]byte("")))
}

func Sha256(data string) string {
	s := sha256.New()
	s.Write([]byte(data))
	return hex.EncodeToString(s.Sum([]byte("")))
}

func RandomFilename(rawFile string) string {
	id := GenerateUID()
	fileEx := path.Ext(path.Base(rawFile))

	if fileEx == "" || fileEx == "." {
		return id
	} else {
		return fmt.Sprintf("%s%s", id, fileEx)
	}
}

func ArrayContains(arr []interface{}, s interface{}) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}

	return false
}
