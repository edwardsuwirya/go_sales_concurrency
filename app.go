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
func main() {
	startTime := time.Now()
	concurrentCalculation()
	//serialCalculation()
	diff := time.Now().Sub(startTime)
	fmt.Printf("Took: %f seconds\n", diff.Seconds())

}
