package storage

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

// Storage defines the interface for storing and retrieving request data.
type Storage interface {
	Load() (map[time.Time]int, error)
	Save(data map[time.Time]int) error
}

// FileStorage implements the Storage interface using file-based storage.
type FileStorage struct {
	filename string
}

func NewFileStorage(filename string) *FileStorage {
	return &FileStorage{filename: filename}
}

func (fs *FileStorage) Load() (map[time.Time]int, error) {
	data := make(map[time.Time]int)
	file, err := os.Open(fs.filename)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty data if file doesn't exist
			return data, nil
		}
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		timestamp, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			return nil, err
		}
		count, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, err
		}
		data[time.Unix(timestamp, 0).UTC()] += count
	}

	return data, nil
}

func (fs *FileStorage) Save(data map[time.Time]int) error {
	file, err := os.Create(fs.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for t, count := range data {
		err := writer.Write([]string{
			strconv.FormatInt(t.UTC().Unix(), 10),
			strconv.Itoa(count),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
