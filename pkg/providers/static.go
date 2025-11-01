package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type staticHTTPProvider struct {
	client   *http.Client
	url      string
	selector func(map[string]any) ([]string, error)
}

func (p *staticHTTPProvider) Fetch(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	cidrs, err := p.selector(payload)
	if err != nil {
		return nil, err
	}
	return sanitize(cidrs)
}

const (
	defaultGoogleEndpoint = "https://www.gstatic.com/ipranges/goog.json"
	defaultAWSEndpoint    = "https://ip-ranges.amazonaws.com/ip-ranges.json"
	defaultGitHubEndpoint = "https://api.github.com/meta"
)

func googleSelector(data map[string]any) ([]string, error) {
	prefixesRaw, ok := data["prefixes"].([]any)
	if !ok {
		return nil, fmt.Errorf("missing prefixes")
	}
	results := make([]string, 0)
	for _, prefix := range prefixesRaw {
		item, _ := prefix.(map[string]any)
		if item == nil {
			continue
		}
		if ipv4, ok := item["ipv4Prefix"].(string); ok {
			if value := strings.TrimSpace(ipv4); value != "" {
				results = append(results, value)
			}
		}
		if ipv6, ok := item["ipv6Prefix"].(string); ok {
			if value := strings.TrimSpace(ipv6); value != "" {
				results = append(results, value)
			}
		}
	}
	return results, nil
}

func awsSelector(data map[string]any) ([]string, error) {
	prefixesRaw, ok := data["prefixes"].([]any)
	if !ok {
		return nil, fmt.Errorf("missing prefixes")
	}
	results := make([]string, 0)
	for _, prefix := range prefixesRaw {
		item, _ := prefix.(map[string]any)
		if item == nil {
			continue
		}
		if service, _ := item["service"].(string); service != "AMAZON" && service != "AMAZON_CONNECT" {
			continue
		}
		// GitHub documents that hook delivery traffic originates from the GLOBAL and us-east-1
		// regions. Restricting the ranges we ingest keeps the resulting NetworkPolicies focused on
		// the documented bot endpoints instead of the full AWS address space.
		if region, _ := item["region"].(string); region == "GLOBAL" || region == "us-east-1" {
			if cidr, ok := item["ip_prefix"].(string); ok {
				if value := strings.TrimSpace(cidr); value != "" {
					results = append(results, value)
				}
			}
		}
	}
	return results, nil
}

func githubSelector(data map[string]any) ([]string, error) {
	hooks, ok := data["hooks"].([]any)
	if !ok {
		return nil, fmt.Errorf("missing hooks field")
	}
	results := make([]string, 0, len(hooks))
	for _, hook := range hooks {
		if cidr, ok := hook.(string); ok {
			if value := strings.TrimSpace(cidr); value != "" {
				results = append(results, value)
			}
		}
	}
	return results, nil
}
