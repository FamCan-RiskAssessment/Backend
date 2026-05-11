package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"
)

type VerificationClientImpl struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewVerificationClient(config *bootstrap.VerificationAPI) *VerificationClientImpl {
	return &VerificationClientImpl{
		baseURL:    config.BaseURL,
		apiKey:     config.APIKey,
		httpClient: &http.Client{},
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

type postalCodeInfoRequest struct {
	PostalCode string `json:"postalCode"`
}

type postalCodeInfoAPIResponse struct {
	Province     *string `json:"province"`
	City         *string `json:"city"`
	Town         *string `json:"town"`
	District     *string `json:"district"`
	Street       *string `json:"street"`
	Street2      *string `json:"street2"`
	Number       *string `json:"number"`
	Floor        *string `json:"floor"`
	SideFloor    *string `json:"sideFloor"`
	BuildingName *string `json:"buildingName"`
	Description  *string `json:"description"`
}

type resultDataOfPostalCodeInfoRes struct {
	Data    *postalCodeInfoAPIResponse `json:"data"`
	Success bool                       `json:"success"`
	Code    int                        `json:"code"`
	Error   string                     `json:"error"`
	Message string                     `json:"message"`
}

func (c *VerificationClientImpl) GetAddressByPostalCode(postalCode string) (*formdto.PostalCodeInfoResponse, error) {
	reqBody := postalCodeInfoRequest{
		PostalCode: postalCode,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api/sw1/PostalCodeInfo", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResp resultDataOfPostalCodeInfoRes
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		if apiResp.Error != "" {
			return nil, fmt.Errorf("postal code API error: %s", apiResp.Error)
		}
		return nil, fmt.Errorf("postal code API call failed with code: %d", apiResp.Code)
	}

	if apiResp.Data == nil {
		return nil, fmt.Errorf("no data returned from postal code API")
	}

	return &formdto.PostalCodeInfoResponse{
		Province:     apiResp.Data.Province,
		City:         apiResp.Data.City,
		Town:         apiResp.Data.Town,
		District:     apiResp.Data.District,
		Street:       apiResp.Data.Street,
		Street2:      apiResp.Data.Street2,
		Number:       apiResp.Data.Number,
		Floor:        apiResp.Data.Floor,
		SideFloor:    apiResp.Data.SideFloor,
		BuildingName: apiResp.Data.BuildingName,
		Description:  apiResp.Data.Description,
	}, nil
}
