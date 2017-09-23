package main

import (
	"fmt"
	"os"
	"encoding/csv"
	"math/rand"
	"strconv"
)

const trainPath = "/Users/plorenzo/dev/uni/rna/final/train.csv"
const validatePath = "/Users/plorenzo/dev/uni/rna/final/validate.csv"
const testPath = "/Users/plorenzo/dev/uni/rna/final/test.csv"

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


func createCSV(train []float64, validate []float64) {

	//TODO Add to existing file instead of creating one each time

	file, _ := os.Create("/Users/plorenzo/dev/uni/rna/errores/total.csv")
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    var strings []string
    var strings1 []string

    for i := range train {
    	strings = append(strings, strconv.FormatFloat(train[i], 'f', 6, 64)) 
    }
    for i := range validate {
    	strings1 = append(strings1, strconv.FormatFloat(validate[i], 'f', 6, 64)) 
    }
    writer.Write(strings)
    writer.Write(strings1)
}

func computeError(data [][]float64, expected []float64, weights []float64) float64 {

	var errors float64
	var errorSum, estimate float64 = 0, 0

	for i := range data {
		estimate = 0
		for j := range data[i] {
			estimate += data[i][j] * weights[j]
		}
		// Cuadratic error E = (Yd - Ye)^2
		errorSum += (expected[i] - estimate) * (expected[i] - estimate)
	}
	errors = errorSum / float64(len(data))

	return errors
}


func main() {

	//Read data from csv file
	data, expectedY := readCSV(trainPath)
	validateData, valExpectedY := readCSV(validatePath)
	testData, testExpectedY := readCSV(testPath)

	//PARAMETERS
	var cylces int = 57
	var learningRate float64 = 0.01
	
	weights := initWeights(len(data[0]))
		
	var estimate float64
	var errorsTrain []float64
	var errorsValidate []float64
	var errorsTest float64
	
	// Learning
	for i := 0; i < cylces; i++ {
		for j := range data {
			//Calculate estimate
			estimate = 0
			for x := range data[j]{
				estimate += data[j][x] * weights[x]
			}
			
			// Update weights (range passes values as a copy)
			for x := 0; x < len(weights); x++ {
				weights[x] += learningRate * (expectedY[j] - estimate) * data[j][x]
			}
		}

		// Compute cylce train error
		errorsTrain = append(errorsTrain, computeError(data, expectedY, weights))
		errorsValidate = append(errorsValidate, computeError(validateData, valExpectedY, weights))

	}

	errorsTest = computeError(testData, testExpectedY, weights)

	fmt.Println("Test error: ")
	fmt.Println(errorsTest)
	

	createCSV(errorsTrain, errorsValidate)
}

