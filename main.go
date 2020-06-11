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

// type brewery struct {
// 	ID             int
// 	name           string
// 	latitude       float64
// 	longitude      float64
// 	beers          []string
// 	distanceToHome float64
// }

type node struct {
	level                             int
	cost, bound                       float64
	includedBrews                     []brew.brewery
	excludedBrews, includedBrewsIndex []int
	beerCnt                           int
}

func (n *node) getLastBrew() brew.brewery {
	var brew brew.brewery
	if len(n.includedBrews) > 0 {
		brew = n.includedBrews[len(n.includedBrews)-1]
	}
	return brew
}

func (n *node) getFirstBrew() brew.brewery {
	var brew brew.brewery
	if len(n.includedBrews) > 0 {
		brew = n.includedBrews[0]
	}
	return brew
}

func (n *node) getTraveledDistance() float64 {
	distance := 0.
	if len(n.includedBrews) == 1 {
		distance = n.includedBrews[0].distanceToHome
	} else {
		distance += n.includedBrews[0].distanceToHome
		for i := 1; i < len(n.includedBrews); i++ {
			distance += haversine(n.includedBrews[i], n.includedBrews[i-1])
		}
		distance += n.includedBrews[len(n.includedBrews)-1].distanceToHome
	}
	return distance
}

var (
	bestBeerCount int
	bestDistance  float64
	bestPath      []brew.brewery
)

