package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/s3store"
)

const UPLOAD_BASE_PATH = "/files/"
const MAX_UPLAOD_SIZE = 500 * 1024 * 1024 // 500 MB

func getS3Client() *s3.Client {
	return s3.New(s3.Options{
		Region: os.Getenv("AWS_REGION"),
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		)),
	})
}

func getTusHandler() (*tusd.Handler, error) {
	store := s3store.New(os.Getenv("AWS_S3_BUCKET_NAME"), getS3Client())
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	return tusd.NewHandler(tusd.Config{
		BasePath:              UPLOAD_BASE_PATH,
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})
}

func getPreSignedUrl(imageId string) (string, error) {
	s3Client := getS3Client()
	presigner := s3.NewPresignClient(s3Client)
	req, err := presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
		Key:    &imageId,
	})

	if err != nil {
		return "", err
	}

	return req.URL, nil
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	handler, err := getTusHandler()
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			event := <-handler.CompleteUploads
			fmt.Printf("========== %s ==========\n", event.Upload.ID)
		}
	}()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowCredentials: true,
	}))

	app.Use(UPLOAD_BASE_PATH, adaptor.HTTPHandler(http.StripPrefix(UPLOAD_BASE_PATH, handler)))

	app.Get("/auth", func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:     "auth",
			Value:    "test",
			HTTPOnly: true,
		})
		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/file/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		// get the auth
		authCookie := c.Cookies("auth", "")
		fmt.Println("======= Auth cookie =============> ", authCookie, "  <==")

		if authCookie == "" || authCookie != "test" {
			return c.SendStatus(fiber.StatusNotFound)
		}

		url, err := getPreSignedUrl(id)
		if err != nil {
			return err
		}
		return c.Redirect(url, http.StatusTemporaryRedirect)
	})

	log.Println("Server is running ...")
	app.Listen(":5000")
}
