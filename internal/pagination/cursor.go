// Package pagination provides cursor-based pagination utilities.
package pagination

import (
	"github.com/cockroachdb/errors"

	keysetpagination "github.com/ory/x/pagination/keysetpagination_v2"
)

// Cursor represents a pagination cursor containing the last seen ID
// for stable, efficient cursor-based pagination using primary key
type Cursor struct {
	// ID is the unique identifier of the last item (primary key)
	ID string `json:"id"`
	// NID is the network ID that issued this cursor.
	// Tokens are valid only in the same network context.
	NID string `json:"nid"`
}

// EncodeCursor encrypts an (id, nid) pair into an opaque page token that
// callers can pass back to resume a list. Returns an empty token when id is
// empty (terminal page) and an error when nid is empty.
func EncodeCursor(key [32]byte, id, nid string) (string, error) {
	if id == "" {
		return "", nil
	}
	if nid == "" {
		return "", errors.Errorf("pagination network ID not configured")
	}

	token := keysetpagination.NewPageToken(
		keysetpagination.Column{Name: "id", Value: id},
		keysetpagination.Column{Name: "nid", Value: nid},
	)
	return token.Encrypt([][32]byte{key}), nil
}

// DecodeCursor decrypts a page token by trying each key in order, which
// supports graceful key rotation: put the current key first and keep
// retired keys until their tokens age out. Returns (nil, nil) for an empty
// token, the decoded cursor if any key succeeds, or an error if all fail.
func DecodeCursor(keys [][32]byte, pageToken string) (*Cursor, error) {
	if pageToken == "" {
		return nil, nil
	}

	if len(keys) == 0 {
		return nil, errors.Errorf("pagination encryption keys not configured")
	}

	parsed, err := keysetpagination.ParsePageToken(keys, pageToken)
	if err != nil {
		return nil, errors.Wrap(err, "decode pagination token")
	}

	var c Cursor
	for _, col := range parsed.Columns() {
		s, _ := col.Value.(string)
		switch col.Name {
		case "id":
			c.ID = s
		case "nid":
			c.NID = s
		}
	}
	return &c, nil
}

// ValidatePageSize clamps a caller-supplied page size to the server bounds:
// non-positive values become the default (50) and values above the cap become
// the maximum (1000). Callers should use the returned value rather than the
// original input.
func ValidatePageSize(pageSize int32) int32 {
	const (
		defaultPageSize = 50
		maxPageSize     = 1000
	)

	if pageSize <= 0 {
		return defaultPageSize
	}
	if pageSize > maxPageSize {
		return maxPageSize
	}

	return pageSize
}

// reviewed - @aeneasr - 2026-03-26
