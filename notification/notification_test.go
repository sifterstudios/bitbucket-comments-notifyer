package notification

import (
	"testing"
)

func TestSendNotification(t *testing.T) {
	testCases := []struct {
		headline string
		message  string
		expected error
	}{
		{
			headline: "Test Notification",
			message:  "This is a test message",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.headline, func(t *testing.T) {
			err := SendNotification(tc.headline, tc.message)

			if err != tc.expected {
				t.Errorf("Expected error: %v, got: %v", tc.expected, err)
			}
		})
	}
}
