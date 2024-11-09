package database_test

import (
	"testing"

	"github.com/baxromumarov/redis-go/database"
)

func TestHandleCommand(t *testing.T) {
	var tests = []struct {
		Input  string
		Output string
	}{
		{"SET key value", ""},
		{"GET key", "value"},
		{"DEL key", ""},
		{"SET key", "SET command requires 2 arguments"},
		{"GET", "GET command requires 1 argument"},
		{"DEL", "DEL command requires 1 argument"},
		{"", "empty command"},
		{"UNKNOWN", "unknown command"},
	}

	db := database.Init()

	for _, test := range tests {
		val, err := database.HandleCommand(db, test.Input)
		if err != nil {
			if err.Error() != test.Output {
				t.Errorf("Expected %v, got %v", test.Output, err.Error())
			}
		} else {
			if val != test.Output {
				t.Errorf("Expected %v, got %v", test.Output, val)
			}
		}
	}

}
