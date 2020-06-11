package brew

//Brewery struct
type Brewery struct {
	ID             int
	name           string
	latitude       float64
	longitude      float64
	beers          []string
	distanceToHome float64
}
