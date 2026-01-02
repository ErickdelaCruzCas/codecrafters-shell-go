package history

import (
	"bufio"
	"os"
)

type Store struct {
	entries   []string
	lastSaved int
}

func New() *Store {
	return &Store{
		entries:   []string{},
		lastSaved: 0,
	}
}

func (s *Store) Add(line string) {
	s.entries = append(s.entries, line)
}

func (s *Store) List() []string {
	out := make([]string, len(s.entries))
	copy(out, s.entries)
	return out
}

func (s *Store) LoadFrom(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	entries := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entries = append(entries, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	s.entries = entries
	s.lastSaved = len(entries)
	return nil
}

func (s *Store) WriteTo(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, entry := range s.entries {
		if _, err := writer.WriteString(entry + "\n"); err != nil {
			return err
		}
	}
	if err := writer.Flush(); err != nil {
		return err
	}

	s.lastSaved = len(s.entries)
	return nil
}

func (s *Store) AppendTo(path string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, entry := range s.entries[s.lastSaved:] {
		if _, err := writer.WriteString(entry + "\n"); err != nil {
			return err
		}
	}
	if err := writer.Flush(); err != nil {
		return err
	}

	s.lastSaved = len(s.entries)
	return nil
}
