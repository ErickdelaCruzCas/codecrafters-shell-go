package history

type Store struct {
	entries []string
}

func New() *Store {
	return &Store{
		entries: []string{},
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
