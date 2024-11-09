package main

import (
	"testing"
	"time"
)

func TestRetailerPoints(t *testing.T) {
	tests := []struct {
		name     string
		retailer string
		expected int
	}{
		{"ValidCharacters", "Wario's Beef and Pork", 17},
		{"EmptyRetailerName", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := calculateRetailerNamePoints(tt.retailer); result != tt.expected {
				t.Errorf("calculateRetailerNamePoints() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRoundTotalPoints(t *testing.T) {
	tests := []struct {
		name     string
		total    float64
		expected int
	}{
		{"RoundTotal", 100.00, 50},
		{"NonRoundTotal", 100.01, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := calculateRoundTotalPoints(tt.total); result != tt.expected {
				t.Errorf("calculateRoundTotalPoints() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestQuarterTotalPoints(t *testing.T) {
	tests := []struct {
		name     string
		total    float64
		expected int
	}{
		{"QuarterTotal", 100.00, 25},
		{"NonQuarterTotal", 100.01, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := calculateQuarterTotalPoints(tt.total); result != tt.expected {
				t.Errorf("calculateQuarterTotalPoints() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestItemPairPoints(t *testing.T) {
	tests := []struct {
		name     string
		items    []Item
		expected int
	}{
		{"TwoPairs",
			[]Item{{ShortDescription: "Chicken Cutlet", Price: 19.00},
				{ShortDescription: "Chicken Cutlet", Price: 19.00},
				{ShortDescription: "Chicken Cutlet", Price: 19.00},
				{ShortDescription: "Chicken Cutlet", Price: 19.00}},
			10,
		},
		{"OnePair",
			[]Item{{ShortDescription: "Chicken Cutlet", Price: 19.00},
				{ShortDescription: "Chicken Cutlet", Price: 19.00},
				{ShortDescription: "Chicken Cutlet", Price: 19.00}},
			5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := calculateItemPairPoints(tt.items); result != tt.expected {
				t.Errorf("calculateItemPairPoints() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestItemNameLengthPoints(t *testing.T) {
	tests := []struct {
		name     string
		items    []Item
		expected int
	}{
		{"ExpectedPoints", []Item{
			{ShortDescription: "Chicken Cutlets", Price: 19.00},
		},
			4},
		{"NoExpectedPoints", []Item{
			{ShortDescription: "Chicken Cutlet", Price: 19.00},
		},
			0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := calculateItemNameLengthPoints(tt.items); result != tt.expected {
				t.Errorf("calculateItemNameLengthPoints() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestOddPurchaseDate(t *testing.T) {
	tests := []struct {
		name     string
		date     PurchaseDate
		expected int
	}{
		{"OddPurchaseDate", PurchaseDate{
			Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)},
			6,
		},
		{"EvenPurchaseDate",
			PurchaseDate{Date: time.Date(2020, 1, 2, 0, 0, 0, 0, time.Local)},
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := calculateOddPurchaseDatePoints(tt.date); result != tt.expected {
				t.Errorf("calculateOddPurchaseDatePoints() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHappyHourPoints(t *testing.T) {
	tests := []struct {
		name     string
		time     PurchaseTime
		expected int
	}{
		{"HappyHour!",
			PurchaseTime{Time: time.Date(2020, 1, 1, 15, 0, 0, 0, time.Local)},
			10},
		{"UnhappyHour",
			PurchaseTime{Time: time.Date(2020, 1, 1, 6, 0, 0, 0, time.Local)},
			0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := calculateHappyHourPoints(tt.time); result != tt.expected {
				t.Errorf("calculateHappyHourPoints() = %v, want %v", result, tt.expected)
			}
		})
	}
}
