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

// NewHISClient เป็น Constructor สำหรับสร้าง HTTP Client ไปต่อกับ Hospital A
func NewHISClient(baseURL string) domain.HISClient {
	return &hisClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // ตั้ง Timeout ป้องกัน API ปลายทางค้างแล้วระบบเราร่มตาม
		},
	}
}

func (c *hisClient) FetchPatientByID(ctx context.Context, id string) (*domain.Patient, error) {
	// 1. สร้าง URL ตามโจทย์: https://hospital-a.api.co.th/patient/search/{id}
	url := fmt.Sprintf("%s/patient/search/%s", c.baseURL, id)

	// 2. สร้าง HTTP Request พร้อมแนบ Context (เผื่อฝั่งผู้ใช้กดยกเลิกกลางคัน)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// 3. สั่งยิง API
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 4. จัดการ Response ตาม HTTP Status Code
	if resp.StatusCode == http.StatusNotFound {
		// ถ้า HIS API ตอบ 404 (หาคนไข้ไม่เจอ) ไม่ต้องคืนค่า Error แต่ให้คืน nil กลับไป
		return nil, nil 
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch from HIS, status: %d", resp.StatusCode)
	}

	// 5. แปลง JSON ที่ได้กลับมาเป็น Struct domain.Patient ของเรา
	var patient domain.Patient
	if err := json.NewDecoder(resp.Body).Decode(&patient); err != nil {
		return nil, err
	}

	return &patient, nil
}