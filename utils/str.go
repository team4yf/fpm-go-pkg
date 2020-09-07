package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/oliveagle/jsonpath"
)

//GetBodyByte 获取io流的字节数据
func GetBodyByte(readCloser io.ReadCloser) ([]byte, error) {
	body, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return nil, err
	}
	defer readCloser.Close()

	return body, nil
}

//GetBodyString 获取io流中的数据并转成string形式
func GetBodyString(readCloser io.ReadCloser) (string, error) {
	body, err := GetBodyByte(readCloser)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

//GetBodyMap 获取io流中的数据，并转成map形式
func GetBodyMap(readCloser io.ReadCloser) (map[string]interface{}, error) {
	var body map[string]interface{}
	err := json.NewDecoder(readCloser).Decode(&body)
	if err != nil {
		return nil, err
	}
	defer readCloser.Close()

	return body, nil
}

//GetBodyStruct 获取io流中的数据，并转成指定的类型
func GetBodyStruct(readCloser io.ReadCloser, obj interface{}) error {
	body, err := GetBodyByte(readCloser)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &obj); err != nil {
		return err
	}
	return nil
}

//GetJSONPathData 获取jsonpath对应的字符
func GetJSONPathData(data, jp string) (interface{}, error) {
	var jsonData interface{}
	json.Unmarshal([]byte(data), &jsonData)
	res, err := jsonpath.JsonPathLookup(jsonData, jp)
	if err != nil {
		return nil, err
	}
	return res, nil
}
