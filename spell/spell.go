package spell

import (
	"fmt"
	"strings"
)

type Spell int

const (
	Empty Spell = iota
	ColdSnap
	GhostWalk
	IceWall
	Tornado
	EMP
	Alacrity
	ChaosMeteor
	ForgeSpirit
	SunStrike
	DeafeningBlast
)

func (s Spell) String() string {
	switch s {
	case Empty:
		return "empty"
	case ColdSnap:
		return "Cold Snap"
	case GhostWalk:
		return "Ghost Walk"
	case IceWall:
		return "Ice Wall"
	case Tornado:
		return "Tornado"
	case EMP:
		return "EMP"
	case Alacrity:
		return "Alacrity"
	case ChaosMeteor:
		return "Chaos Meteor"
	case ForgeSpirit:
		return "Forge Spirit"
	case SunStrike:
		return "Sun Strike"
	case DeafeningBlast:
		return "Deafening Blast"
	default:
		return "unknown"
	}
}

func FromRecipe(recipe string) (Spell, error) {
	normalized, err := NormalizeRecipe(recipe)
	if err != nil {
		return Empty, err
	}

	switch normalized {
	case "QQQ":
		return ColdSnap, nil
	case "QQW":
		return GhostWalk, nil
	case "QQE":
		return IceWall, nil
	case "QWW":
		return Tornado, nil
	case "WWW":
		return EMP, nil
	case "WWE":
		return Alacrity, nil
	case "WEE":
		return ChaosMeteor, nil
	case "QEE":
		return ForgeSpirit, nil
	case "EEE":
		return SunStrike, nil
	case "QWE":
		return DeafeningBlast, nil
	default:
		return Empty, fmt.Errorf("unknown recipe %q", recipe)
	}
}

func NormalizeRecipe(recipe string) (string, error) {
	if len(recipe) != 3 {
		return "", fmt.Errorf("recipe %q must contain exactly three Q/W/E components", recipe)
	}

	var q, w, e int
	for _, ch := range recipe {
		switch ch {
		case 'Q':
			q++
		case 'W':
			w++
		case 'E':
			e++
		default:
			return "", fmt.Errorf("recipe %q contains invalid component %q", recipe, ch)
		}
	}

	return strings.Repeat("Q", q) + strings.Repeat("W", w) + strings.Repeat("E", e), nil
}
