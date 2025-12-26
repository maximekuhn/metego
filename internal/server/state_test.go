package server

import "testing"

func TestNextCity(t *testing.T) {
	tests := []struct {
		title       string
		currentCity string
		cities      []string
		expected    string
	}{
		{
			title:       "No cities provided",
			currentCity: "Paris",
			cities:      []string{},
			expected:    "",
		},
		{
			title:       "Only one city provided",
			currentCity: "Paris",
			cities:      []string{"Paris"},
			expected:    "Paris",
		},
		{
			title:       "Only one city provided but not the current one",
			currentCity: "Berlin",
			cities:      []string{"Paris"},
			expected:    "Paris",
		},
		{
			title:       "Multiple cities, first one is selected",
			currentCity: "Paris",
			cities:      []string{"Paris", "Berlin", "Rome"},
			expected:    "Berlin",
		},
		{
			title:       "Multiple cities, second one is selected",
			currentCity: "Berlin",
			cities:      []string{"Paris", "Berlin", "Rome"},
			expected:    "Rome",
		},
		{
			title:       "Multiple cities, third one is selected",
			currentCity: "Rome",
			cities:      []string{"Paris", "Berlin", "Rome"},
			expected:    "Paris",
		},
		{
			title:       "Multiple cities, none selected",
			currentCity: "New York",
			cities:      []string{"Paris", "Berlin", "Rome"},
			expected:    "Paris",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			state := NewState(
				nil,
				nil,
				nil,
				nil,
				test.cities,
			)
			nextCity := state.nextCity(test.currentCity)
			if nextCity != test.expected {
				t.Errorf("nextCity(): want %s got %s", test.expected, nextCity)
			}
		})
	}

}
