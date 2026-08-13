#!/usr/bin/env node
// generate-collection.js — Generates Postman Collection v2.1 for Go Backend (no PostgREST)
// Usage: node generate-collection.js > FSM.postman_collection.json

const fs = require('fs');

// ============================================================================
// TABLE DEFINITIONS — matches Go backend routes
// ============================================================================

const TABLES = [
  // Core
  { schema: 'core', table: 'users', route: 'users', pk: 'employee_number', pkType: 'string', columns: ['employee_number','email','name','role','is_active','created_at','updated_at'], readOnly: ['employee_number','created_at','updated_at'] },
  { schema: 'core', table: 'tech_roles', route: 'tech-roles', pk: 'id', columns: ['id','code','name','active','created_at'], readOnly: ['id','created_at'] },

  // Partners
  { schema: 'partners', table: 'partners', route: 'partners', pk: 'id', columns: ['id','name','tax_id','status','created_at','updated_at'], readOnly: ['id','created_at','updated_at'] },
  { schema: 'partners', table: 'partner_contacts', route: 'partner-contacts', pk: 'id', parentRoute: 'partners', parentParam: 'partner_id', columns: ['id','partner_id','name','position','phone','email','created_at'], readOnly: ['id','created_at'] },
  { schema: 'partners', table: 'partner_agreements', route: 'partner-agreements', pk: 'id', parentRoute: 'partners', parentParam: 'partner_id', columns: ['id','partner_id','service_type','rate','active','created_at'], readOnly: ['id','created_at'] },
  { schema: 'partners', table: 'partner_agreement_docs', route: 'partner-agreement-docs', pk: 'id', parentRoute: 'partner-agreements', parentParam: 'agreement_id', columns: ['id','agreement_id','work_type','doc_name','required','created_at'], readOnly: ['id','created_at'] },
  { schema: 'partners', table: 'partner_agreement_forms', route: 'partner-agreement-forms', pk: 'id', parentRoute: 'partner-agreements', parentParam: 'agreement_id', columns: ['id','agreement_id','work_type','form_template_id','quantity','created_at'], readOnly: ['id','created_at'] },
  { schema: 'partners', table: 'slas', route: 'slas', pk: 'id', parentRoute: 'partners', parentParam: 'partner_id', columns: ['id','partner_id','name','description','work_type','response_hours','resolution_hours','compliance_target','active','valid_from','valid_until','created_at'], readOnly: ['id','created_at'] },

  // Customers
  { schema: 'customers', table: 'customers', route: 'customers', pk: 'id', columns: ['id','name','phone','email','tax_id','created_at','updated_at'], readOnly: ['id','created_at','updated_at'] },
  { schema: 'customers', table: 'customer_addresses', route: 'customer-addresses', pk: 'id', parentRoute: 'customers', parentParam: 'customer_id', columns: ['id','customer_id','address_id','name','type','created_at'], readOnly: ['id','created_at'] },
  { schema: 'customers', table: 'customer_address_partners', route: 'customer-address-partners', pk: 'id', parentRoute: 'customer-addresses', parentParam: 'customer_address_id', columns: ['id','customer_address_id','partner_id','from_date','to_date','notes','created_at'], readOnly: ['id','created_at'] },
  { schema: 'customers', table: 'properties', route: 'properties', pk: 'id', parentRoute: 'customer-addresses', parentParam: 'customer_address_id', columns: ['id','customer_address_id','name','type','notes','created_at'], readOnly: ['id','created_at'] },
  { schema: 'customers', table: 'technicians', route: 'technicians', pk: 'id', columns: ['id','user_id','name','is_active','specialties','created_at'], readOnly: ['id','created_at'] },
  { schema: 'customers', table: 'tech_certifications', route: 'tech-certifications', pk: 'id', parentRoute: 'technicians', parentParam: 'tech_id', columns: ['id','tech_id','cert_id','number','issuer','issue_date','expiry_date','created_at'], readOnly: ['id','created_at'] },

  // Inventory
  { schema: 'inventory', table: 'materials', route: 'materials', pk: 'id', columns: ['id','sku','name','unit','unit_cost','created_at'], readOnly: ['id','created_at'] },
  { schema: 'inventory', table: 'tools', route: 'tools', pk: 'id', columns: ['id','code','name','status','created_at'], readOnly: ['id','created_at'] },
  { schema: 'inventory', table: 'epp_items', route: 'epp-items', pk: 'id', columns: ['id','name','type','lifecycle','created_at'], readOnly: ['id','created_at'] },
  { schema: 'inventory', table: 'vehicles', route: 'vehicles', pk: 'id', columns: ['id','type','license_plate','name','brand','model','year','status','capacity','created_at'], readOnly: ['id','created_at'] },
  { schema: 'inventory', table: 'rentals', route: 'rentals', pk: 'id', columns: ['id','type','supplier','item_description','daily_cost','hourly_cost','fixed_cost','notes','created_at'], readOnly: ['id','created_at'] },
  { schema: 'inventory', table: 'cost_rates', route: 'cost-rates', pk: 'id', columns: ['id','type','reference_id','reference_type','value','unit','valid_from','valid_until','created_at'], readOnly: ['id','created_at'] },
  { schema: 'inventory', table: 'maintenance_records', route: 'maintenance-records', pk: 'id', columns: ['id','type','reference_id','date','cost','supplier','description','next_date','created_at'], readOnly: ['id','created_at'], noUpdate: true },
  { schema: 'inventory', table: 'maintenance_schedules', route: 'maintenance-schedules', pk: 'id', columns: ['id','type','reference_id','frequency_km','frequency_days','last_service_date','last_service_km','active','created_at'], readOnly: ['id','created_at'] },
  { schema: 'inventory', table: 'checklist_templates', route: 'checklist-templates', pk: 'id', columns: ['id','name','work_type','created_at'], readOnly: ['id','created_at'] },

  // Planning
  { schema: 'planning', table: 'daily_plans', route: 'daily-plans', pk: 'id', columns: ['id','date','notes','created_at'], readOnly: ['id','created_at'] },
  { schema: 'planning', table: 'daily_plan_assignments', route: 'daily-plan-assignments', pk: 'id', parentRoute: 'daily-plans', parentParam: 'daily_plan_id', columns: ['id','daily_plan_id','tech_id','role_id','created_at'], readOnly: ['id','created_at'] },
  { schema: 'planning', table: 'daily_load_materials', route: 'daily-load-materials', pk: 'id', parentRoute: 'daily-plans', parentParam: 'daily_plan_id', columns: ['id','daily_plan_id','material_id','loaded_quantity','returned_quantity','created_at'], readOnly: ['id','created_at'] },
  { schema: 'planning', table: 'daily_load_tools', route: 'daily-load-tools', pk: 'id', parentRoute: 'daily-plans', parentParam: 'daily_plan_id', columns: ['id','daily_plan_id','tool_id','loaded_quantity','returned_quantity','created_at'], readOnly: ['id','created_at'] },
  { schema: 'planning', table: 'daily_load_epps', route: 'daily-load-epps', pk: 'id', parentRoute: 'daily-plans', parentParam: 'daily_plan_id', columns: ['id','daily_plan_id','epp_id','loaded_quantity','returned_quantity','created_at'], readOnly: ['id','created_at'] },
  { schema: 'planning', table: 'routes', route: 'routes', pk: 'id', columns: ['id','type','date','daily_plan_id','notes','status','created_at'], readOnly: ['id','created_at'] },

  // Operations
  { schema: 'operations', table: 'visits', route: 'visits', pk: 'id', columns: ['id','property_id','partner_id','partner_order_id','parent_visit_id','route_id','daily_plan_id','type','status','priority','source','source_reference','billing_to','partner_supervisor','result','rejection_reason_id','rejection_reason_detail','scheduled_at','started_at','completed_at','alternative_location','notes','created_at','updated_at'], readOnly: ['id','created_at','updated_at'] },
  { schema: 'operations', table: 'visit_assignments', route: 'visit-assignments', pk: 'id', parentRoute: 'visits', parentParam: 'visit_id', columns: ['id','visit_id','tech_id','role_id','created_at'], readOnly: ['id','created_at'] },
  { schema: 'operations', table: 'checklist_template_materials', route: 'checklist-template-materials', pk: 'id', parentRoute: 'checklist-templates', parentParam: 'template_id', columns: ['id','template_id','material_id','default_quantity','required','created_at'], readOnly: ['id','created_at'] },
  { schema: 'operations', table: 'checklist_template_tools', route: 'checklist-template-tools', pk: 'id', parentRoute: 'checklist-templates', parentParam: 'template_id', columns: ['id','template_id','tool_id','default_quantity','required','created_at'], readOnly: ['id','created_at'] },
  { schema: 'operations', table: 'checklist_template_epps', route: 'checklist-template-epps', pk: 'id', parentRoute: 'checklist-templates', parentParam: 'template_id', columns: ['id','template_id','epp_id','default_quantity','required','created_at'], readOnly: ['id','created_at'] },
  { schema: 'operations', table: 'visit_checklist_materials', route: 'visit-checklist-materials', pk: 'id', parentRoute: 'visits', parentParam: 'visit_id', columns: ['id','visit_id','material_id','planned_quantity','confirmed','notes','created_at'], readOnly: ['id','created_at'] },
  { schema: 'operations', table: 'visit_checklist_tools', route: 'visit-checklist-tools', pk: 'id', parentRoute: 'visits', parentParam: 'visit_id', columns: ['id','visit_id','tool_id','planned_quantity','confirmed','notes','created_at'], readOnly: ['id','created_at'] },
  { schema: 'operations', table: 'visit_checklist_epps', route: 'visit-checklist-epps', pk: 'id', parentRoute: 'visits', parentParam: 'visit_id', columns: ['id','visit_id','epp_id','planned_quantity','confirmed','notes','created_at'], readOnly: ['id','created_at'] },
  { schema: 'operations', table: 'visit_material_usages', route: 'visit-material-usages', pk: 'id', parentRoute: 'visits', parentParam: 'visit_id', columns: ['id','visit_id','material_id','quantity','notes','created_at'], readOnly: ['id','created_at'] },
  { schema: 'operations', table: 'visit_tool_usages', route: 'visit-tool-usages', pk: 'id', parentRoute: 'visits', parentParam: 'visit_id', columns: ['id','visit_id','tool_id','notes','created_at'], readOnly: ['id','created_at'] },
  { schema: 'operations', table: 'visit_epp_usages', route: 'visit-epp-usages', pk: 'id', parentRoute: 'visits', parentParam: 'visit_id', columns: ['id','visit_id','epp_id','quantity','status','notes','created_at'], readOnly: ['id','created_at'] },
  { schema: 'operations', table: 'visit_measurements', route: 'visit-measurements', pk: 'id', parentRoute: 'visits', parentParam: 'visit_id', columns: ['id','visit_id','type','value','unit','result','measuring_device','notes','photo_url','created_at'], readOnly: ['id','created_at'] },
  { schema: 'operations', table: 'visit_checkpoints', route: 'visit-checkpoints', pk: 'id', parentRoute: 'visits', parentParam: 'visit_id', columns: ['id','visit_id','type','geom','timestamp','device_info','created_at'], readOnly: ['id','created_at'] },
  { schema: 'operations', table: 'visit_photos', route: 'visit-photos', pk: 'id', parentRoute: 'visits', parentParam: 'visit_id', columns: ['id','visit_id','checkpoint_id','url','geom','timestamp','stage','finding_type','created_at'], readOnly: ['id','created_at'] },
  { schema: 'operations', table: 'visit_reports', route: 'visit-reports', pk: 'id', parentRoute: 'visits', parentParam: 'visit_id', columns: ['id','visit_id','report_template_id','source','recorded_at','created_at'], readOnly: ['id','created_at'] },
  { schema: 'operations', table: 'report_images', route: 'report-images', pk: 'id', parentRoute: 'visit-reports', parentParam: 'report_id', columns: ['id','report_id','url','pages','created_at'], readOnly: ['id','created_at'] },
  { schema: 'operations', table: 'report_entries', route: 'report-entries', pk: 'id', parentRoute: 'visit-reports', parentParam: 'report_id', columns: ['id','report_id','data_json','created_at'], readOnly: ['id','created_at'] },
  { schema: 'operations', table: 'vehicle_assignments', route: 'vehicle-assignments', pk: 'id', columns: ['id','vehicle_id','daily_plan_id','visit_id','departure_time','return_time','departure_mileage','return_mileage','notes','created_at'], readOnly: ['id','created_at'] },
  { schema: 'operations', table: 'visit_rentals', route: 'visit-rentals', pk: 'id', parentRoute: 'visits', parentParam: 'visit_id', columns: ['id','visit_id','rental_id','start_time','end_time','hours','total_cost','reason','created_at'], readOnly: ['id','created_at'] },
  { schema: 'operations', table: 'visit_sla_trackings', route: 'visit-sla-trackings', pk: 'id', parentRoute: 'visits', parentParam: 'visit_id', columns: ['id','visit_id','sla_id','requested_at','responded_at','resolved_at','response_time_hours','resolution_time_hours','meets_response_sla','meets_resolution_sla','created_at'], readOnly: ['id','created_at'] },

  // Notifications
  { schema: 'notifications', table: 'notification_templates', route: 'notification-templates', pk: 'id', columns: ['id','type','channel','subject','body','active','created_at'], readOnly: ['id','created_at'] },
  { schema: 'notifications', table: 'notifications', route: 'notifications', pk: 'id', columns: ['id','type','channel','recipient','subject','body','status','entity_type','entity_id','scheduled_for','sent_at','read_at','attempts','error','created_at'], readOnly: ['id','created_at'] },

  // Geocoding
  { schema: 'geocoding', table: 'addresses', route: 'addresses', pk: 'id', columns: ['id','street','number','apartment','neighborhood','city','region','location_references','postal_code','geom','created_at'], readOnly: ['id','created_at'] },

  // Shared
  { schema: 'shared', table: 'report_templates', route: 'report-templates', pk: 'id', columns: ['id','name','fields_json','created_at'], readOnly: ['id','created_at'] },
];

