package main

import "testing"

func TestHaversine(t *testing.T) {

	//Same coordinates
	firstPoint := brewery{0, "firstPoint", 51.355468, 11.100790, []string{}, 0.}
	secondPoint := brewery{0, "secondPoint", 51.355468, 11.100790, []string{}, 0.}
	got := haversine(firstPoint, secondPoint)

	if got != 0 {
		t.Errorf("haversine got %v", got)
	}

	//Should return around 169.5 km
	firstPoint = brewery{0, "firstPoint", 51.355468, 11.100790, []string{}, 0.}
	secondPoint = brewery{0, "secondPoint", 50, 10, []string{}, 0.}
	got = haversine(firstPoint, secondPoint)

	if got < 169.4 || got > 169.6 {
		t.Errorf("haversine got %v", got)
	}

	//Should return around 20015 km
	firstPoint = brewery{0, "firstPoint", 90, 180, []string{}, 0.}
	secondPoint = brewery{0, "secondPoint", -90, -180, []string{}, 0.}
	got = haversine(firstPoint, secondPoint)

	if got < 20014 || got > 20016 {
		t.Errorf("haversine got %v", got)
	}
}

func TestGetBreweriesWithin1000(t *testing.T) {
	//Should return 6 breweries
	home := brewery{0, "home", 51.355468, 11.100790, []string{}, 0.}
	breweries := []brewery{
		brewery{1267, "Tin Whistle Brewing", 49.49440002441406, -119.61000061035156, []string{}, 7869.721023294844},
		brewery{1268, "Titletown Brewing", 44.52000045776367, -88.01730346679688, []string{}, 6839.141532625284},
		brewery{1269, "Tivoli Brewing", 39.739200592041016, -104.98500061035156, []string{}, 8145.293542222566},
		brewery{1271, "Tommyknocker Brewery and Pub", 39.741798400878906, -105.51799774169922, []string{}, 8171.6742897312},
		brewery{1272, "Tomos Watkin and Sons Ltd.", 51.66659927368164, -3.9442999362945557, []string{}, 1039.9077572663819},
		brewery{1274, "Tooheys", -33.850101470947266, 151.0449981689453, []string{}, 16268.64846307816},

		brewery{1056, "Ridgeway Brewing", 51.546199798583984, -1.1354999542236328, []string{}, 847.1914449316964},
		brewery{1083, "Ruppaner-Brauerei", 47.68550109863281, 9.208000183105469, []string{}, 430.3081461406065},
		brewery{1088, "SA Brain & Co. Ltd.", 51.47359848022461, -3.178999900817871, []string{}, 988.8234297974843},
		brewery{1096, "Salopian Brewery", 52.70690155029297, -2.7869999408721924, []string{}, 960.325595317151},
		brewery{1099, "Samuel Smith Old Brewery (Tadcaster)", 53.883399963378906, -1.2625000476837158, []string{}, 879.3105926211241},
		brewery{1111, "Sarah Hughes Brewery", 52.54349899291992, -2.115600109100342, []string{}, 914.0475492655376}}

	got := getBreweriesWithin1000(home, breweries)

	if len(got) != 6 {
		t.Errorf("haversine got %v", got)
	}

	//Should return 0 breweries
	home = brewery{0, "home", 0, 0, []string{}, 0.}
	breweries = []brewery{
		brewery{1267, "Tin Whistle Brewing", 49.49440002441406, -119.61000061035156, []string{}, 7869.721023294844},
		brewery{1268, "Titletown Brewing", 44.52000045776367, -88.01730346679688, []string{}, 6839.141532625284},
		brewery{1269, "Tivoli Brewing", 39.739200592041016, -104.98500061035156, []string{}, 8145.293542222566},
		brewery{1271, "Tommyknocker Brewery and Pub", 39.741798400878906, -105.51799774169922, []string{}, 8171.6742897312},
		brewery{1272, "Tomos Watkin and Sons Ltd.", 51.66659927368164, -3.9442999362945557, []string{}, 1039.9077572663819},
		brewery{1274, "Tooheys", -33.850101470947266, 151.0449981689453, []string{}, 16268.64846307816},

		brewery{1056, "Ridgeway Brewing", 51.546199798583984, -1.1354999542236328, []string{}, 847.1914449316964},
		brewery{1083, "Ruppaner-Brauerei", 47.68550109863281, 9.208000183105469, []string{}, 430.3081461406065},
		brewery{1088, "SA Brain & Co. Ltd.", 51.47359848022461, -3.178999900817871, []string{}, 988.8234297974843},
		brewery{1096, "Salopian Brewery", 52.70690155029297, -2.7869999408721924, []string{}, 960.325595317151},
		brewery{1099, "Samuel Smith Old Brewery (Tadcaster)", 53.883399963378906, -1.2625000476837158, []string{}, 879.3105926211241},
		brewery{1111, "Sarah Hughes Brewery", 52.54349899291992, -2.115600109100342, []string{}, 914.0475492655376}}

	got = getBreweriesWithin1000(home, breweries)

	if len(got) != 0 {
		t.Errorf("haversine got %v", got)
	}

	//Should return 0 breweries
	home = brewery{0, "home", 51.355468, 11.100790, []string{}, 0.}
	breweries = []brewery{}

	got = getBreweriesWithin1000(home, breweries)

	if len(got) != 0 {
		t.Errorf("haversine got %v", got)
	}
}
