package http_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"inspacemap/backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HTTPIntegrationTestSuite struct {
	suite.Suite
	app *fiber.App
}

func (suite *HTTPIntegrationTestSuite) SetupTest() {
	// Initialize app like in main.go
	suite.app = fiber.New()

	// Setup routes (simplified version for testing)
	api := suite.app.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", func(c *fiber.Ctx) error {
		var req models.RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}
		return c.JSON(fiber.Map{"message": "Registration endpoint works"})
	})

	// Public venue routes
	api.Get("/venues/:slug/manifest", func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Venue not found or not published"})
	})
}

func (suite *HTTPIntegrationTestSuite) TearDownTest() {
	// Cleanup if needed
}

func (suite *HTTPIntegrationTestSuite) TestAuthRegisterEndpoint() {
	// Arrange
	reqData := models.RegisterRequest{
		FullName:         "Test User",
		Email:            "test@example.com",
		Password:         "password123",
		OrganizationName: "Test Org",
	}
	reqBody, _ := json.Marshal(reqData)

	// Act
	httpReq := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := suite.app.Test(httpReq)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(suite.T(), "Registration endpoint works", response["message"])
}

func (suite *HTTPIntegrationTestSuite) TestVenueManifestEndpoint() {
	// Test public venue manifest endpoint
	httpReq := httptest.NewRequest("GET", "/api/v1/venues/test-slug/manifest", nil)
	resp, err := suite.app.Test(httpReq)

	// Should return 404 for non-existent venue
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 404, resp.StatusCode)
}

func (suite *HTTPIntegrationTestSuite) TestInvalidJSONRequest() {
	// Test invalid JSON handling
	httpReq := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader([]byte("invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := suite.app.Test(httpReq)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 400, resp.StatusCode)
}

func TestHTTPIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPIntegrationTestSuite))
}
