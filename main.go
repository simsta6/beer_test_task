package main

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	maxDistance := 2000.

	checkArguments()
	lat, err := strconv.ParseFloat(os.Args[1], 64)
	checkForError(err)
	lon, err := strconv.ParseFloat(os.Args[2], 64)
	checkForError(err)

	home := brewery{0, "HOME", lat, lon, []string{}, 0.}
	breweries := getBreweriesWithBeersFromDB(home)

	start := time.Now()
	breweries = getCloseBreweries(breweries, home, maxDistance/2)

	if len(breweries) == 0 {
		fmt.Printf("There is no solution\n")
		os.Exit(1)
	}

	breweries = append([]brewery{home}, breweries...)

	//Beer count descending
	sort.Slice(breweries[:], func(i, j int) bool {
		brew1 := len(breweries[i].beers)
		brew2 := len(breweries[j].beers)
		return brew1 > brew2
	})

	nodes, bestBeerCount := findPaths(breweries, maxDistance)

	//Finds best path
	var bestPath []brewery
	bestDistance := 2000.

	for _, i := range nodes {
		distance := i.getTraveledDistance()
		if i.beerCnt >= bestBeerCount && bestDistance > distance {
			bestPath = []brewery{}
			bestDistance = distance
			for _, j := range i.includedBrews {
				bestPath = append(bestPath, j)
			}
		}
	}

	//Gets slice of beer in best path
	beers := []string{}
	for i := range bestPath {
		beers = append(beers, bestPath[i].beers...)
	}
	beers = getUniqueBeers(beers)

	//Adds home as first and last stop
	bestPath = append([]brewery{home}, bestPath...)
	bestPath = append(bestPath, home)

	elapsed := time.Since(start)

	printResults(beers, elapsed, bestPath, bestDistance)
}

func checkArguments() {
	if len(os.Args) != 3 {
		fmt.Printf("Program usage:\n%v latitude longitude\n", os.Args[0])
		os.Exit(1)
	}

	lat, err := strconv.ParseFloat(os.Args[1], 64)
	checkForError(err)
	lon, err := strconv.ParseFloat(os.Args[2], 64)
	checkForError(err)

	if lat > 90 || lat < -90 {
		fmt.Printf("Latitudes range from -90 to 90")
		os.Exit(1)
	}

	if lon > 180 || lon < -180 {
		fmt.Printf("Longitudes range from -180 to 180")
		os.Exit(1)
	}
}

func printResults(beers []string, elapsed time.Duration, bestPath []brewery, bestDistance float64) {
	fmt.Printf("Found %v beer factories: \n", len(bestPath)-2)

	for i := 0; i < len(bestPath); i++ {
		if i == 0 {
			bestPath[i].printWithDistance(bestPath[i])
			continue
		}
		bestPath[i].printWithDistance(bestPath[i-1])
	}

	fmt.Printf("\nTotal distance traveled: %.1fkm\n", bestDistance)

	fmt.Printf("\n\nCollected %v beer types:\n", len(beers))

	for i := range beers {
		fmt.Printf("\t%v\n", beers[i])
	}

	fmt.Printf("\nProgram took: %s\n", elapsed)
}

func getBreweriesWithBeersFromDB(home brewery) []brewery {
	breweries := make([]brewery, 0)

	var (
		conn  *sql.DB
		err   error
		retry int
	)

	for retry = 0; retry < 10; retry++ {
		conn, err = sql.Open("mysql", "root:@tcp(db)/beer-database")
		if err == nil {
			break // success!
		}
		<-time.After(1 * time.Second)
	}
	if retry == 10 {
		fmt.Println("could not connect to db after 10 retries")
		panic("oh no")
	}

	checkForError(err)

	queryForBreweries := `SELECT b.id, b.name, geo.longitude, geo.latitude 
							FROM breweries b 
							LEFT JOIN geocodes geo ON geo.brewery_id = b.id 
							WHERE geo.latitude IS NOT NULL AND geo.longitude IS NOT NULL 
							ORDER BY b.id ASC`

	statement, err := conn.Prepare(queryForBreweries)
	checkForError(err)

	results, err := statement.Query()
	checkForError(err)

	for results.Next() {
		beers := make([]string, 0)
		brew := brewery{}

		err = results.Scan(&brew.ID, &brew.name, &brew.longitude, &brew.latitude)

		checkForError(err)

		queryForBeers := `SELECT beer.name FROM beer 
							RIGHT JOIN breweries ON breweries.id = beer.brewery_id 
							WHERE beer.name IS NOT NULL AND brewery_id = ` + strconv.Itoa(brew.ID)

		bstatement, err := conn.Prepare(queryForBeers)
		checkForError(err)

		bresults, err := bstatement.Query()
		checkForError(err)

		for bresults.Next() {
			var beerName string

			err = bresults.Scan(&beerName)
			checkForError(err)

			beers = append(beers, beerName)
		}

		brew.beers = beers
		brew.distanceToHome = haversine(home, brew)

		if len(brew.beers) > 0 {
			breweries = append(breweries, brew)
		}
	}

	conn.Close()

	return breweries
}

