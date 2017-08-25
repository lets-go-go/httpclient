package main

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/lets-go-go/httpclient"
)

func main() {

	httpclient.Settings().SetProxy(httpclient.CustomProxy, "http://192.168.16.189:8080")

	// {
	// 	body, _ := httpclient.Get("http://www.baidu.com").Text()
	// 	fmt.Printf("body=%v", body)
	// }

	{

		v := url.Values{
			"appid": []string{"wx782c26e4c19acffb"},
			"fun":   []string{"new"},
			"lang":  []string{"zh_CN"},
			"_":     []string{strconv.FormatInt(time.Now().Unix(), 10)},
		}

		body, _ := httpclient.Post("https://login.weixin.qq.com/jslogin").AddFields(v).Text()
		fmt.Printf("body=\n%v", body)
	}

	// {
	// 	b := map[string]string{
	// 		"appid": "wx782c26e4c19acffb",
	// 		"fun":   "new",
	// 		"lang":  "zh_CN",
	// 		"_":     "1501747870",
	// 	}
	// 	body, _ := httpclient.Post("https://login.weixin.qq.com/jslogin").SetProxy("http://192.168.16.189:8080").SendBody(b).Text()
	// 	fmt.Printf("body=%v", body)
	// }

}
