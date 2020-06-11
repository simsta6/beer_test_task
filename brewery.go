package main

import (
	"fmt"
)

type brewery struct {
	ID             int
	name           string
	latitude       float64
	longitude      float64
	beers          []string
	distanceToHome float64
}

func (brew *brewery) print() {
	fmt.Printf("\t[%v] %v: %.8f, %.8f ", brew.ID, brew.name, brew.latitude, brew.longitude)
}
