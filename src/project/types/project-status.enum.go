package types

type ProjectStatusEnum int

const (
	Active ProjectStatusEnum = iota
	Inactive
)

func (c ProjectStatusEnum) String() string {
	switch c {
	case Active:
		return "active"

	case Inactive:
		return "inactive"
	}

	return ""
}

func (ProjectStatusEnum) Value(str string) ProjectStatusEnum {
	switch str {
	case "active":
		return Active

	case "inactive":
		return Inactive
	}

	return Active
}
