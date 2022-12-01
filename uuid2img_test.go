package uuid2img

import (
	"testing"
)

func TestImageGeneration(test *testing.T) {
	success := GenerateFile("123e4567-e89b-12d3-a456-426614174000", "test_uuid.png")
	if !success {
		test.Errorf("Failed to generate image")
	}
}