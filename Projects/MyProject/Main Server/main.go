package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

func main() {
	router := gin.Default()

	// CORS middleware configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	router.Any("/auth/*path", func(c *gin.Context) {
		proxyRequest(c, "https://localhost:5050")
	})

	router.Any("/crud/*path", func(c *gin.Context) {
		proxyRequest(c, "https://localhost:6000")
	})

	err := router.RunTLS(
		":5000",
		"D:\\Projects\\Employee\\MyProject\\Certificate\\cert.crt", 
		"D:\\Projects\\Employee\\MyProject\\Certificate\\key.pem",
	)
	
	if err != nil {
		log.Fatal("Failed to run server with HTTPS: ", err)
	}
}

// Function to proxy the request to another service
func proxyRequest(c *gin.Context, target string) {
	client := &http.Client{
		Timeout: 10 * 1000000000,
	}

	req, err := http.NewRequest(c.Request.Method, target+c.Param("path"), c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	for k, v := range c.Request.Header {
		if k != "Connection" && k != "Transfer-Encoding" {
			req.Header[k] = v
		}
	}

	// Make the HTTP request to the target service
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reach microservice"})
		return
	}
	defer resp.Body.Close()

	// Set the status code from the response
	c.Status(resp.StatusCode)

	for k, v := range resp.Header {
		if k == "Access-Control-Allow-Origin" || k == "Authorization" {
			c.Header(k, v[0])
		}
	}

	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward response"})
	}
}
