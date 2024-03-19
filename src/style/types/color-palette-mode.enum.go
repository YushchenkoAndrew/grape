package types

type ColorPaletteModeEnum int

const (
	Stroke ColorPaletteModeEnum = iota
	Fill
	Join
)

func (c ColorPaletteModeEnum) String() string {
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

func (ColorPaletteModeEnum) Value(str string) ColorPaletteModeEnum {
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
