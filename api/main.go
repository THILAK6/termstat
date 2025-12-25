package main

import (
	"bytes"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/v1/ping", func(c echo.Context) error {
		var body bytes.Buffer
		if _, err := body.ReadFrom(c.Request().Body); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		}

		toURL := "https://api.us-east.aws.tinybird.co/v0/events?name=cli_pings"
		req, _ := http.NewRequest("POST", toURL, &body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+os.Getenv("TINYBIRD_TOKEN"))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode >= 400 {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "failed to forward request"})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "received"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}