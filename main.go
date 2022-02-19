package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	resp, _ := http.Get("https://www.theiphonewiki.com/wiki/Models")
	defer resp.Body.Close()  // 函数结束时关闭Body
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	//cmd.Execute()
}
