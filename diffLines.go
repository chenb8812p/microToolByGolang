package main

import (
	"sync"
	"os"
	"runtime"
	"fmt"
	"bufio"
)

var wg sync.WaitGroup
var wf *os.File

var lines chan string
var siteMap map[string] bool
var total int
var has int
func main(){
	total = 0
	has = 0
	siteMap = make(map[string]bool)
	maxProcs := runtime.NumCPU()
	runtime.GOMAXPROCS(maxProcs)

	lines = make(chan string, 300000000)
	wf, _ = os.Create("/Users/chenyaoyu/Desktop/fileOut")

	file, err := os.Open("/Users/chenyaoyu/Desktop/file1")
	if err != nil {
		fmt.Println(err)
		return
	}
	siteFile, err := os.Open("/Users/chenyaoyu/Desktop/file2")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()
	defer siteFile.Close()
	fmt.Println("start read")

	scanner := bufio.NewScanner(file)
	siteScanner := bufio.NewScanner(siteFile)

	for siteScanner.Scan(){
		siteMap[siteScanner.Text()] = true
	}

	for scanner.Scan(){
		lines <- scanner.Text()
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		 doCheckMd5()
	}
	wg.Wait()
	fmt.Println("DONE")

}


func doCheckMd5() {
	for {
		if len(lines) == 0  {
			break
		}
		total++
		str := <- lines
		if siteMap[str] {
			has++
			fmt.Println(total, ":", has)
		} else {
			wf.WriteString(str + "\n")
		}

	}
}
