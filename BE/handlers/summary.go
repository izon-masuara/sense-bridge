package handlers

import (
	"context"
	"fmt"
	"net/http"
	"server_lidm/apis"
	"server_lidm/db"
	"server_lidm/helpers"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Summary(c *gin.Context) {
	code := c.Param("code")

	collection := db.GetCollection(db.Client, "sensebridge", code)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"id": code}

	var result bson.M
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("find document err: %s", err.Error()))
		return
	}

	sections, ok := result["sections"].(primitive.M)
	if !ok {
		c.String(http.StatusInternalServerError, "sections tidak ditemukan atau formatnya salah")
		return
	}

	var textBuilder strings.Builder

	for _, section := range sections {
		sectionMap, ok := section.(primitive.M)
		if !ok {
			continue
		}

		audioTexts, ok := sectionMap["audio_text"].(primitive.A)
		if ok {
			for _, audioText := range audioTexts {
				textBuilder.WriteString(audioText.(string) + " ")
			}
		}

		imageText, ok := sectionMap["image_text"].(string)
		if ok {
			textBuilder.WriteString(fmt.Sprintf("Linkupan pembahasan di atas haris berdasarkan topik %s ini \n\n", imageText))
		}
	}

	text := textBuilder.String()
	textSummary, err := apis.GenerateSummary(text)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = helpers.CreatePdf(textSummary, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	update := bson.M{
		"$set": bson.M{
			"link_pdf": "./summary/" + code + ".pdf",
			"summary":  textSummary,
		},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("update document err: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"text":    text,
	})
}
