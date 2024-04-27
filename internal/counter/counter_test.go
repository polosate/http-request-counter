package counter

import (
	"testing"
	"time"
)

func TestRequestCounter_AddRequest(t *testing.T) {
	mockStorage := &mockStorage{
		data: make(map[time.Time]int),
	}

	counter, err := New(mockStorage)
	if err != nil {
		t.Errorf(`Failed to cretae a counter service`)
	}

	counter.AddRequest()
	counter.AddRequest()
	counter.AddRequest()

	if len(mockStorage.data) != 1 {
		t.Errorf(`Unexpected number of records: got %d, want %d`, len(mockStorage.data), 1)
	}
	for _, v := range mockStorage.data {
		if v != 3 {
			t.Errorf(`Unexpected number of requests: got %d, want %d`, v, 3)
		}
	}

}

func TestRequestCounter_CountRequests(t *testing.T) {
	mockStorage := &mockStorage{
		data: map[time.Time]int{
			time.Now().UTC():                        10,
			time.Now().UTC().Add(-10 * time.Second): 5,
		},
	}

	counter, err := New(mockStorage)
	if err != nil {
		t.Errorf(`Failed to cretae a counter service`)
	}

	count := counter.CountRequests()

	if count != 15 {
		t.Errorf(`Unexpected count of requests: got %d, want %d`, count, 15)
	}
}

// Mock storage for testing
type mockStorage struct {
	data map[time.Time]int
}

func (m *mockStorage) Load() (map[time.Time]int, error) {
	return m.data, nil
}

func (m *mockStorage) Save(data map[time.Time]int) error {
	m.data = data
	return nil
}
