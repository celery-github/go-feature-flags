package flags

import "time"

type RolloutType string

const (
	RolloutAll        RolloutType = "all"
	RolloutNone       RolloutType = "none"
	RolloutPercentage RolloutType = "percentage"
)

type Rollout struct {
	Type       RolloutType `json:"type"`
	Percentage int         `json:"percentage,omitempty"` // 0-100 (used if type=percentage)
}

type Flag struct {
	Name        string    `json:"name"`
	Enabled     bool      `json:"enabled"`
	Description string    `json:"description,omitempty"`
	Envs        []string  `json:"envs,omitempty"` // allowed envs, e.g. ["dev","prod"]. empty => all envs
	Rollout     Rollout   `json:"rollout"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type FlagUpsert struct {
	Enabled     *bool    `json:"enabled,omitempty"`
	Description *string  `json:"description,omitempty"`
	Envs        []string `json:"envs,omitempty"` // if provided, replaces entire env list
	Rollout     *Rollout `json:"rollout,omitempty"`
}
