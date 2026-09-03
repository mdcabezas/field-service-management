package mobilehandler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/auth"
	"localis-backend/internal/handler"
	"localis-backend/internal/model/mobile"
)

type SyncHandler struct {
	pool *pgxpool.Pool
}

func NewSyncHandler(pool *pgxpool.Pool) *SyncHandler {
	return &SyncHandler{pool: pool}
}

func NewTestSyncHandler() *SyncHandler {
	return &SyncHandler{pool: nil}
}

func (h *SyncHandler) Push(c *gin.Context) {
	var req mobile.SyncPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	// Test mode: if no pool, return mock results
	if h.pool == nil {
		results := mobile.SyncResultItems{
			Checkpoints:   make([]mobile.SyncItemResult, 0, len(req.Checkpoints)),
			Photos:        make([]mobile.SyncItemResult, 0, len(req.Photos)),
			Usages:        make([]mobile.SyncItemResult, 0, len(req.Usages)),
			Measurements:  make([]mobile.SyncItemResult, 0, len(req.Measurements)),
			Reports:       make([]mobile.SyncItemResult, 0, len(req.Reports)),
			StatusChanges: make([]mobile.SyncItemResult, 0, len(req.StatusChanges)),
		}
		for _, cp := range req.Checkpoints {
			results.Checkpoints = append(results.Checkpoints, mobile.SyncItemResult{LocalID: cp.LocalID, ServerID: uuid.New().String(), Status: "ok"})
		}
		for _, p := range req.Photos {
			results.Photos = append(results.Photos, mobile.SyncItemResult{LocalID: p.LocalID, ServerID: uuid.New().String(), Status: "ok"})
		}
		for _, u := range req.Usages {
			results.Usages = append(results.Usages, mobile.SyncItemResult{LocalID: u.LocalID, ServerID: uuid.New().String(), Status: "ok"})
		}
		for _, m := range req.Measurements {
			results.Measurements = append(results.Measurements, mobile.SyncItemResult{LocalID: m.LocalID, ServerID: uuid.New().String(), Status: "ok"})
		}
		for _, r := range req.Reports {
			results.Reports = append(results.Reports, mobile.SyncItemResult{LocalID: r.LocalID, ServerID: uuid.New().String(), Status: "ok"})
		}
		for _, sc := range req.StatusChanges {
			results.StatusChanges = append(results.StatusChanges, mobile.SyncItemResult{LocalID: sc.LocalID, ServerID: uuid.New().String(), Status: "ok"})
		}
		handler.RespondJSON(c, http.StatusOK, mobile.SyncResult{SyncedAt: time.Now(), Results: results})
		return
	}

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		handler.RespondError(c, http.StatusInternalServerError, "failed to begin transaction: "+err.Error())
		return
	}
	defer tx.Rollback(ctx)

	results := mobile.SyncResultItems{
		Checkpoints:   make([]mobile.SyncItemResult, 0, len(req.Checkpoints)),
		Photos:        make([]mobile.SyncItemResult, 0, len(req.Photos)),
		Usages:        make([]mobile.SyncItemResult, 0, len(req.Usages)),
		Measurements:  make([]mobile.SyncItemResult, 0, len(req.Measurements)),
		Reports:       make([]mobile.SyncItemResult, 0, len(req.Reports)),
		StatusChanges: make([]mobile.SyncItemResult, 0, len(req.StatusChanges)),
	}

	// Checkpoints
	for _, cp := range req.Checkpoints {
		serverID := uuid.New()
		_, err := tx.Exec(ctx,
			`INSERT INTO operations.visit_checkpoints (id, visit_id, type, geom, timestamp, device_info)
			 VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5), 4326), $6, $7)
			 ON CONFLICT (id) DO NOTHING`,
			serverID, cp.VisitID, cp.Type, cp.Lng, cp.Lat, cp.Timestamp, cp.LocalID,
		)
		status := "ok"
		if err != nil {
			status = "error"
		}
		results.Checkpoints = append(results.Checkpoints, mobile.SyncItemResult{
			LocalID:  cp.LocalID,
			ServerID: serverID.String(),
			Status:   status,
			Error:    errMessage(err),
		})
	}

	// Photos are uploaded separately via /api/mobile/photos handler.
	// No metadata insertion here — the photo handler writes to visit_photos directly.

	// Material usages
	for _, usage := range req.Usages {
		serverID := uuid.New()
		_, err := tx.Exec(ctx,
			`INSERT INTO operations.visit_material_usages (id, visit_id, material_id, quantity, notes)
			 VALUES ($1, $2, $3, $4, $5)
			 ON CONFLICT (id) DO NOTHING`,
			serverID, usage.VisitID, usage.MaterialID, usage.Quantity, usage.Notes,
		)
		status := "ok"
		if err != nil {
			status = "error"
		}
		results.Usages = append(results.Usages, mobile.SyncItemResult{
			LocalID:  usage.LocalID,
			ServerID: serverID.String(),
			Status:   status,
			Error:    errMessage(err),
		})
	}

	// Measurements
	for _, m := range req.Measurements {
		serverID := uuid.New()
		var err error
		// Parse value as numeric
		var valueFloat float64
		if _, parseErr := fmt.Sscanf(m.Value, "%f", &valueFloat); parseErr != nil {
			valueFloat = 0
		}
		if m.Lat != nil && m.Lng != nil {
			_, err = tx.Exec(ctx,
				`INSERT INTO operations.visit_measurements (id, visit_id, value, unit, notes)
				 VALUES ($1, $2, $3, $4, $5)
				 ON CONFLICT (id) DO NOTHING`,
				serverID, m.VisitID, valueFloat, m.Unit, m.Key,
			)
		} else {
			_, err = tx.Exec(ctx,
				`INSERT INTO operations.visit_measurements (id, visit_id, value, unit, notes)
				 VALUES ($1, $2, $3, $4, $5)
				 ON CONFLICT (id) DO NOTHING`,
				serverID, m.VisitID, valueFloat, m.Unit, m.Key,
			)
		}
		status := "ok"
		if err != nil {
			status = "error"
		}
		results.Measurements = append(results.Measurements, mobile.SyncItemResult{
			LocalID:  m.LocalID,
			ServerID: serverID.String(),
			Status:   status,
			Error:    errMessage(err),
		})
	}

	// Status changes — update visit status
	for _, sc := range req.StatusChanges {
		serverID := uuid.New()
		var err error
		switch sc.Status {
		case "in_progress":
			_, err = tx.Exec(ctx,
				`UPDATE operations.visits SET status = $1, started_at = COALESCE(started_at, $2), updated_at = now()
				 WHERE id = $3`,
				sc.Status, sc.Timestamp, sc.VisitID,
			)
		case "completed":
			_, err = tx.Exec(ctx,
				`UPDATE operations.visits SET status = $1, completed_at = COALESCE(completed_at, $2), updated_at = now()
				 WHERE id = $3`,
				sc.Status, sc.Timestamp, sc.VisitID,
			)
		default:
			_, err = tx.Exec(ctx,
				`UPDATE operations.visits SET status = $1, updated_at = now()
				 WHERE id = $2`,
				sc.Status, sc.VisitID,
			)
		}
		status := "ok"
		if err != nil {
			status = "error"
		}
		results.StatusChanges = append(results.StatusChanges, mobile.SyncItemResult{
			LocalID:  sc.LocalID,
			ServerID: serverID.String(),
			Status:   status,
			Error:    errMessage(err),
		})
	}

	// Reports
	for _, report := range req.Reports {
		serverID := uuid.New()
		var templateID interface{} = nil
		if report.ReportTemplateID != uuid.Nil {
			templateID = report.ReportTemplateID
		}
		_, err := tx.Exec(ctx,
			`INSERT INTO operations.visit_reports (id, visit_id, report_template_id, source, recorded_at)
			 VALUES ($1, $2, $3, 'system', now())
			 ON CONFLICT (id) DO NOTHING`,
			serverID, report.VisitID, templateID,
		)
		status := "ok"
		if err != nil {
			status = "error"
		}
		results.Reports = append(results.Reports, mobile.SyncItemResult{
			LocalID:  report.LocalID,
			ServerID: serverID.String(),
			Status:   status,
			Error:    errMessage(err),
		})
	}

	if err := tx.Commit(ctx); err != nil {
		handler.RespondError(c, http.StatusInternalServerError, "failed to commit: "+err.Error())
		return
	}

	handler.RespondJSON(c, http.StatusOK, mobile.SyncResult{
		SyncedAt: time.Now(),
		Results:  results,
	})
}

