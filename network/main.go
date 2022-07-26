package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	GetOpenIMToken()
}

// 向本地Flask获取手机号码
func GetOpenIMToken() {
	client := &http.Client{}
	// json序列化
	params := make(map[string]interface{})
	params["secret"] = "tuoyun"
	params["platform"] = 5
	params["uid"] = "163825307760"
	con, err := json.Marshal(params)
	if err != nil {
		fmt.Println("con: ", err)
	}
	var jsonStr = []byte(con)

	req, err := http.NewRequest("POST", "http://106.52.92.189:10000/auth/user_token", bytes.NewBuffer(jsonStr))
	if err != nil {
		// handle error
		fmt.Println("req: ", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		// handle error
		fmt.Println("resp: ", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println("body: ", err)
	}

	fmt.Println(string(body))
}
