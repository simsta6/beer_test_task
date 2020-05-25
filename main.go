package main

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	bestBeerCount int
	bestDistance  float64
	bestPath      []brewery
)

type brewery struct {
	ID             int
	name           string
	latitude       float64
	longitude      float64
	beers          []string
	distanceToHome float64
}

func main() {
	lat1 := 51.74250300
	lon1 := 19.43295600

	breweries := getBreweriesWithBeersFromDB()

	breweries = getBreweriesWithin1000(lat1, lon1, breweries)
	brew := brewery{0, "HOME", lat1, lon1, []string{}, 0.}
	breweries = append([]brewery{brew}, breweries...)

	graph := makeDistancesGraph(breweries)

	visited := make([]bool, len(breweries))

	visited[0] = true

	start := time.Now()

	tspRec(breweries, 0., 0, []brewery{}, graph, visited)

	elapsed := time.Since(start)

	for i := range bestPath {
		fmt.Println(bestPath[i])
	}

	fmt.Printf("%v\n", bestPath)

	fmt.Printf("%v\n", bestBeerCount)

	fmt.Printf("%v\n", bestDistance)

	fmt.Printf("Method took %s\n", elapsed)

	fmt.Println("Done")
}

func getBreweriesWithBeersFromDB() []brewery {
	breweries := make([]brewery, 0)
	conn, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/beer-database")
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

		bconn, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/beer-database")
		checkForError(err)

		queryForBeers := `SELECT beer.name FROM beer 
							RIGHT JOIN breweries ON breweries.id = beer.brewery_id 
							WHERE beer.name IS NOT NULL AND brewery_id = ` + strconv.Itoa(brew.ID)

		bstatement, err := bconn.Prepare(queryForBeers)
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

		bconn.Close()

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

func haversine(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	//Converting to radians
	lat1r := lat1 * math.Pi / 180
	lon1r := lon1 * math.Pi / 180
	lat2r := lat2 * math.Pi / 180
	lon2r := lon2 * math.Pi / 180

	dlat := lat2r - lat1r
	dlon := lon2r - lon1r

	//Haversine formule to calculate ditance between two points
	a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1r)*math.Cos(lat2r)*math.Pow(math.Sin(dlon/2), 2)
	c := 2 * math.Asin(math.Sqrt(a)) * 6371.0

	return c
}

func getBreweriesWithin1000(lat float64, lon float64, breweries []brewery) []brewery {
	breweries1000 := make([]brewery, 0)

	for i := range breweries {
		distance := haversine(lat, lon, breweries[i].latitude, breweries[i].longitude)
		breweries[i].distanceToHome = distance
		if distance <= 1000 {
			breweries1000 = append(breweries1000, breweries[i])
		}
	}

	return breweries1000
}

func makeDistancesGraph(breweries []brewery) [][]float64 {
	graph := make([][]float64, len(breweries))

	for i := range graph {
		graph[i] = make([]float64, len(breweries))
	}

	for i := range graph {
		for j := range graph {
			distance := haversine(breweries[i].latitude, breweries[i].longitude, breweries[j].latitude, breweries[j].longitude)
			graph[i][j] = distance
		}
	}
	return graph
}

func getBeerCnt(path []brewery) int {
	beer := []string{}
	for i := range path {
		for j := range path[i].beers {
			beer = append(beer, path[i].beers[j])
		}
	}

	return len(beer)
}

func getTotalDistance(path []brewery) float64 {
	var distance float64
	for i := 0; i < len(path)-1; i++ {
		distance += haversine(path[i].latitude, path[i].longitude, path[i+1].latitude, path[i+1].longitude)
	}
	return distance
}

func setBestPath(currPath []brewery) {
	currBeerCount := getBeerCnt(currPath)
	currDistance := getTotalDistance(currPath)

	if currBeerCount > bestBeerCount {
		bestBeerCount = currBeerCount
		bestDistance = currDistance
		bestPath = currPath
		fmt.Printf("First brewery of path: %v\n", currPath[1])
		fmt.Printf("Breweries count: %v\n", len(currPath)-2)
		fmt.Printf("Beer count: %v\n", currBeerCount)
		fmt.Printf("Distance: %v\n", currDistance)
		fmt.Println("")
	} else if currBeerCount == bestBeerCount && currDistance < bestDistance {
		bestBeerCount = currBeerCount
		bestDistance = currDistance
		bestPath = currPath
		fmt.Println("Distance was better")
		fmt.Printf("First brewery of path: %v\n", currPath[1])
		fmt.Printf("Breweries count: %v\n", len(currPath)-2)
		fmt.Printf("Beer count: %v\n", currBeerCount)
		fmt.Printf("Distance: %v\n", currDistance)
		fmt.Println("")
	}
}

func tspRec(breweries []brewery, distanceTraveled float64, currPos int, path []brewery, graph [][]float64, visited []bool) {
	maxDistance := 1000.
	if len(path) > 0 {
		newPath := path
		home := breweries[0]

		newPath = append([]brewery{home}, newPath...)
		homeFinal := home
		homeFinal.distanceToHome = haversine(homeFinal.latitude, homeFinal.longitude, path[len(path)-1].latitude, path[len(path)-1].longitude)

		newPath = append(newPath, homeFinal)

		setBestPath(newPath)
	}

	for i := range breweries {
		if !visited[i] {
			distance := distanceTraveled + graph[currPos][i]

			if distance+graph[i][0] > maxDistance {
				continue
			}

			visited[i] = true
			path = append(path, breweries[i])

			tspRec(breweries, distance, i, path, graph, visited)

			path = path[:len(path)-1]
			visited[i] = false
		}
	}
}
