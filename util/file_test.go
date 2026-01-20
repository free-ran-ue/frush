package util_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/free-ran-ue/frush/util"
)

var testFileExistCases = []struct {
	name string
	path string
	want bool
}{
	{
		name: "file exists",
		path: "./testdata/file_exists.txt",
		want: true,
	},
	{
		name: "file not exists",
		path: "./testdata/file_not_exists.txt",
		want: false,
	},
}

func TestFileExists(t *testing.T) {
	testdataDir := "./testdata"
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatalf("failed to create testdata directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(testdataDir); err != nil {
			t.Fatalf("failed to remove testdata directory: %v", err)
		}
	}()

	testFile := filepath.Join(testdataDir, "file_exists.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	for _, testCase := range testFileExistCases {
		t.Run(testCase.name, func(t *testing.T) {
			if util.FileExists(testCase.path) != testCase.want {
				t.Errorf("FileExists(%q) = %v, want %v", testCase.name, testCase.path, testCase.want)
			}
		})
	}
}

var testFileReadCases = []struct {
	name string
	path string
	want []byte
}{
	{
		name: "file read",
		path: "./testdata/file_read.txt",
		want: []byte("test content"),
	},
}

func TestFileRead(t *testing.T) {
	testdataDir := "./testdata"
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatalf("failed to create testdata directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(testdataDir); err != nil {
			t.Fatalf("failed to remove testdata directory: %v", err)
		}
	}()

	testFile := filepath.Join(testdataDir, "file_read.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	for _, testCase := range testFileReadCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := util.FileRead(testCase.path)
			if err != nil {
				t.Errorf("FileRead(%q) = %v, want %v", testCase.name, err, testCase.want)
			}
			if !bytes.Equal(got, testCase.want) {
				t.Errorf("FileRead(%q) = %v, want %v", testCase.name, got, testCase.want)
			}
		})
	}
}
