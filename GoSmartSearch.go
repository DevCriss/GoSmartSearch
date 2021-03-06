package GoSmartSearch

import (
	"fmt"
	"sort"
	"strings"
)

type stringResult struct {
	value    string
	accuracy float32
}

// SearchInMaps returns a map slice formed by the input elements ordered (based on the key) from most to least similar to the input term
func SearchInMaps(elements []map[string]string, term, key string, tolerance float32) ([]map[string]string, error) {

	if err := validateTolerance(tolerance); err != nil {
		return nil, err
	}

	keyValues := make([]string, 0, len(elements))

	for _, item := range elements {
		keyValues = append(keyValues, item[key])
	}

	sortedKeyValues, err := SearchInStrings(keyValues, term, tolerance)

	if err != nil {
		return nil, err
	}

	result := make([]map[string]string, 0, len(sortedKeyValues))

	for _, item := range sortedKeyValues {

		itemMap := findItemInMapSlice(elements, key, item)

		result = append(result, itemMap)
	}

	return result, nil

}

// SearchInStrings returns a slice formed by the input elements ordered from most to least similar to the input term
func SearchInStrings(elements []string, term string, tolerance float32) ([]string, error) {

	if err := validateTolerance(tolerance); err != nil {
		return nil, err
	}

	var tmpResult []stringResult

	for _, currentTerm := range elements {

		var resultObject stringResult
		resultObject.accuracy = calculateAccuracy(term, currentTerm)
		resultObject.value = currentTerm
		if resultObject.accuracy >= tolerance {
			tmpResult = append(tmpResult, resultObject)
		}
	}

	sort.Slice(tmpResult, func(a, b int) bool {
		return tmpResult[a].accuracy > tmpResult[b].accuracy
	})

	result := make([]string, len(tmpResult))
	for i := range tmpResult {
		result[i] = tmpResult[i].value
	}

	return result, nil

}

func calculateAccuracy(original, current string) float32 {

	var hits, hitsExact float32
	var limit int

	if original == current {
		return 1
	}

	original, current = strings.ToLower(original), strings.ToLower(current)

	if original == current {
		return 1
	}

	if len(original) > len(current) {
		limit = len(current)
	} else {
		limit = len(original)
	}

	for i := 0; i < limit; i++ {
		if original[i] == current[i] {
			hitsExact++
		} else {
			for e := 0; e < limit; e++ {
				if (original[i] == current[e]) || (original[e] == current[i]) {
					hits += 0.25
				}
			}
		}
	}

	if int(hitsExact) == len(original) {
		return 1
	}

	hitsExact += hits

	return hitsExact / float32(len(original)) / 4
}

func findItemInMapSlice(elements []map[string]string, key, value string) map[string]string {

	for _, item := range elements {
		if item[key] == value {
			return item
		}
	}
	return nil
}

func validateTolerance(tolerance float32) error {
	if tolerance > 1 || tolerance < 0 {
		return fmt.Errorf("validation error: tolerance (%f) must be in range 0-1", tolerance)
	}
	return nil
}

