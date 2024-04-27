package storage

import (
	"bytes"
	"encoding/csv"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestFileStorage_Load(t *testing.T) {
	testData := map[time.Time]int{
		time.Now().UTC().Truncate(time.Second):                        10,
		time.Now().UTC().Truncate(time.Second).Add(-10 * time.Second): 5,
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	for td, count := range testData {
		err := writer.Write([]string{strconv.FormatInt(td.Unix(), 10), strconv.Itoa(count)})
		if err != nil {
			t.Fatal(`Error writing test data to CSV:`, err)
		}
	}
	writer.Flush()

	tempFile, err := os.CreateTemp(``, `test-requests.csv`)
	if err != nil {
		t.Fatal(`Error creating temporary file:`, err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(buf.Bytes())
	if err != nil {
		t.Fatal(`Error writing test data to temporary file:`, err)
	}

	tempFile.Close()

	fs := &FileStorage{filename: tempFile.Name()}

	loadedData, err := fs.Load()
	if err != nil {
		t.Fatalf(`Error loading data from file: %v`, err)
	}

	for key, val := range testData {
		if loadedData[key] != val {
			t.Errorf(`Loaded data mismatch at %s: got %d, want %d`, key, loadedData[key], val)
		}
	}
}

func TestFileStorage_Save(t *testing.T) {
	testData := map[time.Time]int{
		time.Now().UTC().Truncate(time.Second):                        10,
		time.Now().UTC().Truncate(time.Second).Add(-10 * time.Second): 5,
	}

	tempFile, err := os.CreateTemp(``, `test-requests.csv`)
	if err != nil {
		t.Fatal(`Error creating temporary file:`, err)
	}
	defer os.Remove(tempFile.Name())

	fs := &FileStorage{filename: tempFile.Name()}

	err = fs.Save(testData)
	if err != nil {
		t.Fatalf(`Error saving data to file: %v`, err)
	}

	loadedData, err := fs.Load()
	if err != nil {
		t.Fatalf(`Error loading data from file: %v`, err)
	}

	for key, val := range testData {
		if loadedData[key] != val {
			t.Errorf(`Loaded data mismatch at %s: got %d, want %d`, key, loadedData[key], val)
		}
	}
}