// ============================================================================
// JWT SCRIPT
// ============================================================================

const JWT_SCRIPT = `
(function () {
  function base64url(str) {
    return btoa(str).replace(/=/g, '').replace(/\\+/g, '-').replace(/\\//g, '_');
  }
  function hmacSha256(message, secret) {
    var hash = CryptoJS.HmacSHA256(message, secret);
    return hash.toString(CryptoJS.enc.Base64).replace(/=/g, '').replace(/\\+/g, '-').replace(/\\//g, '_');
  }
  var existingToken = pm.environment.get('jwt_token');
  if (existingToken) {
    try {
      var parts = existingToken.split('.');
      var payload = JSON.parse(atob(parts[1].replace(/-/g, '+').replace(/_/g, '/')));
      var now = Math.floor(Date.now() / 1000);
      if (payload.exp && payload.exp > now + 60) return;
    } catch (e) {}
  }
  var secret = pm.environment.get('JWT_SECRET');
  if (!secret) { console.error('JWT_SECRET not set'); return; }
  var userId = pm.environment.get('employee_number') || '1001';
  var userRole = pm.environment.get('user_role') || 'admin';
  var ttl = parseInt(pm.environment.get('jwt_ttl_minutes')) || 60;
  var now = Math.floor(Date.now() / 1000);
  var header = base64url(JSON.stringify({ alg: 'HS256', typ: 'JWT' }));
  var claims = base64url(JSON.stringify({ sub: userId, role: userRole, iat: now, exp: now + (ttl * 60) }));
  var sig = hmacSha256(header + '.' + claims, secret);
  pm.environment.set('jwt_token', header + '.' + claims + '.' + sig);
})();
`.trim();