func checkForError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func haversine(firstBrewry brewery, secondBrewery brewery) float64 {
	//Converting to radians
	lat1r := firstBrewry.latitude * math.Pi / 180
	lon1r := firstBrewry.longitude * math.Pi / 180
	lat2r := secondBrewery.latitude * math.Pi / 180
	lon2r := secondBrewery.longitude * math.Pi / 180

	dlat := lat2r - lat1r
	dlon := lon2r - lon1r

	//Haversine formule to calculate ditance between two points
	a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1r)*math.Cos(lat2r)*math.Pow(math.Sin(dlon/2), 2)
	c := 2 * math.Asin(math.Sqrt(a)) * 6371.0

	return c
}

func getCloseBreweries(breweries []brewery, currentB brewery, distance float64) []brewery {
	closeBreweries := make([]brewery, 0)

	for i := range breweries {
		currDistance := haversine(currentB, breweries[i])
		if currDistance <= distance && currDistance != 0. {
			closeBreweries = append(closeBreweries, breweries[i])
		}
	}

	return closeBreweries
}

func getUniqueBeers(beers []string) []string {
	keys := make(map[string]bool)
	uBeers := []string{}
	for _, entry := range beers {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uBeers = append(uBeers, entry)
		}
	}
	return uBeers
}

func calcFirstProfitAndBound(breweries []brewery, maxDistance float64) (float64, float64) {
	var (
		profit, bound, tDistance float64
		lastBrew                 brewery
	)

	if len(breweries) > 0 {
		lastBrew = breweries[0]
	}

	for i := 1; i < len(breweries); i++ {
		distance := haversine(breweries[i], lastBrew)
		beerCnt := float64(len(breweries[i].beers))
		if tDistance+distance+breweries[i].distanceToHome <= maxDistance {
			profit += beerCnt
			bound += beerCnt
			tDistance += distance
			lastBrew = breweries[i]
		} else {
			profit += beerCnt / (distance + breweries[i].distanceToHome) * (maxDistance - tDistance)
			break
		}
	}

	if profit != 0 || bound != 0 {
		profit = -profit
		bound = -bound
	}

	return profit, bound
}

func calcProfitAndBound(breweries []brewery, includedBrewsIndex, excludedBrews []int, maxDistance float64) (float64, float64, []brewery) {
	var (
		profit, bound, tDistance float64
		lastBrew                 brewery
		includedBrews            []brewery
	)
	if len(breweries) > 0 {
		lastBrew = breweries[0]
	}

outterloop:
	for i := 1; i < len(breweries); i++ {

		for j := range excludedBrews {
			if excludedBrews[j] == i {
				continue outterloop
			}
		}

		for j := range includedBrewsIndex {
			if includedBrewsIndex[j] == i {
				includedBrews = append(includedBrews, breweries[i])
			}
		}

		distance := haversine(breweries[i], lastBrew)
		beerCnt := float64(len(breweries[i].beers))
		if tDistance+distance+breweries[i].distanceToHome <= maxDistance {
			profit += beerCnt
			bound += beerCnt
			tDistance += distance
			lastBrew = breweries[i]

		} else {
			profit += beerCnt / (distance + breweries[i].distanceToHome) * (maxDistance - tDistance)
			break
		}
	}

	if profit != 0 || bound != 0 {
		profit = -profit
		bound = -bound
	}

	return profit, bound, includedBrews
}

