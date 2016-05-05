package fileprocessor

import (
	"io/ioutil"
	"os"
	"testing"
)

type TestFile struct {
	name string
	data []byte
}

type TestZipFile struct {
	file           TestFile
	targetFileName string
}

var testWritingFiles = []TestFile{
	{"test.txt", []byte("this is text that should end up in file")},
	{"test", []byte("")},
	{"TesT.out", []byte("")},
	{"weird.out", []byte("asdfasdf^^&&$$!!!!   --- @@ &nsp; ")},
}

var testZipFiles = []TestZipFile{
	{testWritingFiles[0], "output.zip"},
	{testWritingFiles[0], "output"},
	{testWritingFiles[0], "output.zip"},
}

var invalidTestZipFiles = []TestZipFile{
	{testWritingFiles[0], testWritingFiles[0].name},
}

func TestWriteFile(t *testing.T) {
	for _, v := range testWritingFiles {
		err := writeFile(v.name, v.data)
		if err != nil {
			t.Fatal("an error occurred while writing file", err)
		}

		dat, err := ioutil.ReadFile(v.name)
		if err != nil {
			t.Fatal("Unable to read file that was written", err)
		}

		if string(dat) != string(v.data) {
			t.Fatal("data was not correctly written to file!")
		}

		deleteFile(v.name)
	}
}

func TestCreateZipFile(t *testing.T) {
	for _, v := range testZipFiles {
		writeFile(v.file.name, v.file.data)
		err := createZipFile(v.file.name, v.targetFileName)
		if err != nil {
			t.Fatalf("Failed to zip file: %v \n %v", v.targetFileName, err)
		}

		s, err := os.Stat(v.targetFileName)
		if err != nil {
			t.Fatal("Something wrong with file!")
		}
		if s.Size() == 0 {
			t.Fatal("data not written to file")
		}

		deleteFile(v.file.name)
		deleteFile(v.targetFileName)
	}
}

func TestCreateZipFile_invalid(t *testing.T) {
	for _, v := range invalidTestZipFiles {
		writeFile(v.file.name, v.file.data)
		err := createZipFile(v.file.name, v.targetFileName)
		if err == nil {
			t.Fatalf("An error should have been thrown for: %v", v.targetFileName)
		}
		deleteFile(v.file.name)
	}
}
