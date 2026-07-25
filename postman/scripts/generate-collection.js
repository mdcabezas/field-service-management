#!/usr/bin/env node
// generate-collection.js — Generates Postman Collection v2.1 from schema definition
// Usage: node generate-collection.js > GAS-FSM.postman_collection.json

const fs = require('fs');
const path = require('path');

// ============================================================================
// SCHEMA DEFINITION
// ============================================================================

const SCHEMAS = {
  core: {
    tables: {
      users: {
        columns: ['employee_number', 'email', 'role', 'name', 'is_active', 'created_at', 'updated_at'],
        readOnly: ['employee_number', 'created_at', 'updated_at'],
        fks: {},
        embeds: {}
      },
      tech_roles: {
        columns: ['id', 'code', 'name', 'active', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: {},
        embeds: {}
      },
      plan_audit_log: {
        columns: ['id', 'entity_type', 'entity_id', 'action', 'old_data', 'new_data', 'user_id', 'reason', 'timestamp'],
        readOnly: ['id', 'timestamp', 'entity_type', 'entity_id', 'action', 'old_data', 'new_data', 'user_id'],
        fks: {},
        embeds: {},
        readOnlyTable: true
      }
    }
  },
  partners: {
    tables: {
      partners: {
        columns: ['id', 'name', 'tax_id', 'status', 'created_at', 'updated_at'],
        readOnly: ['id', 'created_at', 'updated_at'],
        fks: {},
        embeds: {
          partner_contacts: 'partner_id',
          partner_agreements: 'partner_id',
          slas: 'partner_id'
        }
      },
      partner_contacts: {
        columns: ['id', 'partner_id', 'name', 'position', 'phone', 'email', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { partner_id: 'partners.partners' },
        embeds: {}
      },
      partner_agreements: {
        columns: ['id', 'partner_id', 'service_type', 'rate', 'active', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { partner_id: 'partners.partners', service_type: 'domain_gas.partner_service_types' },
        embeds: {
          partner_agreement_docs: 'agreement_id',
          partner_agreement_forms: 'agreement_id'
        }
      },
      partner_agreement_docs: {
        columns: ['id', 'agreement_id', 'work_type', 'doc_name', 'required', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { agreement_id: 'partners.partner_agreements' },
        embeds: {}
      },
      partner_agreement_forms: {
        columns: ['id', 'agreement_id', 'work_type', 'form_template_id', 'quantity', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { agreement_id: 'partners.partner_agreements', form_template_id: 'shared.report_templates' },
        embeds: {}
      },
      slas: {
        columns: ['id', 'partner_id', 'name', 'description', 'work_type', 'response_hours', 'resolution_hours', 'compliance_target', 'active', 'valid_from', 'valid_until', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { partner_id: 'partners.partners', work_type: 'domain_gas.visit_types' },
        embeds: {}
      }
    }
  },
  customers: {
    tables: {
      customers: {
        columns: ['id', 'name', 'phone', 'email', 'tax_id', 'created_at', 'updated_at'],
        readOnly: ['id', 'created_at', 'updated_at'],
        fks: {},
        embeds: {
          customer_addresses: 'customer_id',
          technicians: '*'
        }
      },
      customer_addresses: {
        columns: ['id', 'customer_id', 'address_id', 'name', 'type', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { customer_id: 'customers.customers', address_id: 'geocoding.addresses' },
        embeds: {
          customer_address_partners: 'customer_address_id'
        }
      },
      customer_address_partners: {
        columns: ['id', 'customer_address_id', 'partner_id', 'from_date', 'to_date', 'notes', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { customer_address_id: 'customers.customer_addresses', partner_id: 'partners.partners' },
        embeds: {}
      },
      properties: {
        columns: ['id', 'customer_address_id', 'name', 'type', 'notes', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { customer_address_id: 'customers.customer_addresses', type: 'domain_gas.property_types' },
        embeds: {
          property_assets: 'property_id'
        }
      },
      technicians: {
        columns: ['id', 'user_id', 'name', 'is_active', 'specialties', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { user_id: 'core.users(employee_number)' },
        embeds: {
          tech_certifications: 'tech_id'
        }
      },
      tech_certifications: {
        columns: ['id', 'tech_id', 'cert_id', 'number', 'issuer', 'issue_date', 'expiry_date', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { tech_id: 'customers.technicians', cert_id: 'domain_gas.certifications' },
        embeds: {}
      }
    }
  },
  inventory: {
    tables: {
      materials: {
        columns: ['id', 'sku', 'name', 'unit', 'unit_cost', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: {},
        embeds: {}
      },
      tools: {
        columns: ['id', 'code', 'name', 'status', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: {},
        embeds: {}
      },
      epp_items: {
        columns: ['id', 'name', 'type', 'lifecycle', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: {},
        embeds: {}
      },
      vehicles: {
        columns: ['id', 'type', 'license_plate', 'name', 'brand', 'model', 'year', 'status', 'capacity', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { type: 'domain_gas.vehicle_types' },
        embeds: {}
      },
      rentals: {
        columns: ['id', 'type', 'supplier', 'item_description', 'daily_cost', 'hourly_cost', 'fixed_cost', 'notes', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: {},
        embeds: {}
      },
      cost_rates: {
        columns: ['id', 'type', 'reference_id', 'reference_type', 'value', 'unit', 'valid_from', 'valid_until', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: {},
        embeds: {}
      },
      maintenance_records: {
        columns: ['id', 'type', 'reference_id', 'date', 'cost', 'supplier', 'description', 'next_date', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: {},
        embeds: {}
      },
      maintenance_schedules: {
        columns: ['id', 'type', 'reference_id', 'frequency_km', 'frequency_days', 'last_service_date', 'last_service_km', 'active', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: {},
        embeds: {}
      },
      checklist_templates: {
        columns: ['id', 'name', 'work_type', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { work_type: 'domain_gas.visit_types' },
        embeds: {
          checklist_template_materials: 'template_id',
          checklist_template_tools: 'template_id',
          checklist_template_epps: 'template_id'
        }
      }
    }
  },
  operations: {
    tables: {
      visits: {
        columns: ['id', 'property_id', 'partner_id', 'partner_order_id', 'parent_visit_id', 'route_id', 'daily_plan_id', 'type', 'status', 'priority', 'source', 'source_reference', 'billing_to', 'partner_supervisor', 'result', 'rejection_reason_id', 'rejection_reason_detail', 'scheduled_at', 'started_at', 'completed_at', 'alternative_location', 'notes', 'created_at', 'updated_at'],
        readOnly: ['id', 'created_at', 'updated_at'],
        fks: {
          property_id: 'customers.properties',
          partner_id: 'partners.partners',
          parent_visit_id: 'operations.visits',
          route_id: 'planning.routes',
          daily_plan_id: 'planning.daily_plans',
          type: 'domain_gas.visit_types',
          result: 'domain_gas.pre_visit_results',
          rejection_reason_id: 'domain_gas.rejection_reasons'
        },
        embeds: {
          visit_assignments: 'visit_id',
          visit_checklist_materials: 'visit_id',
          visit_checklist_tools: 'visit_id',
          visit_checklist_epps: 'visit_id',
          visit_material_usages: 'visit_id',
          visit_tool_usages: 'visit_id',
          visit_epp_usages: 'visit_id',
          visit_measurements: 'visit_id',
          visit_checkpoints: 'visit_id',
          visit_photos: 'visit_id',
          visit_reports: 'visit_id',
          vehicle_assignments: 'visit_id',
          visit_rentals: 'visit_id',
          visit_sla_trackings: 'visit_id'
        }
      },
      visit_assignments: {
        columns: ['id', 'visit_id', 'tech_id', 'role_id', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { visit_id: 'operations.visits', tech_id: 'customers.technicians', role_id: 'core.tech_roles' },
        embeds: {}
      },
      checklist_template_materials: {
        columns: ['id', 'template_id', 'material_id', 'default_quantity', 'required', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { template_id: 'inventory.checklist_templates', material_id: 'inventory.materials' },
        embeds: {}
      },
      checklist_template_tools: {
        columns: ['id', 'template_id', 'tool_id', 'default_quantity', 'required', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { template_id: 'inventory.checklist_templates', tool_id: 'inventory.tools' },
        embeds: {}
      },
      checklist_template_epps: {
        columns: ['id', 'template_id', 'epp_id', 'default_quantity', 'required', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { template_id: 'inventory.checklist_templates', epp_id: 'inventory.epp_items' },
        embeds: {}
      },
      visit_checklist_materials: {
        columns: ['id', 'visit_id', 'material_id', 'planned_quantity', 'confirmed', 'notes', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { visit_id: 'operations.visits', material_id: 'inventory.materials' },
        embeds: {}
      },
      visit_checklist_tools: {
        columns: ['id', 'visit_id', 'tool_id', 'planned_quantity', 'confirmed', 'notes', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { visit_id: 'operations.visits', tool_id: 'inventory.tools' },
        embeds: {}
      },
      visit_checklist_epps: {
        columns: ['id', 'visit_id', 'epp_id', 'planned_quantity', 'confirmed', 'notes', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { visit_id: 'operations.visits', epp_id: 'inventory.epp_items' },
        embeds: {}
      },
      visit_material_usages: {
        columns: ['id', 'visit_id', 'material_id', 'quantity', 'notes', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { visit_id: 'operations.visits', material_id: 'inventory.materials' },
        embeds: {}
      },
      visit_tool_usages: {
        columns: ['id', 'visit_id', 'tool_id', 'notes', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { visit_id: 'operations.visits', tool_id: 'inventory.tools' },
        embeds: {}
      },
      visit_epp_usages: {
        columns: ['id', 'visit_id', 'epp_id', 'quantity', 'status', 'notes', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { visit_id: 'operations.visits', epp_id: 'inventory.epp_items' },
        embeds: {}
      },
      visit_measurements: {
        columns: ['id', 'visit_id', 'type', 'value', 'unit', 'result', 'measuring_device', 'notes', 'photo_url', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { visit_id: 'operations.visits', type: 'domain_gas.measurement_types' },
        embeds: {}
      },
      visit_checkpoints: {
        columns: ['id', 'visit_id', 'type', 'geom', 'timestamp', 'device_info', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { visit_id: 'operations.visits' },
        embeds: {}
      },
      visit_photos: {
        columns: ['id', 'visit_id', 'checkpoint_id', 'url', 'geom', 'timestamp', 'stage', 'finding_type', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { visit_id: 'operations.visits', checkpoint_id: 'operations.visit_checkpoints', finding_type: 'domain_gas.photo_findings' },
        embeds: {}
      },
      visit_reports: {
        columns: ['id', 'visit_id', 'report_template_id', 'source', 'recorded_at', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { visit_id: 'operations.visits', report_template_id: 'shared.report_templates' },
        embeds: {
          report_images: 'report_id',
          report_entries: 'report_id'
        }
      },
      report_images: {
        columns: ['id', 'report_id', 'url', 'pages', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { report_id: 'operations.visit_reports' },
        embeds: {}
      },
      report_entries: {
        columns: ['id', 'report_id', 'data_json', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { report_id: 'operations.visit_reports' },
        embeds: {}
      },
      vehicle_assignments: {
        columns: ['id', 'vehicle_id', 'daily_plan_id', 'visit_id', 'departure_time', 'return_time', 'departure_mileage', 'return_mileage', 'notes', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { vehicle_id: 'inventory.vehicles', daily_plan_id: 'planning.daily_plans', visit_id: 'operations.visits' },
        embeds: {}
      },
      visit_rentals: {
        columns: ['id', 'visit_id', 'rental_id', 'start_time', 'end_time', 'hours', 'total_cost', 'reason', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { visit_id: 'operations.visits', rental_id: 'inventory.rentals' },
        embeds: {}
      },
      visit_sla_trackings: {
        columns: ['id', 'visit_id', 'sla_id', 'requested_at', 'responded_at', 'resolved_at', 'response_time_hours', 'resolution_time_hours', 'meets_response_sla', 'meets_resolution_sla', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { visit_id: 'operations.visits', sla_id: 'partners.slas' },
        embeds: {}
      }
    }
  },
  planning: {
    tables: {
      daily_plans: {
        columns: ['id', 'date', 'notes', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: {},
        embeds: {
          daily_plan_assignments: 'daily_plan_id',
          daily_load_materials: 'daily_plan_id',
          daily_load_tools: 'daily_plan_id',
          daily_load_epps: 'daily_plan_id',
          routes: 'daily_plan_id'
        }
      },
      daily_plan_assignments: {
        columns: ['id', 'daily_plan_id', 'tech_id', 'role_id', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { daily_plan_id: 'planning.daily_plans', tech_id: 'customers.technicians', role_id: 'core.tech_roles' },
        embeds: {}
      },
      daily_load_materials: {
        columns: ['id', 'daily_plan_id', 'material_id', 'loaded_quantity', 'returned_quantity', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { daily_plan_id: 'planning.daily_plans', material_id: 'inventory.materials' },
        embeds: {}
      },
      daily_load_tools: {
        columns: ['id', 'daily_plan_id', 'tool_id', 'loaded_quantity', 'returned_quantity', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { daily_plan_id: 'planning.daily_plans', tool_id: 'inventory.tools' },
        embeds: {}
      },
      daily_load_epps: {
        columns: ['id', 'daily_plan_id', 'epp_id', 'loaded_quantity', 'returned_quantity', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { daily_plan_id: 'planning.daily_plans', epp_id: 'inventory.epp_items' },
        embeds: {}
      },
      routes: {
        columns: ['id', 'type', 'date', 'daily_plan_id', 'notes', 'status', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { type: 'domain_gas.route_types', daily_plan_id: 'planning.daily_plans' },
        embeds: {}
      }
    }
  },
  notifications: {
    tables: {
      notification_templates: {
        columns: ['id', 'type', 'channel', 'subject', 'body', 'active', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: {},
        embeds: {}
      },
      notifications: {
        columns: ['id', 'type', 'channel', 'recipient', 'subject', 'body', 'status', 'entity_type', 'entity_id', 'scheduled_for', 'sent_at', 'read_at', 'attempts', 'error', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: {},
        embeds: {}
      }
    }
  },
  geocoding: {
    tables: {
      addresses: {
        columns: ['id', 'street', 'number', 'apartment', 'neighborhood', 'city', 'region', 'location_references', 'postal_code', 'geom', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: {},
        embeds: {}
      }
    }
  },
  shared: {
    tables: {
      report_templates: {
        columns: ['id', 'name', 'fields_json', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: {},
        embeds: {}
      }
    }
  },
  domain_gas: {
    tables: {
      visit_types: {
        columns: ['id', 'code', 'name', 'active'],
        readOnly: ['id'],
        fks: {},
        embeds: {},
        readOnlyTable: true
      },
      measurement_types: {
        columns: ['id', 'code', 'name', 'default_unit', 'active'],
        readOnly: ['id'],
        fks: {},
        embeds: {},
        readOnlyTable: true
      },
      photo_findings: {
        columns: ['id', 'code', 'name', 'active'],
        readOnly: ['id'],
        fks: {},
        embeds: {},
        readOnlyTable: true
      },
      vehicle_types: {
        columns: ['id', 'code', 'name', 'active'],
        readOnly: ['id'],
        fks: {},
        embeds: {},
        readOnlyTable: true
      },
      partner_service_types: {
        columns: ['id', 'code', 'name', 'active'],
        readOnly: ['id'],
        fks: {},
        embeds: {},
        readOnlyTable: true
      },
      pre_visit_results: {
        columns: ['id', 'code', 'name', 'active'],
        readOnly: ['id'],
        fks: {},
        embeds: {},
        readOnlyTable: true
      },
      route_types: {
        columns: ['id', 'code', 'name', 'active'],
        readOnly: ['id'],
        fks: {},
        embeds: {},
        readOnlyTable: true
      },
      property_types: {
        columns: ['id', 'code', 'name', 'columns_config', 'active'],
        readOnly: ['id'],
        fks: {},
        embeds: {},
        readOnlyTable: true
      },
      certifications: {
        columns: ['id', 'code', 'name', 'class', 'description', 'active', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: {},
        embeds: {},
        readOnlyTable: true
      },
      rejection_reasons: {
        columns: ['id', 'code', 'name', 'category', 'active', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: {},
        embeds: {},
        readOnlyTable: true
      },
      property_assets: {
        columns: ['id', 'property_id', 'type', 'value', 'created_at'],
        readOnly: ['id', 'created_at'],
        fks: { property_id: 'customers.properties' },
        embeds: {}
      }
    }
  }
};

// ============================================================================
// JWT GENERATOR SCRIPT (embedded in collection)
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
  var secret = pm.environment.get('PGRST_JWT_SECRET');
  if (!secret) { console.error('PGRST_JWT_SECRET not set'); return; }
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
// REQUEST GENERATORS
// ============================================================================

function makeUrl(schema, table, params) {
  let url = `{{base_url}}/${schema}/${table}`;
  if (params) url += `?${params}`;
  return url;
}

function makeRequest(method, url, body, description) {
  const req = {
    method,
    header: [
      { key: 'Authorization', value: 'Bearer {{jwt_token}}', type: 'text' },
      { key: 'Content-Type', value: 'application/json', type: 'text' }
    ],
    url: { raw: url, host: ['{{base_url}}'], path: url.replace('{{base_url}}/', '').split('/') },
    description
  };
  if (body) {
    req.body = { mode: 'raw', raw: JSON.stringify(body, null, 2) };
  }
  return req;
}

function makePostRequest(schema, table, tableDef) {
  const writableCols = tableDef.columns.filter(c => !tableDef.readOnly.includes(c));
  const example = {};
  writableCols.forEach(c => {
    if (c === 'name' || c === 'doc_name' || c === 'subject') example[c] = `Example ${table}`;
    else if (c === 'email') example[c] = 'test@example.com';
    else if (c === 'phone' || c === 'recipient') example.c = '+56912345678';
    else if (c === 'code' || c === 'sku') example[c] = 'CODE001';
    else if (c === 'status') example[c] = 'active';
    else if (c === 'active') example[c] = true;
    else if (c === 'required') example[c] = true;
    else if (c === 'confirmed') example[c] = false;
    else if (c === 'quantity' || c === 'default_quantity' || c === 'planned_quantity' || c === 'loaded_quantity') example[c] = 1;
    else if (c === 'returned_quantity') example[c] = 0;
    else if (c === 'rate' || c === 'value' || c === 'unit_cost' || c === 'daily_cost' || c === 'hourly_cost' || c === 'fixed_cost' || c === 'total_cost' || c === 'cost') example[c] = 10000;
    else if (c === 'response_hours' || c === 'resolution_hours') example[c] = 24;
    else if (c === 'year') example[c] = 2024;
    else if (c === 'from_date') example[c] = '2024-01-01';
    else if (c === 'frequency_km' || c === 'frequency_days') example[c] = 30;
    else if (c === 'specialties') example[c] = ['gas', 'maintenance'];
    else if (c === 'fields_json') example[c] = [];
    else if (c === 'data_json') example[c] = {};
    else if (c === 'columns_config') example[c] = {};
    else if (c === 'geom') example[c] = 'SRID=4326;POINT(-70.6483 -33.4569)';
    else if (c === 'scheduled_at' || c === 'start_time' || c === 'departure_time' || c === 'timestamp' || c === 'recorded_at' || c === 'requested_at') example[c] = '2024-01-15T08:00:00Z';
    else if (c.endsWith('_at') || c.endsWith('_date')) example[c] = null;
    else if (c.endsWith('_id')) example[c] = '00000000-0000-0000-0000-000000000000';
    else if (c === 'notes' || c === 'description' || c === 'body' || c === 'reason' || c === 'rejection_reason_detail' || c === 'alternative_location' || c === 'partner_supervisor' || c === 'source_reference' || c === 'partner_order_id' || c === 'device_info' || c === 'error') example[c] = null;
    else example[c] = null;
  });
  const url = makeUrl(schema, table);
  return {
    name: `Create ${table}`,
    request: makeRequest('POST', url, example, `POST /${schema}/${table} — Create new ${table} record`),
    response: []
  };
}

function makePatchRequest(schema, table, tableDef) {
  const writableCols = tableDef.columns.filter(c => !tableDef.readOnly.includes(c));
  const example = {};
  if (writableCols.includes('name')) example.name = `Updated ${table}`;
  else if (writableCols.includes('status')) example.status = 'active';
  else if (writableCols.includes('notes')) example.notes = 'Updated notes';
  else if (writableCols.length > 0) example[writableCols[0]] = null;
  const url = makeUrl(schema, table, 'id=eq.{{id}}');
  return {
    name: `Update ${table}`,
    request: makeRequest('PATCH', url, example, `PATCH /${schema}/${table}?id=eq.{{id}} — Update ${table} by ID`),
    response: []
  };
}

function makeDeleteRequest(schema, table) {
  const url = makeUrl(schema, table, 'id=eq.{{id}}');
  return {
    name: `Delete ${table}`,
    request: makeRequest('DELETE', url, null, `DELETE /${schema}/${table}?id=eq.{{id}} — Delete ${table} by ID`),
    response: []
  };
}

function makeListRequest(schema, table) {
  const url = makeUrl(schema, table, 'limit=20&offset=0&order=created_at.desc');
  return {
    name: `List ${table}`,
    request: makeRequest('GET', url, null, `GET /${schema}/${table} — List with pagination, latest first`),
    response: []
  };
}

function makeListSelectRequest(schema, table, tableDef) {
  const cols = tableDef.columns.filter(c => !c.startsWith('password'));
  const url = makeUrl(schema, table, `select=${cols.join(',')}`);
  return {
    name: `List ${table} (select *)`,
    request: makeRequest('GET', url, null, `GET /${schema}/${table}?select=* — List all columns explicitly`),
    response: []
  };
}

function makeEmbedRequest(schema, table, tableDef) {
  const embedParts = Object.entries(tableDef.embeds).map(([child, fk]) => `${child}(*)`);
  if (embedParts.length === 0) return null;
  const selectCols = tableDef.columns.filter(c => !c.startsWith('password') && !tableDef.readOnly.includes(c));
  const select = [...selectCols, ...embedParts].join(',');
  const url = makeUrl(schema, table, `select=${select}`);
  return {
    name: `Embed ${table} → ${Object.keys(tableDef.embeds).join(', ')}`,
    request: makeRequest('GET', url, null, `GET /${schema}/${table}?select=*,<related>(*) — Embed related ${Object.keys(tableDef.embeds).length} tables`),
    response: []
  };
}

function makeGetByIdRequest(schema, table) {
  const url = makeUrl(schema, table, 'id=eq.{{id}}');
  return {
    name: `Get ${table} by ID`,
    request: makeRequest('GET', url, null, `GET /${schema}/${table}?id=eq.{{id}} — Get single ${table} by ID`),
    response: []
  };
}

function makeGeoQueryRequest() {
  const url = makeUrl('geocoding', 'addresses', 'select=*,st_distance(geom,ST_SetSRID(ST_MakePoint(-70.6483,-33.4569),4326)::geography)&geom=st_dwithin(geom,ST_SetSRID(ST_MakePoint(-70.6483,-33.4569),4326)::geography,1000)&order=geom asc');
  return {
    name: 'Geo Query: Within 1km of Santiago',
    request: makeRequest('GET', url, null, 'GET /geocoding/addresses?select=*,st_distance(...)&geom=st_dwithin(...,1000) — PostGIS spatial query'),
    response: []
  };
}

function makeVisitsFullEmbedRequest() {
  const embed = [
    'visit_assignments(*,technician:tech_id(name,specialties))',
    'visit_checklist_materials(*,material:material_id(name,unit,unit_cost))',
    'visit_checklist_tools(*,tool:tool_id(name,code))',
    'visit_checklist_epps(*,epp:epp_id(name,type,lifecycle))',
    'visit_material_usages(*,material:material_id(name))',
    'visit_tool_usages(*,tool:tool_id(name))',
    'visit_epp_usages(*,epp:epp_id(name))',
    'visit_measurements(*,mtype:type(code,name,default_unit))',
    'visit_checkpoints(*)',
    'visit_photos(*,finding:finding_type(code,name))',
    'visit_reports(*,report_images(*),report_entries(*))',
    'vehicle_assignments(*,vehicle:vehicle_id(name,license_plate,vtype:type(code,name)))',
    'visit_rentals(*,rental:rental_id(item_description,daily_cost))',
    'visit_sla_trackings(*,sla:sla_id(name,response_hours,resolution_hours))'
  ].join(',');
  const url = makeUrl('operations', 'visits', `select=*,${embed}&limit=5`);
  return {
    name: 'Visits FULL Embed (14 relationships)',
    request: makeRequest('GET', url, null, 'GET /operations/visits?select=*,<14 embeds> — Full visit with all related data (assignments, checklists, measurements, photos, reports, vehicles, rentals, SLAs)'),
    response: []
  };
}

// ============================================================================
// COLLECTION BUILDER
// ============================================================================

function buildCollection() {
  const collection = {
    info: {
      name: 'GAS-FSM PostgREST API',
      description: 'Field Service Management API for gas installation, maintenance, and certification.\n\n**Stack**: PostgreSQL 16 + PostGIS 3.5 + PostgREST 12 + GLAuth + Authelia + Traefik\n**Mode**: Single-tenant, JWT auth via Traefik → PostgREST\n\n## Quick Start\n1. Select environment: GAS-FSM-dev\n2. Run "Auth: Generate JWT" folder\n3. Start querying endpoints\n\n## Auth\nJWT is generated via pre-request script using PGRST_JWT_SECRET.\nToken is cached in jwt_token variable and auto-refreshed on expiry.\n\n**JWT payload**: {sub: employee_number, role: "admin|operator|technician"}',
      schema: 'https://schema.getpostman.com/json/collection/v2.1.0/collection.json'
    },
    variable: [
      { key: 'base_url', value: '{{base_url}}' },
      { key: 'jwt_token', value: '{{jwt_token}}' }
    ],
    auth: { type: 'bearer', bearer: [{ key: 'token', value: '{{jwt_token}}' }] },
    event: [
      {
        listen: 'prerequest',
        script: { exec: [JWT_SCRIPT], type: 'text/javascript' }
      }
    ],
    item: []
  };

  // Auth folder
  collection.item.push({
    name: 'Auth',
    description: 'Authentication and JWT generation',
    item: [
      {
        name: 'Generate JWT (dev)',
        description: 'Generates a new JWT token using PGRST_JWT_SECRET from environment. Token is cached and auto-refreshed.',
        request: {
          method: 'GET',
          header: [],
          url: { raw: '{{base_url}}/', host: ['{{base_url}}'], path: [''] },
          description: 'Triggers pre-request script to generate JWT. Does not make an actual API call.'
        },
        event: [
          { listen: 'prerequest', script: { exec: [JWT_SCRIPT], type: 'text/javascript' } }
        ],
        response: []
      }
    ]
  });

  // Schema folders
  const schemaOrder = ['core', 'partners', 'customers', 'inventory', 'operations', 'planning', 'notifications', 'geocoding', 'shared', 'domain_gas'];
  const schemaDescriptions = {
    core: 'Core context: users, roles, audit',
    partners: 'Partners context: gas companies and agreements',
    customers: 'Customers context: end clients, addresses, properties, technicians',
    inventory: 'Inventory context: materials, tools, PPE, vehicles, costs, maintenance',
    operations: 'Operations context: visits, assignments, checklists, measurements, checkpoints, photos, reports',
    planning: 'Planning context: daily plans, technician assignments, material loading, routes',
    notifications: 'Notifications context: templates and notification instances',
    geocoding: 'Geocoding context: global georeferenced address catalog',
    shared: 'Shared global catalogs (report_templates, rejection_category)',
    domain_gas: 'Gas Chile specific domain: lookup tables, SEC certifications, rejection reasons'
  };

  schemaOrder.forEach(schemaName => {
    const schema = SCHEMAS[schemaName];
    const schemaFolder = {
      name: schemaName,
      description: schemaDescriptions[schemaName],
      item: []
    };

    Object.entries(schema.tables).forEach(([tableName, tableDef]) => {
      const tableFolder = {
        name: tableName,
        description: `Table: ${schemaName}.${tableName} (${tableDef.columns.length} columns, ${Object.keys(tableDef.embeds).length} embeds)`,
        item: []
      };

      // Always add list and get by ID
      tableFolder.item.push(makeListRequest(schemaName, tableName));
      tableFolder.item.push(makeListSelectRequest(schemaName, tableName, tableDef));
      tableFolder.item.push(makeGetByIdRequest(schemaName, tableName));

      // Add embed request if available
      const embedReq = makeEmbedRequest(schemaName, tableName, tableDef);
      if (embedReq) tableFolder.item.push(embedReq);

      // Add CRUD for non-read-only tables
      if (!tableDef.readOnlyTable) {
        tableFolder.item.push(makePostRequest(schemaName, tableName, tableDef));
        tableFolder.item.push(makePatchRequest(schemaName, tableName, tableDef));
        tableFolder.item.push(makeDeleteRequest(schemaName, tableName));
      }

      schemaFolder.item.push(tableFolder);
    });

    // Special requests
    if (schemaName === 'geocoding') {
      schemaFolder.item.push(makeGeoQueryRequest());
    }

    collection.item.push(schemaFolder);
  });

  // Add special Operations requests at the end
  const opsFolder = collection.item.find(f => f.name === 'operations');
  if (opsFolder) {
    opsFolder.item.push({
      name: 'SPECIAL: Full Visit Embed',
      description: 'Special request with all 14 relationship embeds for visits',
      item: [makeVisitsFullEmbedRequest()]
    });
  }

  return collection;
}

// ============================================================================
// MAIN
// ============================================================================

const collection = buildCollection();
const output = JSON.stringify(collection, null, 2);

// Write to stdout or file
if (process.argv[2]) {
  fs.writeFileSync(process.argv[2], output, 'utf8');
  console.error(`Collection written to ${process.argv[2]}`);
  console.error(`Total items: ${JSON.stringify(collection).split('"name":').length - 1}`);
} else {
  process.stdout.write(output);
}
