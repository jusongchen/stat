// Copyright 2015 Jusong Chen. All rights reserved.
package stat

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

// Scanner is the interface that wraps basic Scan methods.
type Scanner interface {
	Scan() bool
	Text() string
	Err() error
}

type HistoBin struct {
	LowerB    float64
	UpperB    float64
	Frequency int64
}

//HistoAnalyze calls a Scanner to get tokens and convert tokens to numbers
//then,it do analysis to get Top N and Histogram of those numbers
//it returns an error if an token read in cannot be converted to
//a float64 number or the scanner encounter an error
//
// When there is no error detected, it returns
// 		1) number of tokens processed
// 		2) top N numbers
// 		3) histogram data
//		4) nil
func HistoAnalyze(s Scanner, N int, binWidth float64) (int64, []float64, []HistoBin, error) {

	count, topN, histoMap, err := histoAnalyze(s, N, binWidth)
	if err != nil {
		return count, nil, nil, err
	}

	//revers topN
	topNsorted := []float64{}
	topNCount := len(topN)
	for i := range topN {
		topNsorted = append(topNsorted, topN[topNCount-i-1])
	}

	//return histo data in sorted order
	h := []HistoBin{}
	var keys []float64
	for k := range histoMap {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	i := 0
	for _, k := range keys {
		h = append(h, HistoBin{k, k + binWidth, histoMap[k]})
		i++
	}

	//return topN in reversted order, also, histo data sorted
	return count, topNsorted, h, err

}

//HistoAnalyze calls a Scanner to get items and convert items to numbers and do analysis to get Top N and Histogram of those numbers
// it returns an error if a item read in cannot be converted to a float64 number or the scanner encountered an error
// When there is no error detected, it returns
// 		1)number count read in
// 		2) top N numbers as a slice
// 		3) histogram data as a map
//		4) nil (no error)
func histoAnalyze(s Scanner, N int, binWidth float64) (int64, []float64, map[float64]int64, error) {

	//init return parameters
	count := int64(0)
	topN := []float64{}
	histo := map[float64]int64{}
	err := error(nil)

	//check pass in parameter
	if N < 1 || binWidth < 0.0 {
		return -1, nil, nil, fmt.Errorf("HistoAnalyze: negitive N or bin width - %d,%f", N, binWidth)
	}

	for s.Scan() {
		count++
		// fmt.Printf("\n%s", s.Text())

		// func ParseFloat(s string, bitSize int) (f float64, err error)
		num, err := strconv.ParseFloat(s.Text(), 64)
		if err != nil {
			return -1, nil, nil, err
		}

		//calculate the lower boundary
		lowerB := float64(math.Floor(num/binWidth)) * binWidth

		// fmt.Printf("\nlower Boundery %f, number %f, d.binWidth %f, int64(num/d.binWidth) %d ", lowerB, num, d.binWidth, int64(num/d.binWidth))

		// increase frquency of the bin
		histo[lowerB]++

		//Now handling Top N
		switch {
		//append the first TOP_N items when the TopN slice is not long enough
		case len(topN) < N:
			topN = append(topN, num)
			sort.Float64s(topN)
		//the new number is greater than one of the existing TopN
		case num > topN[0]:
			// if new number is bigger than the first one, which is the min value in the TopN as TopN slice is sorted
			// replace the first value
			topN[0] = num
			sort.Float64s(topN)
		}

	}
	//check if the scanner has encountered any error
	if err = s.Err(); err != nil {
		return -1, nil, nil, err
	}

	// fmt.Printf("\n\nNOT SORTED%#v", histo)
	return count, topN, histo, err
}
