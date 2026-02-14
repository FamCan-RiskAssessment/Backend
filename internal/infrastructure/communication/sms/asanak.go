package sms

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/exception"
)

const (
	asanakBaseURL = "https://sms.asanak.ir/webservice/v2rest"
)

// AsanakClient represents the Asanak SMS client
type AsanakClient struct {
	Username   string
	Password   string
	Source     string
	HTTPClient *http.Client
}

// NewAsanakClient creates a new Asanak SMS client
func NewAsanakClient(username, password, source string) *AsanakClient {
	return &AsanakClient{
		Username: username,
		Password: password,
		Source:   source,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendSMSRequest represents the request for sending simple SMS
type SendSMSRequest struct {
	Source          string `json:"source"`
	Destination     string `json:"destination"`
	Message         string `json:"message"`
	SendToBlacklist int    `json:"send_to_blacklist,omitempty"`
}

// SendSMSResponse represents the response from SMS sending
type SendSMSResponse []int64

// AsanakAPIResponse represents the actual API response structure
type AsanakAPIResponse struct {
	Meta struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"meta"`
	Data []int64 `json:"data"`
}

// ErrorResponse represents error response from API
type ErrorResponse struct {
	HTTPStatusCode int    `json:"http_status_code"`
	MetaStatusCode int    `json:"meta_status_code"`
	Message        string `json:"message"`
	Description    string `json:"description"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("SMS API Error: %s (HTTP: %d, Meta: %d)", e.Message, e.HTTPStatusCode, e.MetaStatusCode)
}

// SendSMS sends a simple SMS to one or more recipients
func (c *AsanakClient) SendSMS(req SendSMSRequest) (SendSMSResponse, error) {
	// Prepare form data
	data := url.Values{}
	data.Set("username", c.Username)
	data.Set("password", c.Password)
	data.Set("source", req.Source)
	data.Set("destination", req.Destination)
	data.Set("message", req.Message)
	if req.SendToBlacklist != 0 {
		data.Set("send_to_blacklist", strconv.Itoa(req.SendToBlacklist))
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", asanakBaseURL+"/sendsms", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Accept", "application/json")

	// Send request
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			return nil, exception.NewServiceUnavailableError(
				exception.ServiceSMS,
				exception.ReasonNetwork,
				"failed to reach sms provider",
				err,
			)
		}
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp.StatusCode, body)
	}

	// Parse response
	var apiResp AsanakAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check meta status
	if apiResp.Meta.Status != 200 {
		return nil, fmt.Errorf("API error: %s (status: %d)", apiResp.Meta.Message, apiResp.Meta.Status)
	}

	return apiResp.Data, nil
}

// handleErrorResponse handles API error responses
func (c *AsanakClient) handleErrorResponse(statusCode int, body []byte) error {
	if statusCode == http.StatusForbidden {
		return exception.NewServiceUnavailableError(
			exception.ServiceSMS,
			exception.ReasonNetwork,
			"sms provider blocked the request",
			fmt.Errorf("API error (HTTP %d): %s", statusCode, string(body)),
		)
	}

	// Try to parse as JSON error response
	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err == nil {
		errResp.HTTPStatusCode = statusCode
		return errResp
	}

	// If not JSON, return the raw response
	return fmt.Errorf("API error (HTTP %d): %s", statusCode, string(body))
}

// AsanakSMSService implements the SmsService interface using Asanak API
type AsanakSMSService struct {
	client    *AsanakClient
	templates *bootstrap.SMSTemplates
}

// NewAsanakSMSService creates a new Asanak SMS service
func NewAsanakSMSService(
	providerConfig *bootstrap.SMSGateway,
	smsTemplates *bootstrap.SMSTemplates,
) *AsanakSMSService {
	client := NewAsanakClient(providerConfig.Username, providerConfig.Password, providerConfig.Source)
	return &AsanakSMSService{
		client:    client,
		templates: smsTemplates,
	}
}

// SendOTP sends an OTP SMS using Asanak API
func (s *AsanakSMSService) SendOTP(receptor, token string) error {
	// Create OTP message - you can customize this format
	message := fmt.Sprintf("به فمکن خوش آمدید\n\nکد تایید شما: %s", token)

	req := SendSMSRequest{
		Source:          s.client.Source,
		Destination:     receptor,
		Message:         message,
		SendToBlacklist: 1,
	}

	_, err := s.client.SendSMS(req)
	return err
}
