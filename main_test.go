package main

import (
	"math"
	"testing"
)

func TestHaversine(t *testing.T) {

	//Same coordinates
	firstPoint := brewery{0, "firstPoint", 51.355468, 11.100790, []string{}, 0.}
	secondPoint := brewery{0, "secondPoint", 51.355468, 11.100790, []string{}, 0.}
	expected := 0.
	got := haversine(firstPoint, secondPoint)

	if got != expected {
		t.Errorf("haversine returned value is different. \nExpected: %v\nGot: %v", expected, got)
	}

	//Should return around 169.5 km
	firstPoint = brewery{0, "firstPoint", 51.355468, 11.100790, []string{}, 0.}
	secondPoint = brewery{0, "secondPoint", 50, 10, []string{}, 0.}
	expected = 169.50341118844037
	got = haversine(firstPoint, secondPoint)

	if got != expected {
		t.Errorf("haversine returned value is different. \nExpected: %v\nGot: %v", expected, got)
	}

	//Should return around 20015 km
	firstPoint = brewery{0, "firstPoint", 90, 180, []string{}, 0.}
	secondPoint = brewery{0, "secondPoint", -90, -180, []string{}, 0.}
	expected = 20015.086796020572
	got = haversine(firstPoint, secondPoint)

	if got != expected {
		t.Errorf("haversine returned value is different. \nExpected: %v\nGot: %v", expected, got)
	}
}

func TestGetCloseBreweries(t *testing.T) {
	//Should return 6 breweries
	breweries := []brewery{
		{1267, "Tin Whistle Brewing", 49.49440002441406, -119.61000061035156, []string{}, 7869.721023294844},
		{1268, "Titletown Brewing", 44.52000045776367, -88.01730346679688, []string{}, 6839.141532625284},
		{1269, "Tivoli Brewing", 39.739200592041016, -104.98500061035156, []string{}, 8145.293542222566},
		{1271, "Tommyknocker Brewery and Pub", 39.741798400878906, -105.51799774169922, []string{}, 8171.6742897312},
		{1272, "Tomos Watkin and Sons Ltd.", 51.66659927368164, -3.9442999362945557, []string{}, 1039.9077572663819},
		{1274, "Tooheys", -33.850101470947266, 151.0449981689453, []string{}, 16268.64846307816},

		{1056, "Ridgeway Brewing", 51.546199798583984, -1.1354999542236328, []string{}, 847.1914449316964},
		{1083, "Ruppaner-Brauerei", 47.68550109863281, 9.208000183105469, []string{}, 430.3081461406065},
		{1088, "SA Brain & Co. Ltd.", 51.47359848022461, -3.178999900817871, []string{}, 988.8234297974843},
		{1096, "Salopian Brewery", 52.70690155029297, -2.7869999408721924, []string{}, 960.325595317151},
		{1099, "Samuel Smith Old Brewery (Tadcaster)", 53.883399963378906, -1.2625000476837158, []string{}, 879.3105926211241},
		{1111, "Sarah Hughes Brewery", 52.54349899291992, -2.115600109100342, []string{}, 914.0475492655376}}

	expected := []brewery{
		{1056, "Ridgeway Brewing", 51.546199798583984, -1.1354999542236328, []string{}, 847.1914449316964},
		{1083, "Ruppaner-Brauerei", 47.68550109863281, 9.208000183105469, []string{}, 430.3081461406065},
		{1088, "SA Brain & Co. Ltd.", 51.47359848022461, -3.178999900817871, []string{}, 988.8234297974843},
		{1096, "Salopian Brewery", 52.70690155029297, -2.7869999408721924, []string{}, 960.325595317151},
		{1099, "Samuel Smith Old Brewery (Tadcaster)", 53.883399963378906, -1.2625000476837158, []string{}, 879.3105926211241},
		{1111, "Sarah Hughes Brewery", 52.54349899291992, -2.115600109100342, []string{}, 914.0475492655376}}

	got := getCloseBreweries(breweries, 2000.)

	if len(got) != len(expected) {
		t.Errorf("getCloseBreweries should have returned empty slice. \nExpected: %v\n got: %v", len(expected), len(got))
	}

	for i := range expected {
		if got[i].ID != expected[i].ID {
			t.Errorf("getCloseBreweries slice values are different. \nExpected: %v\n got: %v", expected[i], got[i])
		}
	}

	//When list of breweries is empty
	breweries = []brewery{}

	expected = []brewery{}
	got = getCloseBreweries(breweries, 2000.)

	if len(got) != len(expected) {
		t.Errorf("getCloseBreweries should have returned empty slice. \nExpected: %v\n got: %v", len(expected), len(got))
	}

	for i := range expected {
		if got[i].ID != expected[i].ID {
			t.Errorf("getCloseBreweries slice values are different. \nExpected: %v\n got: %v", expected[i], got[i])
		}
	}
}

