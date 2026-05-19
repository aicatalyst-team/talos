---
title: Secret management
---

# Secret management

Talos uses a single HMAC secret to generate and verify API key checksums. The same secret also
derives the pagination cursor encryption key through domain-separated HMAC-SHA256, so operators
configure one secret family.

## Configuration

```yaml
secrets:
  hmac:
    current: "a-32-byte-or-longer-hmac-secret!"
    retired:
      - "previous-hmac-secret-that-was-rotated"
```

## Secret types

| Secret                 | Purpose                                                                | Required           |
| ---------------------- | ---------------------------------------------------------------------- | ------------------ |
| `secrets.hmac.current` | HMAC secret for API key checksums and pagination cursor key derivation | Yes (min 32 chars) |

## Secret rotation

1. Add the current secret to the `retired` array.
2. Set a new `current` secret.
3. Restart Talos (or wait for config hot-reload).

```yaml
secrets:
  hmac:
    current: "new-hmac-secret-32-chars-minimum!"
    retired:
      - "old-hmac-secret-that-was-previously-current"
```

During verification, Talos tries the current secret first, then each retired secret in order. This
ensures existing API keys remain valid and previously issued pagination tokens continue to decode
after rotation.

## Environment variables

```bash
export TALOS_SECRETS_HMAC_CURRENT="my-hmac-secret-32-chars-minimum"
export TALOS_SECRETS_HMAC_RETIRED="old-hmac-secret-1,old-hmac-secret-2"
```
