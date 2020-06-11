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
		t.Errorf("Same coordinates case haversine distance returned value is different. \nExpected: %v\nGot: %v", expected, got)
	}

	//Should return around 169.5 km
	firstPoint = brewery{0, "firstPoint", 51.355468, 11.100790, []string{}, 0.}
	secondPoint = brewery{0, "secondPoint", 50, 10, []string{}, 0.}
	expected = 169.50341118844037
	got = haversine(firstPoint, secondPoint)

	if got != expected {
		t.Errorf("dinateshaversine distance returned value is different. \nExpected: %v\nGot: %v", expected, got)
	}

	//Should return around 20015 km
	firstPoint = brewery{0, "firstPoint", 90, 180, []string{}, 0.}
	secondPoint = brewery{0, "secondPoint", -90, -180, []string{}, 0.}
	expected = 20015.086796020572
	got = haversine(firstPoint, secondPoint)

	if got != expected {
		t.Errorf("haversine distance returned value is different. \nExpected: %v\nGot: %v", expected, got)
	}
}

func TestGetCloseBreweries(t *testing.T) {

	home := brewery{0, "HOME", 51.355468, 11.100790, []string{}, 0.}
	maxDistance := 2000.
	//Normal case
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

	got := getCloseBreweries(breweries, home, maxDistance/2)

	if len(got) != len(expected) {
		t.Errorf("Normal case getCloseBreweries returned breweries slice are different than expected. \nExpected length: %v\n got length: %v", len(expected), len(got))
	}

	for i := range expected {
		if got[i].ID != expected[i].ID {
			t.Errorf("getCloseBreweries slice values are different. \nExpected: %v\n got: %v", expected[i], got[i])
		}
	}

	//Slice of breweries only has breweries that are futher than maxDistance/2
	breweries = []brewery{{1267, "Tin Whistle Brewing", 49.49440002441406, -119.61000061035156, []string{}, 7869.721023294844},
		{1268, "Titletown Brewing", 44.52000045776367, -88.01730346679688, []string{}, 6839.141532625284},
		{1269, "Tivoli Brewing", 39.739200592041016, -104.98500061035156, []string{}, 8145.293542222566},
		{1271, "Tommyknocker Brewery and Pub", 39.741798400878906, -105.51799774169922, []string{}, 8171.6742897312},
		{1272, "Tomos Watkin and Sons Ltd.", 51.66659927368164, -3.9442999362945557, []string{}, 1039.9077572663819},
		{1274, "Tooheys", -33.850101470947266, 151.0449981689453, []string{}, 16268.64846307816}}

	expected = []brewery{}
	got = getCloseBreweries(breweries, home, maxDistance/2)

	if len(got) != len(expected) {
		t.Errorf("Slice of breweries only has breweries that are futher than maxDistance/2 case \ngetCloseBreweries returned breweries slice are different than expected. \nExpected length: %v\n got length: %v", len(expected), len(got))
	}

	for i := range expected {
		if got[i].ID != expected[i].ID {
			t.Errorf("Slice of breweries only has breweries that are futher than maxDistance/2 case \ngetCloseBreweries slice values are different. \nExpected: %v\n got: %v", expected[i], got[i])
		}
	}
}

func TestGetUniqueBeers(t *testing.T) {

	//Normal case
	beers := []string{"one", "two", "one"}
	expected := []string{"one", "two"}
	got := getUniqueBeers(beers)

	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("Normal case getUniqueBeers returned slice values are different. \nExpected: %v\n got: %v", expected[i], got[i])
		}
	}

	//Slice consists of only same element
	beers = []string{"one", "one", "one"}
	expected = []string{"one"}
	got = getUniqueBeers(beers)

	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("Slice consists of only same element case getUniqueBeers returned slice values are different. \nExpected: %v\n got: %v", expected[i], got[i])
		}
	}

	//Beers slice is empty
	beers = []string{}
	expected = []string{}
	got = getUniqueBeers(beers)

	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("Beers slice is empty case getUniqueBeers returned slice values are different. \nExpected: %v\n got: %v", expected[i], got[i])
		}
	}
}

