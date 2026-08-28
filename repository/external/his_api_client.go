package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"agnos-assignment/domain"
)

type hisClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHISClient(baseURL string) domain.HISClient {
	return &hisClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second, 
		},
	}
}

func (c *hisClient) FetchPatientByID(ctx context.Context, id string) (*domain.Patient, error) {
	url := fmt.Sprintf("%s/patient/search/%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil 
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch from HIS, status: %d", resp.StatusCode)
	}

	var patient domain.Patient
	if err := json.NewDecoder(resp.Body).Decode(&patient); err != nil {
		return nil, err
	}

	return &patient, nil
}