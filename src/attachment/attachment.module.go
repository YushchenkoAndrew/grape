package attachment

// import (
// 	c "grape/controllers"
// 	"grape/interfaces"
// 	i "grape/interfaces/controller"
// 	m "grape/middleware"
// 	s "grape/service"

// 	"github.com/gin-gonic/gin"
// 	"github.com/go-redis/redis/v8"
// 	"gorm.io/gorm"
// )

// type fileRouter struct {
// 	route *gin.RouterGroup
// 	auth  *gin.RouterGroup
// 	file  i.Default
// }

// func NewFileRouter(rg *gin.RouterGroup, db *gorm.DB, client *redis.Client) interfaces.Router {
// 	return &fileRouter{
// 		route: rg.Group("/file"),
// 		auth:  rg.Group("/file", m.GetMiddleware().Auth()),
// 		file:  c.NewFileController(s.NewFileService(db, client)),
// 	}
// }

// func (c *fileRouter) Init() {
// 	c.auth.POST("/list/:id", c.file.CreateAll)
// 	c.auth.POST("/:id", c.file.CreateOne)

// 	c.route.GET("/:id", c.file.ReadOne)
// 	c.route.GET("", c.file.ReadAll)

// 	c.auth.PUT("/:id", c.file.UpdateOne)
// 	c.auth.PUT("", c.file.UpdateAll)

// 	c.auth.DELETE("/:id", c.file.DeleteOne)
// 	c.auth.DELETE("", c.file.DeleteAll)
// }
