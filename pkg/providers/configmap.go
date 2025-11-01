package providers

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1alpha1 "github.com/sugaf1204/botnetworkpolicy/api/v1alpha1"
)

type configMapProvider struct {
	client    client.Reader
	namespace string
	name      string
	key       string
}

func (p *configMapProvider) Fetch(ctx context.Context) ([]string, error) {
	var cfg corev1.ConfigMap
	if err := p.client.Get(ctx, client.ObjectKey{Name: p.name, Namespace: p.namespace}, &cfg); err != nil {
		return nil, err
	}
	payload, ok := cfg.Data[p.key]
	if !ok {
		return nil, errMissingKey(p.key)
	}
	return sanitize(v1alpha1.ExtractCIDRs(payload))
}

type errMissingKey string

func (e errMissingKey) Error() string {
	return "configmap missing key: " + string(e)
}
