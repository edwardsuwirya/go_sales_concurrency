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

var workDir = "/Users/edwardsuwirya/Desktop/sample_data"

func salesDataFileList() []fs.FileInfo {
	files, err := ioutil.ReadDir(workDir)
	if err != nil {
		panic(err)
	}
	return files
}

func readSalesData(fileName string) []salesData {
	file, err := os.Open(filepath.Join(workDir, fileName))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var listSales []salesData
	for scanner.Scan() {
		str := scanner.Text()
		split := strings.Split(str, ",")
		outletCode := split[0]
		sales, _ := strconv.Atoi(split[1])
		data := salesData{
			outletCode,
			sales,
		}
		listSales = append(listSales, data)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return listSales
}

func calculateSales(fileName string) int {
	var totalSales int
	for _, s := range readSalesData(fileName) {
		totalSales = totalSales + s.sales
	}
	return totalSales
}

func serialCalculation() {
	for _, f := range salesDataFileList() {
		calculateSales(f.Name())
		//fmt.Println(f.Name(), totalSales)
	}
}

func concurrentCalculation() {
	var wg sync.WaitGroup
	for _, f := range salesDataFileList() {
		wg.Add(1)
		go func(f fs.FileInfo) {
			defer wg.Done()
			calculateSales(f.Name())
			//fmt.Println(f.Name(), totalSales)
		}(f)
	}
	wg.Wait()
}

// go run -race app.go untuk mendeteksi apakah akan terjadi race condition
func concurrentMutexCalculation() {
	var wg sync.WaitGroup
	var grandTotalSales int
	var mtx sync.Mutex
	for _, f := range salesDataFileList() {
		wg.Add(1)
		go func(f fs.FileInfo) {
			defer wg.Done()
			res := calculateSales(f.Name())
			mtx.Lock()
			grandTotalSales = grandTotalSales + res
			mtx.Unlock()
			//fmt.Println(f.Name(), totalSales)
		}(f)
	}
	wg.Wait()
	fmt.Println(grandTotalSales)
}

func channelCalculation() {
	fileList := salesDataFileList()
	lenFileList := len(fileList)

	var results []chan string
	var jobs []chan salesData
	for i := 0; i < lenFileList; i++ {
		results = append(results, make(chan string))
		jobs = append(jobs, make(chan salesData))
	}

	for i, f := range fileList {
		go calculateSalesChannel(jobs[i], results[i])
		go readSalesDataChannel(f.Name(), jobs[i])
	}
	for _, r := range results {
		//fmt.Println(<-r)
		<-r
	}
}

func calculateSalesChannel(jobs chan salesData, result chan string) {
	var totalSales int
	for s := range jobs {
		totalSales = totalSales + s.sales
	}
	result <- fmt.Sprintf("%d", totalSales)
}

func readSalesDataChannel(fileName string, jobs chan salesData) {
	file, err := os.Open(filepath.Join(workDir, fileName))
	if err != nil {
		panic(err)
	}
	defer file.Close()

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

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	close(jobs)
}
func main() {
	startTime := time.Now()
	concurrentMutexCalculation()
	diff := time.Now().Sub(startTime)
	fmt.Printf("Took: %f seconds\n", diff.Seconds())

}
