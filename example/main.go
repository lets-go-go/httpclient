package main

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/lets-go-go/httpclient"
)

func main() {

	// httpclient.Settings().SetProxy(httpclient.NoProxy, "")
	// httpclient.Settings().SetProxy(httpclient.DefaultProxy, "")
	// httpclient.Settings().SetProxy(httpclient.CustomProxy, "http://192.168.16.232:8080")
	testPostWithFiled()
	// testDownload()
	// testPost()
	// testGoogle()
}

func testPost1() {

	// v := url.Values{
	// 	"scKey":          []string{"wx782c26e4c19acffb"},
	// 	"currentVersion": []string{"2.0.0.0"},
	// 	"fixVersion":     []string{"2.0.0.0"},
	// }

	b := map[string]string{
		"scKey":          "wx782c26e4c19acffb",
		"currentVersion": "2.0.0.0",
		"fixVersion":     "2.0.0.0",
	}

	// b := []int{1, 2, 3}

	if body, err := httpclient.Post("http://cloudcn.focusteach.com/pc/sc/equipnum/update", 3*time.Second).SendBody(b).Text(); err != nil {
		fmt.Printf("err=%v\n", err)
	} else {
		fmt.Printf("body=%v\n", body)
	}

	fmt.Println("-------------------------------------")
}

func testPost() {

	v := url.Values{
		"appid": []string{"wx782c26e4c19acffb"},
		"fun":   []string{"new"},
		"lang":  []string{"zh_CN"},
		"_":     []string{strconv.FormatInt(time.Now().Unix(), 10)},
	}

	// b := []int{1, 2, 3}

	body, _ := httpclient.Post("http://httpbin.org/post", 3*time.Second).AddFields(v).AttachFile("file1", "d:/1.jpg", "1.jpg").Text()
	fmt.Printf("body=%v\n", body)

	fmt.Println("-------------------------------------")
}

func testGet() {

	rsp, _ := httpclient.Head("http://httpbin.org/get", 3*time.Second).Execute()

	fmt.Printf("body=%+v\n", rsp.Header)
	fmt.Println("-------------------------------------")
}

func testHead() {
	rsp, _ := httpclient.Options("http://httpbin.org/get", 3*time.Second).Execute()

	fmt.Printf("body=%+v\n", rsp.Header)
	fmt.Println("-------------------------------------")
}

func testGetBody() {

	body, _ := httpclient.Get("http://httpbin.org/get", 3*time.Second).Text()
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

	body, _ := httpclient.Post("http://httpbin.org/post", 3*time.Second).SendBody(b).Text()
	fmt.Printf("body=%v\n", body)

	fmt.Println("-------------------------------------")
}

func testPut() {
	rsp, _ := httpclient.Put("http://httpbin.org/put", 3*time.Second).Text()

	fmt.Printf("body=%+v\n", rsp)
	fmt.Println("-------------------------------------")
}

func testDelete() {
	rsp, _ := httpclient.Delete("http://httpbin.org/delete", 3*time.Second).Text()

	fmt.Printf("body=%+v\n", rsp)
	fmt.Println("-------------------------------------")
}

func testPatch() {
	rsp, _ := httpclient.Patch("http://httpbin.org/patch", 3*time.Second).Text()

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

	c := httpclient.Post("https://login.weixin.qq.com/jslogin", 3*time.Second).AddFields(v)
	body, _ := c.Execute()
	fmt.Printf("body=%v\n", body)
	dmp, err := c.DumpRequest()
	fmt.Printf("dmp req=%v,err=%v\n", string(dmp), err)

	dmp, err = c.DumpResponse()
	fmt.Printf("dmp rsp=%v,err=%v\n", string(dmp), err)
}

func testPostWithBody2() {
	b := map[string]string{
		"appid": "wx782c26e4c19acffb",
		"fun":   "new",
		"lang":  "zh_CN",
		"_":     "1501747870",
	}
	body, _ := httpclient.Post("https://login.weixin.qq.com/jslogin", 3*time.Second).SendBody(b).Text()
	fmt.Printf("body=%v", body)
}

func testDownload() {
	filePath := "d:/"
	url := "https://github.com/henrylee2cn/goutil/blob/master/pool/GoPool.go"

	err := httpclient.Get(url, 0).ToFile(filePath, "")

	fmt.Printf("err:%+v", err)
}

func testGoogle() {
	filePath := "d:/"
	url := "https://www.google.com"

	err := httpclient.Get(url, 10*time.Second).ToFile(filePath, "")

	fmt.Printf("err:%+v", err)
}
