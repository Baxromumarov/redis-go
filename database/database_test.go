package database_test

import (
	"testing"
	"time"

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
		{"SET key", "SET command requires at least two arguments: key and value"},
		{"GET", "GET command requires 1 argument"},
		{"DEL", "DEL command requires 1 argument"},
		{"", "empty command"},
		{"UNKNOWN", "unknown command"},
		{"SET key value 4", ""},
		{"GET key", "value"},
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

	time.Sleep(4 * time.Second)

	_, err := database.HandleCommand(db, "GET key")
	if err == nil || err.Error() != "key has expired" {
		t.Errorf("Expected key to expire after TTL, but it was found")
	}

	val, err := database.HandleCommand(db, "GET key")
	if err != nil && err.Error() != "key has expired" {
		t.Errorf("Expected error for expired key, got %v", err)
	} else if val != "" {
		t.Errorf("Expected empty value for expired key, got %v", val)
	}
	
	db.Set("key", "value", 10)
	_, err = database.HandleCommand(db, "TTL key")
	if err != nil {
		t.Errorf("Expected TTL value, but got error: %v", err)
	}
}
