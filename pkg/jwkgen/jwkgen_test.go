package jwkgen_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/talos/pkg/jwkgen"
)

func TestGenerateSigningKeyJWKS_DefaultKidIsThumbprint(t *testing.T) {
	t.Parallel()

	raw, kid, err := jwkgen.GenerateSigningKeyJWKS("", "")
	require.NoError(t, err)
	assert.NotEmpty(t, kid, "empty kid argument must fall back to thumbprint kid")

	var parsed struct {
		Keys []map[string]any `json:"keys"`
	}
	require.NoError(t, json.Unmarshal([]byte(raw), &parsed))
	require.Len(t, parsed.Keys, 1)
	assert.Equal(t, "EdDSA", parsed.Keys[0]["alg"])
	assert.Equal(t, "sig", parsed.Keys[0]["use"])
	assert.Equal(t, kid, parsed.Keys[0]["kid"])
}

func TestGenerateSigningKeyJWKS_ExplicitKidIsPreserved(t *testing.T) {
	t.Parallel()

	raw, kid, err := jwkgen.GenerateSigningKeyJWKS("", "my-custom-kid")
	require.NoError(t, err)
	assert.Equal(t, "my-custom-kid", kid)

	var parsed struct {
		Keys []map[string]any `json:"keys"`
	}
	require.NoError(t, json.Unmarshal([]byte(raw), &parsed))
	require.Len(t, parsed.Keys, 1)
	assert.Equal(t, "my-custom-kid", parsed.Keys[0]["kid"])
}

func TestGenerateSigningKeyJWKS_RS256(t *testing.T) {
	t.Parallel()

	raw, kid, err := jwkgen.GenerateSigningKeyJWKS("RS256", "")
	require.NoError(t, err)
	assert.NotEmpty(t, kid)

	var parsed struct {
		Keys []map[string]any `json:"keys"`
	}
	require.NoError(t, json.Unmarshal([]byte(raw), &parsed))
	require.Len(t, parsed.Keys, 1)
	assert.Equal(t, "RS256", parsed.Keys[0]["alg"])
	assert.Equal(t, kid, parsed.Keys[0]["kid"])
}

func TestGenerateSigningKeyJWKS_UnsupportedAlgorithmErrors(t *testing.T) {
	t.Parallel()

	_, _, err := jwkgen.GenerateSigningKeyJWKS("ES256", "")
	require.Error(t, err)
}

func TestExtractSigningKeyID(t *testing.T) {
	t.Parallel()

	t.Run("returns the thumbprint kid", func(t *testing.T) {
		t.Parallel()

		raw, _, err := jwkgen.GenerateSigningKeyJWKS("", "")
		require.NoError(t, err)

		kid, err := jwkgen.ExtractSigningKeyID(raw)
		require.NoError(t, err)
		assert.NotEmpty(t, kid)
	})

	t.Run("returns the explicit kid", func(t *testing.T) {
		t.Parallel()

		raw, _, err := jwkgen.GenerateSigningKeyJWKS("", "explicit-kid")
		require.NoError(t, err)

		kid, err := jwkgen.ExtractSigningKeyID(raw)
		require.NoError(t, err)
		assert.Equal(t, "explicit-kid", kid)
	})

	t.Run("errors on empty JWKS", func(t *testing.T) {
		t.Parallel()

		_, err := jwkgen.ExtractSigningKeyID(`{"keys":[]}`)
		require.Error(t, err)
	})

	t.Run("errors on invalid JSON", func(t *testing.T) {
		t.Parallel()

		_, err := jwkgen.ExtractSigningKeyID(`not-json`)
		require.Error(t, err)
	})
}
