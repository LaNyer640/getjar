package main

import (
	"flag"
	"fmt"
	"getjar/getjar"
	"getjar/lib"
	"math"
	"sync"
	"time"
)

var (
	startTime = time.Now()
	url       = flag.String("url", "https://repo1.maven.org/maven2/org/apache/", "想要递归下载的url 例如:-url=https://repo1.maven.org/maven2/org/apache/")
	path      = flag.String("path", "D:/jar/", "想要保存的文件的地址 例如:-path=D:/jar/")
	process   = flag.Int("n", 10, "进程数 例如:-n=10")
	h         = flag.Bool("h", false, "帮助信息")
)

func main() {
	var (
		downloadUrls []string
		pageCount    int
	)

	startTime := time.Now()
	flag.Parse()

	if *h == true {
		lib.Usage("getjar version: getjar/1.0\n Usage: getjar [-h] [-url ip地址] [-n 进程数]")
		return
	}

	if getjar.Isjarurl(*url) {
		downloadUrls = append(downloadUrls, *url)
	}
	downloadUrls = getjar.Geturl(*url)
	fmt.Printf("花费了%.3fs获得了所有url\n", time.Since(startTime).Seconds())
	startTime2 := time.Now()
	total := len(downloadUrls)

	fmt.Printf("一共有%d条url\n", total)
	if total < *process {
		pageCount = total
	} else {
		pageCount = *process
	}

	num := int(math.Ceil(float64(total) / float64(pageCount)))
	all := map[int][]string{}
	for i := 1; i <= pageCount; i++ {
		for j := 0; j < num; j++ {
			tmp := (i-1)*num + j
			if tmp < total {
				all[i] = append(all[i], downloadUrls[tmp])
			}
		}
	}

	wg := sync.WaitGroup{}
	for _, v := range all {
		wg.Add(1)
		go func(value []string) {
			defer wg.Done()
			for i := 0; i < len(value); i++ {
				//fmt.Printf("正在下载%s\n", value[i])
				getjar.Download(value[i], *path)
			}
		}(v)
	}
	wg.Wait()
	fmt.Printf("花费了%.3fs下载了所有jar\n", time.Since(startTime2).Seconds())
}
