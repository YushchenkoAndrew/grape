package entities

import (
	"time"
)

type SvgPatternEntity struct {
	ID        uint32    `gorm:"primaryKey" json:"id" xml:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" xml:"created_at" example:"2021-08-06"`
	Mode      string    `gorm:"not null" json:"mode" xml:"mode" example:"fill"`
	Colors    uint8     `gorm:"not null" json:"colors" xml:"colors" example:"5"`

	MaxStroke   float32 `gorm:"not null" json:"max_stroke" xml:"max_stroke" example:"5"`
	MaxScale    uint8   `gorm:"not null" json:"max_scale" xml:"max_scale" example:"5"`
	MaxSpacingX float32 `gorm:"not null" json:"max_spacing_x" xml:"max_spacing_x" example:"5"`
	MaxSpacingY float32 `gorm:"not null" json:"max_spacing_y" xml:"max_spacing_y" example:"5"`

	Width  float32 `gorm:"not null" json:"width" xml:"width" example:"70.00"`
	Height float32 `gorm:"not null" json:"height" xml:"height" example:"8.00"`

	Path string `gorm:"not null" json:"path" xml:"path" example:"<path d='M-.02 22c8.373 0 11.938-4.695 16.32-9.662C20.785 7.258 25.728 2 35 2c9.272 0 14.215 5.258 18.7 10.338C58.082 17.305 61.647 22 70.02 22M-.02 14.002C8.353 14 11.918 9.306 16.3 4.339 20.785-.742 25.728-6 35-6 44.272-6 49.215-.742 53.7 4.339c4.382 4.967 7.947 9.661 16.32 9.664M70 6.004c-8.373-.001-11.918-4.698-16.3-9.665C49.215-8.742 44.272-14 35-14c-9.272 0-14.215 5.258-18.7 10.339C11.918 1.306 8.353 6-.02 6.002'/>"`
}

func (*SvgPatternEntity) TableName() string {
	return "svg_patterns"
}

// func (c *Pattern) Migrate(db *gorm.DB, forced bool) error {
// 	if forced {
// 		db.Migrator().DropTable(c)
// 	}

// 	if err := db.AutoMigrate(c); err != nil {
// 		return err
// 	}

// 	var nSize int64
// 	if db.Model(c).Count(&nSize); nSize == 0 {

// 		// The most quick and easiest way !!!
// 		db.Exec(fmt.Sprintf("copy pattern(created_at, mode, colors, max_stroke, max_scale, max_spacing_x, max_spacing_y, width, height, path) from '%s/Patterns.csv' delimiter ',' csv header;", config.ENV.MigrationPath))
// 	}

// 	return nil
// }

// func (c *Pattern) Copy() *Pattern {
// 	return &Pattern{
// 		ID:          c.ID,
// 		CreatedAt:   c.CreatedAt,
// 		Mode:        c.Mode,
// 		Colors:      c.Colors,
// 		MaxStroke:   c.MaxStroke,
// 		MaxScale:    c.MaxScale,
// 		MaxSpacingX: c.MaxSpacingX,
// 		MaxSpacingY: c.MaxSpacingY,
// 		Width:       c.Width,
// 		Height:      c.Height,
// 		Path:        c.Path,
// 	}
// }

// func (c *Pattern) Fill(updated *Pattern) *Pattern {
// 	if updated.ID != 0 {
// 		c.ID = updated.ID
// 	}

// 	if !updated.CreatedAt.IsZero() {
// 		c.CreatedAt = updated.CreatedAt
// 	}

// 	if updated.Mode != "" {
// 		c.Mode = updated.Mode
// 	}

// 	if updated.Colors != 0 {
// 		c.Colors = updated.Colors
// 	}

// 	if updated.MaxStroke != 0 {
// 		c.MaxStroke = updated.MaxStroke
// 	}

// 	if updated.MaxSpacingX != 0 {
// 		c.MaxSpacingX = updated.MaxSpacingX
// 	}

// 	if updated.MaxSpacingY != 0 {
// 		c.MaxSpacingY = updated.MaxSpacingY
// 	}

// 	if updated.Width != 0 {
// 		c.Width = updated.Width
// 	}

// 	if updated.Height != 0 {
// 		c.Height = updated.Height
// 	}

// 	if updated.Path != "" {
// 		c.Path = updated.Path
// 	}

// 	return c
// }

// type PatternDto struct {
// 	Mode   string `json:"mode" xml:"mode" example:"fill"`
// 	Colors uint8  `json:"colors" xml:"colors" example:"5"`

// 	MaxStroke   float32 `json:"max_stroke" xml:"max_stroke" example:"5"`
// 	MaxScale    uint8   `json:"max_scale" xml:"max_scale" example:"5"`
// 	MaxSpacingX float32 `json:"max_spacing_x" xml:"max_spacing_x" example:"5"`
// 	MaxSpacingY float32 `json:"max_spacing_y" xml:"max_spacing_y" example:"5"`

// 	Width   float32 `json:"width" xml:"width" example:"70.00"`
// 	Height  float32 `json:"height" xml:"height" example:"8.00"`
// 	VHeight uint16  `json:"v_height" xml:"v_height" example:"8"`

// 	Path string `json:"path" xml:"path" example:"<path d='M-.02 22c8.373 0 11.938-4.695 16.32-9.662C20.785 7.258 25.728 2 35 2c9.272 0 14.215 5.258 18.7 10.338C58.082 17.305 61.647 22 70.02 22M-.02 14.002C8.353 14 11.918 9.306 16.3 4.339 20.785-.742 25.728-6 35-6 44.272-6 49.215-.742 53.7 4.339c4.382 4.967 7.947 9.661 16.32 9.664M70 6.004c-8.373-.001-11.918-4.698-16.3-9.665C49.215-8.742 44.272-14 35-14c-9.272 0-14.215 5.258-18.7 10.339C11.918 1.306 8.353 6-.02 6.002'/>"`
// }

// func (c *PatternDto) IsOK() bool {
// 	decoder := xml.NewDecoder(strings.NewReader(c.Path))

// 	decoder.Strict = false
// 	decoder.AutoClose = xml.HTMLAutoClose
// 	decoder.Entity = xml.HTMLEntity

// 	return c.Mode != "" && c.Colors > 1 && c.MaxStroke > 0 && c.MaxScale > 0 && c.Width > 0 && c.Height > 0 && v.ValidateHTML(decoder, []v.ValidationHTMLCondition{
// 		{
// 			Name: "path",
// 			El:   "start",
// 			Err:  nil,
// 		},
// 		{
// 			Name: "path",
// 			El:   "end",
// 			Err:  nil,
// 		},
// 		{
// 			Name: "",
// 			El:   "",
// 			Err:  io.EOF,
// 		},
// 	})
// }

// type PatternQueryDto struct {
// 	ID     uint32 `form:"id,omitempty"`
// 	Mode   string `form:"mode,omitempty" example:"fill"`
// 	Colors uint8  `form:"colors,omitempty" example:"5"`

// 	Page  int `form:"page,omitempty,default=-1" example:"1"`
// 	Limit int `form:"limit,omitempty" example:"10"`
// }

// func (c *PatternQueryDto) IsOK(model *Pattern) bool {
// 	return (c.Mode == "" || c.Mode == model.Mode) && (c.Colors == 0 || c.Colors == model.Colors)
// }
