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

type brewery struct {
	ID             int
	name           string
	latitude       float64
	longitude      float64
	beers          []string
	distanceToHome float64
}

var (
	bestBeerCount int
	bestDistance  float64
	bestPath      []brewery
)

func main() {
	lat1 := 51.355468
	lon1 := 11.100790

	start := time.Now()

	breweries := getBreweriesWithBeersFromDB()

	home := brewery{0, "HOME", lat1, lon1, []string{}, 0.}
	breweries = getBreweriesWithin1000(home, breweries)
	breweries = append([]brewery{home}, breweries...)

	graph := makeDistancesGraph(breweries)

	visited := make([]bool, len(breweries))

	visited[0] = true

	tspRec2(breweries, 0., 0, []brewery{}, graph, visited, 0, 0.)

	elapsed := time.Since(start)

	beers := []string{}

	for i := range bestPath {
		beers = append(beers, bestPath[i].beers...)
	}

	beers = getUniqueBeers(beers)

	printResults(len(beers), elapsed)

}

func printResults(beersCnt int, elapsed time.Duration) {
	fmt.Printf("Breweries: \n")

	for i := range bestPath {
		fmt.Printf("%v\n", bestPath[i])
	}

	fmt.Printf("Beers: %v\n", beersCnt)

	fmt.Printf("Distance: %v\n", bestDistance)

	fmt.Printf("Program took: %s\n", elapsed)

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

func getBreweriesWithin1000(home brewery, breweries []brewery) []brewery {
	breweries1000 := make([]brewery, 0)

	for i := range breweries {
		distance := haversine(home, breweries[i])
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

// func tspRec(breweries []brewery, distanceTraveled float64, currPos int, path []brewery, graph [][]float64, visited []bool, cBeerCnt int, cDistance float64) {
// 	maxDistance := 2000.
// 	if cBeerCnt > bestBeerCount || (cBeerCnt == bestBeerCount && cDistance+graph[currPos][0] < bestDistance) {
// 		newPath := path
// 		home := breweries[0]
//
// 		newPath = append([]brewery{home}, newPath...)
// 		homeFinal := home
// 		homeFinal.distanceToHome = haversine(homeFinal, path[len(path)-1])
//
// 		newPath = append(newPath, homeFinal)
//
// 		bestBeerCount = cBeerCnt
// 		bestDistance = cDistance + graph[currPos][0]
// 		bestPath = newPath
// 	}
//
// 	for i := 0; i < len(breweries); i++ {
// 		if !visited[i] {
// 			distance := distanceTraveled + graph[currPos][i]
//
// 			if distance+graph[i][0] > maxDistance {
// 				continue
// 			}
//
// 			visited[i] = true
// 			path = append(path, breweries[i])
// 			cBeerCnt += len(breweries[i].beers)
//
// 			tspRec(breweries, distance, i, path, graph, visited, cBeerCnt, distance)
//
// 			cBeerCnt -= len(breweries[i].beers)
// 			path = path[:len(path)-1]
// 			visited[i] = false
// 		}
// 	}
// }

func tspRec2(breweries []brewery, distanceTraveled float64, currPos int, path []brewery, graph [][]float64, visited []bool, cBeerCnt int, cDistance float64) {
	if len(breweries) == 0 {
		fmt.Printf("There is no solution\n")
	}
	maxDistance := 2000.

	for i := 0; i < len(breweries); i++ {
		if !visited[i] {
			distance := distanceTraveled + graph[currPos][i]

			if distance+graph[i][0] > maxDistance {
				continue
			}

			visited[i] = true
			path = append(path, breweries[i])
			cBeerCnt += len(breweries[i].beers)

			if cBeerCnt > bestBeerCount || (cBeerCnt == bestBeerCount && distance+graph[i][0] < bestDistance) {
				newPath := path
				home := breweries[0]

				newPath = append([]brewery{home}, newPath...)
				homeFinal := home
				homeFinal.distanceToHome = haversine(homeFinal, path[len(path)-1])

				newPath = append(newPath, homeFinal)

				bestBeerCount = cBeerCnt
				bestDistance = distance + graph[i][0]
				bestPath = newPath

				tspRec2(breweries, distance, i, path, graph, visited, cBeerCnt, distance)

				cBeerCnt -= len(breweries[i].beers)
				path = path[:len(path)-1]
				visited[i] = false
			} else {
				cBeerCnt -= len(breweries[i].beers)
				path = path[:len(path)-1]
				visited[i] = false
			}
		}
	}
}
