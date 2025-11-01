package providers

import (
	"testing"
)

func TestNavigateField(t *testing.T) {
	payload := map[string]any{
		"outer": map[string]any{
			"inner": []any{"10.0.0.0/24", "2001:db8::/32"},
		},
	}

	value, err := navigateField(payload, "outer.inner")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	arr, ok := value.([]any)
	if !ok || len(arr) != 2 {
		t.Fatalf("unexpected value: %#v", value)
	}
}

func TestInterpretCIDRs(t *testing.T) {
	items := []any{"10.0.0.0/24", "10.0.1.0/24"}
	cidrs, err := interpretCIDRs(items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cidrs) != 2 {
		t.Fatalf("expected 2 cidrs, got %d", len(cidrs))
	}

	if _, err := interpretCIDRs([]any{"10.0.0.0/24", 5}); err == nil {
		t.Fatalf("expected error for non-string entry")
	}
}
