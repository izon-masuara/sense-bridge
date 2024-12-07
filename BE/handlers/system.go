package handlers

import (
	"context"
	"net/http"
	"server_lidm/db"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GenerateCode(c *gin.Context) {
	uniqueCode := uuid.New().String()

	collection := db.GetCollection(db.Client, "sensebridge", uniqueCode)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	document := map[string]interface{}{
		"id": uniqueCode,
		"sections": map[string]interface{}{
			"1": map[string]interface{}{
				"audio_text": []string{},
				"image_text": "",
			},
		},
		"link_pdf": "",
		"summary":  "",
	}

	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": uniqueCode,
		"sec":  1,
	})
}

func FindCode(c *gin.Context) {
	code := c.Param("code")

	filter := bson.M{"id": code}
	collection := db.GetCollection(db.Client, "sensebridge", code)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var result bson.M
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sections, ok := result["sections"].(bson.M)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid sections format"})
		return
	}

	var lastSectionID string
	var maxSectionID int
	for sectionID := range sections {
		intID, err := strconv.Atoi(sectionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid section ID format"})
			return
		}
		if intID > maxSectionID {
			maxSectionID = intID
			lastSectionID = sectionID
		}
	}

	response := map[string]interface{}{
		"id":           result["id"],
		"sections":     result["sections"],
		"link_pdf":     result["link_pdf"],
		"summary":      result["summary"],
		"last_section": lastSectionID,
	}

	c.JSON(http.StatusOK, response)
}
