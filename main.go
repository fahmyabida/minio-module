package main

import (
	configEnv "github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"minio/config"
	"minio/minio"
	"net/http"
	"time"
)

func main() {
	Echo()
}

func Echo(){
	err := configEnv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	minioClient := config.MinioClient()
	minioEngine := minio.NewMinioHelper(minioClient)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/bucket", func(c echo.Context) error {
		nameBucket := c.FormValue("name")
		err := minioEngine.AddBucket(nameBucket)
		if err != nil {
			return c.JSON(http.StatusBadRequest, `{"message":"`+err.Error()+`"}`)
		}
		return c.JSON(200, `{"message":"success"}`)
	})
	e.POST("/upload", func(c echo.Context) error {
		fileHeader, err := c.FormFile("upload")
		if err != nil {
			return c.JSON(http.StatusBadRequest, `{"message":"`+err.Error()+`"}`)
		}
		file, err := fileHeader.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, `{"message":"`+err.Error()+`"}`)
		}
		res, err := minioEngine.UploadFileWithFile("mymusic", "photo/file2.jpeg", file)
		if err != nil {
			return c.JSON(http.StatusBadRequest, `{"message":"`+err.Error()+`"}`)
		}
		return c.JSON(200, res)
	})

	e.GET("/path", func(c echo.Context) error {
		duration, _ := time.ParseDuration("10m")
		res, err := minioEngine.GetFile("mymusic", "photo/file2.jpeg", duration)
		if err != nil {
			return c.JSON(http.StatusBadRequest, `{"message":"`+err.Error()+`"}`)
		}
		return c.JSON(200, res)
	})

	e.GET("/file", func(c echo.Context) error {
		duration, _ := time.ParseDuration("10m")
		res, err := minioEngine.DownloadFile("mymusic", "photo/file2.jpeg", duration)
		if err != nil {
			return c.JSON(http.StatusBadRequest, `{"message":"`+err.Error()+`"}`)
		}
		return c.JSON(200, res)
	})

	e.Logger.Fatal(e.Start(":1234"))
}