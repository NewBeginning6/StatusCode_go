package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var urllist []string

func httpres(url string, c chan string, wgscan *sync.WaitGroup) {
	rsps, err := http.Get(url)
	/*
		如果请求成功err == nil
	*/
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer rsps.Body.Close()
	Code := strconv.Itoa(rsps.StatusCode)
	c <- (url + "\t" + Code)
	wgscan.Done()
}

func fileread(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("open file err:", err.Error())
		return
	}
	defer file.Close()
	r := bufio.NewReader(file) //建立缓冲区，把文件内容放到缓冲区中
	for {
		// 分行读取文件  ReadLine返回单个行，不包括行尾字节(\n  或 \r\n)
		data, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read err", err.Error())
			break
		}
		// 打印出内容
		//fmt.Printf("%v", string(data))
		urllist = append(urllist, string(data))
	}
}

func main() {
	ip := flag.String("u", "", "ip or url")
	file := flag.String("r", "", "url list file")
	flag.Parse()
	var wg sync.WaitGroup
	c := make(chan string, 500) //通道定义
	start := time.Now()
	if *ip != "" {
		urllist = append(urllist, *ip)
		wg.Add(len(urllist)) //计数器，只有带那个计数器为0才执行某个操作
		/*
			遍历url数组进行扫描
		*/
		for i := range urllist {
			if strings.Contains(urllist[i], "://") {
			} else {
				urllist[i] = "http://" + urllist[i]
			}
			go func(url string) {
				httpres(url, c, &wg)
			}(urllist[i])
		}
		go func() {
			wg.Wait() //计数器等待，为0即关闭通道
			close(c)
		}()
		for i := range c {
			fmt.Println(i)
		}
		end := time.Since(start)
		fmt.Println("花费时间为:", end)
	}
	if *file != "" {
		fileread(*file)
		wg.Add(len(urllist)) //计数器，只有带那个计数器为0才执行某个操作
		/*
			遍历url数组进行扫描
		*/
		for i := range urllist {
			if strings.Contains(urllist[i], "://") {
			} else {
				urllist[i] = "http://" + urllist[i]
			}
			go func(url string) {
				httpres(url, c, &wg)
			}(urllist[i])
		}
		go func() {
			wg.Wait() //计数器等待，为0即关闭通道
			close(c)
		}()
		for i := range c {
			fmt.Println(i)
		}
		end := time.Since(start)
		fmt.Println("花费时间为:", end)
	}
	if *file == "" && *ip == "" {
		flag.Usage()
	}
}