// ============================================================================
// REQUEST BUILDERS
// ============================================================================

function url(base, ...parts) {
  return [base, ...parts].filter(Boolean).join('/');
}

function makeRequest(method, urlStr, body, description) {
  const pathParts = urlStr.replace('{{base_url}}/', '').split('/');
  const req = {
    method,
    header: [
      { key: 'Authorization', value: 'Bearer {{jwt_token}}', type: 'text' },
      { key: 'Content-Type', value: 'application/json', type: 'text' }
    ],
    url: { raw: urlStr, host: ['{{base_url}}'], path: pathParts },
    description
  };
  if (body) {
    req.body = { mode: 'raw', raw: JSON.stringify(body, null, 2) };
  }
  return req;
}

function exampleValue(col) {
  if (col === 'name' || col === 'doc_name' || col === 'subject') return 'Example';
  if (col === 'email') return 'test@example.com';
  if (col === 'phone' || col === 'recipient') return '+56912345678';
  if (col === 'code' || col === 'sku') return 'CODE001';
  if (col === 'status') return 'active';
  if (col === 'active' || col === 'required') return true;
  if (col === 'confirmed') return false;
  if (col === 'is_active') return true;
  if (['quantity','default_quantity','planned_quantity','loaded_quantity','attempts','pages','year'].includes(col)) return 1;
  if (col === 'returned_quantity') return 0;
  if (['rate','value','unit_cost','daily_cost','hourly_cost','fixed_cost','total_cost','cost','compliance_target'].includes(col)) return 10000;
  if (['response_hours','resolution_hours'].includes(col)) return 24;
  if (col === 'from_date') return '2024-01-01';
  if (col === 'frequency_km' || col === 'frequency_days') return 30;
  if (col === 'specialties') return ['maintenance', 'repair'];
  if (col === 'fields_json') return [];
  if (col === 'data_json') return {};
  if (col === 'geom') return 'SRID=4326;POINT(-70.6483 -33.4569)';
  if (['scheduled_at','start_time','departure_time','timestamp','recorded_at','requested_at','responded_at','resolved_at','sent_at','read_at','scheduled_for'].includes(col)) return '2024-01-15T08:00:00Z';
  if (['valid_from','valid_until','issue_date','expiry_date','date','next_date','last_service_date'].includes(col)) return '2024-01-15';
  if (col.endsWith('_id')) return '00000000-0000-0000-0000-000000000000';
  if (['notes','description','body','reason','rejection_reason_detail','alternative_location','partner_supervisor','source_reference','partner_order_id','device_info','error','doc_name','supplier','item_description','measuring_device','photo_url','url'].includes(col)) return null;
  if (col === 'type') return 'other';
  if (col === 'priority') return 'normal';
  if (col === 'billing_to') return 'customer';
  if (col === 'source') return 'system';
  if (col === 'result') return 'approved';
  if (col === 'lifecycle') return 'reusable';
  if (col === 'channel') return 'email';
  return null;
}

