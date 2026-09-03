package operations

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestResolvePartner_NilPool(t *testing.T) {
	m := newVisitService(t)
	// vPool is nil by default in test setup
	partnerID, err := m.service.ResolvePartner(context.Background(), uuid.New())
	require.NoError(t, err)
	require.Nil(t, partnerID)
}

func TestResolveSLA_NilPool(t *testing.T) {
	m := newVisitService(t)
	slaID, err := m.service.ResolveSLA(context.Background(), uuid.New(), "electrical")
	require.NoError(t, err)
	require.Nil(t, slaID)
}