func TestCalcFirstProfitAndBound(t *testing.T) {

	// Normal case

	breweries := []brewery{
		{0, "HOME", 57, 35, []string{}, 0},
		{957, "Ostankinskij Pivovarennij Zavod", 55.75579833984375, 37.61759948730469, []string{"Beer"}, 212.37640863461093},
		{961, "Oy Sinebrychoff AB", 60.38100051879883, 25.110200881958008, []string{"Porter IV", "Koff Special III"}, 682.979146237569},
		{999, "Pivzavod Baltika /", 59.93899917602539, 30.315799713134766, []string{"Baltika 6 Porter", "Baltika #5", "Baltika #8", "Baltika #9"}, 425.22565488237495},
		{1094, "Saku lletehas", 59.30139923095703, 24.66790008544922, []string{"Porter"}, 657.168864470112},
		{1329, "Vivungs Bryggeri", 57.49850082397461, 18.458999633789062, []string{"Romakloster", "DragÃ¶l"}, 994.0936874184431}}

	eprofit, eBound := -7.61, -7.

	gProfit, gBound := calcFirstProfitAndBound(breweries, 2000.)
	gProfit = math.Round(gProfit*100) / 100

	if eprofit != gProfit {
		t.Errorf("Normal case countProfitandBound returned profit value is different than expected. \nExpected: %v\n got: %v", eprofit, gProfit)
	}

	if eBound != gBound {
		t.Errorf("Normal case countProfitandBound returned bound value is different than expected. \nExpected: %v\n got: %v", eBound, gBound)
	}

	// Edge case when only Profit should be returned

	breweries = []brewery{
		{0, "HOME", 57, 35, []string{}, 0},
		{957, "Ostankinskij Pivovarennij Zavod", -55, -37, []string{"Beer"}, 2000}}

	eprofit, eBound = -0.12, 0.

	gProfit, gBound = calcFirstProfitAndBound(breweries, 2000.)
	gProfit = math.Round(gProfit*100) / 100

	if eprofit != gProfit {
		t.Errorf("Edge case when only Profit should be returned countProfitandBound returned profit value is different than expected. \nExpected: %v\n got: %v", eprofit, gProfit)
	}

	if eBound != gBound {
		t.Errorf("Edge case when only Profit should be returned countProfitandBound returned bound value is different than expected. \nExpected: %v\n got: %v", eBound, gBound)
	}
}

func TestCalcProfitAndBound(t *testing.T) {

	// Normal case

	breweries := []brewery{
		{0, "HOME", 57, 35, []string{}, 0},
		{957, "Ostankinskij Pivovarennij Zavod", 55.75579833984375, 37.61759948730469, []string{"Beer"}, 212.37640863461093},
		{961, "Oy Sinebrychoff AB", 60.38100051879883, 25.110200881958008, []string{"Porter IV", "Koff Special III"}, 682.979146237569},
		{999, "Pivzavod Baltika /", 59.93899917602539, 30.315799713134766, []string{"Baltika 6 Porter", "Baltika #5", "Baltika #8", "Baltika #9"}, 425.22565488237495},
		{1094, "Saku lletehas", 59.30139923095703, 24.66790008544922, []string{"Porter"}, 657.168864470112},
		{1329, "Vivungs Bryggeri", 57.49850082397461, 18.458999633789062, []string{"Romakloster", "DragÃ¶l"}, 994.0936874184431}}

	eIncludedBrews := []brewery{{957, "Ostankinskij Pivovarennij Zavod", 55.75579833984375, 37.61759948730469, []string{"Beer"}, 212.37640863461093}}
	eProfit, eBound := -7.18, -6.

	gProfit, gBound, includedBrews := calcProfitAndBound(breweries, []int{1}, []int{2}, 2000.)
	gProfit = math.Round(gProfit*100) / 100

	if eProfit != gProfit {
		t.Errorf("Normal case countProfitandBound returned profit value is different than expected. \nExpected: %v\n got: %v", eProfit, gProfit)
	}

	if eBound != gBound {
		t.Errorf("Normal case countProfitandBound returned bound value is different than expected. \nExpected: %v\n got: %v", eBound, gBound)
	}

	for i := range eIncludedBrews {
		if includedBrews[i].ID != eIncludedBrews[i].ID {
			t.Errorf("Normal case makeDistancesGraph slice values are different. \nExpected: %v\n got: %v", eIncludedBrews[i], includedBrews[i])
		}
	}

	// Edge case when only Profit should be returned

	breweries = []brewery{
		{0, "HOME", 57, 35, []string{}, 0},
		{957, "Ostankinskij Pivovarennij Zavod", -55, -37, []string{"Beer"}, 2000}}
	eIncludedBrews = []brewery{}
	eProfit, eBound = -0.12, 0.

	gProfit, gBound, includedBrews = calcProfitAndBound(breweries, []int{1}, []int{}, 2000.)
	gProfit = math.Round(gProfit*100) / 100

	if eProfit != gProfit {
		t.Errorf("Edge case when only Profit should be returned countProfitandBound returned profit value is different than expected. \nExpected: %v\n got: %v", eProfit, gProfit)
	}

	if eBound != gBound {
		t.Errorf("Edge case when only Profit should be returned countProfitandBound returned bound value is different than expected. \nExpected: %v\n got: %v", eBound, gBound)
	}

	for i := range eIncludedBrews {
		if includedBrews[i].ID != eIncludedBrews[i].ID {
			t.Errorf("Edge case when only Profit should be returned makeDistancesGraph slice values are different. \nExpected: %v\n got: %v", eIncludedBrews[i], includedBrews[i])
		}
	}
}

