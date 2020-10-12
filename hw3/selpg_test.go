package main

import "testing"

func TestUsuage(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Usage()
		})
	}
}

func TestProcessArgs(t *testing.T) {

	tests := []struct {
		name string
		args selpgArgs
	}{
		// TODO: Add test cases.
		{
			name: "Nil options",
			args: selpgArgs{2, 2, "test.txt", 7, false, ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ProcessArgs(&(tt.args))
		})
	}
}

func TestProcessInput(t *testing.T) {
	tests := []struct {
		name string
		args selpgArgs
	}{
		// TODO: Add test cases.
		{
			name: "Nil options",
			args: selpgArgs{2, 2, "test.txt", 7, false, ""},
		},
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ProcessInput(&(tt.args))
		})
	}
}
