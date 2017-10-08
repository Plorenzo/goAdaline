package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strconv"

	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/mat"
)

func main() {

	var (
		folderPath   = flag.String("path", ".", "Path to the datasets")
		cycles       = flag.Int("cycles", 1, "Nº of training cycles")
		learningRate = flag.Float64("lr", 0.1, "Learning rate of the neuron")
		outputPath   = flag.String("out", ".", "Path to save the output file")
		doCPU        = flag.Bool("cpu", false, "enable CPU profiling")
	)
	flag.Parse()

	if *doCPU {
		f, err := os.Create("cpu.prof")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		err = pprof.StartCPUProfile(f)
		if err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	trainPath := filepath.Join(*folderPath, "train.csv")
	validatePath := filepath.Join(*folderPath, "validate.csv")
	testPath := filepath.Join(*folderPath, "test.csv")

	// Read data from csv file
	data, expectedY := readCSV(trainPath)
	validateData, valExpectedY := readCSV(validatePath)
	testData, testExpectedY := readCSV(testPath)

	nrows, ncols := data.Dims()

	// Init weights randomly [-1,1]
	weights := initWeights(ncols)

	var (
		errorsTrain    = make([]float64, 0, *cycles)
		errorsValidate = make([]float64, 0, *cycles)
		errorsTest     float64
	)

	rate := *learningRate
	row := blas64.Vector{
		Inc: 1,
	}
	raw := data.RawMatrix()
	exp := expectedY.RawVector()
	ws := weights.RawVector()
	// Learning
	for cycle := 0; cycle < *cycles; cycle++ {
		for i := 0; i < nrows; i++ {
			row.Data = raw.Data[i*raw.Stride:]
			iexp := exp.Data[i]
			// Calculate estimate
			estimate := blas64.Dot(ncols, row, ws)
			_ = row.Data[len(ws.Data)-1]
			for j := range ws.Data {
				ws.Data[j] += rate * (iexp - estimate) * row.Data[j]
			}
		}

		// Compute cycle train error
		errorsTrain = append(errorsTrain, computeError(data, expectedY, weights))
		errorsValidate = append(errorsValidate, computeError(validateData, valExpectedY, weights))
	}

	errorsTest = computeError(testData, testExpectedY, weights)

	var estimates mat.VecDense
	estimates.MulVec(testData, weights)

	fmt.Println("Test error: ")
	fmt.Println(errorsTest)
	fmt.Println("Weights:")
	fmt.Println(weights.RawVector().Data)

	createCSV(*outputPath, errorsTrain, errorsValidate, weights.RawVector().Data, estimates.RawVector().Data)
}

func readCSV(filepath string) (*mat.Dense, *mat.VecDense) {

	csvfile, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("could not open %q: %v", filepath, err)
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("could not decode CSV file: %v", err)
	}

	nrows := len(records)
	ncols := len(records[0])

	m := mat.NewDense(nrows, ncols, nil)
	y := mat.NewVecDense(nrows, nil)

	// Parse string matrix into float64
	for i, record := range records {
		for j, str := range record {
			var v float64 = 1
			if j < 8 {
				v, err = strconv.ParseFloat(str, 64)
				if err != nil {
					log.Fatalf("could not parse float %q: %v", str, err)
				}
			} else {
				var yv float64
				// Extract expected output date from file (last column)
				yv, err = strconv.ParseFloat(str, 64)
				if err != nil {
					log.Fatalf("could not parse float %q: %v", str, err)
				}
				y.SetVec(i, yv)
			}
			m.Set(i, j, v)

		}
	}
	return m, y
}

// This also inits the threshold
func initWeights(length int) *mat.VecDense {

	weights := make([]float64, length)
	// Inits the slice with random numbers between [-1, 1]
	for index := range weights {
		w := 2*rand.Float64() - 1
		weights[index] = w
	}
	return mat.NewVecDense(length, weights)
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
			log.Fatalf("could not write back sample %q: %v", v[0], err)
		}
	}

	writer.Flush()
	err = file.Close()
	if err != nil {
		log.Fatalf("could not write back data to file: %v", err)
	}
}

func computeError(data *mat.Dense, expected, weights *mat.VecDense) float64 {

	var errs mat.VecDense
	errs.MulVec(data, weights)
	errs.SubVec(expected, &errs)
	errs.MulElemVec(&errs, &errs)

	return mat.Sum(&errs) / float64(errs.Len())
}
