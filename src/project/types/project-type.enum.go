package types

type ProjectTypeEnum int

const (
	P5js ProjectTypeEnum = iota + 1
	Emscripten
	Html
	Markdown
	Link
	K3s
)

func (c ProjectTypeEnum) String() string {
	switch c {
	case P5js:
		return "p5js"

	case Emscripten:
		return "emscripten"

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
	case "p5js":
		return P5js

	case "emscripten":
		return Emscripten

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
