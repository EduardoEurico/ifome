package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

var secretKey = []byte(os.Getenv("SECRET_KEY")) // Corrigido o nome do env para uma chave mais segura

type Restaurante struct {
	RestID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RestName          string             `bson:"restname" json:"restname"`
	Email             string             `bson:"email" json:"email"`
	Password          string             `bson:"password" json:"password"`
	RestauranteCodigo string             `bson:"restauranteCodigo" json:"restauranteCodigo"`
	Listas            []Lista            `bson:"listas" json:"listas"`
	Image             string             `bson:"image" json:"image"`
	Token             string             `bson:"token" json:"token"`
}
type RestaurantePersonalizacao struct {
	RestauranteCodigo string `bson:"restauranteCodigo" json:"restauranteCodigo"`
	RestName          string `bson:"restname" json:"restname"`
	Image             string `bson:"image" json:"image"`
}

func contaRest(c *gin.Context) {
	var rest Restaurante
	if err := c.ShouldBindJSON(&rest); err != nil {
		log.Printf("Erro ao fazer bind do JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	const RestauranteCodigoFixo = "567"
	if rest.RestauranteCodigo != RestauranteCodigoFixo {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant code"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Erro ao criptografar senha: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while hashing password"})
		return
	}
	rest.Password = string(hashedPassword)

	collection := DB.Collection("restaurante")
	doc := bson.D{
		{"restname", rest.RestName},
		{"email", rest.Email},
		{"password", rest.Password},
		{"restauranteCodigo", rest.RestauranteCodigo},
		{"listas", []Lista{}},
		{"image", rest.Image},
	}

	_, err = collection.InsertOne(context.Background(), doc)
	if err != nil {
		log.Printf("Erro ao inserir restaurante no MongoDB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new restaurant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Restaurant created successfully"})
}

func loginRest(c *gin.Context) {
	var loginInfo Restaurante
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		log.Printf("Erro ao fazer bind do JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := DB.Collection("restaurante")
	filter := bson.D{{"restname", loginInfo.RestName}}
	var result Restaurante

	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Restaurante não encontrado: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
			return
		}
		log.Printf("Erro ao encontrar restaurante: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while finding restaurant"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(loginInfo.Password))
	if err != nil {
		log.Printf("Erro ao comparar senha: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.MapClaims{
		"restid":   result.RestID.Hex(),
		"restname": result.RestName,
		"exp":      expirationTime.Unix(),
		"type":     "restaurante",
		"token":    result.Token,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("Erro ao gerar token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while generating token"})
		return
	}

	// Atualizar o token no banco de dados
	update := bson.D{{"$set", bson.D{{"token", tokenString}}}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Erro ao atualizar token no banco de dados: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating token in the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "message": "Restaurant logged in successfully", "restname": result.RestName})
}

func findRest(c *gin.Context) {
	restname := c.Param("restname")

	collection := DB.Collection("restaurante")
	filter := bson.D{{"restname", restname}}
	var result Restaurante

	err := collection.FindOne(context.Background(), filter).Decode(&result)
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

	c.JSON(http.StatusOK, gin.H{
		"message":  "Restaurant found",
		"restname": restname,
	})
}

func getAllRestaurantes(c *gin.Context) {
	collection := DB.Collection("restaurante")
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Printf("Erro ao buscar restaurantes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching restaurants"})
		return
	}
	defer cur.Close(context.Background())

	var restaurantes []RestaurantePersonalizacao
	for cur.Next(context.Background()) {
		var rest RestaurantePersonalizacao
		err := cur.Decode(&rest)
		if err != nil {
			log.Printf("Erro ao decodificar restaurante: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while decoding restaurant"})
			return
		}
		restaurantes = append(restaurantes, rest)
	}

	if err := cur.Err(); err != nil {
		log.Printf("Erro no cursor: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
		return
	}

	c.JSON(http.StatusOK, restaurantes)
}

func updateRestPersonalization(c *gin.Context) {
	collection := DB.Collection("restaurante")
	restname := c.Param("userName")

	var requestData struct {
		ImagemBase64 string `bson:"image" json:"image"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	filter := bson.D{{"restname", restname}}
	update := bson.D{
		{"$set", bson.D{
			{"image", requestData.ImagemBase64},
		}},
	}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Erro ao atualizar a imagem do restaurante: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update restaurant image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Restaurant image updated successfully"})
}
func getImage(c *gin.Context) {
	collection := DB.Collection("restaurante")
	restname := c.Param("userName")
	filter := bson.D{{"restname", restname}}
	var result RestaurantePersonalizacao

	err := collection.FindOne(context.Background(), filter).Decode(&result)
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

	log.Printf("Imagem encontrada para o restaurante %s: %s", restname, result.Image)
	c.JSON(http.StatusOK, gin.H{"image": result.Image})
}
