package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

const path = "/Users/plorenzo/dev/uni/rna/p1"

func readCSV(filepath string) [][]string {

	csvfile, err := os.Open(filepath)

	if err != nil {
		return nil
	}

	reader := csv.NewReader(csvfile)
	fields, err := reader.ReadAll()

	csvfile.Close()

	return fields
}

func main() {

	data := readCSV(path)
	fmt.Println(data[0][0])

	

}
