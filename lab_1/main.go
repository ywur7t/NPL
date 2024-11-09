package main

// добавлен вывод медиан и отклонений
import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
)

type Data struct {
	category string
	value    float64
}

func median(values []float64) float64 {
	sort.Float64s(values)
	n := len(values)
	if n%2 == 1 {
		return values[n/2]
	}
	return (values[n/2-1] + values[n/2]) / 2
}


func stdDev(values []float64) float64 {
	mean := 0.0
	for _, v := range values {
		mean += v
	}
	mean /= float64(len(values))

	var variance float64
	for _, v := range values {
		variance += (v - mean) * (v - mean)
	}
	variance /= float64(len(values) - 1)
	return math.Sqrt(variance)
}

func processFile(filename string, wg *sync.WaitGroup, mediansCh, stdDevsCh chan<- float64) {
	defer wg.Done()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Не удалось открыть файл: %s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	categories := make(map[string][]float64)
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		category := record[0]
		value, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			continue
		}
		categories[category] = append(categories[category], value)

	}

	for _, values := range categories {

		median := median(values)
		mediansCh <- median
		fmt.Println(filename,"median ", median)

		stdDev := stdDev(values)
		stdDevsCh <- stdDev
		fmt.Println(filename, "stdDev ", stdDev)
	}


}

func main() {

        files, err := filepath.Glob("*.csv")
	if err != nil { panic(err) }
	fmt.Println("Найденные CSV файлы:", files)


	var wg sync.WaitGroup
	mediansCh := make(chan float64, len(files)*4)
	stdDevsCh := make(chan float64, len(files)*4)


	for _, file := range files {
		wg.Add(1)
		go processFile(file, &wg, mediansCh, stdDevsCh)
	}

        wg.Wait()
        close(mediansCh)
        close(stdDevsCh)



	var medians, stdDevs []float64
	for median := range mediansCh {
		medians = append(medians, median)
                // fmt.Println(median)
	}
	for stdDev := range stdDevsCh {
		stdDevs = append(stdDevs, stdDev)
		// fmt.Println(stdDevsCh)
	}

	finalMedian := median(medians)
	finalStdDev := stdDev(stdDevs)

	fmt.Printf("Медиана %.2f\n", finalMedian)
	fmt.Printf("Ср кв отклонение %.2f\n", finalStdDev)

}
