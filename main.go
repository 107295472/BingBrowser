package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"yin/AdminVue/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
	"github.com/tidwall/gjson"
)

var po = "8099"

func main() {
	p := utils.GetCurrentPath() //当前路径
	r := gin.Default()
	r.Static("/static", p+"/public") //公开资源文件
	//模版示例
	r.LoadHTMLGlob(p + "templates/*")
	r.GET("/", func(c *gin.Context) {
		a, b := Next("0")
		c.HTML(http.StatusOK, "bing.html", gin.H{
			"title": a,
			"url":   b,
		})
	})
	r.GET("/exit", func(c *gin.Context) {
		os.Exit(3)
	})
	type imgJson struct {
		Name string
		Url  string
	}
	r.GET("/bing", func(c *gin.Context) {
		var id = c.Query("idx")
		a, b := Next(id)
		//组合json
		joinStr := imgJson{a, b}
		result, err := json.Marshal(joinStr)
		if err != nil {
			fmt.Println("encoding failed...")
		}
		c.Writer.WriteString(string(result))
	})
	openPro()
	r.Run(":" + po)
}
func Next(idx string) (string, string) {
	var htt = "http://cn.bing.com"
	var url = "http://cn.bing.com/HPImageArchive.aspx?format=js&idx=" + idx + "&n=1&nc=" + string(time.Now().Unix()) + "&video=1"
	json := httpGet(url)
	img := gjson.Get(json, "images.#.url")
	title := gjson.Get(json, "images.#.copyright")
	var title2 string
	var img1 string
	for _, name := range title.Array() {
		title2 = name.String()
	}
	for _, name := range img.Array() {
		img1 = name.String()
	}
	//截取字符串
	str1 := strings.Index(title2, "©")
	var s = string([]byte(title2)[:str1-1])
	return s, htt + img1
}
func httpGet(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	return string(body)
}
func openPro() error {
	exeStr := filepath.ToSlash("C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe")
	t, err := PathExists(exeStr)
	check(err)
	if !t {
		exeStr = filepath.ToSlash("C:\\Program Files\\Internet Explorer\\iexplore.exe")
	}
	arg := "http://localhost:" + po
	cmd := exec.Command(exeStr, arg)
	return cmd.Start()
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
