package main

type node struct {
	level                             int
	profit, bound                     float64
	includedBrews                     []brewery
	excludedBrews, includedBrewsIndex []int
	beerCnt                           int
}

func (n *node) getLastBrew() brewery {
	var brew brewery
	if len(n.includedBrews) > 0 {
		brew = n.includedBrews[len(n.includedBrews)-1]
	}
	return brew
}

func (n *node) getFirstBrew() brewery {
	var brew brewery
	if len(n.includedBrews) > 0 {
		brew = n.includedBrews[0]
	}
	return brew
}

func (n *node) getTraveledDistance() float64 {
	distance := 0.
	if len(n.includedBrews) == 1 {
		distance = n.includedBrews[0].distanceToHome
	} else if len(n.includedBrews) == 0 {
		distance = 0
	} else {
		distance += n.includedBrews[0].distanceToHome
		for i := 1; i < len(n.includedBrews); i++ {
			distance += haversine(n.includedBrews[i], n.includedBrews[i-1])
		}
		distance += n.includedBrews[len(n.includedBrews)-1].distanceToHome
	}
	return distance
}
