package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strings"
)

type Claims struct {
	RestaurantID string `json:"restname" bson:"restname"`
	Token        string `json:"token" bson:"token"`
	jwt.RegisteredClaims
}

type PermissionResponse struct {
	IsOwner bool `json:"isOwner" bson:"isOwner"`
}

func filterUrl(c *gin.Context) {
	tokenUser := c.GetHeader("Authorization")
	if tokenUser == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
		return
	}

	tokenUser = strings.TrimSpace(tokenUser)
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenUser, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar token"})
		return
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}

	var requestData struct {
		UrlUserName string `json:"urlUserName" bson:"urlUserName"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	urlUserName := requestData.UrlUserName

	collection := DB.Collection("restaurante")
	filter := bson.D{{"restname", urlUserName}}
	var tokenOwner Restaurante

	err = collection.FindOne(context.Background(), filter).Decode(&tokenOwner)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Restaurante não encontrado: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
			return
		}
		log.Printf("Erro ao encontrar o restaurante: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while finding restaurant"})
		return
	}

	// Verificar se o Token existe antes de defini-lo no contexto
	if tokenOwner.Token == "" {
		log.Println("Token do restaurante não encontrado")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Restaurant token not found"})
		return
	}

	// Definindo o token no contexto
	c.Set("Token", tokenOwner.Token)
	log.Printf("Token definido no contexto: %s", tokenOwner.Token)

	isOwner := tokenUser == tokenOwner.Token
	response := gin.H{
		"message": "Restaurant found",
		"Token":   tokenOwner.Token,
		"IsOwner": isOwner,
	}
	c.JSON(http.StatusOK, response)
}

func verifyToken(c *gin.Context) (*Claims, error) {
	tokenUser := c.GetHeader("Authorization")
	if tokenUser == "" {
		return nil, fmt.Errorf("Token não fornecido")
	}

	tokenUser = strings.TrimSpace(tokenUser)
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenUser, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, fmt.Errorf("Token inválido")
		}
		return nil, fmt.Errorf("Erro ao processar token")
	}

	if !token.Valid {
		return nil, fmt.Errorf("Token inválido")
	}

	return claims, nil
}

func getIsOwner(c *gin.Context) {
	authToken := c.GetHeader("Authorization")
	_, err := verifyToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	urlUserName := c.Query("urlUserName")
	if urlUserName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "urlUserName não fornecido"})
		return
	}

	collection := DB.Collection("restaurante")
	filter := bson.D{{"restname", urlUserName}}
	var tokenOwner Restaurante

	err = collection.FindOne(context.Background(), filter).Decode(&tokenOwner)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Restaurante não encontrado: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
			return
		}
		log.Printf("Erro ao encontrar o restaurante: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while finding restaurant"})
		return
	}
	log.Printf("Valor de claims.Token: %s", authToken)

	isOwner := authToken == tokenOwner.Token
	response := gin.H{"IsOwner": isOwner}
	log.Printf("Token do proprietário: %s \n , Token do usuário: %s \n, É proprietário: %v \n", tokenOwner.Token, authToken, isOwner)
	c.JSON(http.StatusOK, response)
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if ok && token.Valid {
			c.Set("urlUserName", claims.RestaurantID)
			log.Printf("Token válido, urlUserName definido no contexto: %s", claims.RestaurantID)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
		}

		c.Next()
	}
}
