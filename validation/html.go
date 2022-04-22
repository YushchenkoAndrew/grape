package validation

import "encoding/xml"

func IsHtmlTag(tt xml.Token, name, el string) bool {
	// A Token is an interface holding one of the token types:
	// StartElement, EndElement, CharData, Comment, ProcInst, or Directive.
	switch tt.(type) {
	case xml.StartElement:
		return el == "start" && tt.(xml.StartElement).Name.Local == name

	case xml.EndElement:
		return el == "end" && tt.(xml.EndElement).Name.Local == name

	default:
		return false
	}
}

type ValidationHTMLCondition struct {
	Name string
	El   string
	Err  error
}

func ValidateHTML(d *xml.Decoder, validation []ValidationHTMLCondition) bool {
	for _, item := range validation {
		tt, err := d.Token()
		if err != item.Err || !(item.Err != nil || IsHtmlTag(tt, item.Name, item.El)) {
			return false
		}
	}

	return true
}
