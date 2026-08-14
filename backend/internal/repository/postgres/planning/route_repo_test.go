package postgresplanning

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/planning"
	"localis-backend/internal/model/shared"
)

func TestRouteRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewRouteRepo(pool)

	route := &planning.Route{
		ID:          uuid.New(),
		Type:        "7290d8d5-ab3a-4465-8bb9-e328e5363afb",
		Date:        time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour),
		DailyPlanID: uuidPtr(uuid.MustParse("70000000-0000-0000-0000-000000000001")),
		Status:      shared.RouteStatusScheduled,
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, route)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, route.ID)

	got, err := repo.GetByID(ctx, route.ID)
	require.NoError(t, err)
	require.Equal(t, route.Type, got.Type)
	require.Equal(t, route.Status, got.Status)

	result, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 1)
	found := false
	for _, r := range result.Items {
		if r.ID == route.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	byPlan, err := repo.ListByPlan(ctx, uuid.MustParse("70000000-0000-0000-0000-000000000001"))
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(byPlan), 1)

	route.Status = shared.RouteStatusInProgress
	err = repo.Update(ctx, route.ID, route)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, route.ID)
	require.NoError(t, err)
	require.Equal(t, shared.RouteStatusInProgress, updated.Status)

	err = repo.Delete(ctx, route.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, route.ID)
	require.Error(t, err)
}

func TestRouteRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewRouteRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}
