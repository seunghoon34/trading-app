package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ForwardToAuthService(c *gin.Context) {
	path := c.Param("path")
	targetUrl := "http://user-management:8080" + path

	// Use the same method as original request
	req, err := http.NewRequest(c.Request.Method, targetUrl, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy headers from original request
	req.Header = c.Request.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	// Read response and send back to client
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Forward the response back to client
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func ForwardToMarketDataService(c *gin.Context) {
	path := c.Param("path")
	targetUrl := "http://market-data:8082" + path
	if c.Request.URL.RawQuery != "" {
		targetUrl += "?" + c.Request.URL.RawQuery
	}

	// Use the same method as original request
	req, err := http.NewRequest(c.Request.Method, targetUrl, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy headers from original request
	req.Header = c.Request.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	// Read response and send back to client
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Forward the response back to client
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func ForwardToTradingService(c *gin.Context) {
	path := c.Param("path")
	targetUrl := "http://trading-engine:8083" + path
	if c.Request.URL.RawQuery != "" {
		targetUrl += "?" + c.Request.URL.RawQuery
	}

	req, err := http.NewRequest(c.Request.Method, targetUrl, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header = c.Request.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func ForwardToPortfolioService(c *gin.Context) {
	path := c.Param("path")
	targetUrl := "http://portfolio:8084/portfolio" + path
	if c.Request.URL.RawQuery != "" {
		targetUrl += "?" + c.Request.URL.RawQuery
	}

	req, err := http.NewRequest(c.Request.Method, targetUrl, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header = c.Request.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}