func TestUpdateNodes(t *testing.T) {

	// Normal case

	nodes := []node{
		{1, -6.1, -6, []brewery{}, []int{}, []int{}, 0},
		{2, -5.1, -5, []brewery{}, []int{}, []int{}, 0},
		{3, -7.1, -5, []brewery{}, []int{}, []int{}, 0}}

	expected := []node{
		{1, -6.1, -6, []brewery{}, []int{}, []int{}, 0},
		{3, -7.1, -5, []brewery{}, []int{}, []int{}, 0}}

	got := updateNodes(nodes, -6)

	for i := range expected {
		if got[i].level != expected[i].level {
			t.Errorf("Normal case updateNodes slice values are different. \nExpected: %v\n got: %v", expected[i], got[i])
		}
	}

	// Given nodes slice is empty

	nodes = []node{}

	expected = []node{}

	got = updateNodes(nodes, -6)

	for i := range expected {
		if got[i].level != expected[i].level {
			t.Errorf("Given nodes slice is empty case updateNodes slice values are different. \nExpected: %v\n got: %v", expected[i], got[i])
		}
	}
}

func TestAddNode(t *testing.T) {

	// Normal case

	breweries := []brewery{
		{0, "HOME", 57, 35, []string{}, 0},
		{957, "Ostankinskij Pivovarennij Zavod", 55.75579833984375, 37.61759948730469, []string{"Beer"}, 212.37640863461093},
		{961, "Oy Sinebrychoff AB", 60.38100051879883, 25.110200881958008, []string{"Porter IV", "Koff Special III"}, 682.979146237569},
		{999, "Pivzavod Baltika /", 59.93899917602539, 30.315799713134766, []string{"Baltika 6 Porter", "Baltika #5", "Baltika #8", "Baltika #9"}, 425.22565488237495},
		{1094, "Saku lletehas", 59.30139923095703, 24.66790008544922, []string{"Porter"}, 657.168864470112},
		{1329, "Vivungs Bryggeri", 57.49850082397461, 18.458999633789062, []string{"Romakloster", "DragÃ¶l"}, 994.0936874184431}}

	currentNode := node{1, -7.1, -7,
		[]brewery{
			{957, "Ostankinskij Pivovarennij Zavod", 55.75579833984375, 37.61759948730469, []string{"Beer"}, 212.37640863461093},
			{961, "Oy Sinebrychoff AB", 60.38100051879883, 25.110200881958008, []string{"Porter IV", "Koff Special III"}, 682.979146237569}}, []int{}, []int{}, 0}

	eBNodes := []node{currentNode}

	eNodes := []node{currentNode}

	eBBeerCnt := 3

	gNodes, gBNodes, gBBeerCnt := addNode(breweries, []node{}, []node{}, currentNode, 2000., -7., 2)

	if gBBeerCnt != eBBeerCnt {
		t.Errorf("normal case addNode returned beerCount values were different. \nExpected: %v\n got: %v", eBBeerCnt, gBBeerCnt)
	}

	for i := range eBNodes {
		if eBNodes[i].level != gBNodes[i].level {
			t.Errorf("normal case addNode returned bestNodes slice values are different. \nExpected: %v\n got: %v", eBNodes[i], gBNodes[i])
		}
	}

	for i := range eNodes {
		if eNodes[i].level != gNodes[i].level {
			t.Errorf("normal case addNode returned nodes slice values are different. \nExpected: %v\n got: %v", eNodes[i], gNodes[i])
		}
	}

	// Profit is less than upper bound

	currentNode = node{1, -7.1, -7,
		[]brewery{
			{957, "Ostankinskij Pivovarennij Zavod", 55.75579833984375, 37.61759948730469, []string{"Beer"}, 212.37640863461093},
			{961, "Oy Sinebrychoff AB", 60.38100051879883, 25.110200881958008, []string{"Porter IV", "Koff Special III"}, 682.979146237569}}, []int{}, []int{}, 0}

	eBNodes = []node{}

	eNodes = []node{}

	eBBeerCnt = 3

	gNodes, gBNodes, gBBeerCnt = addNode(breweries, []node{}, []node{}, currentNode, 2000., -8., 2)

	if gBBeerCnt != eBBeerCnt {
		t.Errorf("Profit is less than upper bound case addNode returned beerCount values were different. \nExpected: %v\n got: %v", eBBeerCnt, gBBeerCnt)
	}

	for i := range eBNodes {
		if eBNodes[i].level != gBNodes[i].level {
			t.Errorf("Profit is less than upper case bound addNode returned bestNodes slice values are different. \nExpected: %v\n got: %v", eBNodes[i], gBNodes[i])
		}
	}

	for i := range eNodes {
		if eNodes[i].level != gNodes[i].level {
			t.Errorf("Profit is less than upper bound case addNode returned nodes slice values are different. \nExpected: %v\n got: %v", eNodes[i], gNodes[i])
		}
	}

	// Profit is less than upper bound

	currentNode = node{1, -7.1, -7,
		[]brewery{
			{957, "Ostankinskij Pivovarennij Zavod", 55.75579833984375, 37.61759948730469, []string{"Beer"}, 212.37640863461093},
			{961, "Oy Sinebrychoff AB", 60.38100051879883, 25.110200881958008, []string{"Porter IV", "Koff Special III"}, 682.979146237569}}, []int{}, []int{}, 0}

	eBNodes = []node{}

	eNodes = []node{}

	eBBeerCnt = 3

	gNodes, gBNodes, gBBeerCnt = addNode(breweries, []node{}, []node{}, currentNode, 2000., -8., 2)

	if gBBeerCnt != eBBeerCnt {
		t.Errorf("Profit is less than upper bound case addNode returned beerCount values were different. \nExpected: %v\n got: %v", eBBeerCnt, gBBeerCnt)
	}

	for i := range eBNodes {
		if eBNodes[i].level != gBNodes[i].level {
			t.Errorf("Profit is less than upper bound case addNode returned bestNodes slice values are different. \nExpected: %v\n got: %v", eBNodes[i], gBNodes[i])
		}
	}

	for i := range eNodes {
		if eNodes[i].level != gNodes[i].level {
			t.Errorf("Profit is less than upper bound case addNode returned nodes slice values are different. \nExpected: %v\n got: %v", eNodes[i], gNodes[i])
		}
	}

	// There is no breweries to visit after. Should return empty nodes slice, but current path distance is under 2000 km

	currentNode = node{1, -7.1, -7,
		[]brewery{
			{957, "Ostankinskij Pivovarennij Zavod", 55.75579833984375, 37.61759948730469, []string{"Beer"}, 212.37640863461093},
			{961, "Oy Sinebrychoff AB", 60.38100051879883, 25.110200881958008, []string{"Porter IV", "Koff Special III"}, 682.979146237569},
			{999, "Pivzavod Baltika /", 59.93899917602539, 30.315799713134766, []string{"Baltika 6 Porter", "Baltika #5", "Baltika #8", "Baltika #9"}, 425.22565488237495}}, []int{}, []int{}, 0}

	eBNodes = []node{currentNode}

	eNodes = []node{}

	eBBeerCnt = 7

	gNodes, gBNodes, gBBeerCnt = addNode(breweries, []node{}, []node{}, currentNode, 2000., -7., 2)

	if gBBeerCnt != eBBeerCnt {
		t.Errorf("There is no breweries to visit after case addNode returned beerCount values were different. \nExpected: %v\n got: %v", eBBeerCnt, gBBeerCnt)
	}

	for i := range eBNodes {
		if eBNodes[i].level != gBNodes[i].level {
			t.Errorf("There is no breweries to visit after case addNode returned bestNodes slice values are different. \nExpected: %v\n got: %v", eBNodes[i], gBNodes[i])
		}
	}

	for i := range eNodes {
		if eNodes[i].level != gNodes[i].level {
			t.Errorf("There is no breweries to visit after case addNode returned nodes slice values are different. \nExpected: %v\n got: %v", eNodes[i], gNodes[i])
		}
	}

	// Total distance of given breweries is more than 2000. Should return empty nodes, bestNodes slices

	currentNode = node{1, -7.1, -7,
		[]brewery{
			{957, "Ostankinskij Pivovarennij Zavod", 55.75579833984375, 37.61759948730469, []string{"Beer"}, 212.37640863461093},
			{961, "Oy Sinebrychoff AB", 60.38100051879883, 25.110200881958008, []string{"Porter IV", "Koff Special III"}, 682.979146237569},
			{1329, "Vivungs Bryggeri", 57.49850082397461, 18.458999633789062, []string{"Romakloster", "DragÃ¶l"}, 994.0936874184431}}, []int{}, []int{}, 0}

	eBNodes = []node{}

	eNodes = []node{}

	eBBeerCnt = 3

	gNodes, gBNodes, gBBeerCnt = addNode(breweries, []node{}, []node{}, currentNode, 2000., -7., 3)

	if gBBeerCnt != eBBeerCnt {
		t.Errorf("Total distance of given breweries is more than 2000 case addNode returned beerCount values were different. \nExpected: %v\n got: %v", eBBeerCnt, gBBeerCnt)
	}

	for i := range eBNodes {
		if eBNodes[i].level != gBNodes[i].level {
			t.Errorf("Total distance of given breweries is more than 2000 case addNode returned bestNodes slice values are different. \nExpected: %v\n got: %v", eBNodes[i], gBNodes[i])
		}
	}

	for i := range eNodes {
		if eNodes[i].level != gNodes[i].level {
			t.Errorf("Total distance of given breweries is more than 2000 case addNode returned nodes slice values are different. \nExpected: %v\n got: %v", eNodes[i], gNodes[i])
		}
	}
}

