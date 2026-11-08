package route

import "github.com/gin-gonic/gin"

func auth(Engine *gin.Engine) {
	authGroup := Engine.Group("/auth")
	authGroup.POST("/login", login)
	authGroup.POST("/logout", logout)
}

func login(c *gin.Context) {
	// todo chamar o controller nesse ponto
}

func logout(c *gin.Context) {
	// todo: chamar o controller nesse ponto
}
