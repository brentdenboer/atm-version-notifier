package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Message represents a Discord webhook message
type Message struct {
	Content string `json:"content"`
}

// SendVersionUpdateNotification sends a notification about a version update to Discord
func SendVersionUpdateNotification(webhookURL, modpackName, oldVersion, newVersion string) error {
	message := Message{
		Content: fmt.Sprintf("🎮 **%s Update Detected!**\n\nThe modpack version has changed:\n- Old version: `%s`\n- New version: `%s`",
			modpackName, oldVersion, newVersion),
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal discord message: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send discord webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("discord webhook returned non-2xx status code: %d", resp.StatusCode)
	}

	return nil
}
