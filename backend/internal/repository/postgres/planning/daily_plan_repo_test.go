package postgresplanning

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/planning"
)

func TestDailyPlanAssignmentRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewDailyPlanAssignmentRepo(pool)

	planID := uuid.MustParse("70000000-0000-0000-0000-000000000001")
	techID := uuid.MustParse("60000000-0000-0000-0000-000000000001")
	roleID := uuid.MustParse("2757c3b2-012a-451e-b5d4-47f11bd65437")

	assignment := &planning.DailyPlanAssignment{
		ID:          uuid.New(),
		DailyPlanID: planID,
		TechID:      techID,
		RoleID:      roleID,
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, assignment)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, assignment.ID)

	got, err := repo.GetByID(ctx, assignment.ID)
	require.NoError(t, err)
	require.Equal(t, assignment.TechID, got.TechID)
	require.Equal(t, assignment.RoleID, got.RoleID)

	result, err := repo.ListByPlan(ctx, planID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(result), 1)
	found := false
	for _, a := range result {
		if a.ID == assignment.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	listResult, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, listResult.Total, 1)

	err = repo.Delete(ctx, assignment.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, assignment.ID)
	require.Error(t, err)
}

func TestDailyPlanAssignmentRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewDailyPlanAssignmentRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestDailyPlanRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewDailyPlanRepo(pool)

	plan := &planning.DailyPlan{
		ID:        uuid.New(),
		Name:      "Test Plan " + uuid.New().String(),
		Date:      time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour).UTC(),
		Notes:     strPtr("Test notes"),
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, plan)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, plan.ID)

	got, err := repo.GetByID(ctx, plan.ID)
	require.NoError(t, err)
	require.Equal(t, plan.Name, got.Name)
	require.Equal(t, plan.Date, got.Date)

	result, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 1)
	found := false
	for _, p := range result.Items {
		if p.ID == plan.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	plan.Name = "Updated " + plan.Name
	err = repo.Update(ctx, plan.ID, plan)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, plan.ID)
	require.NoError(t, err)
	require.Equal(t, plan.Name, updated.Name)

	err = repo.Delete(ctx, plan.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, plan.ID)
	require.Error(t, err)
}

func TestDailyPlanRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewDailyPlanRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}
