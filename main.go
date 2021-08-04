package main

import (
	"os"

	"github.com/gin-gonic/gin"
	as "github.com/petrolax/shorten-api/AbbreviationStorage"
	"github.com/petrolax/shorten-api/handler"
)

func main() {

	file, err := os.OpenFile("test.json", os.O_CREATE, 0644)
	if err != nil {
		panic("Can't open file")
	}
	defer file.Close()

	storage := as.NewAbbreviationStorage(file)
	h := handler.NewHandler(storage)

	router := gin.Default()

	router.POST("/", h.NewShorten)               // Получить из URL-параметра url= и вернуть сокращённый вариант, example - 'localhost:8080/wda12eda'
	router.DELETE("/delete", h.RemoveAllShorten) // Удалить весь список сокращений
	router.GET("/list/:page", h.GetList)         // Передать в json список первых 10-ти сокращений

	router.GET("/:shorturl", h.RedirectFromShorten)    // Перейти на зашифрованный сайт
	router.GET("/:shorturl/lengthen", h.GetLengthen)   // Получить из json сокращённый URL и вернуть полноценный вариант
	router.DELETE(":shorturl/delete", h.RemoveShorten) // Удалить конкретную ссылку из списка сокращений

	router.Run(":8080")
}
