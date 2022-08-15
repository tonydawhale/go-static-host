package server

import (
	"bytes"
	"net/http"
	"os"
	"log"
	"strings"

	"go-static-host/s3utils"
	"go-static-host/mongoutils"
	"go-static-host/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Init() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLFiles("public/index.html")

	r.GET("/", getIndex)
	r.POST("/", postIndex)
	r.POST("/api/upload", postApi)
	r.GET("/favicon.ico", func (c *gin.Context) { c.Status(http.StatusAccepted) } )
	r.GET("/:id", getFile)

	log.Println("Server Listening to Port " + os.Getenv("PORT"))
	r.Run()
}

func getIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func postIndex(c *gin.Context) {
	form, _ := c.MultipartForm()

	file := form.File["file-input"][0]

	if (!(strings.HasPrefix(file.Header.Get("content-type"), "image/") && 
		strings.HasPrefix(file.Header.Get("content-type"), "video/"))) {
		log.Println("Someone has attempted to upload an invalid file type")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid file type",
		})
		return;
	}

	extractedFile, err := file.Open()
	
	if (err != nil) {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error Uploading File",
		})
		return
	}

	newName := uuid.New().String()
	shortId := util.GenerateId(7)

	_, err = s3utils.UploadS3Object(newName, extractedFile)

	if (err != nil) {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error Uploading File",
		})
		return
	}

	_, mongoErr := mongoutils.CreateItemMetaData(newName, shortId, file.Header.Get("content-type"))

	if (mongoErr != nil) {
		log.Println(mongoErr.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error Uploading File",
		})
		return
	}

	c.Redirect(http.StatusFound, "/" + shortId)
}

func postApi(c *gin.Context) {
	form, _ := c.MultipartForm()

	file := form.File["file-input"][0]

	if (!(strings.HasPrefix(file.Header.Get("content-type"), "image/") && 
		strings.HasPrefix(file.Header.Get("content-type"), "video/"))) {
		log.Println("Someone has attempted to upload an invalid file type")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid file type",
		})
		return;
	}

	extractedFile, err := file.Open()
	
	if (err != nil) {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error Uploading File",
		})
		return
	}

	newName := uuid.New().String()
	shortId := util.GenerateId(7)

	_, err = s3utils.UploadS3Object(newName, extractedFile)

	if (err != nil) {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error Uploading File",
		})
		return
	}

	_, mongoErr := mongoutils.CreateItemMetaData(newName, shortId, file.Header.Get("content-type"))

	if (mongoErr != nil) {
		log.Println(mongoErr.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error Uploading File",
		})
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"url": os.Getenv("URI") + "/" + shortId,
	})
}

func getFile(c *gin.Context) {
	data, dataErr := mongoutils.FetchItemMetaData(c.Param("id"))
	if (dataErr != nil) {
		log.Println(dataErr.Error())
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Item Not Found",
		})
		return
	}

	rawObject, err := s3utils.GetS3Object(data.Uuid)

	if (err != nil) {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error Downloading File",
		})
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(rawObject.Body)

	c.Data(
		http.StatusOK,
		data.ContentType,  
		buf.Bytes(), 
	)
}