package util

// CalculateLevelFromExperience calculates user level based on experience points
// Level formula based on API spec:
// - Level 1: 0-199 XP
// - Level 2: 200-499 XP
// - Level 3: 500-1499 XP
// - Level 4: 1500-2999 XP
// - Level 5: 3000-4999 XP
// - Level 6: 5000-6999 XP
// - Level 7+: +2000 XP per level (5000 + (level - 6) * 2000)
// - Level 36 (max): 65000+ XP
func CalculateLevelFromExperience(exp int32) int32 {
	if exp < 200 {
		return 1
	}
	if exp < 500 {
		return 2
	}
	if exp < 1500 {
		return 3
	}
	if exp < 3000 {
		return 4
	}
	if exp < 5000 {
		return 5
	}
	if exp < 7000 {
		return 6
	}

	// Level 7+: 5000 + (level - 6) * 2000
	// Solve: exp >= 5000 + (level - 6) * 2000
	// level = 6 + (exp - 5000) / 2000
	level := 6 + (exp-5000)/2000
	if level > 36 {
		return 36 // Max level cap
	}
	return level
}
