package types

type PatternColorModeEnum int

const (
	Stroke PatternColorModeEnum = iota + 1
	Fill
	Join
)

func (c PatternColorModeEnum) String() string {
	switch c {
	case Stroke:
		return "stroke"

	case Fill:
		return "fill"

	case Join:
		return "Join"
	}

	return ""
}

func (PatternColorModeEnum) Value(str string) PatternColorModeEnum {
	switch str {
	case "stroke":
		return Stroke

	case "fill":
		return Fill

	case "join":
		return Join
	}

	return Join
}
