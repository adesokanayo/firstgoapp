package main

import "testing"

func TestOne(t *testing.T) {
	one := false
	if !one {
		t.Errorf("test failed")
	}
}