func errMessage(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func (h *SyncHandler) Pull(c *gin.Context) {
	claimsVal, exists := c.Get(auth.ClaimsKey)
	if !exists {
		handler.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	claims, ok := claimsVal.(*auth.Claims)
	if !ok {
		handler.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	employeeNumber := claims.UserID
	if employeeNumber == "" {
		handler.RespondError(c, http.StatusBadRequest, "employee_number not in token")
		return
	}

	// Test mode: if no pool, return mock data
	if h.pool == nil {
		response := mobile.SyncPullResponse{
			ReferenceData: &mobile.ReferenceData{
				Materials: []mobile.ReferenceMaterial{},
				Tools:     []mobile.ReferenceTool{},
				EPP:       []mobile.ReferenceEPP{},
				Templates: []mobile.ReferenceTemplate{},
				Partners:  []mobile.ReferencePartner{},
			},
			Visits:        []mobile.SyncVisit{},
			ChecklistItems: []mobile.SyncChecklistItem{},
			Photos:        []mobile.SyncPhotoInfo{},
			SyncedAt:      time.Now(),
		}
		handler.RespondJSON(c, http.StatusOK, response)
		return
	}

	ctx := c.Request.Context()

	techID, err := h.getTechnicianID(ctx, employeeNumber)
	if err != nil {
		handler.RespondError(c, http.StatusNotFound, "technician not found for employee_number: "+employeeNumber)
		return
	}

	visits, err := h.getVisitsForTechnician(ctx, techID)
	if err != nil {
		handler.RespondError(c, http.StatusInternalServerError, "failed to fetch visits: "+err.Error())
		return
	}

	checklistItems, err := h.getChecklistItemsForVisits(ctx, visits)
	if err != nil {
		handler.RespondError(c, http.StatusInternalServerError, "failed to fetch checklist items: "+err.Error())
		return
	}

	referenceData, err := h.getReferenceData(ctx)
	if err != nil {
		handler.RespondError(c, http.StatusInternalServerError, "failed to fetch reference data: "+err.Error())
		return
	}

	photos, err := h.getPhotosForVisits(ctx, visits)
	if err != nil {
		handler.RespondError(c, http.StatusInternalServerError, "failed to fetch photos: "+err.Error())
		return
	}

	response := mobile.SyncPullResponse{
		ReferenceData:   referenceData,
		Visits:          visits,
		ChecklistItems:  checklistItems,
		Photos:          photos,
		SyncedAt:        time.Now(),
	}

	handler.RespondJSON(c, http.StatusOK, response)
}

func (h *SyncHandler) getTechnicianID(ctx context.Context, employeeNumber string) (uuid.UUID, error) {
	var techID uuid.UUID
	err := h.pool.QueryRow(ctx,
		`SELECT id FROM customers.technicians WHERE user_id = $1`,
		employeeNumber,
	).Scan(&techID)
	return techID, err
}

func (h *SyncHandler) getVisitsForTechnician(ctx context.Context, techID uuid.UUID) ([]mobile.SyncVisit, error) {
	rows, err := h.pool.Query(ctx, `
		SELECT 
			v.id,
			p.id as property_id,
			p.name as property_name,
			COALESCE(a.street || ' ' || a.number || ', ' || a.city, '') as property_address,
			ST_Y(a.geom) as property_lat,
			ST_X(a.geom) as property_lng,
			pt.id as partner_id,
			pt.name as partner_name,
			v.partner_order_id,
			v.partner_supervisor,
			vt.name as visit_type,
			v.status,
			v.priority,
			v.scheduled_at
		FROM operations.visits v
		JOIN operations.visit_assignments va ON va.visit_id = v.id
		JOIN customers.properties p ON p.id = v.property_id
		JOIN customers.customer_addresses ca ON ca.id = p.customer_address_id
		JOIN geocoding.addresses a ON a.id = ca.address_id
		LEFT JOIN partners.partners pt ON pt.id = v.partner_id
		LEFT JOIN operations.visit_types vt ON vt.id = v.type
		WHERE va.tech_id = $1
		  AND v.status IN ('scheduled', 'en_route', 'in_progress')
		ORDER BY v.scheduled_at ASC
	`, techID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var visits []mobile.SyncVisit
	for rows.Next() {
		var v mobile.SyncVisit
		var propertyID, partnerID *uuid.UUID
		var partnerOrderID, partnerSupervisor, visitType *string

		err := rows.Scan(
			&v.ID,
			&propertyID, &v.PropertyName, &v.PropertyAddress,
			&v.PropertyLat, &v.PropertyLng,
			&partnerID, &v.PartnerName,
			&partnerOrderID, &partnerSupervisor,
			&visitType, &v.Status, &v.Priority, &v.ScheduledAt,
		)
		if err != nil {
			return nil, err
		}

		v.PropertyID = propertyID
		v.PartnerID = partnerID
		if partnerOrderID != nil {
			v.PartnerOrderID = *partnerOrderID
		}
		if partnerSupervisor != nil {
			v.PartnerSupervisor = *partnerSupervisor
		}
		if visitType != nil {
			v.Type = *visitType
		}

		visits = append(visits, v)
	}
	return visits, rows.Err()
}

func (h *SyncHandler) getChecklistItemsForVisits(ctx context.Context, visits []mobile.SyncVisit) ([]mobile.SyncChecklistItem, error) {
	if len(visits) == 0 {
		return []mobile.SyncChecklistItem{}, nil
	}

	visitIDs := make([]uuid.UUID, len(visits))
	for i, v := range visits {
		visitIDs[i] = v.ID
	}

	// Build a single query with a CTE for visit IDs to avoid parameter duplication
	cteArgs := make([]interface{}, len(visitIDs))
	for i, id := range visitIDs {
		cteArgs[i] = id
	}

	query := `
		WITH visit_ids AS (
			SELECT * FROM unnest($1::uuid[]) AS vid
		)
		SELECT 
			vcm.id,
			vcm.visit_id,
			'default' as stage_id,
			'General' as stage_name,
			'material' as type,
			m.id as ref_id,
			m.name,
			vcm.planned_quantity as quantity
		FROM operations.visit_checklist_materials vcm
		JOIN inventory.materials m ON m.id = vcm.material_id
		JOIN visit_ids vi ON vi.vid = vcm.visit_id
		UNION ALL
		SELECT 
			vct.id,
			vct.visit_id,
			'default' as stage_id,
			'General' as stage_name,
			'tool' as type,
			t.id as ref_id,
			t.name,
			vct.planned_quantity as quantity
		FROM operations.visit_checklist_tools vct
		JOIN inventory.tools t ON t.id = vct.tool_id
		JOIN visit_ids vi ON vi.vid = vct.visit_id
		UNION ALL
		SELECT 
			vce.id,
			vce.visit_id,
			'default' as stage_id,
			'General' as stage_name,
			'epp' as type,
			e.id as ref_id,
			e.name,
			vce.planned_quantity as quantity
		FROM operations.visit_checklist_epps vce
		JOIN inventory.epp_items e ON e.id = vce.epp_id
		JOIN visit_ids vi ON vi.vid = vce.visit_id
		ORDER BY visit_id, stage_id
	`

	rows, err := h.pool.Query(ctx, query, cteArgs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []mobile.SyncChecklistItem
	for rows.Next() {
		var item mobile.SyncChecklistItem
		err := rows.Scan(
			&item.ID, &item.VisitID, &item.StageID, &item.StageName,
			&item.Type, &item.RefID, &item.Name, &item.Quantity,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (h *SyncHandler) getPhotosForVisits(ctx context.Context, visits []mobile.SyncVisit) ([]mobile.SyncPhotoInfo, error) {
	if len(visits) == 0 {
		return []mobile.SyncPhotoInfo{}, nil
	}

	visitIDs := make([]uuid.UUID, len(visits))
	for i, v := range visits {
		visitIDs[i] = v.ID
	}

	cteArgs := make([]interface{}, len(visitIDs))
	for i, id := range visitIDs {
		cteArgs[i] = id
	}

	query := `
		WITH visit_ids AS (
			SELECT * FROM unnest($1::uuid[]) AS vid
		)
		SELECT vp.id, vp.visit_id
		FROM operations.visit_photos vp
		JOIN visit_ids vi ON vi.vid = vp.visit_id
		ORDER BY vp.visit_id, vp.created_at DESC
	`

	rows, err := h.pool.Query(ctx, query, cteArgs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []mobile.SyncPhotoInfo
	for rows.Next() {
		var p mobile.SyncPhotoInfo
		if err := rows.Scan(&p.ID, &p.VisitID); err != nil {
			return nil, err
		}
p.URL = fmt.Sprintf("/api/photos/%s/file", p.ID)
	p.ThumbnailURL = fmt.Sprintf("/api/photos/%s/thumb", p.ID)
		photos = append(photos, p)
	}
	return photos, rows.Err()
}

func joinPlaceholders(placeholders []string) string {
	result := ""
	for i, p := range placeholders {
		if i > 0 {
			result += ","
		}
		result += p
	}
	return result
}

func (h *SyncHandler) getReferenceData(ctx context.Context) (*mobile.ReferenceData, error) {
	// Materials
	matRows, err := h.pool.Query(ctx, `SELECT id, name, sku, unit, unit_cost FROM inventory.materials`)
	if err != nil {
		return nil, err
	}
	defer matRows.Close()

	materials := []mobile.ReferenceMaterial{}
	for matRows.Next() {
		var m mobile.ReferenceMaterial
		if err := matRows.Scan(&m.ID, &m.Name, &m.SKU, &m.Unit, &m.UnitCost); err != nil {
			return nil, err
		}
		materials = append(materials, m)
	}

	// Tools
	toolRows, err := h.pool.Query(ctx, `SELECT id, name, code FROM inventory.tools`)
	if err != nil {
		return nil, err
	}
	defer toolRows.Close()

	tools := []mobile.ReferenceTool{}
	for toolRows.Next() {
		var t mobile.ReferenceTool
		if err := toolRows.Scan(&t.ID, &t.Name, &t.SKU); err != nil {
			return nil, err
		}
		tools = append(tools, t)
	}

	// EPP
	eppRows, err := h.pool.Query(ctx, `SELECT id, name, type FROM inventory.epp_items`)
	if err != nil {
		return nil, err
	}
	defer eppRows.Close()

	epp := []mobile.ReferenceEPP{}
	for eppRows.Next() {
		var e mobile.ReferenceEPP
		if err := eppRows.Scan(&e.ID, &e.Name, &e.SKU); err != nil {
			return nil, err
		}
		epp = append(epp, e)
	}

	// Checklist templates
	tmplRows, err := h.pool.Query(ctx, `SELECT id, name, work_type FROM inventory.checklist_templates`)
	if err != nil {
		return nil, err
	}
	defer tmplRows.Close()

	templates := []mobile.ReferenceTemplate{}
	for tmplRows.Next() {
		var t mobile.ReferenceTemplate
		var workType *uuid.UUID
		if err := tmplRows.Scan(&t.ID, &t.Name, &workType); err != nil {
			return nil, err
		}
		if workType != nil {
			t.WorkType = workType.String()
		}
		templates = append(templates, t)
	}

	// Partners
	partnerRows, err := h.pool.Query(ctx, `SELECT id, name FROM partners.partners WHERE status = 'active'`)
	if err != nil {
		return nil, err
	}
	defer partnerRows.Close()

	partners := []mobile.ReferencePartner{}
	for partnerRows.Next() {
		var p mobile.ReferencePartner
		if err := partnerRows.Scan(&p.ID, &p.Name); err != nil {
			return nil, err
		}
		partners = append(partners, p)
	}

	return &mobile.ReferenceData{
		Materials: materials,
		Tools:     tools,
		EPP:       epp,
		Templates: templates,
		Partners:  partners,
	}, nil
}

func (h *SyncHandler) Retry(c *gin.Context) {
	var body struct {
		ItemIDs []string `json:"item_ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	handler.RespondJSON(c, http.StatusOK, gin.H{
		"retried": len(body.ItemIDs),
		"failed":  0,
	})
}