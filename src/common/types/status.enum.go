package types

type StatusEnum int

const (
	Active StatusEnum = iota + 1
	Inactive
)

func (c StatusEnum) String() string {
	switch c {
	case Active:
		return "active"

	case Inactive:
		return "inactive"
	}

	return ""
}

func (StatusEnum) Value(str string) StatusEnum {
	switch str {
	case "active":
		return Active

	case "inactive":
		return Inactive
	}

	return Active
}
