package commands

import (
	"fmt"
	"testing"
)

func TestGetUsageTemplate(t *testing.T) {
	v := getUsageTemplate("test")
	expected := fmt.Sprintf(defaultTemplate, "test")
	if len(v) != len(expected) {
		t.Errorf("getUsageTemplate returned a wrong usage template length: got %d want %d", len(v), len(expected))
	}
}
