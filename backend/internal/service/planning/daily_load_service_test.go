package planning

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/inventory"
	"localis-backend/internal/model/planning"
	"localis-backend/internal/service"
	coreservice "localis-backend/internal/service/core"
	"localis-backend/internal/testutil/mocks"
)

type dailyLoadSvcMocks struct {
	loadMatRepo  *mocks.DailyLoadMaterialRepository
	loadToolRepo *mocks.DailyLoadToolRepository
	loadEPPRepo  *mocks.DailyLoadEPPRepository
	matRepo      *mocks.MaterialRepository
	toolRepo     *mocks.ToolRepository
	eppRepo      *mocks.EPPItemRepository
	auditRepo    *mocks.CorePlanAuditLogRepository
	svc          *DailyLoadService
}

func newDailyLoadService(t *testing.T) *dailyLoadSvcMocks {
	t.Helper()
	loadMatRepo := mocks.NewDailyLoadMaterialRepository(t)
	loadToolRepo := mocks.NewDailyLoadToolRepository(t)
	loadEPPRepo := mocks.NewDailyLoadEPPRepository(t)
	matRepo := mocks.NewMaterialRepository(t)
	toolRepo := mocks.NewToolRepository(t)
	eppRepo := mocks.NewEPPItemRepository(t)
	auditRepo := mocks.NewCorePlanAuditLogRepository(t)
	svc := NewDailyLoadService(loadMatRepo, loadToolRepo, loadEPPRepo, matRepo, toolRepo, eppRepo, coreservice.NewAuditService(auditRepo))
	return &dailyLoadSvcMocks{loadMatRepo: loadMatRepo, loadToolRepo: loadToolRepo, loadEPPRepo: loadEPPRepo, matRepo: matRepo, toolRepo: toolRepo, eppRepo: eppRepo, auditRepo: auditRepo, svc: svc}
}

func TestLoadMaterial_Success(t *testing.T) {
	ctx := context.Background()
	m := newDailyLoadService(t)

	mat := &planning.DailyLoadMaterial{MaterialID: uuid.New(), LoadedQuantity: 5}
	m.matRepo.On("GetByID", ctx, mat.MaterialID).Return(newInventoryMaterial(), nil).Once()
	m.loadMatRepo.On("Create", ctx, mat).Return(nil).Once()
	m.auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := m.svc.LoadMaterial(ctx, mat)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, mat.ID)
	require.Equal(t, float64(0), mat.ReturnedQuantity)
}

func TestLoadMaterial_NotFound(t *testing.T) {
	ctx := context.Background()
	m := newDailyLoadService(t)

	mat := &planning.DailyLoadMaterial{MaterialID: uuid.New()}
	m.matRepo.On("GetByID", ctx, mat.MaterialID).Return(nil, &service.NotFoundError{Resource: "material", ID: mat.MaterialID.String()}).Once()

	err := m.svc.LoadMaterial(ctx, mat)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
	require.Equal(t, "material", nf.Resource)
}

func TestLoadTool_Success(t *testing.T) {
	ctx := context.Background()
	m := newDailyLoadService(t)

	tool := &planning.DailyLoadTool{ToolID: uuid.New(), LoadedQuantity: 2}
	m.toolRepo.On("GetByID", ctx, tool.ToolID).Return(newInventoryTool(), nil).Once()
	m.loadToolRepo.On("Create", ctx, tool).Return(nil).Once()
	m.auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := m.svc.LoadTool(ctx, tool)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, tool.ID)
	require.Equal(t, float64(0), tool.ReturnedQuantity)
}

func TestLoadTool_NotFound(t *testing.T) {
	ctx := context.Background()
	m := newDailyLoadService(t)

	tool := &planning.DailyLoadTool{ToolID: uuid.New()}
	m.toolRepo.On("GetByID", ctx, tool.ToolID).Return(nil, &service.NotFoundError{Resource: "tool", ID: tool.ToolID.String()}).Once()

	err := m.svc.LoadTool(ctx, tool)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
	require.Equal(t, "tool", nf.Resource)
}

func TestLoadEPP_Success(t *testing.T) {
	ctx := context.Background()
	m := newDailyLoadService(t)

	epp := &planning.DailyLoadEPP{EPPID: uuid.New(), LoadedQuantity: 10}
	m.eppRepo.On("GetByID", ctx, epp.EPPID).Return(newInventoryEPP(), nil).Once()
	m.loadEPPRepo.On("Create", ctx, epp).Return(nil).Once()
	m.auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := m.svc.LoadEPP(ctx, epp)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, epp.ID)
	require.Equal(t, float64(0), epp.ReturnedQuantity)
}

