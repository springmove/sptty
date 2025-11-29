package sptty

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"runtime"
	"strings"

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
	_, _ = s.Write([]byte(data))
	return hex.EncodeToString(s.Sum([]byte("")))
}

func Sha256(data string) string {
	s := sha256.New()
	_, _ = s.Write([]byte(data))
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

// param1: content body
// param2: mime type
func GetUrlImage(url string) ([]byte, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	mime := resp.Header.Get("content-type")
	vals := strings.Split(mime, "/")
	if len(vals) > 1 {
		mime = vals[1]
	}

	return body, mime, nil
}

func CurrentFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	vals := strings.Split(f.Name(), ".")
	return vals[len(vals)-1]
}
