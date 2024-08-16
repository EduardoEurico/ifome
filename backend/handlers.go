package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func handler() {
	// Conectar ao banco de dados
	connectDB()
	defer disconnectDB()

	r := gin.Default()
	// Configuração do CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	config.AllowHeaders = append(config.AllowHeaders, "params") // Add 'params' to the list of allowed headers

	r.Use(cors.New(config))

	// Definindo as rotas

	//User normal
	r.POST("/conta", createUser)
	r.GET("/login/user/:username", findUser)
	r.POST("/login", login)
	r.GET("/auth/protected", func(c *gin.Context) {
		userID := c.MustGet("userID").(string)
		username := c.MustGet("username").(string)
		c.JSON(http.StatusOK, gin.H{"userID": userID, "username": username})
	})

	//User Restaurante
	r.POST("/login/restaurante", loginRest)
	r.POST("/conta/restaurante", contaRest)
	r.GET("/restaurante", getAllRestaurantes)

	r.GET("/login/rest/:restname", findRest)
	r.POST("/saveLists/:userName", saveLists)
	r.DELETE("/restaurante/:userName/lista/:nome", deleteList)
	r.GET("/restaurante/:userName/lista/get", getLists)

	r.POST("/restaurante/:userName/imagem", updateRestPersonalization)
	r.GET("/restaurante/:userName/pegarimagem", getImage)

	r.POST("/filterUrl", filterUrl)

	// Checar permissões
	r.GET("/isOwner", TokenAuthMiddleware(), getIsOwner)

	if err := r.Run(":8080"); err != nil {
		log.Fatalln(err)
	}
}
