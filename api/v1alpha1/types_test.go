package v1alpha1

import "testing"

func TestExtractCIDRs(t *testing.T) {
	payload := "10.0.0.0/24\n10.0.1.0/24, 2001:db8::/32"
	cidrs := ExtractCIDRs(payload)
	if len(cidrs) != 3 {
		t.Fatalf("expected 3 cidrs, got %d", len(cidrs))
	}
}
