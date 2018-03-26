pckage main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
)

var urlChan chan string
var wg sync.WaitGroup
var count int
var wf *os.File

func main() {
	wf, _ = os.Create("C:\\url\\out.txt")

	count = 0
	maxProcs := runtime.NumCPU()
	runtime.GOMAXPROCS(maxProcs)

	fmt.Println("start read")
	urlChan = make(chan string, 50000)
	file, err := os.Open("C:\\url\\urlList.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	fmt.Println("start read")

	for scanner.Scan() {
		urlChan <- scanner.Text()
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go fetchUrl()
	}
	wg.Wait()
	fmt.Println("DONE")
}

func fetchUrl() {
	for {
		count++

		if len(urlChan) <= 0 {
			break
		}
		url := <-urlChan

		response, err := http.Get(url)
		if err != nil {
			fmt.Printf("%s   %s", url, err)
			return
		}
		body, _ := ioutil.ReadAll(response.Body) //转换byte数组
		resultArr := make([]string, 0, 10)
		resultArr = append(resultArr, url)
		defer response.Body.Close()
		code := response.StatusCode
		if code == 404 {
			resultArr = append(resultArr, "404")
			res := fmt.Sprintf("%s", resultArr)
			fmt.Println(count, ":", res)

			wf.WriteString(res + "\n")
			continue
		}
		bodystr := string(body) //转换字符串
		if strings.Contains(bodystr, "pageLevel = 2") {
			resultArr = append(resultArr, "二级页面")
			res := fmt.Sprintf("%s", resultArr)
			fmt.Println(res)

			wf.WriteString(res + "\n")
		}

		ind := strings.LastIndex(url, "/")
		if ind != -1 {
			tmpUrl := url[:ind]
			ind := strings.LastIndex(tmpUrl, "/")
			if len(tmpUrl[ind:]) == 7 {
				resultArr = append(resultArr, "六位随机")
			}
		}

		if !strings.Contains(bodystr, "siteId") {
			resultArr = append(resultArr, "没有站点定义")
		}
		if !strings.Contains(bodystr, "sunlandsLog_online.js") {
			resultArr = append(resultArr, "没有sunlandsjs")
		}
		if !strings.Contains(bodystr, "interactive_online.js") {
			resultArr = append(resultArr, "没有interactivejs")
		}

		if !strings.Contains(bodystr, "interactive.init") {
			if !strings.Contains(bodystr, "var stPhone") && !strings.Contains(bodystr, "var qq") && !strings.Contains(bodystr, "var wechat") {
				resultArr = append(resultArr, "没有init且没有留言")
			}
			if strings.Contains(bodystr, "openMeiqia") && strings.Contains(bodystr, "openNtkf") {
				resultArr = append(resultArr, "没有init有聊天")
			}
		}

		if !strings.Contains(bodystr, "openMeiqia") && !strings.Contains(bodystr, "openNtkf") {
			if !strings.Contains(bodystr, "var stPhone") && !strings.Contains(bodystr, "var qq") && !strings.Contains(bodystr, "var wechat") && !strings.Contains(bodystr, "submitMsg") {
				resultArr = append(resultArr, "没有发起且没有留言")
			}
		}
		res := fmt.Sprintf("%s", resultArr)
		fmt.Println(count, ":", res)

		wf.WriteString(res + "\n")
	}
	wg.Done()
}