function buildCollection() {
  const collection = {
    info: {
      name: 'FSM Go Backend API',
      description: 'Field Service Management API.\n\n**Stack**: PostgreSQL 16 + PostGIS 3.5 + Go Backend (Gin) + GLAuth + Traefik\n**Mode**: Single-tenant, JWT auth via Go Backend → direct PostgreSQL\n\n## Quick Start\n1. Select environment: FSM-dev\n2. Run "Auth: Login" folder\n3. Start querying endpoints\n\n## Auth\nPOST /auth/login with employee_number + password to get JWT.\nToken is cached in jwt_token variable and auto-refreshed on expiry.\n\n**JWT payload**: {sub: employee_number, role: "admin|operator|technician"}',
      schema: 'https://schema.getpostman.com/json/collection/v2.1.0/collection.json'
    },
    variable: [
      { key: 'base_url', value: '{{base_url}}' },
      { key: 'jwt_token', value: '{{jwt_token}}' }
    ],
    auth: { type: 'bearer', bearer: [{ key: 'token', value: '{{jwt_token}}' }] },
    event: [
      { listen: 'prerequest', script: { exec: [JWT_SCRIPT], type: 'text/javascript' } }
    ],
    item: []
  };

  // Auth folder
  collection.item.push({
    name: 'Auth',
    description: 'Authentication (LDAP + JWT)',
    item: [
      {
        name: 'Login',
        description: 'POST /auth/login — Authenticate via LDAP, returns JWT',
        request: makeRequest('POST', '{{base_url}}/auth/login', { employee_number: '1001', password: 'admin123' }, 'Login with LDAP credentials'),
        response: []
      },
      {
        name: 'Refresh Token',
        description: 'POST /auth/refresh — Refresh JWT token',
        request: makeRequest('POST', '{{base_url}}/auth/refresh', { refresh_token: '{{jwt_token}}' }, 'Refresh JWT'),
        response: []
      },
      {
        name: 'Me',
        description: 'GET /auth/me — Get current user info',
        request: makeRequest('GET', '{{base_url}}/auth/me', null, 'Validate token and get user info'),
        response: []
      }
    ]
  });

  // Group tables by schema
  const schemaOrder = ['core','partners','customers','inventory','planning','operations','notifications','geocoding','shared'];
  const schemaDesc = {
    core: 'Users, roles, audit log',
    partners: 'Partners, contacts, agreements, SLAs',
    customers: 'Customers, addresses, properties, technicians',
    inventory: 'Materials, tools, EPP, vehicles, rentals, costs, maintenance, templates',
    planning: 'Daily plans, assignments, material loading, routes',
    operations: 'Visits, assignments, checklists, measurements, checkpoints, photos, reports, vehicles, rentals, SLAs',
    notifications: 'Templates and notification instances',
    geocoding: 'Georeferenced addresses',
    shared: 'Global catalogs'
  };

  for (const schemaName of schemaOrder) {
    const tables = TABLES.filter(t => t.schema === schemaName);
    if (tables.length === 0) continue;

    const schemaFolder = {
      name: schemaName,
      description: schemaDesc[schemaName] || '',
      item: []
    };

    for (const t of tables) {
      const tableFolder = {
        name: t.table,
        description: `${t.schema}.${t.table} — ${t.columns.length} columns`,
        item: []
      };

      // List
      tableFolder.item.push({
        name: `List ${t.table}`,
        request: makeRequest('GET', `{{base_url}}/api/${t.route}?limit=20&offset=0`, null, `GET /api/${t.route} — List with pagination`),
        response: []
      });

      // Get by ID
      if (t.pkType === 'string') {
        tableFolder.item.push({
          name: `Get ${t.table} by employee_number`,
          request: makeRequest('GET', `{{base_url}}/api/${t.route}/:id`, null, `GET /api/${t.route}/:id — Get by employee number`),
          response: []
        });
      } else {
        tableFolder.item.push({
          name: `Get ${t.table} by ID`,
          request: makeRequest('GET', `{{base_url}}/api/${t.route}/:id`, null, `GET /api/${t.route}/:id — Get by UUID`),
          response: []
        });
      }

      // Create (skip read-only tables)
      if (t.table !== 'plan_audit_log') {
        const writableCols = t.columns.filter(c => !t.readOnly.includes(c));
        const body = {};
        writableCols.forEach(c => { body[c] = exampleValue(c); });

        tableFolder.item.push({
          name: `Create ${t.table}`,
          request: makeRequest('POST', `{{base_url}}/api/${t.route}`, body, `POST /api/${t.route} — Create new record`),
          response: []
        });
      }

      // Update (skip read-only and noUpdate tables)
      if (t.table !== 'plan_audit_log' && !t.noUpdate) {
        const writableCols = t.columns.filter(c => !t.readOnly.includes(c));
        const body = {};
        writableCols.forEach(c => { body[c] = exampleValue(c); });

        tableFolder.item.push({
          name: `Update ${t.table}`,
          request: makeRequest('PUT', `{{base_url}}/api/${t.route}/:id`, body, `PUT /api/${t.route}/:id — Update record`),
          response: []
        });
      }

      // Delete (skip read-only tables)
      if (t.table !== 'plan_audit_log') {
        tableFolder.item.push({
          name: `Delete ${t.table}`,
          request: makeRequest('DELETE', `{{base_url}}/api/${t.route}/:id`, null, `DELETE /api/${t.route}/:id — Delete record`),
          response: []
        });
      }

      schemaFolder.item.push(tableFolder);
    }

    collection.item.push(schemaFolder);
  }

  return collection;
}

// ============================================================================
// MAIN
// ============================================================================

const collection = buildCollection();
const output = JSON.stringify(collection, null, 2);

if (process.argv[2]) {
  fs.writeFileSync(process.argv[2], output, 'utf8');
  const count = output.split('"name":').length - 1;
  console.error(`Collection written to ${process.argv[2]}`);
  console.error(`Total named items: ${count}`);
} else {
  process.stdout.write(output);
}
