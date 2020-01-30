package sptty

import (
	"encoding/json"
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
