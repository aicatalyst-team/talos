// Package jwkgen generates JSON Web Keys for signing and HMAC secrets.
// It produces keys in the canonical format expected by talos: each key includes
// a thumbprint-based key ID, the "sig" usage, and the correct algorithm.
//
// This package is safe to import from other services (e.g. backoffice) that
// need to generate keys in the same format talos consumes.
package jwkgen

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"

	"github.com/ory/x/randx"
)

// GenerateSigningKeyJWKS generates a new signing key and returns it as a JWKS
// JSON string containing the private key material, along with the assigned
// key ID. The key includes the "sig" usage, the specified algorithm, and the
// returned key ID.
//
// If alg is empty, EdDSA is used. Supported values: "EdDSA", "RS256".
//
// If kid is empty, the key ID is derived from a SHA-256 thumbprint of the
// public key. If kid is non-empty, it is used verbatim with no format
// validation, matching talos's lookup rules.
func GenerateSigningKeyJWKS(alg, kid string) (jwks, assignedKid string, err error) {
	if alg == "" {
		alg = "EdDSA"
	}

	var (
		rawKey any
		sigAlg jwa.SignatureAlgorithm
	)
	switch alg {
	case "EdDSA":
		_, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return "", "", errors.Wrap(err, "generate Ed25519 key pair")
		}
		rawKey = priv
		sigAlg = jwa.EdDSA()
	case "RS256":
		priv, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return "", "", errors.Wrap(err, "generate RSA key pair")
		}
		rawKey = priv
		sigAlg = jwa.RS256()
	default:
		return "", "", errors.Errorf("unsupported signing algorithm: %s (must be EdDSA or RS256)", alg)
	}

	key, err := jwk.Import(rawKey)
	if err != nil {
		return "", "", errors.Wrap(err, "import key as JWK")
	}

	finalKid := kid
	if finalKid == "" {
		finalKid, err = computeThumbprintKeyID(key)
		if err != nil {
			return "", "", err
		}
	}
	if err := key.Set(jwk.KeyIDKey, finalKid); err != nil {
		return "", "", errors.Wrap(err, "set key ID")
	}
	if err := key.Set(jwk.AlgorithmKey, sigAlg); err != nil {
		return "", "", errors.Wrap(err, "set algorithm")
	}
	if err := key.Set(jwk.KeyUsageKey, "sig"); err != nil {
		return "", "", errors.Wrap(err, "set key usage")
	}

	keySet := jwk.NewSet()
	if err := keySet.AddKey(key); err != nil {
		return "", "", errors.Wrap(err, "add key to set")
	}

	data, err := json.Marshal(keySet)
	if err != nil {
		return "", "", errors.Wrap(err, "marshal JWKS")
	}

	withCreatedAt, err := injectCreatedAt(string(data), time.Now().UTC())
	if err != nil {
		return "", "", err
	}
	return withCreatedAt, finalKid, nil
}

// ExtractSigningKeyID returns the "kid" of the first key in a JWKS JSON
// string. It lets callers recover the kid assigned by GenerateSigningKeyJWKS
// without reimplementing JWKS parsing.
func ExtractSigningKeyID(jwksJSON string) (string, error) {
	var jwks struct {
		Keys []struct {
			Kid string `json:"kid"`
		} `json:"keys"`
	}
	if err := json.Unmarshal([]byte(jwksJSON), &jwks); err != nil {
		return "", errors.Wrap(err, "parse JWKS JSON")
	}
	if len(jwks.Keys) == 0 {
		return "", errors.New("JWKS contains no keys")
	}
	return jwks.Keys[0].Kid, nil
}

// HMACSecret is the JSON-serializable format for HMAC secrets stored in
// project revisions. It wraps the raw secret string with an ID and timestamp.
type HMACSecret struct {
	ID        string `json:"id"`
	Secret    string `json:"secret"`
	CreatedAt string `json:"created_at"`
}

// GenerateHMACSecret generates a new HMAC secret with a UUID identifier and
// a 32-character alphanumeric random secret.
func GenerateHMACSecret() HMACSecret {
	return HMACSecret{
		ID:        uuid.Must(uuid.NewV4()).String(),
		Secret:    randx.MustString(32, randx.AlphaLowerNum),
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// computeThumbprintKeyID returns a SHA-256 thumbprint of the public key,
// encoded as a URL-safe base64 string. This matches the format used by talos's
// CLI key generation.
func computeThumbprintKeyID(key jwk.Key) (string, error) {
	pubKey, err := key.PublicKey()
	if err != nil {
		return "", errors.Wrap(err, "extract public key for thumbprint")
	}

	jsonBytes, err := json.Marshal(pubKey)
	if err != nil {
		return "", errors.Wrap(err, "marshal key for thumbprint")
	}

	hash := sha256.Sum256(jsonBytes)
	return base64.URLEncoding.EncodeToString(hash[:]), nil
}

// injectCreatedAt adds a created_at timestamp to each key in a JWKS JSON string.
func injectCreatedAt(jwksJSON string, createdAt time.Time) (string, error) {
	var jwks struct {
		Keys []map[string]any `json:"keys"`
	}
	if err := json.Unmarshal([]byte(jwksJSON), &jwks); err != nil {
		return "", errors.Wrap(err, "parse JWKS JSON")
	}

	ts := createdAt.UTC().Format(time.RFC3339)
	for i := range jwks.Keys {
		jwks.Keys[i]["created_at"] = ts
	}

	raw, err := json.Marshal(jwks)
	if err != nil {
		return "", errors.Wrap(err, "marshal JWKS JSON")
	}
	return string(raw), nil
}
