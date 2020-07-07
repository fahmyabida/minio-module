package config

import (
	"fmt"
	"github.com/minio/minio-go/v6"
	"log"
	"os"
	"strconv"
	"time"
)

func MinioClient() *minio.Client{
	sslMode, err := strconv.ParseBool(os.Getenv("MINIO_SSL_MODE"))
	if err != nil {
		fmt.Println(err)
		time.Sleep(7*time.Second); os.Exit(1)
	}
	minioClient, err := minio.New(
		os.Getenv("MINIO_CLIENT_URL"),
		os.Getenv("MINIO_ACCESS_KEY_ID"),
		os.Getenv("MINIO_SECRET_ACCESS_KEY"),
		sslMode,
		)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to minio!")
	return minioClient
}
