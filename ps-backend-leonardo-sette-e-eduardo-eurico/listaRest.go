package main

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"log"
	"net/http"
)

type ItemLista struct {
	Nome      string  `bson:"nome" json:"nome"`
	Valor     float64 `bson:"valor" json:"valor"`
	Descricao string  `bson:"descricao" json:"descricao"`
}

type Lista struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Nome      string             `bson:"nome" json:"nome"`
	Itens     []ItemLista        `bson:"itens" json:"itens"`
	Categoria string             `bson:"categoria" json:"categoria"`
}

func saveLists(c *gin.Context) {
	userName := c.Param("userName")
	var newList Lista

	if err := c.ShouldBindJSON(&newList); err != nil {
		log.Printf("Erro ao fazer bind do JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	log.Printf("Dados recebidos: %+v\n", newList)

	collection := DB.Collection("restaurante")
	filter := bson.D{{"restname", userName}}
	var restaurante Restaurante

	err := collection.FindOne(context.Background(), filter).Decode(&restaurante)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding restaurant"})
		return
	}

	listaEncontrada := false
	for i, lista := range restaurante.Listas {
		if lista.Nome == newList.Nome {
			restaurante.Listas[i].Itens = newList.Itens
			restaurante.Listas[i].Categoria = newList.Categoria
			listaEncontrada = true
			break
		}
	}

	if !listaEncontrada {
		restaurante.Listas = append(restaurante.Listas, newList)
	}

	_, err = collection.ReplaceOne(context.Background(), filter, restaurante)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving list"})
		return
	}

	if listaEncontrada {
		c.JSON(http.StatusOK, gin.H{"message": "List updated successfully"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "List added successfully"})
	}
}

func getLists(c *gin.Context) {
	collection := DB.Collection("restaurante")
	restname := c.Param("userName")
	filter := bson.D{{"restname", restname}}
	var result Restaurante
	err := collection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
		return
	}

	err = collection.FindOne(context.Background(), filter).Decode(&result)
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

	c.JSON(http.StatusOK, gin.H{"listas": result.Listas})
}

func deleteList(c *gin.Context) {
	userName := c.Param("userName")
	nomeLista := c.Param("nome")

	collection := DB.Collection("restaurante")
	filter := bson.D{{"restname", userName}}
	var restaurante Restaurante

	err := collection.FindOne(context.Background(), filter).Decode(&restaurante)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding restaurant"})
		return
	}

	listaIndex := -1
	for i, lista := range restaurante.Listas {
		if lista.Nome == nomeLista {
			listaIndex = i
			break
		}
	}

	if listaIndex == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "List not found"})
		return
	}

	restaurante.Listas = append(restaurante.Listas[:listaIndex], restaurante.Listas[listaIndex+1:]...)

	_, err = collection.ReplaceOne(context.Background(), filter, restaurante)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "List deleted successfully"})
}
