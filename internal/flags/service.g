package flags

import (
	"encoding/json"
	"fmt"
	"os"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) List() ([]Flag, error) {
	return s.store.List()
}

func (s *Service) Get(name string) (Flag, error) {
	return s.store.Get(name)
}

func (s *Service) Put(flag Flag) (Flag, error) {
	if flag.Name == "" {
		return Flag{}, fmt.Errorf("name is required")
	}
	// sensible defaults
	if flag.Rollout.Type == "" {
		flag.Rollout.Type = RolloutAll
	}
	if flag.Rollout.Type == RolloutPercentage {
		if flag.Rollout.Percentage < 0 || flag.Rollout.Percentage > 100 {
			return Flag{}, fmt.Errorf("rollout.percentage must be 0..100")
		}
	}
	return s.store.Put(flag)
}

func (s *Service) Patch(name string, patch FlagUpsert) (Flag, error) {
	if patch.Rollout != nil && patch.Rollout.Type == RolloutPercentage {
		if patch.Rollout.Percentage < 0 || patch.Rollout.Percentage > 100 {
			return Flag{}, fmt.Errorf("rollout.percentage must be 0..100")
		}
	}
	return s.store.Patch(name, patch)
}

func (s *Service) Delete(name string) error {
	return s.store.Delete(name)
}

func (s *Service) Evaluate(name, env, userKey string) (bool, Flag, error) {
	f, err := s.store.Get(name)
	if err != nil {
		return false, Flag{}, err
	}
	return Evaluate(f, env, userKey), f, nil
}

type seedFile struct {
	Flags []Flag `json:"flags"`
}

func (s *Service) LoadFromFile(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var sf seedFile
	if err := json.Unmarshal(b, &sf); err != nil {
		return fmt.Errorf("parse seed json: %w", err)
	}
	for _, f := range sf.Flags {
		if _, err := s.Put(f); err != nil {
			return fmt.Errorf("seed put %q: %w", f.Name, err)
		}
	}
	return nil
}
