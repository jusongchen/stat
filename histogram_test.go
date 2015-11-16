// Copyright 2015 Jusong Chen. All rights reserved.
package stat

import (
	"bufio"
	"fmt"
	"reflect"

	"strings"
	"testing"
)

func TestHistoAnalyze(t *testing.T) {

	bw1 := 0.5
	s1 := `1.2`
	t1 := []float64{1.2}
	h1 := []HistoBin{
		{1, 1 + bw1, 1},
	}

	bw2 := 2.0
	s2 := `0.0
1.2
1.3
-2.3
10
`

	t2 := []float64{10.0, 1.3, 1.2, 0.0, -2.3}
	h2 := []HistoBin{
		{-4.0, -4.0 + bw2, 1},
		{0.0, 0.0 + bw2, 3},
		{10.0, 10.0 + bw2, 1},
	}

	bw3 := 5.0
	s3 := `-29.66339874
32.12055105
18.41178493
2.938947794
19.54987592
46.416079
46.55036662
4.830790915
56.30215256
15.81995617
69.95266781
28.3679655
29.04550694
23.09620432
-18.5052939
62.89160911
60.14851663
`
	t3 := []float64{69.95266781,
		62.89160911,
		60.14851663,
		56.30215256,
		46.55036662,
		46.416079,
		32.12055105,
		29.04550694,
		28.3679655,
		23.09620432,
	}

	h3 := []HistoBin{
		{-30, -30 + bw3, 1},
		{-20, -20 + bw3, 1},
		{0, 0 + bw3, 2},
		{15, 15 + bw3, 3},
		{20, 20 + bw3, 1},
		{25, 25 + bw3, 2},
		{30, 30 + bw3, 1},
		{45, 45 + bw3, 2},
		{55, 55 + bw3, 1},
		{60, 60 + bw3, 2},
		{65, 65 + bw3, 1},
	}

	var HistoTests = []struct {
		binWidth float64
		inStr    string // input
		topN     []float64
		h        []HistoBin
		err      error
	}{

		{-0.5, "", nil, nil, fmt.Errorf("error expected")}, //BAD - negitive bin width
		{0.5, `1.2
invaidFloatNumber
232`, nil, nil, fmt.Errorf("error expected")}, //BAD - invalid float number
		{bw1, s1, t1, h1, nil},
		{bw2, s2, t2, h2, nil},
		{bw3, s3, t3, h3, nil},
	}

	for i, tt := range HistoTests {

		buf := strings.NewReader(tt.inStr)
		s := bufio.NewScanner(buf)
		// func HistoAnalyze(s Scanner, N int, binWidth float64) (int64, []float64, []HistoBin, error) {
		count, topN, histo, err := HistoAnalyze(s, 10, tt.binWidth)

		if tt.err != nil {
			if err == nil {
				t.Errorf("expected an error but no error returned")
			} else if err != nil {
				//we are OK here
				fmt.Printf("\ncase #%d:error returned as expected:\n%s", i+1, err)
			}
		} else {
			//should return results
			if !reflect.DeepEqual(topN, tt.topN) {
				t.Errorf("\nHistoAnalyze: input strint\n %s\n expected TopN %#v,\n actual %#v", tt.inStr, tt.topN, topN)
			}

			if !reflect.DeepEqual(histo, tt.h) {
				t.Errorf("\nHistoAnalyze input strint \n %s\n expected Histograms %#v,\n\n actual %#v", tt.inStr, tt.h, histo)
			}

			fmt.Printf("\n\ncase #%d:Numbers processed -%d", i+1, count)

		}
	}
	fmt.Println()
}