func TestLoadEPP_NotFound(t *testing.T) {
	ctx := context.Background()
	m := newDailyLoadService(t)

	epp := &planning.DailyLoadEPP{EPPID: uuid.New()}
	m.eppRepo.On("GetByID", ctx, epp.EPPID).Return(nil, &service.NotFoundError{Resource: "epp_item", ID: epp.EPPID.String()}).Once()

	err := m.svc.LoadEPP(ctx, epp)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
	require.Equal(t, "epp_item", nf.Resource)
}

func TestReturnMaterial_Success(t *testing.T) {
	ctx := context.Background()
	m := newDailyLoadService(t)

	id := uuid.New()
	old := &planning.DailyLoadMaterial{ID: id, ReturnedQuantity: 1}
	m.loadMatRepo.On("GetByID", ctx, id).Return(old, nil).Once()
	m.loadMatRepo.On("Update", ctx, id, old).Return(nil).Once()
	m.auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := m.svc.ReturnMaterial(ctx, id, 3)
	require.NoError(t, err)
	require.Equal(t, float64(3), old.ReturnedQuantity)
}

func TestReturnMaterial_Negative(t *testing.T) {
	ctx := context.Background()
	m := newDailyLoadService(t)

	err := m.svc.ReturnMaterial(ctx, uuid.New(), -1)
	var ve *service.ValidationError
	require.ErrorAs(t, err, &ve)
	require.Equal(t, "returned_quantity", ve.Field)
}

func TestReturnTool_Negative(t *testing.T) {
	ctx := context.Background()
	m := newDailyLoadService(t)

	err := m.svc.ReturnTool(ctx, uuid.New(), -1)
	var ve *service.ValidationError
	require.ErrorAs(t, err, &ve)
}

func TestReturnEPP_Negative(t *testing.T) {
	ctx := context.Background()
	m := newDailyLoadService(t)

	err := m.svc.ReturnEPP(ctx, uuid.New(), -1)
	var ve *service.ValidationError
	require.ErrorAs(t, err, &ve)
}

func TestReturnMaterial_NotFound(t *testing.T) {
	ctx := context.Background()
	m := newDailyLoadService(t)

	id := uuid.New()
	m.loadMatRepo.On("GetByID", ctx, id).Return(nil, &service.NotFoundError{Resource: "daily_load_material", ID: id.String()}).Once()

	err := m.svc.ReturnMaterial(ctx, id, 1)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
}

func TestListMaterialsByPlan_Success(t *testing.T) {
	ctx := context.Background()
	m := newDailyLoadService(t)

	planID := uuid.New()
	expected := []planning.DailyLoadMaterial{{ID: uuid.New(), DailyPlanID: planID}}
	m.loadMatRepo.On("ListByPlan", ctx, planID).Return(expected, nil).Once()

	got, err := m.svc.ListMaterialsByPlan(ctx, planID)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestListMaterialsByPlan_Error(t *testing.T) {
	ctx := context.Background()
	m := newDailyLoadService(t)

	m.loadMatRepo.On("ListByPlan", ctx, uuid.Nil).Return(nil, errors.New("db down")).Once()

	_, err := m.svc.ListMaterialsByPlan(ctx, uuid.Nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list daily load materials")
}

func TestUpdateMaterial_NotFound(t *testing.T) {
	ctx := context.Background()
	m := newDailyLoadService(t)

	id := uuid.New()
	m.loadMatRepo.On("GetByID", ctx, id).Return(nil, &service.NotFoundError{Resource: "daily_load_material", ID: id.String()}).Once()

	err := m.svc.UpdateMaterial(ctx, id, &planning.DailyLoadMaterial{})
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
}

func TestDeleteMaterial_NotFound(t *testing.T) {
	ctx := context.Background()
	m := newDailyLoadService(t)

	id := uuid.New()
	m.loadMatRepo.On("GetByID", ctx, id).Return(nil, &service.NotFoundError{Resource: "daily_load_material", ID: id.String()}).Once()

	err := m.svc.DeleteMaterial(ctx, id)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
}

func newInventoryMaterial() *inventory.Material {
	return &inventory.Material{ID: uuid.New()}
}

func newInventoryTool() *inventory.Tool {
	return &inventory.Tool{ID: uuid.New()}
}

func newInventoryEPP() *inventory.EPPItem {
	return &inventory.EPPItem{ID: uuid.New()}
}