func TestMakeDistancesGraph(t *testing.T) {
	//With 3 breweries
	breweries := []brewery{
		{1056, "Ridgeway Brewing", 51.546199798583984, -1.1354999542236328, []string{}, 847.1914449316964},
		{1083, "Ruppaner-Brauerei", 47.68550109863281, 9.208000183105469, []string{}, 430.3081461406065},
		{1088, "SA Brain & Co. Ltd.", 51.47359848022461, -3.178999900817871, []string{}, 988.8234297974843}}

	expected := [][]float64{
		{0, 858.8600028094968, 141.6468895179824},
		{858.8600028094968, 0, 985.756263328248},
		{141.6468895179824, 985.756263328248, 0}}

	got := makeDistancesGraph(breweries)

	if len(got) != len(expected) {
		t.Errorf("makeDistancesGraph returned slices is different length. \nExpected: %v\n got: %v", len(expected), len(got))
	}

	for i := range expected {
		for j := range expected[i] {
			if got[i][j] != expected[i][j] {
				t.Errorf("makeDistancesGraph graph values are different. \nExpected: %v\n got: %v", expected[i][j], got[i][j])
			}
		}
	}

	//Passing empty breweries slice
	breweries = []brewery{}

	expected = [][]float64{}

	got = makeDistancesGraph(breweries)

	if len(got) != len(expected) {
		t.Errorf("makeDistancesGraph should have returned 0 length slice. \nExpected: %v\n got: %v", len(expected), len(got))
	}

	//Passing 1 brewery slice
	breweries = []brewery{{1056, "Ridgeway Brewing", 51.546199798583984, -1.1354999542236328, []string{}, 847.1914449316964}}

	expected = [][]float64{{0}}

	got = makeDistancesGraph(breweries)

	if len(got) != len(expected) {
		t.Errorf("makeDistancesGraph should have returned 1 length slice. \nExpected: %v\n got: %v", len(expected), len(got))
	}

	for i := range expected {
		for j := range expected[i] {
			if got[i][j] != expected[i][j] {
				t.Errorf("makeDistancesGraph graph values are different. \nExpected: %v\n got: %v", expected[i][j], got[i][j])
			}
		}
	}
}

func TestGetUniqueBeers(t *testing.T) {

	//Normal conditions
	beers := []string{"one", "two", "one"}
	expected := []string{"one", "two"}
	got := getUniqueBeers(beers)

	if len(got) != len(expected) {
		t.Errorf("getUniqueBeers should have returned 2 length slice. \nExpected: %v\n got: %v", len(expected), len(got))
	}

	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("makeDistancesGraph slice values are different. \nExpected: %v\n got: %v", expected[i], got[i])
		}
	}

	//Slice consists of only same element
	beers = []string{"one", "one", "one"}
	expected = []string{"one"}
	got = getUniqueBeers(beers)

	if len(got) != len(expected) {
		t.Errorf("getUniqueBeers should have returned 1 length slice. \nExpected: %v\n got: %v", len(expected), len(got))
	}

	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("makeDistancesGraph slice values are different. \nExpected: %v\n got: %v", expected[i], got[i])
		}
	}

	//Slice is empty
	beers = []string{}
	expected = []string{}
	got = getUniqueBeers(beers)

	if len(got) != len(expected) {
		t.Errorf("getUniqueBeers should have returned empty slice. \nExpected: %v\n got: %v", len(expected), len(got))
	}

	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("makeDistancesGraph slice values are different. \nExpected: %v\n got: %v", expected[i], got[i])
		}
	}
}

func TestCountCostandBound(t *testing.T) {

	// Normal case

	breweries := []brewery{
		{0, "HOME", 57, 35, []string{}, 0},
		{957, "Ostankinskij Pivovarennij Zavod", 55.75579833984375, 37.61759948730469, []string{"Beer"}, 212.37640863461093},
		{961, "Oy Sinebrychoff AB", 60.38100051879883, 25.110200881958008, []string{"Porter IV", "Koff Special III"}, 682.979146237569},
		{999, "Pivzavod Baltika /", 59.93899917602539, 30.315799713134766, []string{"Baltika 6 Porter", "Baltika #5", "Baltika #8", "Baltika #9"}, 425.22565488237495},
		{1094, "Saku lletehas", 59.30139923095703, 24.66790008544922, []string{"Porter"}, 657.168864470112},
		{1329, "Vivungs Bryggeri", 57.49850082397461, 18.458999633789062, []string{"Romakloster", "DragÃ¶l"}, 994.0936874184431}}

	eCost, eBound := -7.18, -6.

	gCost, gBound := countCostandBound(breweries, 2000., []int{2})
	gCost = math.Round(gCost*100) / 100

	if eCost != gCost {
		t.Errorf("countCostandBound returned cost value is different than expected. \nExpected: %v\n got: %v", eCost, gCost)
	}

	if eBound != gBound {
		t.Errorf("countCostandBound returned bound value is different than expected. \nExpected: %v\n got: %v", eBound, gBound)
	}

	// Given slice is empty

	breweries = []brewery{}

	eCost, eBound = 0, 0.

	gCost, gBound = countCostandBound(breweries, 2000., []int{2})
	gCost = math.Round(gCost*100) / 100

	if eCost != gCost {
		t.Errorf("countCostandBound returned cost value is different than expected. \nExpected: %v\n got: %v", eCost, gCost)
	}

	if eBound != gBound {
		t.Errorf("countCostandBound returned bound value is different than expected. \nExpected: %v\n got: %v", eBound, gBound)
	}

	// Edge case when only Cost should be returned

	breweries = []brewery{
		{0, "HOME", 57, 35, []string{}, 0},
		{957, "Ostankinskij Pivovarennij Zavod", -55, -37, []string{"Beer"}, 2000}}

	eCost, eBound = -0.12, 0.

	gCost, gBound = countCostandBound(breweries, 2000., []int{2})
	gCost = math.Round(gCost*100) / 100

	if eCost != gCost {
		t.Errorf("countCostandBound returned cost value is different than expected. \nExpected: %v\n got: %v", eCost, gCost)
	}

	if eBound != gBound {
		t.Errorf("countCostandBound returned bound value is different than expected. \nExpected: %v\n got: %v", eBound, gBound)
	}
}
