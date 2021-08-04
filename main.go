package main

import (
	"os"

	"github.com/gin-gonic/gin"
	a "github.com/petrolax/shorten-api/abbreviation"
	"github.com/petrolax/shorten-api/handler"
)

func main() {

	file, err := os.OpenFile("data/test.json", os.O_CREATE, 0644)
	if err != nil {
		panic("Can't open file")
	}
	defer file.Close()

	storage := a.NewAbbreviationStorage(file)
	h := handler.NewHandler(storage)

	router := gin.Default()

	router.POST("/", h.NewShorten)
	router.DELETE("/delete", h.RemoveAllShorten)
	router.GET("/list/:page", h.GetList)

	router.GET("/:shorturl", h.RedirectFromShorten)
	router.GET("/:shorturl/original", h.GetOriginal)
	router.DELETE(":shorturl/delete", h.RemoveShorten)

	router.Run(":8080")
}
