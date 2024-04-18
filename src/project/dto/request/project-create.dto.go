package request

type ProjectCreateDto struct {
	Name        string `json:"name" xml:"name" binding:"required"`
	Description string `json:"description" xml:"description" binding:"omitempty"`
	Type        string `json:"type" xml:"type" binding:"required,oneof=html markdown link k3s"`
	Footer      string `json:"footer" xml:"footer" binding:"omitempty"`
	Link        string `json:"link" xml:"link" binding:"required_if=Type link,excluded_unless=Type link,url"`
	README      bool   `json:"readme" xml:"readme" form:"readme,default=false" binding:"omitempty,excluded_if=Type link"`
}
