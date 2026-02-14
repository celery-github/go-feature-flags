package flags

import "testing"

func TestEvaluate_DisabledAlwaysFalse(t *testing.T) {
	f := Flag{
		Name:    "x",
		Enabled: false,
		Rollout: Rollout{Type: RolloutAll},
	}

	if got := Evaluate(f, "prod", "user1"); got {
		t.Fatalf("expected false when flag is disabled, got true")
	}
}

func TestEvaluate_EnvTargeting(t *testing.T) {
	f := Flag{
		Name:    "env_flag",
		Enabled: true,
		Envs:    []string{"dev", "prod"},
		Rollout: Rollout{Type: RolloutAll},
	}

	if got := Evaluate(f, "dev", "u"); !got {
		t.Fatalf("expected true for allowed env dev, got false")
	}
	if got := Evaluate(f, "prod", "u"); !got {
		t.Fatalf("expected true for allowed env prod, got false")
	}
	if got := Evaluate(f, "staging", "u"); got {
		t.Fatalf("expected false for disallowed env staging, got true")
	}
}

func TestEvaluate_EnvTargeting_IsCaseInsensitive(t *testing.T) {
	f := Flag{
		Name:    "case_env",
		Enabled: true,
		Envs:    []string{"Dev"},
		Rollout: Rollout{Type: RolloutAll},
	}

	if got := Evaluate(f, "dev", "u"); !got {
		t.Fatalf("expected true for case-insensitive env match, got false")
	}
}

func TestEvaluate_RolloutAll(t *testing.T) {
	f := Flag{
		Name:    "all_flag",
		Enabled: true,
		Rollout: Rollout{Type: RolloutAll},
	}

	if got := Evaluate(f, "prod", "user1"); !got {
		t.Fatalf("expected true for rollout=all, got false")
	}
}

func TestEvaluate_RolloutNone(t *testing.T) {
	f := Flag{
		Name:    "none_flag",
		Enabled: true,
		Rollout: Rollout{Type: RolloutNone},
	}

	if got := Evaluate(f, "prod", "user1"); got {
		t.Fatalf("expected false for rollout=none, got true")
	}
}

func TestEvaluate_RolloutPercentage_Bounds(t *testing.T) {
	base := Flag{
		Name:    "pct_flag",
		Enabled: true,
		Rollout: Rollout{Type: RolloutPercentage},
	}

	t.Run("0_percent_is_always_false", func(t *testing.T) {
		f := base
		f.Rollout.Percentage = 0
		if got := Evaluate(f, "prod", "user1"); got {
			t.Fatalf("expected false for 0%%, got true")
		}
	})

	t.Run("100_percent_is_always_true", func(t *testing.T) {
		f := base
		f.Rollout.Percentage = 100
		if got := Evaluate(f, "prod", "user1"); !got {
			t.Fatalf("expected true for 100%%, got false")
		}
	})
}

func TestEvaluate_RolloutPercentage_IsDeterministic(t *testing.T) {
	f := Flag{
		Name:    "deterministic",
		Enabled: true,
		Rollout: Rollout{Type: RolloutPercentage, Percentage: 25},
	}

	// same inputs should always produce same result
	first := Evaluate(f, "prod", "user1")
	for i := 0; i < 50; i++ {
		if got := Evaluate(f, "prod", "user1"); got != first {
			t.Fatalf("expected deterministic result, got flip at iter %d: first=%v now=%v", i, first, got)
		}
	}
}

func TestEvaluate_RolloutPercentage_ChangesWithUserOrEnv(t *testing.T) {
	f := Flag{
		Name:    "vary_inputs",
		Enabled: true,
		Rollout: Rollout{Type: RolloutPercentage, Percentage: 25},
	}

	// Not guaranteed to differ, but overwhelmingly likely across different inputs.
	// We'll check "at least one differs" across a small set to avoid flaky tests.
	base := Evaluate(f, "prod", "userA")

	diffCount := 0
	cases := []struct {
		env  string
		user string
	}{
		{"prod", "userB"},
		{"prod", "userC"},
		{"dev", "userA"},
		{"dev", "userB"},
		{"staging", "userA"},
	}

	for _, c := range cases {
		if Evaluate(f, c.env, c.user) != base {
			diffCount++
		}
	}

	if diffCount == 0 {
		t.Fatalf("expected at least one different evaluation when env/user changes, got none")
	}
}

func TestEvaluate_RolloutPercentage_AnonymousUserKey(t *testing.T) {
	f := Flag{
		Name:    "anon_user",
		Enabled: true,
		Rollout: Rollout{Type: RolloutPercentage, Percentage: 50},
	}

	// Ensure empty userKey doesn't panic and is deterministic
	first := Evaluate(f, "prod", "")
	for i := 0; i < 25; i++ {
		if got := Evaluate(f, "prod", ""); got != first {
			t.Fatalf("expected deterministic for empty userKey, first=%v now=%v", first, got)
		}
	}
}

func TestEvaluate_UnknownRolloutType_SafestOff(t *testing.T) {
	f := Flag{
		Name:    "unknown_rollout",
		Enabled: true,
		Rollout: Rollout{Type: "mystery"},
	}

	if got := Evaluate(f, "prod", "user1"); got {
		t.Fatalf("expected false for unknown rollout type, got true")
	}
}

func TestEvaluate_DefaultRolloutTypeEmptyTreatsAsAll(t *testing.T) {
	f := Flag{
		Name:    "empty_rollout",
		Enabled: true,
		Rollout: Rollout{Type: ""}, // Evaluate treats "" like "all"
	}

	if got := Evaluate(f, "prod", "user1"); !got {
		t.Fatalf(`expected true for empty rollout type (treated as "all"), got false`)
	}
}
