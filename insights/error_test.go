package insights

import "testing"

/*
	This file contains the tests for the error utility
*/

func TestError_Error(t *testing.T) {
	e := &Error{"TestError", 0}
	if e.Error() != "C-0 TestError" {
		t.Fatal("Expected C-0 TestError. Got", e.Error())
	}
}

func TestError_String(t *testing.T) {
	e := &Error{"TestError", 0}
	if e.String() != "C-0 TestError" {
		t.Fatal("Expected C-0 TestError. Got", e.String())
	}
}
