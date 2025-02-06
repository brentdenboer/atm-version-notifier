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

// NotificationFunc is a function type for sending version update notifications
type NotificationFunc func(webhookURL, modpackName, oldVersion, newVersion string) error

// defaultSendVersionUpdateNotification is the default implementation for sending notifications
func defaultSendVersionUpdateNotification(webhookURL, modpackName, oldVersion, newVersion string) error {
	message := Message{
		Content: fmt.Sprintf("🎮 **%s Update Detected!**\n\nThe modpack version has changed:\n- Old version: `%s`\n- New version: `%s`",
			modpackName, oldVersion, newVersion),
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal discord message: %w", err)
	}

	// Create a client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create request
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send discord webhook: %w", err)
	}
	defer resp.Body.Close()

	// Read response body for error details
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("discord webhook returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// SendVersionUpdateNotification is the current function used for sending notifications
// It can be replaced with a mock for testing
var SendVersionUpdateNotification NotificationFunc = defaultSendVersionUpdateNotification
