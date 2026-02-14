package flags

import (
	"errors"
	"sort"
	"sync"
	"time"
)

var ErrNotFound = errors.New("flag not found")

type Store interface {
	List() ([]Flag, error)
	Get(name string) (Flag, error)
	Put(flag Flag) (Flag, error)
	Patch(name string, patch FlagUpsert) (Flag, error)
	Delete(name string) error
}

type InMemoryStore struct {
	mu    sync.RWMutex
	flags map[string]Flag
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{flags: map[string]Flag{}}
}

func (s *InMemoryStore) List() ([]Flag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Flag, 0, len(s.flags))
	for _, f := range s.flags {
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func (s *InMemoryStore) Get(name string) (Flag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	f, ok := s.flags[name]
	if !ok {
		return Flag{}, ErrNotFound
	}
	return f, nil
}

func (s *InMemoryStore) Put(flag Flag) (Flag, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	flag.UpdatedAt = time.Now().UTC()
	s.flags[flag.Name] = flag
	return flag, nil
}

func (s *InMemoryStore) Patch(name string, patch FlagUpsert) (Flag, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, ok := s.flags[name]
	if !ok {
		return Flag{}, ErrNotFound
	}

	if patch.Enabled != nil {
		f.Enabled = *patch.Enabled
	}
	if patch.Description != nil {
		f.Description = *patch.Description
	}
	if patch.Envs != nil {
		// if user sends [] it's allowed (means all envs? or none). We'll treat empty as "all envs" at eval time.
		f.Envs = patch.Envs
	}
	if patch.Rollout != nil {
		f.Rollout = *patch.Rollout
	}
	f.UpdatedAt = time.Now().UTC()

	s.flags[name] = f
	return f, nil
}

func (s *InMemoryStore) Delete(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.flags[name]; !ok {
		return ErrNotFound
	}
	delete(s.flags, name)
	return nil
}
