package entities

import (
	"time"
)

type ColorPalateEntity struct {
	ID        uint32    `gorm:"primaryKey" json:"id" xml:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" xml:"created_at" example:"2021-08-06"`
	Colors    string    `gorm:"not null,type:varchar(34)" json:"colors" xml:"colors" example:"FF0000-FF0000-FF0000-FF0000-FF0000"`
}

func (*ColorPalateEntity) TableName() string {
	return "color_palates"
}

// func (c *Color) Migrate(db *gorm.DB, forced bool) error {
// 	if forced {
// 		db.Migrator().DropTable(c)
// 	}

// 	if err := db.AutoMigrate(c); err != nil {
// 		return err
// 	}

// 	var nSize int64
// 	if db.Model(c).Count(&nSize); nSize == 0 {

// 		// The most quick and easiest way !!!
// 		db.Exec(fmt.Sprintf("copy colors(created_at, colors) from '%s/Colors.csv' delimiter ',' csv header;", config.ENV.MigrationPath))
// 	}

// 	return nil
// }

// func (c *Color) Copy() *Color {
// 	return &Color{
// 		ID:        c.ID,
// 		CreatedAt: c.CreatedAt,
// 		Colors:    c.Colors,
// 	}
// }

// func (c *Color) Fill(updated *Color) *Color {
// 	if updated.ID != 0 {
// 		c.ID = updated.ID
// 	}

// 	if !updated.CreatedAt.IsZero() {
// 		c.CreatedAt = updated.CreatedAt
// 	}

// 	if updated.Colors != "" {
// 		c.Colors = updated.Colors
// 	}

// 	return c
// }

// type ColorDto struct {
// 	Colors string `json:"colors" xml:"colors" example:"[#ff0000, #0ff000, #00ff00, #000ff0, #0000ff]"`
// }

// func (c *ColorDto) IsOK() bool {
// 	// TODO: !!

// 	return c.Colors != ""
// }

// type ColorQueryDto struct {
// 	ID    uint32 `form:"id,omitempty"`
// 	Color string `form:"colors,omitempty" example:"#520055"`

// 	Page  int `form:"page,omitempty,default=-1" example:"1"`
// 	Limit int `form:"limit,omitempty" example:"10"`
// }

// func (c *ColorQueryDto) IsOK(model *Pattern) bool {
// 	// return (c.Mode == "" || c.Mode == model.Mode) && (c.Colors == 0 || c.Colors == model.Colors)
// 	return true
// }
