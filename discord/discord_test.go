package discord

import (
	"testing"
)

// TestSendVersionUpdateNotification tests the message formatting
func TestSendVersionUpdateNotification(t *testing.T) {
	// Store original function
	original := SendVersionUpdateNotification
	defer func() {
		SendVersionUpdateNotification = original
	}()

	// Create mock function
	var called bool
	var gotWebhook, gotName, gotOld, gotNew string

	SendVersionUpdateNotification = NotificationFunc(func(webhookURL, modpackName, oldVersion, newVersion string) error {
		called = true
		gotWebhook = webhookURL
		gotName = modpackName
		gotOld = oldVersion
		gotNew = newVersion
		return nil
	})

	// Test the function
	err := SendVersionUpdateNotification("test-webhook", "Test Pack", "1.0", "2.0")
	if err != nil {
		t.Errorf("SendVersionUpdateNotification() error = %v", err)
	}

	// Verify the mock was called with correct parameters
	if !called {
		t.Error("Mock function was not called")
	}
	if gotWebhook != "test-webhook" {
		t.Errorf("Got webhook = %v, want %v", gotWebhook, "test-webhook")
	}
	if gotName != "Test Pack" {
		t.Errorf("Got modpack name = %v, want %v", gotName, "Test Pack")
	}
	if gotOld != "1.0" {
		t.Errorf("Got old version = %v, want %v", gotOld, "1.0")
	}
	if gotNew != "2.0" {
		t.Errorf("Got new version = %v, want %v", gotNew, "2.0")
	}
}
