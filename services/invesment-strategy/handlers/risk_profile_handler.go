package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seunghoon34/trading-app/services/investment-strategy/internal/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

// Add these to your handlers/investment_handler.go file

// RiskProfile represents a user's risk assessment
type RiskProfile struct {
	ID                   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AccountID            string             `json:"account_id" bson:"alpaca_id"`
	RiskTolerance        string             `json:"risk_tolerance" bson:"risk_tolerance" binding:"required,oneof=conservative moderate aggressive"`
	InvestmentTimeline   string             `json:"investment_timeline" bson:"investment_timeline" binding:"required,oneof=short_term medium_term long_term"`
	FinancialGoals       []string           `json:"financial_goals" bson:"financial_goals" binding:"required,dive,oneof=retirement wealth_building income_generation capital_preservation education home_purchase"`
	AgeBracket           string             `json:"age_bracket" bson:"age_bracket" binding:"required,oneof=18-25 26-35 36-45 46-55 56-65 65+"`
	AnnualIncomeBracket  string             `json:"annual_income_bracket" bson:"annual_income_bracket" binding:"required,oneof=0-25000 25000-50000 50000-75000 75000-100000 100000-150000 150000+"`
	InvestmentExperience string             `json:"investment_experience" bson:"investment_experience" binding:"required,oneof=beginner intermediate advanced"`
	RiskCapacity         string             `json:"risk_capacity" bson:"risk_capacity" binding:"required,oneof=low medium high"`
	CreatedAt            time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at" bson:"updated_at"`
}

// RiskProfileRequest represents the request body for creating/updating risk profiles
type RiskProfileRequest struct {
	RiskTolerance        string   `json:"risk_tolerance" binding:"required,oneof=conservative moderate aggressive"`
	InvestmentTimeline   string   `json:"investment_timeline" binding:"required,oneof=short_term medium_term long_term"`
	FinancialGoals       []string `json:"financial_goals" binding:"required,dive,oneof=retirement wealth_building income_generation capital_preservation education home_purchase"`
	AgeBracket           string   `json:"age_bracket" binding:"required,oneof=18-25 26-35 36-45 46-55 56-65 65+"`
	AnnualIncomeBracket  string   `json:"annual_income_bracket" binding:"required,oneof=0-25000 25000-50000 50000-75000 75000-100000 100000-150000 150000+"`
	InvestmentExperience string   `json:"investment_experience" binding:"required,oneof=beginner intermediate advanced"`
	RiskCapacity         string   `json:"risk_capacity" binding:"required,oneof=low medium high"`
}

// CreateRiskProfile creates a new risk profile (prevents duplicates)
func CreateRiskProfile(c *gin.Context) {
	// Get account_id from header
	accountID := c.GetHeader("X-Account-ID")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Account-ID header is required"})
		return
	}

	var req RiskProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if risk profile already exists
	var existingProfile RiskProfile
	err := mongo.RiskProfileCollection.FindOne(ctx, bson.M{"alpaca_id": accountID}).Decode(&existingProfile)

	if err == nil {
		// Risk profile exists
		c.JSON(http.StatusConflict, gin.H{
			"error":                    "Risk profile already exists for this user",
			"existing_risk_profile_id": existingProfile.ID,
			"message":                  "Use update endpoint to modify existing risk profile",
		})
		return
	} else if err != mongodriver.ErrNoDocuments {
		// Some other error occurred
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing risk profile: " + err.Error()})
		return
	}

	// Create new risk profile (no existing profile found)
	riskProfile := RiskProfile{
		AccountID:            accountID,
		RiskTolerance:        req.RiskTolerance,
		InvestmentTimeline:   req.InvestmentTimeline,
		FinancialGoals:       req.FinancialGoals,
		AgeBracket:           req.AgeBracket,
		AnnualIncomeBracket:  req.AnnualIncomeBracket,
		InvestmentExperience: req.InvestmentExperience,
		RiskCapacity:         req.RiskCapacity,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	result, err := mongo.RiskProfileCollection.InsertOne(ctx, riskProfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create risk profile"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":         "risk profile created successfully",
		"account_id":      accountID,
		"risk_profile_id": result.InsertedID,
	})
}

// UpdateRiskProfile updates an existing risk profile
func UpdateRiskProfile(c *gin.Context) {
	accountID := c.GetHeader("X-Account-ID")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Account-ID header is required"})
		return
	}

	var req RiskProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Update risk profile (keep original created_at, update updated_at)
	update := bson.M{
		"$set": bson.M{
			"risk_tolerance":        req.RiskTolerance,
			"investment_timeline":   req.InvestmentTimeline,
			"financial_goals":       req.FinancialGoals,
			"age_bracket":           req.AgeBracket,
			"annual_income_bracket": req.AnnualIncomeBracket,
			"investment_experience": req.InvestmentExperience,
			"risk_capacity":         req.RiskCapacity,
			"updated_at":            time.Now(),
		},
	}

	result, err := mongo.RiskProfileCollection.UpdateOne(
		ctx,
		bson.M{"alpaca_id": accountID},
		update,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update risk profile"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Risk profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "risk profile updated successfully",
		"account_id": accountID,
	})
}

// GetRiskProfile retrieves a risk profile by account_id
func GetRiskProfile(c *gin.Context) {
	accountID := c.GetHeader("X-Account-ID")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Account-ID header is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var riskProfile RiskProfile
	err := mongo.RiskProfileCollection.FindOne(ctx, bson.M{"alpaca_id": accountID}).Decode(&riskProfile)

	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Risk profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch risk profile"})
		return
	}

	c.JSON(http.StatusOK, riskProfile)
}
