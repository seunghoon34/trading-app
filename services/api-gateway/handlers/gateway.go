// file: api-gateway/handlers/gateway.go
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

	// Copy cookies from response back to client
	for _, cookie := range resp.Cookies() {
		c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
	}

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
	accountID, _ := c.Get("account_id")
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
	req.Header.Set("X-Account-ID", accountID.(string))
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
	accountID, _ := c.Get("account_id")
	targetUrl := "http://portfolio:8084" + path
	if c.Request.URL.RawQuery != "" {
		targetUrl += "?" + c.Request.URL.RawQuery
	}

	req, err := http.NewRequest(c.Request.Method, targetUrl, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header = c.Request.Header
	req.Header.Set("X-Account-ID", accountID.(string))

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

func ForwardToInvestmentStrategy(c *gin.Context) {
	path := c.Param("path")
	accountID, _ := c.Get("account_id")
	targetUrl := "http://investment-strategy:8089" + path
	if c.Request.URL.RawQuery != "" {
		targetUrl += "?" + c.Request.URL.RawQuery
	}

	req, err := http.NewRequest(c.Request.Method, targetUrl, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header = c.Request.Header
	req.Header.Set("X-Account-ID", accountID.(string))

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

func ForwardToPaymentService(c *gin.Context) {
	path := c.Param("path")
	accountID, _ := c.Get("account_id")
	targetUrl := "http://payment:8090" + path
	if c.Request.URL.RawQuery != "" {
		targetUrl += "?" + c.Request.URL.RawQuery
	}

	req, err := http.NewRequest(c.Request.Method, targetUrl, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header = c.Request.Header
	req.Header.Set("X-Account-ID", accountID.(string))

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

func ForwardToCrewAIPortfolio(c *gin.Context) {
	path := c.Param("path")
	targetUrl := "http://crewai-portfolio:8000" + path
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

func ForwardToZeusBackend(c *gin.Context) {
	path := c.Param("path")
	accountID, _ := c.Get("account_id")
	targetUrl := "http://zeus-backend-service:3002" + path
	if c.Request.URL.RawQuery != "" {
		targetUrl += "?" + c.Request.URL.RawQuery
	}

	req, err := http.NewRequest(c.Request.Method, targetUrl, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header = c.Request.Header
	req.Header.Set("X-Account-ID", accountID.(string))

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
