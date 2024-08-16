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
	"time"
)

type User struct {
	UserID   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Erro ao fazer bind do JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Erro ao criptografar senha: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while hashing password"})
		return
	}
	user.Password = string(hashedPassword)

	collection := DB.Collection("users")

	// Create a BSON document
	doc := bson.D{
		{"username", user.Username},
		{"email", user.Email},
		{"password", user.Password},
	}

	_, err = collection.InsertOne(context.Background(), doc)
	if err != nil {
		log.Printf("Erro ao inserir usuário no MongoDB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func findUser(c *gin.Context) {
	username := c.Param("username")

	collection := DB.Collection("users")
	filter := bson.D{{"username", username}}
	var result User

	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Usuário não encontrado: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		log.Printf("Erro ao encontrar usuário: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while finding user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "User found",
		"username": username,
	})
}

func login(c *gin.Context) {
	var loginInfo User
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		log.Printf("Erro ao fazer bind do JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := DB.Collection("users")
	filter := bson.D{{"username", loginInfo.Username}}
	var result User

	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Usuário não encontrado: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
			return
		}
		log.Printf("Erro ao encontrar usuário: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while finding user"})
		return
	}

	// Compare the password with the hashed password stored in the database
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(loginInfo.Password))
	if err != nil {
		log.Printf("Erro ao comparar senha: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid":   result.UserID.Hex(),
		"username": result.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expira em 24 horas
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("Erro ao gerar token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while generating token"})
		return
	}
	update := bson.D{{"$set", bson.D{{"token", tokenString}}}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Erro ao atualizar token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Logged in successfully",
		"token":    tokenString,
		"userName": result.Username, // Inclui o username na resposta
	})
}
