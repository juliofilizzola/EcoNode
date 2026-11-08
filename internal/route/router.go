package route

import "github.com/gin-gonic/gin"

func SetupEngine() *gin.Engine {
	engine := gin.Default()
	auth(engine)
	health(engine)
	return engine
}
