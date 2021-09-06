package mapbox

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildURL(t *testing.T) {
	expected := "https://api.mapbox.com/geocoding/v5/mapbox.places/0.449912,51.569032.json?types=postcode&limit=1&access_token=pk.eyabcd"
	require.Equal(t, expected, buildRequestURL(0.449912, 51.569032, "pk.eyabcd"))
}

func TestGetText(t *testing.T) {
	tests := []struct {
		name string
		in   SimplifiedResp
		out  string
	}{
		{
			"Valid",
			SimplifiedResp{Features: []SimplifiedFeature{{Text: "HANOI"}}},
			"HANOI",
		},
		{
			"Empty text",
			SimplifiedResp{Features: []SimplifiedFeature{{Text: ""}}},
			"",
		},
		{
			"Empty features",
			SimplifiedResp{Features: []SimplifiedFeature{}},
			"",
		},
		{
			"Nil features",
			SimplifiedResp{Features: nil},
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.in.GetText()
			if test.out == "" {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
			require.Equal(t, test.out, actual)
		})
	}
}

func TestGetJsonResp(t *testing.T) {
	tests := []struct {
		name   string
		sample string
		out    *SimplifiedResp
	}{
		{
			"Valid",
			"{\"type\":\"FeatureCollection\",\"query\":[0.449912,51.569032],\"features\":[{\"id\":\"postcode.9916238034211520\",\"type\":\"Feature\",\"place_type\":[\"postcode\"],\"relevance\":1,\"properties\":{},\"text\":\"SS15 5AS\",\"place_name\":\"SS15 5AS, Basildon, Essex, England, United Kingdom\",\"bbox\":[0.449224,51.568259,0.449937,51.569776],\"center\":[0.4496737,51.569431],\"geometry\":{\"type\":\"Point\",\"coordinates\":[0.4496737,51.569431]},\"context\":[{\"id\":\"place.13976826799562680\",\"wikidata\":\"Q216649\",\"text\":\"Basildon\"},{\"id\":\"district.18801805364378390\",\"wikidata\":\"Q23240\",\"text\":\"Essex\"},{\"id\":\"region.13483278848453920\",\"wikidata\":\"Q21\",\"short_code\":\"GB-ENG\",\"text\":\"England\"},{\"id\":\"country.12405201072814600\",\"wikidata\":\"Q145\",\"short_code\":\"gb\",\"text\":\"United Kingdom\"}]}],\"attribution\":\"NOTICE: © 2021 Mapbox and its suppliers. All rights reserved. Use of this data is subject to the Mapbox Terms of Service (https://www.mapbox.com/about/maps/). This response and the information it contains may not be retained. POI(s) provided by Foursquare.\"}",
			&SimplifiedResp{Features: []SimplifiedFeature{{Text: "SS15 5AS"}}},
		},
		{
			"No features list",
			"{\"type\":\"FeatureCollection\",\"query\":[0.449912,51.569032]}",
			&SimplifiedResp{Features: nil},
		},
		{
			"No text field",
			"{\"type\":\"FeatureCollection\",\"query\":[0.449912,51.569032],\"features\":[{\"id\":\"postcode.9916238034211520\",\"type\":\"Feature\",\"place_type\":[\"postcode\"],\"relevance\":1,\"properties\":{},\"place_name\":\"SS15 5AS, Basildon, Essex, England, United Kingdom\",\"bbox\":[0.449224,51.568259,0.449937,51.569776],\"center\":[0.4496737,51.569431],\"geometry\":{\"type\":\"Point\",\"coordinates\":[0.4496737,51.569431]},\"context\":[{\"id\":\"place.13976826799562680\",\"wikidata\":\"Q216649\",\"text\":\"Basildon\"},{\"id\":\"district.18801805364378390\",\"wikidata\":\"Q23240\",\"text\":\"Essex\"},{\"id\":\"region.13483278848453920\",\"wikidata\":\"Q21\",\"short_code\":\"GB-ENG\",\"text\":\"England\"},{\"id\":\"country.12405201072814600\",\"wikidata\":\"Q145\",\"short_code\":\"gb\",\"text\":\"United Kingdom\"}]}],\"attribution\":\"NOTICE: © 2021 Mapbox and its suppliers. All rights reserved. Use of this data is subject to the Mapbox Terms of Service (https://www.mapbox.com/about/maps/). This response and the information it contains may not be retained. POI(s) provided by Foursquare.\"}",
			&SimplifiedResp{Features: []SimplifiedFeature{{Text: ""}}},
		},
		{
			"Unrelevant JSON",
			"{\"names\": [\"laputa\", \"totoro\"]}",
			&SimplifiedResp{},
		},
		{
			"Empty JSON",
			"{}",
			&SimplifiedResp{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := getJSONResp(strings.NewReader(test.sample))
			t.Logf("%+v", actual)
			if test.out == nil {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
			require.Equal(t, test.out, actual)
		})
	}
}
