package kutil

import (
	crand "crypto/rand"
	"math/rand"
	"encoding/base64"

	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(length int) (str string, err error) {

	bytes := make([]byte, length)
	_, err = crand.Read(bytes)
	if nil != err {
		return
	}

	str = base64.URLEncoding.EncodeToString(bytes)
	return
}

func RandomRange(min, max int) int {
	return rand.Intn(max - min) + min
}