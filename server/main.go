package main

import (
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/echo/v4"
)

func main() {
	// Create an AWS session
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2"),
	}))

	// Create an S3 client1
	svc := s3.New(sess)

	e := echo.New()

	  // Add custom middleware to handle CORS
	  e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
		  c.Response().Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5173")
		  c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		  c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, X-Requested-With, remember-me")
	
		  if c.Request().Method == http.MethodOptions {
			return c.NoContent(http.StatusOK)
		  }
	
		  return next(c)
		}
	  })

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/request-presigned-url", func(c echo.Context) error {
		type Object struct {
			Key string `json:"key"`
		}

		var objects []Object
		if err := c.Bind(&objects); err != nil {
			return c.String(400, "bad request")
		}

		bucket := "presigned-poc"

		// 1 minute duration for PUT URL
		expiresPut := time.Minute * 1

		// 5 minutes duration for GET URL
		expiresGet := time.Minute * 5

		// I want to make urls that are like this [ {key: objectName, put: url, get: url}]
		urls := make([]map[string]string, 0)
		for _, obj := range objects {
			objectName := obj.Key
		
			// Get the URL that allows only a PUT operation
			req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(objectName),
			})
			urlStr, err := req.Presign(expiresPut)
			if err != nil {
				return c.String(500, "internal server error")
			}
		
			// Get the URL that allows only a GET operation
			req, _ = svc.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(objectName),
			})
			urlStr2, err := req.Presign(expiresGet)
			if err != nil {
				return c.String(500, "internal server error")
			}
		
			urls = append(urls, map[string]string{
				"key": objectName,
				"put": urlStr,
				"get": urlStr2,
			})
		}

		return c.JSON(200, urls)
	})

	e.Logger.Fatal(e.Start(":1323"))
}