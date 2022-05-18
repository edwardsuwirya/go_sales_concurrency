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

func channelCalculation() {
	var wg sync.WaitGroup
	for _, f := range salesDataFileList() {
		wg.Add(1)
		go func(f fs.FileInfo) {
			readSalesDataChannel(&wg, f.Name())
			//fmt.Println(f.Name(), totalSales)
		}(f)
	}
	wg.Wait()
}

func calculateSalesChannel(wg *sync.WaitGroup, jobs chan salesData) {
	var totalSales int
	for s := range jobs {
		totalSales = totalSales + s.sales
	}
	//fmt.Println(fileName, totalSales)
	wg.Done()
}
func readSalesDataChannel(wg *sync.WaitGroup, fileName string) {
	jobs := make(chan salesData)
	file, err := os.Open(filepath.Join(workDir, fileName))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	go calculateSalesChannel(wg, jobs)

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
	channelCalculation()
	diff := time.Now().Sub(startTime)
	fmt.Printf("Took: %f seconds\n", diff.Seconds())

}
