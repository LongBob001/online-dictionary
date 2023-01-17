package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type DictRequest struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
	UserID    string `json:"user_id"`
}

type DictResponse struct {
	Rc   int `json:"rc"`
	Wiki struct {
		KnownInLaguages int `json:"known_in_laguages"`
		Description     struct {
			Source string      `json:"source"`
			Target interface{} `json:"target"`
		} `json:"description"`
		ID   string `json:"id"`
		Item struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"item"`
		ImageURL  string `json:"image_url"`
		IsSubject string `json:"is_subject"`
		Sitelink  string `json:"sitelink"`
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string      `json:"explanations"`
		Synonym      []string      `json:"synonym"`
		Antonym      []string      `json:"antonym"`
		WqxExample   [][]string    `json:"wqx_example"`
		Entry        string        `json:"entry"`
		Type         string        `json:"type"`
		Related      []interface{} `json:"related"`
		Source       string        `json:"source"`
	} `json:"dictionary"`
}

func main() {
	//直接用编译器运行的代码
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n') //从该流中读取一行，读取数据会默认带有换行符
	if err != nil {
		fmt.Println("错误", err)
		return
	}
	input = strings.TrimSuffix(input, "\r\n") //去掉换行符
	query(input)

	//用cmd命令行输入参数时的代码
	//检查输入是否是一个单词
	//用os.args的len检查，该函数返回一个数组，第一个值为文件的绝对路径，第二个值为单词，如果len！=2，表示没有输入单词
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, `usage:simpleDict WORD
	example:simpleDict hello`)
		os.Exit(1)
	}
	//如果没出错，则数组的第二个值就是要查询的单词
	word := os.Args[1]
	query(word)
}

func query(word string) {
	client := &http.Client{}
	//创建结构体并初始化
	request := DictRequest{TransType: "en2zh", Source: word}
	//将结构体序列化
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	//读取流，序列化后的buf为字节数组
	var data = bytes.NewReader(buf)
	//创建请求
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
	if err != nil {
		log.Fatal(err)
	}
	//设置请求头
	req.Header.Set("authority", "api.interpreter.caiyunai.com")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("app-name", "xy")
	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("device-id", "")
	req.Header.Set("dnt", "1")
	req.Header.Set("origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("os-type", "web")
	req.Header.Set("os-version", "")
	req.Header.Set("referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("sec-ch-ua", `"Not?A_Brand";v="8", "Chromium";v="108", "Microsoft Edge";v="108"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36 Edg/108.0.1462.76")
	req.Header.Set("x-authorization", "token:qgemv4jr1y38jyq6vhvi")
	//发起请求
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//用defer手动关闭响应流
	defer resp.Body.Close()
	//读取响应，bodyText是字符串数组
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//对response进行安全检查，如果不是200，要把状态码和错误原因log出来以便检查
	if resp.StatusCode != 200 {
		log.Fatal("bad statusCode：", resp.StatusCode, "body", string(bodyText))
	}
	//定义结构体
	var dictResponse DictResponse
	//将响应的json数组反序列化到结构体里，注意这里需要用指针操作
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(word, "UK:", dictResponse.Dictionary.Prons.En, "US:", dictResponse.Dictionary.Prons.EnUs)
	for _, item := range dictResponse.Dictionary.Explanations {
		fmt.Println(item)
	}
}
