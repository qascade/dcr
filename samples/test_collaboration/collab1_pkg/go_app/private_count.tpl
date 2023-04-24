//nolint
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/google/differential-privacy/go/v2/dpagg"
	"github.com/google/differential-privacy/go/v2/noise"
)

//nolint
func extractUniqueId(loc string, UniqueID string) ([]string, error) {
	unique_id_list := []string{}
	// Open the csv file
	csvFile, err := os.Open(loc)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()
	
	// Read the csv file
	r := csv.NewReader(csvFile)
	
	// Iterate through the csv file
	firstLine := true
	uniqueIdx := -1 // Index of the column with name UniqueID
	for {
		record, err := r.Read()
		if err == io.EOF {
			return unique_id_list, nil
		}
		if firstLine {
			// Get the index of the column with name UniqueID
			for i, col := range record {
				if col == UniqueID {
					uniqueIdx = i
					break
				}
			}
			if uniqueIdx == -1 {
				return nil, fmt.Errorf("column with name %s not found", UniqueID)
			}
			firstLine = false
			continue
		}
		// Add the value of the column with name UniqueID to the list
		unique_id_list = append(unique_id_list, record[uniqueIdx])
	}
	return unique_id_list, nil
}

// This function just combines all the unique ids from the two csv files
//nolint
func joinUniqueIds(loc1 string, loc2 string, UniqueID string) ([]string, []string, error) {
	// Open the first csv file
	// Combine The values of col with name UniqueId from both the csv files
	list1, err := extractUniqueId(loc1, UniqueID)
	if err != nil {
		return nil, nil,err
	}
	list2, err := extractUniqueId(loc2, UniqueID)
	if err != nil {
		return nil, nil, err
	}
	return list1, list2,nil
}

//nolint
func CalculatePrivateCount(unique_id_list1, unique_id_list2 []string) (int64, int64, error){
	var count int64 = 0
	privateCount, err := dpagg.NewCount(&dpagg.CountOptions{
		Noise: noise.{{noiseType}},
		Epsilon: {{epsilon}},
		MaxPartitionsContributed: {{maxPartitionsContributed}},
	})
	if err != nil {
		return -1, -1, err
	}
	// no. of common ids
	for _, x := range unique_id_list1 {
		for _, y := range unique_id_list2 {
			if x == y {
				count++;
				privateCount.Increment()
			}
		}
	}
	result, err  := privateCount.Result()
	if err != nil {
		return -1, -1, err
	}
	return count, result,  nil
}

//nolint
func writeToCSV(count int64, privateCountResult int64, outputFolderLocation string) {
	// Create a new file
	file, err := os.Create(outputFolderLocation + "/output.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Write the header
	file.WriteString("count, privateCountResult\n")
	// Write the data
	file.WriteString(fmt.Sprintf("%d, %d\n", count, privateCountResult))
}

//nolint
func main() {
	csvlocation1 := "{{csvLocation1}}"
	csvlocation1 := "{{csvLocation2}}}"
	//outputFolderLocation := "{{outputFolderLocation}}"
	unique_Id := "{{uniqueId}}}"
	unique_id_list1, unique_id_list2, err := joinUniqueIds(csvlocation1, csvlocation2, unique_Id)
	if err != nil {
		panic(err)
	}
	count, privateCountResult, err := CalculatePrivateCount(unique_id_list1, unique_id_list2)
	if err != nil {
		panic(err)
	}
	fmt.Println(count, privateCountResult)
	//writeToCSV(count, privateCountResult, outputFolderLocation)
}