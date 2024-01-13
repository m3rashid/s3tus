package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/joho/godotenv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/s3store"
)

const UPLOAD_BASE_PATH = "/files/"
const MAX_UPLAOD_SIZE = 500 * 1024 * 1024 // 500 MB

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	store := s3store.New(os.Getenv("AWS_S3_BUCKET_NAME"), s3.New(s3.Options{
		Region: os.Getenv("AWS_REGION"),
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		)),
	}))

	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:              UPLOAD_BASE_PATH,
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})

	if err != nil {
		panic(fmt.Errorf("unable to create handler: %s", err))
	}

	go func() {
		for {
			event := <-handler.CompleteUploads
			fmt.Printf("========== %s ==========\n", event.Upload.ID)
		}
	}()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Use(UPLOAD_BASE_PATH, adaptor.HTTPHandler(http.StripPrefix(UPLOAD_BASE_PATH, handler)))
	app.Listen(":5000")
}
