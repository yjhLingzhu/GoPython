package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	data := `{"errcode":0,"errmsg":"ok","phone_info":{"phoneNumber":"17875513387","purePhoneNumber":"17875513387","countryCode":"86","watermark":{"timestamp":1640739648,"appid":"wx62b400658529b5a1"}}}`
	//json str 转map
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(data), &dat); err == nil {
		fmt.Println("==============json str 转map=======================")
		fmt.Println(dat, dat["errcode"])
		fmt.Printf("%T", dat["errcode"])
		if value, flag := dat["errcode"].(float64); value == 0 && flag {
			a, err := dat["phone_info"].(map[string]interface{})
			fmt.Println(err)
			// fmt.Println(dat["phone_info"].(map[string]string))
			fmt.Println(a)
			fmt.Println(a["purePhoneNumber"])
		}
	}

	var s interface{} = ""
	if s == "" {
		fmt.Println("哈哈哈")
	}
}
