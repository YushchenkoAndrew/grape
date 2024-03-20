package types

type ProjectTypeEnum int

const (
	Html ProjectTypeEnum = iota + 1
	Markdown
	Link
	K3s
)

func (c ProjectTypeEnum) String() string {
	switch c {
	case Html:
		return "html"

	case Markdown:
		return "markdown"

	case Link:
		return "link"

	case K3s:
		return "k3s"
	}

	return ""
}

func (ProjectTypeEnum) Value(str string) ProjectTypeEnum {
	switch str {
	case "html":
		return Html

	case "markdown":
		return Markdown

	case "link":
		return Link

	case "k3s":
		return K3s
	}

	return Html
}
