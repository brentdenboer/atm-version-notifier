package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Message represents a Discord webhook message
type Message struct {
	Content string `json:"content"`
}

// SendVersionUpdateNotification sends a version update notification to Discord
func SendVersionUpdateNotification(webhookURL, modpackName, oldVersion, newVersion string) error {
	if webhookURL == "" {
		return fmt.Errorf("empty webhook URL")
	}

	message := Message{
		Content: fmt.Sprintf(
			`🎮 **%s Update Detected!**

The modpack version has changed:
- Old version: %q
- New version: %q`,
			modpackName, oldVersion, newVersion),
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message %+v: %w", message, err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("unexpected status %d: %q", resp.StatusCode, string(body))
	}

	return nil
}
