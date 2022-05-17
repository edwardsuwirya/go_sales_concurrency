package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type salesData struct {
	outletCode string
	sales      int
}

func main() {
	startTime := time.Now()
	workDir := "/Users/edwardsuwirya/Desktop/sample_data"
	files, err := ioutil.ReadDir(workDir)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		file, err := os.Open(filepath.Join(workDir, f.Name()))
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

		var totalSales int
		for _, s := range listSales {
			totalSales = totalSales + s.sales
		}
		fmt.Println(f.Name(), totalSales)
	}
	diff := time.Now().Sub(startTime)
	fmt.Printf("Took: %f seconds\n", diff.Seconds())

}
