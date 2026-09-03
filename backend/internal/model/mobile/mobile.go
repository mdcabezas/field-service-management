package mobile

import (
	"time"

	"github.com/google/uuid"
)

type MobilePhoto struct {
	ID            uuid.UUID `json:"id" db:"id"`
	VisitID       uuid.UUID `json:"visit_id" db:"visit_id"`
	URL           string    `json:"url" db:"url"`
	ThumbnailURL  string    `json:"thumbnail_url" db:"thumbnail_url"`
	OriginalPath  string    `json:"-" db:"original_path"`
	ThumbnailPath string    `json:"-" db:"thumbnail_path"`
	Lat           *float64  `json:"lat,omitempty" db:"lat"`
	Lng           *float64  `json:"lng,omitempty" db:"lng"`
	Timestamp     time.Time `json:"timestamp" db:"timestamp"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type MobileConfig struct {
	PhotoRetentionDays      int `json:"photo_retention_days" db:"photo_retention_days"`
	PhotoMaxStorageMB       int `json:"photo_max_storage_mb" db:"photo_max_storage_mb"`
	PhotoWarningThresholdMB int `json:"photo_warning_threshold_mb" db:"photo_warning_threshold_mb"`
}

type GPSWaiver struct {
	ID         uuid.UUID `json:"id" db:"id"`
	VisitID    uuid.UUID `json:"visit_id" db:"visit_id"`
	Type       string    `json:"type" db:"type"`
	AcceptedBy string    `json:"accepted_by" db:"accepted_by"`
	Timestamp  time.Time `json:"timestamp" db:"timestamp"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type SyncPushRequest struct {
	Checkpoints    []SyncCheckpoint    `json:"checkpoints"`
	Photos         []SyncPhoto         `json:"photos"`
	Usages         []SyncUsage         `json:"usages"`
	Measurements   []SyncMeasurement   `json:"measurements"`
	Reports        []SyncReport        `json:"reports"`
	StatusChanges  []SyncStatusChange  `json:"status_changes"`
}

type SyncCheckpoint struct {
	VisitID    uuid.UUID `json:"visit_id"`
	Type       string    `json:"type"`
	Timestamp  time.Time `json:"timestamp"`
	Lat        float64   `json:"lat"`
	Lng        float64   `json:"lng"`
	Accuracy   float64   `json:"accuracy"`
	LocalID    string    `json:"local_id"`
}

type SyncPhoto struct {
	VisitID   uuid.UUID `json:"visit_id"`
	LocalID   string    `json:"local_id"`
	Timestamp time.Time `json:"timestamp"`
	Lat       *float64  `json:"lat,omitempty"`
	Lng       *float64  `json:"lng,omitempty"`
}

type SyncUsage struct {
	VisitID    uuid.UUID `json:"visit_id"`
	MaterialID uuid.UUID `json:"material_id"`
	Quantity   float64   `json:"quantity"`
	Notes      string    `json:"notes"`
	LocalID    string    `json:"local_id"`
}

type SyncMeasurement struct {
	VisitID uuid.UUID `json:"visit_id"`
	Key     string    `json:"key"`
	Value   string    `json:"value"`
	Unit    string    `json:"unit"`
	Lat     *float64  `json:"lat,omitempty"`
	Lng     *float64  `json:"lng,omitempty"`
	LocalID string    `json:"local_id"`
}

type SyncReport struct {
	VisitID          uuid.UUID              `json:"visit_id"`
	ReportTemplateID uuid.UUID              `json:"report_template_id"`
	Payload          map[string]interface{} `json:"payload"`
	LocalID          string                 `json:"local_id"`
}

type SyncStatusChange struct {
	VisitID   uuid.UUID `json:"visit_id"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Notes     string    `json:"notes"`
	LocalID   string    `json:"local_id"`
}

type SyncResult struct {
	SyncedAt time.Time       `json:"synced_at"`
	Results  SyncResultItems `json:"results"`
}

type SyncResultItems struct {
	Checkpoints   []SyncItemResult `json:"checkpoints"`
	Photos        []SyncItemResult `json:"photos"`
	Usages        []SyncItemResult `json:"usages"`
	Measurements  []SyncItemResult `json:"measurements"`
	Reports       []SyncItemResult `json:"reports"`
	StatusChanges []SyncItemResult `json:"status_changes"`
}

type SyncItemResult struct {
	LocalID  string `json:"local_id"`
	ServerID string `json:"server_id"`
	Status   string `json:"status"`
	Conflict bool   `json:"conflict,omitempty"`
	Error    string `json:"error,omitempty"`
}

type SyncPullResponse struct {
	ReferenceData *ReferenceData      `json:"reference_data,omitempty"`
	Visits        []SyncVisit         `json:"visits"`
	ChecklistItems []SyncChecklistItem `json:"checklist_items"`
	Photos        []SyncPhotoInfo     `json:"photos"`
	SyncedAt      time.Time           `json:"synced_at"`
}

type ReferenceData struct {
	Materials []ReferenceMaterial `json:"materials"`
	Tools     []ReferenceTool     `json:"tools"`
	EPP       []ReferenceEPP      `json:"epp"`
	Templates []ReferenceTemplate `json:"templates"`
	Partners  []ReferencePartner  `json:"partners"`
}

type ReferenceMaterial struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	SKU      string    `json:"sku"`
	Unit     string    `json:"unit"`
	UnitCost float64   `json:"unit_cost"`
}

type ReferenceTool struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	SKU         string    `json:"sku"`
	Description string    `json:"description"`
}

type ReferenceEPP struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	SKU         string    `json:"sku"`
	Description string    `json:"description"`
}

type ReferenceTemplate struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	VisitType      string    `json:"visit_type"`
	WorkType       string    `json:"work_type"`
	PartnerID      *uuid.UUID `json:"partner_id,omitempty"`
	StagesEnabled  bool      `json:"stages_enabled"`
}

type ReferencePartner struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type SyncVisit struct {
	ID                   uuid.UUID              `json:"id"`
	PropertyID           *uuid.UUID             `json:"property_id,omitempty"`
	PropertyName         string                 `json:"property_name"`
	PropertyAddress      string                 `json:"property_address"`
	PropertyLat          *float64               `json:"property_lat,omitempty"`
	PropertyLng          *float64               `json:"property_lng,omitempty"`
	PartnerID            *uuid.UUID             `json:"partner_id,omitempty"`
	PartnerName          string                 `json:"partner_name"`
	PartnerOrderID       string                 `json:"partner_order_id"`
	PartnerSupervisor    string                 `json:"partner_supervisor"`
	Type                 string                 `json:"type"`
	Status               string                 `json:"status"`
	Priority             string                 `json:"priority"`
	ScheduledAt          time.Time              `json:"scheduled_at"`
	ChecklistTemplateID  *uuid.UUID             `json:"checklist_template_id,omitempty"`
	StagesEnabled        bool                   `json:"stages_enabled"`
	SelectedStages       []string               `json:"selected_stages"`
	SLAResponseHours     *int                   `json:"sla_response_hours,omitempty"`
	SLAResolutionHours   *int                   `json:"sla_resolution_hours,omitempty"`
}

type SyncChecklistItem struct {
	ID        uuid.UUID `json:"id"`
	VisitID   uuid.UUID `json:"visit_id"`
	StageID   string    `json:"stage_id"`
	StageName string    `json:"stage_name"`
	Type      string    `json:"type"`
	RefID     uuid.UUID `json:"ref_id"`
	Name      string    `json:"name"`
	Quantity  float64   `json:"quantity"`
	Unit      string    `json:"unit"`
	Confirmed bool      `json:"confirmed"`
	Required  bool      `json:"required"`
}

type SyncPhotoInfo struct {
	ID           uuid.UUID `json:"id"`
	VisitID      uuid.UUID `json:"visit_id"`
	URL          string    `json:"url"`
	ThumbnailURL string    `json:"thumbnail_url"`
}

type PropertyHistory struct {
	Visits    []PropertyHistoryVisit `json:"visits"`
	Total     int                    `json:"total"`
}

type PropertyHistoryVisit struct {
	ID         uuid.UUID `json:"id"`
	Type       string    `json:"type"`
	Status     string    `json:"status"`
	Technician string    `json:"technician"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type AuditEntry struct {
	ID        uuid.UUID `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	UserName  string    `json:"user_name"`
	Action    string    `json:"action"`
	FieldName string    `json:"field_name"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
}

type AuditResponse struct {
	Entries []AuditEntry `json:"entries"`
	Total   int          `json:"total"`
	Page    int          `json:"page"`
	Limit   int          `json:"limit"`
}
