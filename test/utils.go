package test

import "testing"

// AssertEqualInt checks if two integers of any type are equal and throws an error if they are not.
func AssertEqualInt[T ~int | ~int8](t *testing.T, expected, actual T, message string) {
	if actual != expected {
		t.Errorf("actual %d != expected %d: %s", actual, expected, message)
	}
}
