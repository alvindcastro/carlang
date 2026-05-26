package spell

import "testing"

func TestFromRecipeNormalizesRecipes(t *testing.T) {
	tests := []struct {
		recipe string
		want   Spell
	}{
		{"QQQ", ColdSnap},
		{"WQQ", GhostWalk},
		{"EQW", DeafeningBlast},
		{"EEE", SunStrike},
	}

	for _, tt := range tests {
		t.Run(tt.recipe, func(t *testing.T) {
			got, err := FromRecipe(tt.recipe)
			if err != nil {
				t.Fatalf("FromRecipe returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("FromRecipe(%q) = %s, want %s", tt.recipe, got, tt.want)
			}
		})
	}
}

func TestNormalizeRecipeRejectsInvalidRecipes(t *testing.T) {
	tests := []string{
		"QQ",
		"QQQQ",
		"QXR",
	}

	for _, recipe := range tests {
		t.Run(recipe, func(t *testing.T) {
			if _, err := NormalizeRecipe(recipe); err == nil {
				t.Fatalf("NormalizeRecipe(%q) returned nil error", recipe)
			}
		})
	}
}

func TestSpellString(t *testing.T) {
	tests := []struct {
		spell Spell
		want  string
	}{
		{Empty, "empty"},
		{ColdSnap, "Cold Snap"},
		{GhostWalk, "Ghost Walk"},
		{IceWall, "Ice Wall"},
		{Tornado, "Tornado"},
		{EMP, "EMP"},
		{Alacrity, "Alacrity"},
		{ChaosMeteor, "Chaos Meteor"},
		{ForgeSpirit, "Forge Spirit"},
		{SunStrike, "Sun Strike"},
		{DeafeningBlast, "Deafening Blast"},
		{Spell(99), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.spell.String(); got != tt.want {
			t.Fatalf("%d.String() = %q, want %q", tt.spell, got, tt.want)
		}
	}
}