func main() {
	maxDistance := 2000.
	//checkArguments()

	// lat, err := strconv.ParseFloat(os.Args[1], 64)
	// checkForError(err)
	// lon, err := strconv.ParseFloat(os.Args[2], 64)
	// checkForError(err)

	// lat := 57.
	// lon := 35.

	//lat := 51.74250300
	//lon := 19.43295600

	lat := 51.355468
	lon := 11.100790

	home := brew.brewery{0, "HOME", lat, lon, []string{}, 0.}
	breweries := getBreweriesWithBeersFromDB(home)

	start := time.Now()
	breweries = getCloseBreweries(breweries, home, maxDistance/2)

	if len(breweries) == 0 {
		fmt.Printf("There is no solution\n")
		os.Exit(1)
	}

	breweries = append([]brew.brewery{home}, breweries...)

	sort.Slice(breweries[:], func(i, j int) bool {
		return len(breweries[i].beers) > len(breweries[j].beers)
	})

	nodes := findPaths(breweries, maxDistance)

	elapsed := time.Since(start)

	fmt.Printf("\n\n\n")

	var bestPath []brew.brewery
	bestDistance := 2000.

	for _, i := range nodes {
		distance := i.getTraveledDistance()
		if i.beerCnt == bestBeerCount {
			bestDistance = distance
			for _, j := range i.includedBrews {
				fmt.Printf("%v\n", j)
				bestPath = append(bestPath, j)
			}
			fmt.Printf("%v\n", i.beerCnt)
			fmt.Printf("\n\n\n")
		}
	}

	fmt.Printf("%v\n", elapsed)

	beers := []string{}
	for i := range bestPath {
		beers = append(beers, bestPath[i].beers...)
	}
	beers = getUniqueBeers(beers)

	bestPath = append([]brew.brewery{home}, bestPath...)
	bestPath = append(bestPath, home)

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

func printResults(beers []string, elapsed time.Duration, bestPath []brew.brewery, bestDistance float64) {
	fmt.Printf("Found %v beer factories: \n", len(bestPath)-2)

	for i := 0; i < len(bestPath); i++ {
		if i == 0 {
			fmt.Printf("\t[%v] %v: %.8f, %.8f distance %.0fkm\n", bestPath[i].ID, bestPath[i].name, bestPath[i].latitude, bestPath[i].longitude, 0.)
			continue
		}
		fmt.Printf("\t[%v] %v: %.8f, %.8f distance %.1fkm\n", bestPath[i].ID, bestPath[i].name, bestPath[i].latitude, bestPath[i].longitude, haversine(bestPath[i], bestPath[i-1]))
	}

	fmt.Printf("\nTotal distance traveled: %.1fkm\n", bestDistance)

	fmt.Printf("\n\nCollected %v beer types:\n", len(beers))

	// for i := range beers {
	// 	fmt.Printf("\t%v\n", beers[i])
	// }

	fmt.Printf("\nProgram took: %s\n", elapsed)
}

func getBreweriesWithBeersFromDB(home brew.brewery) []brew.brewery {
	breweries := make([]brew.brewery, 0)

	var (
		conn  *sql.DB
		err   error
		retry int
	)

	for retry = 0; retry < 10; retry++ {
		conn, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/beer-database")
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
		brew := brew.brewery{}

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

func haversine(firstBrewry brew.brewery, secondBrewery brew.brewery) float64 {
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

func getCloseBreweries(breweries []brew.brewery, currentB brew.brewery, distance float64) []brew.brewery {
	closeBreweries := make([]brew.brewery, 0)

	for i := range breweries {
		currDistance := haversine(currentB, breweries[i])
		if currDistance <= distance && currDistance != 0. {
			closeBreweries = append(closeBreweries, breweries[i])
		}
	}

	return closeBreweries
}

func breweriesCanGetHome(breweries []brew.brewery, currentB brew.brewery, distance float64) []brew.brewery {
	closeBreweries := make([]brew.brewery, 0)

	for i := range breweries {
		if haversine(currentB, breweries[i])+breweries[i].distanceToHome <= distance {
			closeBreweries = append(closeBreweries, breweries[i])
		}
	}

	return closeBreweries
}

func makeDistancesGraph(breweries []brew.brewery) [][]float64 {
	graph := make([][]float64, len(breweries))

	for i := range graph {
		graph[i] = make([]float64, len(breweries))
	}

	for i := range graph {
		for j := range graph[i] {
			distance := haversine(breweries[i], breweries[j])
			graph[i][j] = distance
		}
	}
	return graph
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

func countCostandBound(breweries []brew.brewery, maxDistance float64) (float64, float64) {
	var (
		cost, bound, tDistance float64
		lastBrew               brew.brewery
	)

	if len(breweries) > 0 {
		lastBrew = breweries[0]
	}

	for i := 1; i < len(breweries); i++ {
		distance := haversine(breweries[i], lastBrew)
		beerCnt := float64(len(breweries[i].beers))
		if tDistance+distance+breweries[i].distanceToHome <= maxDistance {
			cost += beerCnt
			bound += beerCnt
			tDistance += distance
			lastBrew = breweries[i]
		} else {
			cost += beerCnt / (distance + breweries[i].distanceToHome) * (maxDistance - tDistance)
			break
		}
	}

	if cost != 0 || bound != 0 {
		cost = -cost
		bound = -bound
	}

	return cost, bound
}

func countCostandBound2(breweries []brew.brewery, includedBrewsIndex, excludedBrews []int, maxDistance float64) (float64, float64, []brew.brewery) {
	var (
		cost, bound, tDistance float64
		lastBrew               brew.brewery
		includedBrews          []brew.brewery
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
			cost += beerCnt
			bound += beerCnt
			tDistance += distance
			lastBrew = breweries[i]

		} else {
			cost += beerCnt / (distance + breweries[i].distanceToHome) * (maxDistance - tDistance)
			break
		}
	}

	if cost != 0 || bound != 0 {
		cost = -cost
		bound = -bound
	}

	return cost, bound, includedBrews
}

func updateNodes(nodes []node, upper float64) []node {
	for i := 0; i < len(nodes); i++ {
		if nodes[i].cost < upper && len(nodes) > i {
			nodes[i] = nodes[len(nodes)-1]
			nodes = nodes[:len(nodes)-1]
		}
	}
	return nodes
}

func addNode(breweries []brew.brewery, nodes, bestNodes []node, currentNode node, maxDistance, upper float64) ([]node, []node) {
	distance := currentNode.getTraveledDistance()

	if len(getCloseBreweries(breweries, currentNode.getLastBrew(), maxDistance-distance)) > 0 {

		beers := []string{}
		for _, i := range currentNode.includedBrews {
			beers = append(beers, i.beers...)
		}
		beers = getUniqueBeers(beers)
		currentNode.beerCnt = len(beers)

		if currentNode.cost == currentNode.bound || currentNode.beerCnt >= bestBeerCount {
			bestNodes = append(bestNodes, currentNode)
			bestBeerCount = currentNode.beerCnt
		}

		if currentNode.cost < upper {
			nodes = append(nodes, currentNode)
		}
	}

	return nodes, bestNodes
}

func findPaths(breweries []brew.brewery, maxDistance float64) []node {
	var (
		oldNode, leftNode, rightNode node
		bestNodes                    []node
		includedBrewsIndex           []int
		upper, cost, bound           float64
		includedBrews                []brew.brewery
	)
	cost, bound = countCostandBound(breweries, maxDistance)
	upper = bound

	nodes := []node{
		{0, cost, bound, []brew.brewery{}, []int{}, []int{}, 0}}

	for len(nodes) > 0 {
		if len(nodes) > 1 {
			var highestCostNode node
			index := -1
			highestCostNode.cost = 0
			for i := range nodes {
				if nodes[i].cost < highestCostNode.cost {
					highestCostNode = nodes[i]
					index = i
				}
			}
			oldNode = highestCostNode
			nodes[index] = nodes[len(nodes)-1]
			nodes = nodes[:len(nodes)-1]
		} else if len(nodes) == 1 {
			oldNode = nodes[len(nodes)-1]
			nodes = nodes[:len(nodes)-1]
		} else {
			break
		}

		newLevel := oldNode.level + 1

		leftNode.level, rightNode.level = newLevel, newLevel

		//Left node calc

		includedNodeID := leftNode.level

		excludedBrews := oldNode.excludedBrews
		includedBrewsIndex = append(oldNode.includedBrewsIndex, includedNodeID)

		cost, bound, includedBrews = countCostandBound2(breweries, includedBrewsIndex, excludedBrews, maxDistance)

		leftNode.cost = cost
		leftNode.bound = bound
		leftNode.includedBrews = includedBrews
		leftNode.includedBrewsIndex = includedBrewsIndex
		leftNode.excludedBrews = excludedBrews

		if leftNode.bound < upper {
			upper = leftNode.bound
			updateNodes(nodes, upper)
		}

		//Right node calc

		excludedBrews = append(excludedBrews, includedNodeID)
		includedBrewsIndex = oldNode.includedBrewsIndex

		cost, bound, includedBrews = countCostandBound2(breweries, includedBrewsIndex, excludedBrews, maxDistance)

		rightNode.cost = cost
		rightNode.bound = bound
		rightNode.includedBrews = includedBrews
		rightNode.includedBrewsIndex = includedBrewsIndex
		rightNode.excludedBrews = excludedBrews

		if rightNode.bound < upper {
			upper = rightNode.bound
			updateNodes(nodes, upper)
		}

		//Adding nodes to future consideration

		//Left node
		nodes, bestNodes = addNode(breweries, nodes, bestNodes, leftNode, maxDistance, upper)

		//Right node
		if rightNode.getLastBrew().ID > 0 {
			nodes, bestNodes = addNode(breweries, nodes, bestNodes, rightNode, maxDistance, upper)
		} else if rightNode.excludedBrews[len(rightNode.excludedBrews)-1] <= len(breweries)-1 {
			if rightNode.cost < upper {
				nodes = append(nodes, rightNode)
			}
		}
	}

	return bestNodes
}
