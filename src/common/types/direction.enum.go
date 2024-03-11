package types

type DirectionEnum int

const (
	Asc DirectionEnum = iota
	Desc
)

func (c DirectionEnum) String() string {
	switch c {
	case Asc:
		return "asc"

	case Desc:
		return "desc"
	}

	return ""
}

func (c DirectionEnum) Bool() bool {
	return c == Desc
}

func (DirectionEnum) Value(str string) DirectionEnum {
	switch str {
	case "asc":
		return Asc

	case "dsc":
		return Desc
	}

	return Desc
}
