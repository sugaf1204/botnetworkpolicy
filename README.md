# Bot Network Policy Operator

The Bot Network Policy Operator manages Kubernetes `NetworkPolicy` objects that allow ingress/egress traffic exclusively from known bot IP ranges published by popular cloud platforms and custom sources.  It periodically refreshes provider allowlists and keeps a deterministic NetworkPolicy in sync for each `BotNetworkPolicy` custom resource.

## Features

- Built-in providers for Google, AWS, and GitHub bot/metadata endpoints.
- ConfigMap provider to supply custom CIDR ranges managed within the cluster.
- JSON endpoint provider that retrieves CIDRs from an arbitrary HTTP endpoint and extracts them via a JSON field path.
- Deterministic NetworkPolicy generation with optional ingress/egress toggles and custom CIDR overrides.
- Periodic re-sync with configurable intervals per resource.

## Custom Resource Overview

```yaml
apiVersion: bot.networking.dev/v1alpha1
kind: BotNetworkPolicy
metadata:
  name: example
spec:
  podSelector:
    matchLabels:
      app: web
  syncPeriod: 30m
  providers:
    - name: google
    - name: aws
    - name: github
    - name: configMap
      configMap:
        name: extra-bot-ips
        key: cidrs
    - name: jsonEndpoint
      jsonEndpoint:
        url: https://example.com/bots.json
        fieldPath: data.cidrs
        headers:
          Accept: application/json
        headerSecretRefs:
          - name: Authorization
            secretKeyRef:
              name: bot-endpoint-token
              key: token
  customCidrs:
    - 192.0.2.0/24
```

The operator will create or update a `NetworkPolicy` named `<metadata.name>-allow-bots` (or a custom name specified via the `bot.networking.dev/networkpolicy-name` annotation) in the same namespace. The generated policy contains ingress rules (and optional egress rules) limited to the merged set of CIDRs.

## Getting Started

1. Deploy the CRD and controller manifests (to be generated via controller-tools) to your cluster.
2. Apply a `BotNetworkPolicy` resource in the target namespace.
3. Confirm that a `NetworkPolicy` with the `botnetworkpolicy.bot.networking.dev/owner` label appears and contains the expected IP blocks.

See [`docs/development.md`](docs/development.md) for development workflows and testing guidance.
