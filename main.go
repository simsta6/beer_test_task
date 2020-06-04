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
	checkArguments()

	lat, err := strconv.ParseFloat(os.Args[1], 64)
	checkForError(err)
	lon, err := strconv.ParseFloat(os.Args[2], 64)
	checkForError(err)

	home := brewery{0, "HOME", lat, lon, []string{}, 0.}
	breweries := getBreweriesWithBeersFromDB(home)

	start := time.Now()

	breweries = getBreweriesWithin1000(breweries)

	if len(breweries) == 0 {
		fmt.Printf("There is no solution\n")
		os.Exit(1)
	}

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

	printResults(beers, elapsed)
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

func printResults(beers []string, elapsed time.Duration) {
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

	for i := range beers {
		fmt.Printf("\t%v\n", beers[i])
	}

	fmt.Printf("\nProgram took: %s\n", elapsed)
}

func getBreweriesWithBeersFromDB(home brewery) []brewery {
	breweries := make([]brewery, 0)

	var (
		conn *sql.DB
		err error
		retry int
	)

	for retry = 0; retry <100; retry++ {
		conn, err = sql.Open("mysql", "root:@tcp(golang_db)/beer-database")
		if err == nil {
			break // success!
		}
		<- time.After(1 * time.Second)
	}
	if retry == 100 {
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

func getBreweriesWithin1000(breweries []brewery) []brewery {
	breweries1000 := make([]brewery, 0)

	for i := range breweries {
		if breweries[i].distanceToHome <= 1000 {
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
	maxDistance := 2000.

	for i := 0; i < len(breweries); i++ {
		if !visited[i] {
			distance := distanceTraveled + graph[currPos][i]
			cBeerCnt += len(breweries[i].beers)

			if distance+graph[i][0] < maxDistance && (cBeerCnt > bestBeerCount || (cBeerCnt == bestBeerCount && distance+graph[i][0] < bestDistance)) {
				visited[i] = true
				path = append(path, breweries[i])

				newPath := path
				home := breweries[0]

				newPath = append([]brewery{home}, newPath...)
				homeFinal := home
				homeFinal.distanceToHome = haversine(homeFinal, path[len(path)-1])

				newPath = append(newPath, homeFinal)

				bestBeerCount = cBeerCnt
				bestDistance = distance + graph[i][0]
				bestPath = newPath

				tspRec2(breweries, distance, i, path, graph, visited, cBeerCnt, distance+graph[i][0])

				path = path[:len(path)-1]
				visited[i] = false
			} else {
				continue
			}
		}
	}
}
