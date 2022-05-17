package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type salesData struct {
	outletCode string
	sales      int
}

func salesCaluction(wg *sync.WaitGroup, fileName string) {
	jobs := make(chan salesData)
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	go func() {
		var totalSales int
		for s := range jobs {
			totalSales = totalSales + s.sales
		}
		fmt.Println(fileName, totalSales)
		wg.Done()
	}()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := scanner.Text()
		split := strings.Split(str, ",")
		outletCode := split[0]
		sales, _ := strconv.Atoi(split[1])
		data := salesData{
			outletCode,
			sales,
		}
		jobs <- data
	}
	close(jobs)

	if err := scanner.Err(); err != nil {
		panic(err)
	}

}
func main() {
	startTime := time.Now()
	workDir := "/Users/edwardsuwirya/Desktop/sample_data"
	files, err := ioutil.ReadDir(workDir)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for _, f := range files {
		wg.Add(1)
		go func(f fs.FileInfo) {
			salesCaluction(&wg, filepath.Join(workDir, f.Name()))
		}(f)
	}
	wg.Wait()

	diff := time.Now().Sub(startTime)
	fmt.Printf("Took: %f seconds\n", diff.Seconds())

}
