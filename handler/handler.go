package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Storage interface {
	CreateNewShortenUrl(url string) (string, error)
	GetLengthenUrl(shorturl string) (string, error)
	GetListOfAbbreviations(page int) (map[string]string, error)
	DeleteAllShortenUrl() error
	DeleteShortenUrl(shorturl string) error
}

type Handler struct {
	storage Storage
	log     *log.Logger
}

func NewHandler(storage Storage) *Handler {
	return &Handler{
		storage: storage,
		log:     log.New(os.Stdout, "", log.LstdFlags),
	}
}

func serverResponse(c *gin.Context, code int, message string, result interface{}, err string) {
	c.JSON(code, map[string]interface{}{
		"StatusCode": code,
		"Message":    message,
		"Result":     result,
		"Error":      err,
	})
}

func (h *Handler) NewShorten(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		h.log.Println("NewShorten:Error: URL is empty")
		serverResponse(c, http.StatusBadRequest, "", "", "URL is empty")
		return
	}

	// Не лучший способ проверить существование сайта, но самый простой
	resp, err := http.Get(url)
	if err != nil {
		h.log.Printf("NewShorten:Error: URL - '%s' is like a simple string", url)
		serverResponse(c, http.StatusBadRequest, fmt.Sprintf("NewShorten:Error: URL - '%s' is like a simple string", url), "", err.Error())
		return
	}
	resp.Body.Close()

	shorturl, err := h.storage.CreateNewShortenUrl(url)
	if err != nil {
		h.log.Printf("NewShorten:Error: %s\n", err.Error())
		serverResponse(c, http.StatusNotAcceptable, "", "", err.Error())
		return
	}

	h.log.Printf("NewShorten: param - %s\n", url)
	serverResponse(c, http.StatusOK, "New shorten url created", shorturl, "")
}

func (h *Handler) RemoveAllShorten(c *gin.Context) {
	err := h.storage.DeleteAllShortenUrl()
	if err != nil {
		h.log.Printf("RemoveAllShorten:Error: %s\n", err.Error())
		serverResponse(c, http.StatusBadRequest, "", "", err.Error())
		return
	}
	serverResponse(c, http.StatusOK, "Every abbreviations was delete", "", "")
	h.log.Println("RemoveAllShorten: Every abbreviations was delete")
}

func (h *Handler) GetList(c *gin.Context) {
	page := c.Param("page")

	npage, err := strconv.Atoi(page)
	if err != nil {
		h.log.Printf("GetList:Error: %s\n", err.Error())
		serverResponse(c, http.StatusBadRequest, "Page number isn't number type", "", err.Error())
		return
	}
	list, err := h.storage.GetListOfAbbreviations(npage)
	if err != nil {
		h.log.Printf("GetList:Error: %s\n", err.Error())
		serverResponse(c, http.StatusBadRequest, "", "", err.Error())
		return
	}
	h.log.Printf("GetList: page %d available\n", npage)
	serverResponse(c, http.StatusOK, "List of abbreviations", list, "")
}

func (h *Handler) RedirectFromShorten(c *gin.Context) {
	shorturl := c.Param("shorturl")
	mainurl, err := h.storage.GetLengthenUrl(shorturl)
	if err != nil {
		h.log.Printf("RedirectFromShorten:Error: %s\n", err.Error())
		serverResponse(c, http.StatusBadRequest, "", "", err.Error())
		return
	}
	h.log.Printf("RedirectFromShorten: Redirect from shorten: %s to url: %s\n", shorturl, mainurl)
	c.Redirect(http.StatusMovedPermanently, mainurl)
}

func (h *Handler) GetOriginal(c *gin.Context) {
	shorturl := c.Param("shorturl")
	mainurl, err := h.storage.GetLengthenUrl(shorturl)
	if err != nil {
		h.log.Printf("GetOriginal:Error: %s\n", err.Error())
		serverResponse(c, http.StatusBadRequest, "", "", err.Error())
		return
	}
	h.log.Printf("GetOriginal: Short URL %s: %s\n", shorturl, mainurl)
	serverResponse(c, http.StatusOK, "Main url of shorten "+shorturl, mainurl, "")
}

func (h *Handler) RemoveShorten(c *gin.Context) {
	shorturl := c.Param("shorturl")
	err := h.storage.DeleteShortenUrl(shorturl)
	if err != nil {
		h.log.Printf("RemoveShorten:Error: %s\n", err.Error())
		serverResponse(c, http.StatusBadRequest, "", "", err.Error())
		return
	}
	serverResponse(c, http.StatusOK, "Short url "+shorturl+" was delete", "", "")
	h.log.Printf("RemoveShorten: Short URL - %s was delete\n", shorturl)
}
