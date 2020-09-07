package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/teris-io/shortid"
	tnet "github.com/toolkits/net"
)

//RespJSON the common json
type RespJSON struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var (
	once     sync.Once
	clientIP = "127.0.0.1"
)

//CheckErr panic if err is not nil
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// GetLocalIP 获取本地内网IP
func GetLocalIP() string {
	once.Do(func() {
		ips, _ := tnet.IntranetIP()
		if len(ips) > 0 {
			clientIP = ips[0]
		} else {
			clientIP = "127.0.0.1"
		}
	})
	return clientIP
}

//JSON2String convert the json object to string
func JSON2String(j interface{}) (str string) {
	bytes, err := json.Marshal(j)
	if err != nil {
		return "{}"
	}
	str = (string)(bytes)
	return
}

//StringToStruct convert the string to struct
func StringToStruct(data string, desc interface{}) (err error) {
	if err = json.Unmarshal(([]byte)(data), desc); err != nil {
		return
	}
	return
}

// GenShortID 生成一个id
func GenShortID() string {
	sid, _ := shortid.Generate()
	return sid
}

// GenUUID 生成随机字符串，eg: 76d27e8c-a80e-48c8-ad20-e5562e0f67e4
func GenUUID() string {
	u, _ := uuid.NewRandom()
	return u.String()
}

// Sha256Encode sha256 加密
func Sha256Encode(origin string) string {
	sum := sha256.Sum256(([]byte)(origin))
	return fmt.Sprintf("%x", sum)
}
