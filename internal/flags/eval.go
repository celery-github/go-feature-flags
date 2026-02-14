package flags

import (
	"crypto/sha256"
	"encoding/binary"
	"strings"
)

// Evaluate returns whether a flag is ON for a given environment and user key.
// - env: "dev"|"prod"|anything
// - userKey: stable identifier (email, uuid, etc.). If empty and rollout=percentage, evaluation uses "anonymous".
func Evaluate(flag Flag, env string, userKey string) bool {
	// hard off
	if !flag.Enabled {
		return false
	}

	// env targeting: if Envs is empty => all envs allowed
	if len(flag.Envs) > 0 {
		allowed := false
		for _, e := range flag.Envs {
			if strings.EqualFold(e, env) {
				allowed = true
				break
			}
		}
		if !allowed {
			return false
		}
	}

	switch flag.Rollout.Type {
	case RolloutAll, "":
		return true
	case RolloutNone:
		return false
	case RolloutPercentage:
		p := flag.Rollout.Percentage
		if p <= 0 {
			return false
		}
		if p >= 100 {
			return true
		}
		if userKey == "" {
			userKey = "anonymous"
		}
		// Deterministic bucketing 0..99 using hash(flagName|env|userKey)
		h := sha256.Sum256([]byte(flag.Name + "|" + env + "|" + userKey))
		n := binary.BigEndian.Uint32(h[:4])
		bucket := int(n % 100)
		return bucket < p
	default:
		// unknown rollout type => safest is off
		return false
	}
}
