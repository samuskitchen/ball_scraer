package analyze

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Analyzer struct to hold the file path
type Analyzer struct {
	filePath string
}

// NewAnalyzer is the constructor for the Analyzer struct
func NewAnalyzer(filePath string) *Analyzer {
	return &Analyzer{filePath: filePath}
}

// AnalyzeResults reads the CSV file, counts the appearances of each number,
// and calculates the percentage of appearance for each number.
func (a *Analyzer) AnalyzeResults() error {
	// Open the CSV file
	file, err := os.Open(a.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Map to count the appearances of each number
	numberCount := make(map[int]int)
	totalDraws := len(records) - 1 // Exclude the header row

	// Iterate over the records and count the appearances
	for _, record := range records[1:] {
		// Count the appearances in the "Result" field
		results := strings.Split(record[2], "-")
		for _, numStr := range results {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return err
			}
			numberCount[num]++
		}

		// Count the appearances in the "Super Ball" field
		superBall, err := strconv.Atoi(record[3])
		if err != nil {
			return err
		}
		numberCount[superBall]++
	}

	// Create a slice to store the results for sorting
	type result struct {
		number      int
		appearances int
		percentage  float64
	}
	var resultsList []result

	// Calculate the percentages and add to the slice
	for number, count := range numberCount {
		percentage := (float64(count) / float64(totalDraws*6)) * 100 // totalDraws*6 because each draw has 5 numbers + 1 super ball
		resultsList = append(resultsList, result{number, count, percentage})
	}

	// Sort the results by percentage in descending order
	sort.Slice(resultsList, func(i, j int) bool {
		return resultsList[i].percentage > resultsList[j].percentage
	})

	// Print the sorted results
	fmt.Println()
	fmt.Println("Number\tAppearances\tPercentage")
	for _, res := range resultsList {
		fmt.Printf("%d\t%d\t\t%.2f%%\n", res.number, res.appearances, res.percentage)
	}

	return nil
}
