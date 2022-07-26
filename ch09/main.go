package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"os"
)

type Data struct {
	Status int `json:"status"`
	Msg    int `json:"msg"`
}

// 画渐变色 从蓝到白，从上到下
func main() {
	// gradientRamp()
	access_token := HttpGet()
	code := "70c4d26e36ad530fea5451037b737f5d4cc66396139221f63dd90892a9897e7d"
	// fmt.Println(access_token)
	GetPhoneNumber(access_token, code)

}

func gradientRamp() {
	RGB1 := []int{255, 255, 255} // 白色
	RGB2 := []int{102, 181, 242} //淡蓝色
	width := 413
	height := 531

	a1, a2, a3, b1, b2, b3 := RGB2[0], RGB2[1], RGB2[2], RGB1[0], RGB1[1], RGB1[2]
	// 相差的rgb
	r, g, b := (b1 - a1), (b2 - a2), (b3 - a3)
	fmt.Println(r, g, b)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < img.Bounds().Dy(); y++ { // 更换背景
		x1 := y
		t := (float64(x1)) / (float64(height))
		rgb := []uint8{uint8(float64(a1) + float64(r)*t), uint8(float64(a2) + float64(g)*t), uint8(float64(a3) + float64(b)*t)}
		fmt.Println(rgb)
		for x := 0; x < img.Bounds().Dx(); x++ {

			// fmt.Println(t, rgb)
			img.Set(x, y, color.RGBA{rgb[0], rgb[1], rgb[2], 255})
		}
	}
	outFile, err := os.Create("./1.jpg") // 创建一张空的图片文件
	defer outFile.Close()
	if err != nil {
		panic(err)
	}

	err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 100}) // 使用jpeg的编码方式存储
	if err != nil {
		panic(err)
	}
}

func HttpGet() string {
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	url = fmt.Sprintf(url, "wx62b400658529b5a1", "db3a95ff2e8abc96067f6cbe0efe578d")
	// fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	type Respone struct {
		Access_token string `json:"access_token"`
		Expires_in   int64  `json:"expires_in"`
	}
	var res Respone
	err = json.Unmarshal(s, &res)
	if err != nil {
		panic(err)
	}
	fmt.Printf(string(s))
	fmt.Println()
	// fmt.Println(res.Access_token)
	return res.Access_token
}

// 获取手机号码
func GetPhoneNumber(access_token string, code string) interface{} {
	Url := "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s"
	Url = fmt.Sprintf(Url, access_token)
	// fmt.Println("Url: ", Url)
	// json序列化
	postData := "{\"code\":" + "\"" + code + "\"" + "}"
	fmt.Println(postData)
	var jsonStr = []byte(postData)
	req, err := http.NewRequest("POST", Url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	//json str 转map
	var dat map[string]interface{}
	err = json.Unmarshal(body, &dat)
	if err != nil {
		panic(err)
	}
	phone_info, isSuccess := dat["phone_info"].(map[string]interface{})
	if isSuccess {
		return phone_info["purePhoneNumber"]
	}

	// fmt.Println(string(body))
	return ""
}
