package hiscores

// SkillEntry holds hiscores rank, level, and XP for one OSRS skill.
type SkillEntry struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Rank  int    `json:"rank"`
	Level int    `json:"level"`
	XP    int64  `json:"xp"`
}

// HiscoresResponse is the top-level response from the OSRS hiscores JSON endpoint.
// The "activities" array (boss kills, minigames) is present in the API but not modelled.
type HiscoresResponse struct {
	Skills []SkillEntry `json:"skills"`
}
