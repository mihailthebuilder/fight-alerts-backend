package main

import (
	"testing"
)

func Test_buildConnectionString(t *testing.T) {
	got := buildConnectionString("localhost", "password", "username", "test", 5432)

	expected := "host=localhost port=5432 user=username password=password dbname=test sslmode=disable"
	if got != expected {
		t.Errorf("buildConnectionString() = %v, want = %v", got, expected)
	}
}

func Test_setUpDatastore(t *testing.T) {
	got := setUpDatastore()

	if got.Db == nil {
		t.Errorf("setUpDatastore() does not set sql.DB driver")
	}
}