func TestFindPaths(t *testing.T) {

	// Normal case

	breweries := []brewery{
		{0, "HOME", 57, 35, []string{}, 0},
		{957, "Ostankinskij Pivovarennij Zavod", 55.75579833984375, 37.61759948730469, []string{"Beer"}, 212.37640863461093},
		{961, "Oy Sinebrychoff AB", 60.38100051879883, 25.110200881958008, []string{"Porter IV", "Koff Special III"}, 682.979146237569},
		{999, "Pivzavod Baltika /", 59.93899917602539, 30.315799713134766, []string{"Baltika 6 Porter", "Baltika #5", "Baltika #8", "Baltika #9"}, 425.22565488237495},
		{1094, "Saku lletehas", 59.30139923095703, 24.66790008544922, []string{"Porter"}, 657.168864470112},
		{1329, "Vivungs Bryggeri", 57.49850082397461, 18.458999633789062, []string{"Romakloster", "DragÃ¶l"}, 994.0936874184431}}

	eBNodes := []node{
		{0, 0, 0, []brewery{}, []int{}, []int{}, 0},
		{0, 0, 0, []brewery{}, []int{}, []int{}, 0},
		{0, 0, 0, []brewery{}, []int{}, []int{}, 0},
		{0, 0, 0, []brewery{}, []int{}, []int{}, 0},
		{0, 0, 0, []brewery{}, []int{}, []int{}, 0},
		{0, 0, 0, []brewery{}, []int{}, []int{}, 0},
		{0, 0, 0, []brewery{}, []int{}, []int{}, 0}}

	eBBeerCnt := 7

	gBNodes, gBBeerCnt := findPaths(breweries, 2000.)

	if gBBeerCnt != eBBeerCnt {
		t.Errorf("findPaths normal case returned beerCount values were different. \nExpected: %v\n got: %v", eBBeerCnt, gBBeerCnt)
	}

	if len(eBNodes) != len(gBNodes) {
		t.Errorf("findPaths normal case returned bestNodes slice are different than expected. Expected length: %v\n got length: %v", len(eBNodes), len(gBNodes))
	}
}
