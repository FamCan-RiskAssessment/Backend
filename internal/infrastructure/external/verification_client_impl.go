package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
)

type VerificationClientImpl struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewVerificationClient(config *bootstrap.VerificationAPI) *VerificationClientImpl {
	return &VerificationClientImpl{
		baseURL: config.BaseURL,
		apiKey:  config.APIKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type shahkarRequest struct {
	NationalCode string `json:"nationalCode"`
	Mobile       string `json:"mobile"`
	IsCompany    bool   `json:"isCompany"`
}

type shahkarResponse struct {
	Data    bool   `json:"data"`
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (c *VerificationClientImpl) VerifyPhoneAndSSN(phone string, ssn string) (bool, error) {
	reqBody := shahkarRequest{
		NationalCode: ssn,
		Mobile:       phone,
		IsCompany:    false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return false, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api/sw1/Shahkar", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var shahkarResp shahkarResponse
	if err := json.NewDecoder(resp.Body).Decode(&shahkarResp); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	if !shahkarResp.Success {
		if shahkarResp.Error != "" {
			return false, fmt.Errorf("shahkar API error: %s", shahkarResp.Error)
		}
		return false, fmt.Errorf("shahkar API call failed with code: %d", shahkarResp.Code)
	}

	return shahkarResp.Data, nil
}
