package model

import (
	"testing"
)

func TestTestSuite_UnmarshalFromFile(t *testing.T) {
	testsuite := TestSuite{}

	err := testsuite.UnmarshalFromFile("test.xml")
	if err != nil {
		t.Error(err)
	}
}
