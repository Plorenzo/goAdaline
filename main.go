package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
)

func main() {

	var (
		folderPath   = flag.String("path", ".", "Path to the datasets")
		cycles       = flag.Int("cycles", 1, "Nº of training cycles")
		learningRate = flag.Float64("lr", 0.1, "Learning rate of the neuron")
		outputPath   = flag.String("out", ".", "Path to save the output file")
	)
	flag.Parse()

	trainPath := filepath.Join(*folderPath, "train.csv")
	validatePath := filepath.Join(*folderPath, "validate.csv")
	testPath := filepath.Join(*folderPath, "test.csv")

	// Read data from csv file
	data, expectedY := readCSV(trainPath)
	validateData, valExpectedY := readCSV(validatePath)
	testData, testExpectedY := readCSV(testPath)

	// Init weights randomly [-1,1]
	weights := initWeights(len(data[0]))

	var estimate float64
	var estimates []float64
	var errorsTrain []float64
	var errorsValidate []float64
	var errorsTest float64

	// Learning
	for i := 0; i < *cycles; i++ {
		for j := range data {
			//Calculate estimate
			estimate = 0
			for x := range data[j] {
				estimate += data[j][x] * weights[x]
			}

			// Update weights (range passes values as a copy)
			for x := 0; x < len(weights); x++ {
				weights[x] += *learningRate * (expectedY[j] - estimate) * data[j][x]
			}
		}

		// Compute cycle train error
		errorsTrain = append(errorsTrain, computeError(data, expectedY, weights))
		errorsValidate = append(errorsValidate, computeError(validateData, valExpectedY, weights))
	}

	errorsTest = computeError(testData, testExpectedY, weights)

	for i := range testData {
		estimate = 0
		for j := range testData[i] {
			estimate += testData[i][j] * weights[j]
		}
		estimates = append(estimates, estimate)
	}

	fmt.Println("Test error: ")
	fmt.Println(errorsTest)
	fmt.Println("Weights:")
	fmt.Println(weights)

	createCSV(*outputPath, errorsTrain, errorsValidate, weights, estimates)
}

// TODO: refactor to improve performance (S.O. QUESTION)
func readCSV(filepath string) ([][]float64, []float64) {

	csvfile, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("could not open %q: %v", filepath, err)
	}

	reader := csv.NewReader(csvfile)
	reader.Comma = ';'
	stringMatrix, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("could not decode CSV file: %v", err)
	}

	csvfile.Close()

	matrix := make([][]float64, len(stringMatrix))
	expectedY := make([]float64, len(stringMatrix))

	//Parse string matrix into float64
	for i := range stringMatrix {
		matrix[i] = make([]float64, len(stringMatrix[0]))
		for j := range stringMatrix[i] {
			if j < 8 {
				matrix[i][j], err = strconv.ParseFloat(stringMatrix[i][j], 64)
				if err != nil {
					log.Fatalf("could not parse float %q: %v", stringMatrix[i][j], err)
				}
			} else {
				//Extract expected output date from file (last column)
				expectedY[i], err = strconv.ParseFloat(stringMatrix[i][j], 64)
				if err != nil {
					log.Fatalf("could not parse float %q: %v", stringMatrix[i][j], err)
				}
				matrix[i][j] = 1
			}

		}
	}
	return matrix, expectedY
}

//This also inits the threshold
func initWeights(length int) []float64 {

	weights := make([]float64, length)
	//Inits the slice with random numbers between [-1, 1]
	for index := range weights {
		w := 2*rand.Float64() - 1
		weights[index] = w
	}
	return weights
}

func createCSV(path string, train []float64, validate []float64, weights []float64, estimates []float64) {

	var filePath string

	if path == "." {
		filePath = "errors.csv"
	} else {
		filePath = path
	}

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	trainS := []string{"Train"}
	validateS := []string{"Validate"}
	weightsS := []string{"Weights"}
	estimatesS := []string{"Estimates"}

	for i := range train {
		trainS = append(trainS, strconv.FormatFloat(train[i], 'f', 6, 64))
	}
	for i := range validate {
		validateS = append(validateS, strconv.FormatFloat(validate[i], 'f', 6, 64))
	}
	for i := range estimates {
		estimatesS = append(estimatesS, strconv.FormatFloat(estimates[i], 'f', 6, 64))
	}
	for i := range weights {
		weightsS = append(weightsS, strconv.FormatFloat(weights[i], 'f', 6, 64))
	}

	for _, v := range [][]string{trainS, validateS, estimatesS, weightsS} {
		err = writer.Write(v)
		if err != nil {
			log.Fatalf("could not write back sample: %v", err)
		}
	}

	writer.Flush()
	err = file.Close()
	if err != nil {
		log.Fatalf("could not write back data to file: %v", err)
	}
}

func computeError(data [][]float64, expected []float64, weights []float64) float64 {

	var errors float64
	var errorSum, estimate float64 = 0, 0

	for i := range data {
		estimate = 0
		for j := range data[i] {
			estimate += data[i][j] * weights[j]
		}
		// Squared error E = (Yd - Ye)^2
		errorSum += (expected[i] - estimate) * (expected[i] - estimate)
	}
	errors = errorSum / float64(len(data))

	return errors
}
