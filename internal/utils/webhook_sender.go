package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"medods_test_task/internal/config"
)

func SendWarningToWebhook(userID uuid.UUID, ip, newIp, userAgent string) (err error) {
	payload := map[string]interface{}{
		"user_id":    userID.String(),
		"ip":         ip,
		"new_ip":     newIp,
		"user_agent": userAgent,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", config.Load().WebHook, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	log.Printf("Webhook sent. Status: %s", resp.Status)
	return nil
}
