package main

import (
	"testing"
)

func TestGetFirstBrew(t *testing.T) {

	//Normal case
	currentNode := node{1, -7.1, -7,
		[]brewery{
			{957, "Ostankinskij Pivovarennij Zavod", 55.75579833984375, 37.61759948730469, []string{"Beer"}, 212.37640863461093},
			{961, "Oy Sinebrychoff AB", 60.38100051879883, 25.110200881958008, []string{"Porter IV", "Koff Special III"}, 682.979146237569}}, []int{}, []int{}, 0}

	expected := brewery{957, "Ostankinskij Pivovarennij Zavod", 55.75579833984375, 37.61759948730469, []string{"Beer"}, 212.37640863461093}
	got := currentNode.getFirstBrew()

	if got.ID != expected.ID {
		t.Errorf("Normal case getFirstBrew returned brewery value is different. \nExpected: %v\nGot: %v", expected, got)
	}

	//Breweries slice is empty
	currentNode = node{1, -7.1, -7, []brewery{}, []int{}, []int{}, 0}

	expected = brewery{}
	got = currentNode.getFirstBrew()

	if got.ID != expected.ID {
		t.Errorf("Breweries slice is empty case getFirstBrew returned brewery value is different. \nExpected: %v\nGot: %v", expected, got)
	}
}

func TestGetLastBrew(t *testing.T) {

	//Normal case
	currentNode := node{1, -7.1, -7,
		[]brewery{
			{957, "Ostankinskij Pivovarennij Zavod", 55.75579833984375, 37.61759948730469, []string{"Beer"}, 212.37640863461093},
			{961, "Oy Sinebrychoff AB", 60.38100051879883, 25.110200881958008, []string{"Porter IV", "Koff Special III"}, 682.979146237569}}, []int{}, []int{}, 0}

	expected := brewery{961, "Oy Sinebrychoff AB", 60.38100051879883, 25.110200881958008, []string{"Porter IV", "Koff Special III"}, 682.979146237569}
	got := currentNode.getLastBrew()

	if got.ID != expected.ID {
		t.Errorf("Normal case getLastBrew returned brewery value is different. \nExpected: %v\nGot: %v", expected, got)
	}

	//Breweries slice is empty
	currentNode = node{1, -7.1, -7, []brewery{}, []int{}, []int{}, 0}

	expected = brewery{}
	got = currentNode.getLastBrew()

	if got.ID != expected.ID {
		t.Errorf("Breweries slice is empty case getLastBrew returned brewery value is different. \nExpected: %v\nGot: %v", expected, got)
	}
}

func TestGetTraveledDistance(t *testing.T) {

	//Normal case
	currentNode := node{1, -7.1, -7,
		[]brewery{
			{957, "Ostankinskij Pivovarennij Zavod", 55.75579833984375, 37.61759948730469, []string{"Beer"}, 212.37640863461093},
			{961, "Oy Sinebrychoff AB", 60.38100051879883, 25.110200881958008, []string{"Porter IV", "Koff Special III"}, 682.979146237569}}, []int{}, []int{}, 0}

	expected := 1790.607461383326
	got := currentNode.getTraveledDistance()

	if got != expected {
		t.Errorf("Normal case getTraveledDistance returned brewery value is different. \nExpected: %v\nGot: %v", expected, got)
	}

	//Breweries slice is empty
	currentNode = node{1, -7.1, -7, []brewery{}, []int{}, []int{}, 0}

	expected = 0.
	got = currentNode.getTraveledDistance()

	if got != expected {
		t.Errorf("Breweries slice is empty case getTraveledDistance returned brewery value is different. \nExpected: %v\nGot: %v", expected, got)
	}
}
