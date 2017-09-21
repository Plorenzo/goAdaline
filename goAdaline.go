package main

import (
	"fmt"
	"os"
	"encoding/csv"
	"math/rand"
	"strconv"
	"reflect"
	"time"
)

const trainPath = "/Users/plorenzo/dev/uni/rna/final/train.csv"
const validatePath = "/Users/plorenzo/dev/uni/rna/final/train.csv"
const testPath = "/Users/plorenzo/dev/uni/rna/final/train.csv"

// TODO: refacto to improve performance (S.O. QUESTION)
func readCSV(filepath string) ([][]float64, []float64) {

	csvfile, err := os.Open(filepath)
	if err != nil {
		return nil, nil
	}

	reader := csv.NewReader(csvfile)
	reader.Comma = ';'
	stringMatrix, err := reader.ReadAll()

	csvfile.Close()

	matrix := make([][]float64, len(stringMatrix))
	expectedY := make([]float64, len(stringMatrix))
	
	//Parse string matrix into float64
	for i := range stringMatrix {
		matrix[i] = make([]float64, len(stringMatrix[0]))
		for j := range stringMatrix[i] {
			if j < 8 {
				matrix[i][j], err = strconv.ParseFloat(stringMatrix[i][j], 64)
			} else {
				//Extract expected output date from file (last column)
				expectedY[i], err = strconv.ParseFloat(stringMatrix[i][j], 64)
				matrix[i][j] = 1
			}
			
		}
	}
	return matrix, expectedY
}

//This also inits the threshold
func initWeights(length int) []float64 {

	weights := make([]float64, length)

	//Inits the slice with random numbers betwen [-1, 1]
	for index := range weights {
		w := rand.Float64()
		s := rand.Float64()

		if s < 0.5 {
			weights[index] = w	
		} else {
			weights[index] = w * -1
		}
		
	}
	return weights
}

func main() {

	start := time.Now()

	//Read data from csv file
	data, expectedY := readCSV(trainPath)

	//PARAMETERS
	var cylces int = 10000
	var learningRate float64 = 0.1

	// Sanity checking 
	fmt.Println("Size data: ")
	fmt.Println(len(data[0]))
	fmt.Println(len(data))
	fmt.Println("First data value: ")
	fmt.Println(data[0][0])
	fmt.Println("Data type: ")
	fmt.Println(reflect.TypeOf(data))
	fmt.Println("Expected Y: ")
	fmt.Println(expectedY[0])

	inputLength := len(data[0])
	weights := initWeights(inputLength)
	
	var errors []float64
	
	// Learning
	for i := 0; i < cylces; i++ {
		for j := range data {
			//Calculate estimate
			var estimate float64 = 0
			for x := range data[j]{
				estimate += data[j][x] * weights[x]
			}
		
			// Update weights (range passes calues as a copy)
			for x := 0 ; x < len(weights); x++ {
				weights[x] = learningRate * (expectedY[j] - estimate) * data[j][x]
			}
		}
		
		// Measure cylce error
		var totalErr float64 = 0
		
		for j := range data {
				var estimate float64 = 0
				for x := range data[j] {
					estimate += data[j][x] * weights[x]
				} 
				// Cuadratic error E = (Yd - Ye)^2
				totalErr += (expectedY[j] - estimate) * (expectedY[j] - estimate)
		}
		errors = append(errors, totalErr / float64(len(errors)))
	}
	

	//TODO save global error to file

	elapsed := time.Since(start)
	fmt	.Println(elapsed)
	//fmt.Println(errors)


}
