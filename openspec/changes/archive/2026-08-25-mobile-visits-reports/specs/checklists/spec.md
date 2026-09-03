# Checklists by Configurable Stages Spec

## Delta: New capability — configurable stages per checklist template

## Context
Checklist templates currently have flat lists of materials, tools, and EPPs. We add a `checklist_stages` system where templates define ordered stages (e.g., "Pre-departure", "On-site", "Post-work"). When creating a visit, the planner selects which stages to include. Required items in a stage block status transitions.

## Functional Requirements

### FR-1: Checklist Stages
- New table: `inventory.checklist_stages` (id, template_id, name, sort_order, required, created_at)
- New fields on templates: `stages_enabled` (boolean, default false)
- If stages_enabled = false: flat checklist (current behavior)
- If stages_enabled = true: items grouped by stage

### FR-2: Stage Management
- Create/update/delete stages on template
- Reorder stages via sort_order
- Each stage has a name (e.g., "Pre-departure", "On-site", "Post-work")
- Each stage can be marked as required (blocks status transition if incomplete)

### FR-3: Stage Selection on Visit Creation
- When template has stages_enabled = true, planner selects which stages to include
- Default: all stages selected
- Emergency visits can skip pre-departure stage
- Selected stages stored on visit (not template)

### FR-4: Checklist Population
- `PopulateChecklist` filtered by selected stages
- If no stages selected: populate all (backward compatible)
- If stages selected: only populate items belonging to selected stages
- Stage information included in checklist items

### FR-5: Status Transition Validation
- When visiting status (e.g., "en_curso" → "completada"):
  - Check required stages for completeness
  - If any required stage has unchecked required items → block transition
  - Return specific error: "Stage 'Post-work' has 3 incomplete required items"
- Optional stages do NOT block transitions
- Emergency visits: all stages become optional

### FR-6: Checklist Display by Stage
- Mobile shows checklist grouped by stage
- Stage header with name and progress (X/Y completed)
- Collapsible stages
- Required indicator on items that block transitions
- Color coding: green (complete), yellow (in progress), grey (not started)

## Non-Functional Requirements

### NFR-1: Backward Compatibility
- Existing templates continue to work (stages_enabled = false)
- No data migration required for existing checklists
- Flat checklists remain functional

### NFR-2: Performance
- Stage filtering is O(1) per item (indexed)
- No performance impact on checklist population

## API Contract

### POST /api/checklist-stages
```json
// Request
{
  "template_id": "...",
  "name": "Pre-departure",
  "sort_order": 0,
  "required": true
}

// Response
{ "id": "...", "template_id": "...", "name": "Pre-departure", "sort_order": 0, "required": true }
```

### PUT /api/checklist-stages/:id
```json
// Request
{ "name": "Pre-work", "sort_order": 1, "required": false }

// Response
{ "id": "...", "name": "Pre-work", "sort_order": 1, "required": false }
```

### DELETE /api/checklist-stages/:id
```json
// Response 204 No Content
```

### GET /api/checklist-stages?template_id=X
```json
// Response
[
  { "id": "...", "name": "Pre-departure", "sort_order": 0, "required": true },
  { "id": "...", "name": "On-site", "sort_order": 1, "required": true },
  { "id": "...", "name": "Post-work", "sort_order": 2, "required": false }
]
```

### POST /api/visits/:id/populate-checklist
```json
// Request (new field)
{
  "template_id": "...",
  "selected_stages": ["stage-id-1", "stage-id-3"]
}

// Response
{
  "items": [
    { "id": "...", "stage_id": "stage-id-1", "stage_name": "Pre-departure", "type": "material", "name": "...", "confirmed": false }
  ]
}
```

### PUT /api/visits/:id/status
```json
// Request
{ "status": "completada" }

// Response (if blocked)
{
  "error": "Checklist incomplete",
  "blocking_stages": [
    { "stage_id": "...", "name": "Post-work", "incomplete_required": 3 }
  ]
}
```
