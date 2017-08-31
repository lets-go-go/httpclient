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
	testDownload()
}

func testPost() {

	v := url.Values{
		"appid": []string{"wx782c26e4c19acffb"},
		"fun":   []string{"new"},
		"lang":  []string{"zh_CN"},
		"_":     []string{strconv.FormatInt(time.Now().Unix(), 10)},
	}

	// b := []int{1, 2, 3}

	body, _ := httpclient.Post("http://httpbin.org/post").AddFields(v).AttachFile("file1", "d:/1.jpg", "1.jpg").Text()
	fmt.Printf("body=%v\n", body)

	fmt.Println("-------------------------------------")
}

func testGet() {

	rsp, _ := httpclient.Head("http://httpbin.org/get").Execute()

	fmt.Printf("body=%+v\n", rsp.Header)
	fmt.Println("-------------------------------------")
}

func testHead() {
	rsp, _ := httpclient.Options("http://httpbin.org/get").Execute()

	fmt.Printf("body=%+v\n", rsp.Header)
	fmt.Println("-------------------------------------")
}

func testGetBody() {

	body, _ := httpclient.Get("http://httpbin.org/get").Text()
	fmt.Printf("body=%v\n", body)

	fmt.Println("-------------------------------------")
}

func testPostWithBody() {

	// b := map[string]string{
	// 	"appid": "wx782c26e4c19acffb",
	// 	"fun":   "new",
	// 	"lang":  "zh_CN",
	// 	"_":     "1501747870",
	// }

	b := []int{1, 2, 3}

	body, _ := httpclient.Post("http://httpbin.org/post").SendBody(b).Text()
	fmt.Printf("body=%v\n", body)

	fmt.Println("-------------------------------------")
}

func testPut() {
	rsp, _ := httpclient.Put("http://httpbin.org/put").Text()

	fmt.Printf("body=%+v\n", rsp)
	fmt.Println("-------------------------------------")
}

func testDelete() {
	rsp, _ := httpclient.Delete("http://httpbin.org/delete").Text()

	fmt.Printf("body=%+v\n", rsp)
	fmt.Println("-------------------------------------")
}

func testPatch() {
	rsp, _ := httpclient.Patch("http://httpbin.org/patch").Text()

	fmt.Printf("body=%+v\n", rsp)
	fmt.Println("-------------------------------------")
}

func testPostWithFiled() {

	v := url.Values{
		"appid": []string{"wx782c26e4c19acffb"},
		"fun":   []string{"new"},
		"lang":  []string{"zh_CN"},
		"_":     []string{strconv.FormatInt(time.Now().Unix(), 10)},
	}

	body, _ := httpclient.Post("https://login.weixin.qq.com/jslogin").AddFields(v).Text()
	fmt.Printf("body=\n%v", body)
}

func testPostWithBody2() {
	b := map[string]string{
		"appid": "wx782c26e4c19acffb",
		"fun":   "new",
		"lang":  "zh_CN",
		"_":     "1501747870",
	}
	body, _ := httpclient.Post("https://login.weixin.qq.com/jslogin").SendBody(b).Text()
	fmt.Printf("body=%v", body)
}

func testDownload() {
	// SetTimeout(5 * time.Second)
	filePath := "d:/"
	err := httpclient.Get("https://ss0.bdstatic.com/5aV1bjqh_Q23odCf/static/superman/img/logo/bd_logo1_31bdc765.png").ToFile(filePath, "bd_logo1_31bdc765.png")

	fmt.Printf("err:%+v", err)
}