func updateNodes(nodes []node, upper float64) []node {
	for i := 0; i < len(nodes); i++ {
		if nodes[i].profit > upper && len(nodes) > i {
			nodes[i] = nodes[len(nodes)-1]
			nodes = nodes[:len(nodes)-1]
		}
	}
	return nodes
}

func addNode(breweries []brewery, nodes, bestNodes []node, currentNode node, maxDistance, upper float64, bestBeerCount int) ([]node, []node, int) {
	distance := currentNode.getTraveledDistance()

	beers := []string{}
	for _, i := range currentNode.includedBrews {
		beers = append(beers, i.beers...)
	}
	beers = getUniqueBeers(beers)
	currentNode.beerCnt = len(beers)

	if len(getCloseBreweries(breweries, currentNode.getLastBrew(), maxDistance-distance)) > 0 {

		if currentNode.profit == currentNode.bound {
			bestNodes = append(bestNodes, currentNode)
		}

		if currentNode.beerCnt >= bestBeerCount {
			bestNodes = append(bestNodes, currentNode)
			bestBeerCount = currentNode.beerCnt
		}

		if currentNode.profit < upper {
			nodes = append(nodes, currentNode)
		}
	} else if currentNode.beerCnt >= bestBeerCount && currentNode.getTraveledDistance() <= 2000. {
		bestNodes = append(bestNodes, currentNode)
		bestBeerCount = currentNode.beerCnt
	}

	return nodes, bestNodes, bestBeerCount
}

func findPaths(breweries []brewery, maxDistance float64) ([]node, int) {
	var (
		includedBrews                []brewery
		bestNodes                    []node
		oldNode, leftNode, rightNode node
		includedBrewsIndex           []int
		upper, profit, bound         float64
		bestBeerCount                int
	)
	profit, bound = calcFirstProfitAndBound(breweries, maxDistance)
	upper = bound

	nodes := []node{
		{0, profit, bound, []brewery{}, []int{}, []int{}, 0}}

	for len(nodes) > 0 {

		//Take node with highest value and deletes it from slice
		if len(nodes) > 1 {
			var highestProfitNode node
			index := -1
			highestProfitNode.profit = 0
			for i := range nodes {
				if nodes[i].profit < highestProfitNode.profit {
					highestProfitNode = nodes[i]
					index = i
				}
			}
			oldNode = highestProfitNode
			nodes[index] = nodes[len(nodes)-1]
			nodes = nodes[:len(nodes)-1]
		} else if len(nodes) == 1 {
			oldNode = nodes[len(nodes)-1]
			nodes = nodes[:len(nodes)-1]
		}

		//Goes to a new level
		newLevel := oldNode.level + 1
		leftNode.level, rightNode.level = newLevel, newLevel
		includedNodeID := newLevel

		//Left node calc
		excludedBrews := oldNode.excludedBrews
		includedBrewsIndex = append(oldNode.includedBrewsIndex, includedNodeID)
		profit, bound, includedBrews = calcProfitAndBound(breweries, includedBrewsIndex, excludedBrews, maxDistance)
		leftNode = node{newLevel, profit, bound, includedBrews, excludedBrews, includedBrewsIndex, 0}

		if leftNode.bound < upper {
			upper = leftNode.bound
			updateNodes(nodes, upper)
		}

		//Right node calc
		excludedBrews = append(excludedBrews, includedNodeID)
		includedBrewsIndex = oldNode.includedBrewsIndex
		profit, bound, includedBrews = calcProfitAndBound(breweries, includedBrewsIndex, excludedBrews, maxDistance)
		rightNode = node{newLevel, profit, bound, includedBrews, excludedBrews, includedBrewsIndex, 0}

		if rightNode.bound < upper {
			upper = rightNode.bound
			updateNodes(nodes, upper)
		}

		//Adding nodes to future consideration
		//Left node
		nodes, bestNodes, bestBeerCount = addNode(breweries, nodes, bestNodes, leftNode, maxDistance, upper, bestBeerCount)

		//Right node
		if rightNode.getLastBrew().ID > 0 {
			nodes, bestNodes, bestBeerCount = addNode(breweries, nodes, bestNodes, rightNode, maxDistance, upper, bestBeerCount)
		} else if rightNode.excludedBrews[len(rightNode.excludedBrews)-1] <= len(breweries)-1 {
			if rightNode.profit < upper {
				nodes = append(nodes, rightNode)
			}
		}
	}

	return bestNodes, bestBeerCount
}
