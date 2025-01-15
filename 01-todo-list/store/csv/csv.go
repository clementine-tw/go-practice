package csv

import (
	"encoding/csv"
	"os"
	"slices"
	"strconv"
	"time"
)

type CsvStore struct {
	file *os.File
	data [][]string
}

const defaultFilePath = "./default_tasks.csv"

// New CsvStore instance with default csv file path.
func NewDefaultCsvStore() *CsvStore {
	return NewCsvStore(defaultFilePath)
}

// New CsvStore instance with custom stored csv file path.
func NewCsvStore(path string) *CsvStore {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.FileMode(0644))
	if err != nil {
		panic(err)
	}
	data := readAll(file)
	if len(data) < 1 {
		data = CreateTable()
		writeAll(file, data)
	}
	return &CsvStore{file, data}
}

// Create a table with default headers.
func CreateTable() [][]string {
	return [][]string{{"ID", "Task", "CreatedAt"}}
}

// Close the opened csv file.
func (c *CsvStore) Close() {
	c.file.Close()
}

// Read all content from a csv file.
func readAll(file *os.File) [][]string {
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	return data
}

// Write all to csv file.
func writeAll(file *os.File, data [][]string) {
	writer := csv.NewWriter(file)
	err := writer.WriteAll(data)
	if err != nil {
		panic(err)
	}
}

// Write a line to csv file.
func write(file *os.File, record []string) {
	writer := csv.NewWriter(file)
	err := writer.Write(record)
	if err != nil {
		panic(err)
	}
	writer.Flush()
}

// Clean content of the file.
func cleanFile(file *os.File) {
	file.Truncate(0)
	file.Seek(0, 0)
}
// Add a record to store.
func (c *CsvStore) Add(task string) {
	id := 0
	if len(c.data) > 1 {
		var err error
		id, err = strconv.Atoi(c.data[len(c.data)-1][0])
		if err != nil {
			panic(err)
		}
		id++
	}
	record := []string{strconv.Itoa(id), task, strconv.FormatInt(time.Now().Unix(), 10)}
	c.data = append(c.data, record)
	write(c.file, record)
}

// Delete a record from store.
func (c *CsvStore) Delete(id string) {
	c.data = slices.DeleteFunc(c.data, func(elem []string) bool {
		return elem[0] == id
	})
	cleanFile(c.file)
	writeAll(c.file, c.data)
}

// Update a record in store.
func (c *CsvStore) Update(id string, task string) {
	index := slices.IndexFunc(c.data, func(elem []string) bool {
		return elem[0] == id
	})
	if index == -1 {
		return
	}
	c.data[index][1] = task
	cleanFile(c.file)
	writeAll(c.file, c.data)
}

// Get records from store.
func (c *CsvStore) List() [][]string {
	return c.data
}

// Remove all records from store and create a new table.
func (c *CsvStore) Clean() {
	if len(c.data) <= 1 {
		return
	}
	cleanFile(c.file)
	c.data = CreateTable()
	writeAll(c.file, c.data)
}
