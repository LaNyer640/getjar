package getjar

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

/*
func NewGetjar(timeout int, process int) *Getjar {
	return &Getjar{
		timeout: timeout,
		process: process,
	}
}*/

func Download(url string, pathbasic string) {
	full_name := strings.Split(url, "//")
	filename1 := full_name[len(full_name)-1]
	full_name2 := strings.Split(url, "/")
	filename2 := full_name2[len(full_name2)-1]
	path := pathbasic + filename1
	if !Exists(path) {
		path2 := strings.Trim(path, filename2)
		f, err := os.Stat(path2)
		if err != nil || f.IsDir() == false {
			if err := os.MkdirAll(path2, os.ModePerm); err != nil {
				fmt.Println("创建目录失败", err)
				return
			}
		}

		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		out, err := os.Create(path)

		if err != nil {
			panic(err)
		}

		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Printf("已下载: %s\n", path)
	}
}

func Geturl(url string) []string {
	var (
		Urls  []string
		Urls2 []string
	)
	//Getjar := NewGetjar(200, 10)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "访问url出错 %v\n", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	s := strings.Split(string(body), "\n")
	//获得最新url
	Urls = getreally(s, url)
	for _, value := range Urls {
		if Isjarurl(value) {
			Urls2 = append(Urls2, value)
		} else if value != "" {
			str := Geturl(value)
			Urls2 = append(Urls2, str...)
		}
	}
	return Urls2
}

func getreally(str []string, url string) []string {
	var (
		url2   string
		Urlmap map[time.Time]string
		Urls   []string
		times  []time.Time
		max    time.Time
	)
	Urlmap = make(map[time.Time]string)
	r, _ := regexp.Compile("href=\"")
	for _, value := range str {
		match := r.MatchString(value)
		if match {
			url1 := getUrl2(value)
			url2 = url + url1
			//这段根据时间戳来放，以便进行时间戳对比得到最新的时间戳
			if Isjarurl(url2) {
				Urls = append(Urls, url2)
			} else {
				if isShuZiChuo(value) && isbanben(url1) {
					Time, _ := time.Parse("2006-01-02 15:04", pipei(value))
					Urlmap[Time] = url2
					times = append(times, Time)
				} else if isUrl(url2) {
					Urls = append(Urls, url2)
				}
			}
		}
	}
	//根据时间戳对比获得最新url
	if times != nil {
		max = times[0]
		for i := 0; i < len(times); i++ {
			if max.Before(times[i]) {
				max = times[i]
			}
		}
	}
	Urls = append(Urls, Urlmap[max])
	return Urls
}

func isUrl(url string) bool {
	if strings.HasSuffix(url, "/") && !strings.HasSuffix(url, "../") {
		return true
	} else {
		return false
	}
}

// 正则匹配得到网页上的时间戳
func getUrl2(url string) string {
	compileRegex := regexp.MustCompile("href=\"(.*?)\"")
	machArr_1 := compileRegex.FindStringSubmatch(url)
	return machArr_1[len(machArr_1)-1]
}

// 判断是否有时间戳
func isShuZiChuo(url string) bool {
	match, _ := regexp.MatchString("[0-9]{4}-[0-9]{2}-[0-9]{2}\\s[0-9]{2}:[0-9]{2}", url)
	return match
}

// 匹配时间戳，得到形如2006-01-02 15:04的时间戳
func pipei(url string) string {
	compileRegex := regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2}\\s[0-9]{2}:[0-9]{2}")
	machArr := compileRegex.FindStringSubmatch(url)
	return machArr[len(machArr)-1]
}

// 判断是不是具有版本号的url，有版本号说明该url下一级是有jar包的
func isbanben(str string) bool {
	var match bool
	if strings.HasSuffix(str, "/") {
		match, _ = regexp.MatchString(".", str)
	}
	return match
}

func Isjarurl(url string) bool {
	if strings.HasSuffix(url, "jar") && !strings.HasSuffix(url, "javadoc.jar") && !strings.HasSuffix(url, "sources.jar") && !strings.HasSuffix(url, "dependencies.jar") {
		return true
	} else {
		return false
	}
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
