// Copyright 2017 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strconv"
)

// RandID return a rand string used in frp.
func RandID() (id string, err error) {
	return RandIDWithLen(16)
}

// RandIDWithLen return a rand string with idLen length.
func RandIDWithLen(idLen int) (id string, err error) {
	if idLen <= 0 {
		return "", nil
	}
	b := make([]byte, idLen/2+1)
	_, err = rand.Read(b)
	if err != nil {
		return
	}

	id = fmt.Sprintf("%x", b)
	return id[:idLen], nil
}

func GetAuthKey(token string, timestamp int64) (key string) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(token))
	md5Ctx.Write([]byte(strconv.FormatInt(timestamp, 10)))
	data := md5Ctx.Sum(nil)
	return hex.EncodeToString(data)
}

func GenerateResponseErrorString(summary string, err error, detailed bool) string {
	if detailed {
		return err.Error()
	}
	return summary
}

func ConstantTimeEqString(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
