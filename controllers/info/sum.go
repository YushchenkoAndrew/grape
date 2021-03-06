package info

import (
	"api/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type sumController struct{}

func NewSumController() interface{} {
	return &sumController{}
}

// @Tags Info
// @Summary Get Info Sum
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Success 200 {object} models.Success{result=[]models.StatInfo}
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /info/sum [get]
func (o *sumController) Read(c *gin.Context) {
	helper.ErrHandler(c, http.StatusInternalServerError, "Not implimented")

	// var stat models.StatInfo
	// ctx := context.Background()

	// if data, err := db.Redis.Get(ctx, "Info:Sum").Result(); err == nil {
	// 	json.Unmarshal([]byte(data), &stat)
	// } else {
	// 	result := db.DB.Table("info").
	// 		Select("SUM(views) as views, SUM(clicks) AS clicks, SUM(media) as media, SUM(visitors) as visitors").
	// 		Scan(&stat)

	// 	if result.Error != nil {
	// 		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
	// 		go logs.SendLogs(&models.LogMessage{
	// 			Stat:    "ERR",
	// 			Name:    "API",
	// 			Url:     "/api/world",
	// 			File:    "/controllers/info/sum.go",
	// 			Message: "It's not an error Karl; It's a bug!!",
	// 			Desc:    result.Error,
	// 		})
	// 		return
	// 	}

	// 	// Encode json to str
	// 	if str, err := json.Marshal(&stat); err == nil {
	// 		go db.Redis.Set(ctx, "Info:Sum", str, 0)
	// 	}
	// }

	// items, err := db.Redis.Get(ctx, "nINFO").Int64()
	// if err != nil {
	// 	items = -1
	// 	go (&models.Info{}).Redis(db.DB, db.Redis)
	// 	go logs.SendLogs(&models.LogMessage{
	// 		Stat:    "ERR",
	// 		Name:    "API",
	// 		File:    "/controllers/info/sum.go",
	// 		Message: "Ohh nooo Cache is broken; Anyway...",
	// 		Desc:    err.Error(),
	// 	})
	// }

	// helper.ResHandler(c, http.StatusOK, models.Success{
	// 	Status:     "OK",
	// 	Result:     []models.StatInfo{stat},
	// 	Items:      1,
	// 	TotalItems: items,
	// })
}
