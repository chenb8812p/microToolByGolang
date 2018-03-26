package main

import (
	"sync"
	"os"
	"runtime"
	"fmt"
	"bufio"
	"crypto/md5"
)

var wg sync.WaitGroup
var wf *os.File
var lines chan string

func main(){
	wf, _ = os.Create("/Users/chenyaoyu/Desktop/md5Out")

	maxProcs := runtime.NumCPU()
	runtime.GOMAXPROCS(maxProcs)

	lines = make(chan string, 2000000)

	file, err := os.Open("/Users/chenyaoyu/Desktop/egurl")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Println("start read")

	scanner := bufio.NewScanner(file)

	for scanner.Scan(){
		lines <- scanner.Text()
	}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go doMd5()
	}
	wg.Wait()
	fmt.Println("DONE")

}


func doMd5() {
	for {
		if len(lines) == 0  {
			break
		}

		str := <- lines
		cipStr := md5.Sum([]byte(str))
		md5Str := fmt.Sprintf("%x", cipStr)
		fmt.Println(md5Str)
		wf.WriteString(md5Str + "\n")
	}
}
