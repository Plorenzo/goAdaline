package main

import (
	"fmt"
	"os"
	"encoding/csv"
	"math/rand"
	"strconv"
	"reflect"
)

const path = "/Users/plorenzo/dev/uni/rna/p1"

func readCSV(filepath string) [][]float64 {

	csvfile, err := os.Open(filepath)
	if err != nil {
		return nil
	}

	reader := csv.NewReader(csvfile)
	stringMatrix, err := reader.ReadAll()

	csvfile.Close()

	matrix := make([][]float64, len(stringMatrix))

	for i := range stringMatrix {
		matrix[i] = make([]float64, len(stringMatrix[0]))
		for y := range stringMatrix[i] {
			matrix[i][y], err = strconv.ParseFloat(stringMatrix[i][y], 64)
		}
	}

	return matrix
}

func initWeights(length int) []float64 {

	weights := make([]float64, length)
	for index := range weights {
		weights[index] = rand.Float64()
	}

	return weights
}

func main() {

	data := readCSV(path)

}
