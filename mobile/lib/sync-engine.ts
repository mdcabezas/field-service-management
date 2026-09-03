import { getDatabase } from './database';
import { apiRequest } from './auth';
import { uploadPhoto } from './photo-upload';

export { addToOutbox } from './outbox';

interface SyncState {
  lastSync: string | null;
  isOnline: boolean;
  isSyncing: boolean;
  pendingItems: number;
}

export async function getSyncState(): Promise<SyncState> {
  console.log('[Sync] getSyncState called');
  const db = getDatabase();
  const stmt = db.prepareSync('SELECT value FROM sync_meta WHERE key = ?');
  const result = stmt.executeSync('last_sync');
  const row = result.getFirstSync() as { value: string } | null;
  stmt.finalizeSync();
  const lastSync = row?.value || null;

  const countStmt = db.prepareSync('SELECT COUNT(*) as count FROM outbox WHERE status = ?');
  const countResult = countStmt.executeSync('pending');
  const countRow = countResult.getFirstSync() as { count: number } | null;
  countStmt.finalizeSync();
  const pendingItems = countRow?.count || 0;

  console.log('[Sync] getSyncState:', { lastSync, pendingItems });
  return {
    lastSync,
    isOnline: true,
    isSyncing: false,
    pendingItems,
  };
}

export async function fullSync(): Promise<void> {
  console.log('[Sync] fullSync: START');
  const db = getDatabase();

  console.log('[Sync] fullSync: pushing outbox before pull');
  await pushOutbox(db);

  console.log('[Sync] fullSync: pulling from /api/mobile/sync');
  let refData: any;
  try {
    refData = await apiRequest<any>('/api/mobile/sync');
    console.log('[Sync] fullSync: pull OK', {
      hasReferenceData: !!refData.reference_data,
      visitsCount: refData.visits?.length || 0,
      checklistCount: refData.checklist_items?.length || 0,
    });
  } catch (error) {
    console.error('[Sync] fullSync: pull FAILED', error);
    throw error;
  }

  if (refData.reference_data) {
    console.log('[Sync] fullSync: saving reference data');

    for (const mat of refData.reference_data.materials || []) {
      await db.runAsync(
        'INSERT OR REPLACE INTO materials (id, name, sku, unit, unit_cost) VALUES (?, ?, ?, ?, ?)',
        mat.id, mat.name, mat.sku, mat.unit, mat.unit_cost
      );
    }
    console.log('[Sync] fullSync: materials saved:', refData.reference_data.materials?.length || 0);

    for (const tool of refData.reference_data.tools || []) {
      await db.runAsync(
        'INSERT OR REPLACE INTO tools (id, name, sku, description) VALUES (?, ?, ?, ?)',
        tool.id, tool.name, tool.sku, tool.description
      );
    }
    console.log('[Sync] fullSync: tools saved:', refData.reference_data.tools?.length || 0);

    for (const epp of refData.reference_data.epp || []) {
      await db.runAsync(
        'INSERT OR REPLACE INTO epp (id, name, sku, description) VALUES (?, ?, ?, ?)',
        epp.id, epp.name, epp.sku, epp.description
      );
    }
    console.log('[Sync] fullSync: epp saved:', refData.reference_data.epp?.length || 0);

    for (const tmpl of refData.reference_data.templates || []) {
      await db.runAsync(
        'INSERT OR REPLACE INTO checklist_templates (id, name, description, visit_type, work_type, partner_id, stages_enabled) VALUES (?, ?, ?, ?, ?, ?, ?)',
        tmpl.id, tmpl.name, tmpl.description, tmpl.visit_type, tmpl.work_type, tmpl.partner_id, tmpl.stages_enabled ? 1 : 0
      );
    }
    console.log('[Sync] fullSync: templates saved:', refData.reference_data.templates?.length || 0);

    for (const partner of refData.reference_data.partners || []) {
      await db.runAsync(
        'INSERT OR REPLACE INTO partners (id, name) VALUES (?, ?)',
        partner.id, partner.name
      );
    }
    console.log('[Sync] fullSync: partners saved:', refData.reference_data.partners?.length || 0);
  }

  for (const visit of refData.visits || []) {
    await db.runAsync(
      `INSERT OR REPLACE INTO visits (id, property_id, property_name, property_address, property_lat, property_lng,
       partner_id, partner_name, partner_order_id, partner_supervisor, type, status, priority, scheduled_at,
       checklist_template_id, stages_enabled, selected_stages, sla_response_hours, sla_resolution_hours)
       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
      visit.id, visit.property_id, visit.property_name, visit.property_address,
      visit.property_lat, visit.property_lng, visit.partner_id, visit.partner_name,
      visit.partner_order_id, visit.partner_supervisor, visit.type, visit.status,
      visit.priority, visit.scheduled_at, visit.checklist_template_id,
      visit.stages_enabled ? 1 : 0, JSON.stringify(visit.selected_stages),
      visit.sla_response_hours, visit.sla_resolution_hours
    );
  }
  console.log('[Sync] fullSync: visits saved:', refData.visits?.length || 0);

  for (const item of refData.checklist_items || []) {
    await db.runAsync(
      'INSERT OR REPLACE INTO checklist_items (id, visit_id, stage_id, stage_name, type, ref_id, name, quantity, unit, confirmed, required) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)',
      item.id, item.visit_id, item.stage_id, item.stage_name, item.type, item.ref_id, item.name, item.quantity, item.unit, item.confirmed ? 1 : 0, item.required ? 1 : 0
    );
  }
  console.log('[Sync] fullSync: checklist_items saved:', refData.checklist_items?.length || 0);

  await db.runAsync('INSERT OR REPLACE INTO sync_meta (key, value) VALUES (?, ?)', 'last_sync', new Date().toISOString());
  console.log('[Sync] fullSync: DONE');
}

export async function deltaSync(): Promise<void> {
  console.log('[Sync] deltaSync: START');
  const db = getDatabase();

  console.log('[Sync] deltaSync: pushing outbox');
  await pushOutbox(db);

  const lastSyncStmt = db.prepareSync('SELECT value FROM sync_meta WHERE key = ?');
  const lastSyncResult = lastSyncStmt.executeSync('last_sync');
  const lastSyncRow = lastSyncResult.getFirstSync() as { value: string } | null;
  lastSyncStmt.finalizeSync();
  const lastSync = lastSyncRow?.value || null;
  console.log('[Sync] deltaSync: lastSync =', lastSync);

  const pullPath = lastSync ? `/api/mobile/sync?since=${encodeURIComponent(lastSync)}` : '/api/mobile/sync';
  console.log('[Sync] deltaSync: pulling from', pullPath);

  let pullData: any;
  try {
    pullData = await apiRequest<any>(pullPath);
    console.log('[Sync] deltaSync: pull OK', {
      visitsCount: pullData.visits?.length || 0,
      checklistCount: pullData.checklist_items?.length || 0,
    });
  } catch (error) {
    console.error('[Sync] deltaSync: pull FAILED', error);
    throw error;
  }

  for (const visit of pullData.visits || []) {
    await db.runAsync(
      `INSERT OR REPLACE INTO visits (id, property_id, property_name, property_address, property_lat, property_lng,
       partner_id, partner_name, partner_order_id, partner_supervisor, type, status, priority, scheduled_at,
       checklist_template_id, stages_enabled, selected_stages, sla_response_hours, sla_resolution_hours)
       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
      visit.id, visit.property_id, visit.property_name, visit.property_address,
      visit.property_lat, visit.property_lng, visit.partner_id, visit.partner_name,
      visit.partner_order_id, visit.partner_supervisor, visit.type, visit.status,
      visit.priority, visit.scheduled_at, visit.checklist_template_id,
      visit.stages_enabled ? 1 : 0, JSON.stringify(visit.selected_stages),
      visit.sla_response_hours, visit.sla_resolution_hours
    );
  }

  for (const item of pullData.checklist_items || []) {
    await db.runAsync(
      'INSERT OR REPLACE INTO checklist_items (id, visit_id, stage_id, stage_name, type, ref_id, name, quantity, unit, confirmed, required) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)',
      item.id, item.visit_id, item.stage_id, item.stage_name, item.type, item.ref_id, item.name, item.quantity, item.unit, item.confirmed ? 1 : 0, item.required ? 1 : 0
    );
  }

  await db.runAsync('INSERT OR REPLACE INTO sync_meta (key, value) VALUES (?, ?)', 'last_sync', new Date().toISOString());
  console.log('[Sync] deltaSync: DONE');
}

async function pushOutbox(db: ReturnType<typeof getDatabase>): Promise<void> {
  console.log('[Sync] pushOutbox: START');
  const stmt = db.prepareSync(
    'SELECT * FROM outbox WHERE status = ? ORDER BY CASE entity_type WHEN ? THEN 1 WHEN ? THEN 2 WHEN ? THEN 3 WHEN ? THEN 4 WHEN ? THEN 5 WHEN ? THEN 6 WHEN ? THEN 7 WHEN ? THEN 8 END, created_at'
  );
  const result = stmt.executeSync('pending', 'status', 'checkpoint', 'gps_waiver', 'photo', 'usage', 'measurement', 'checklist', 'report');
  const items = result.getAllSync() as Record<string, unknown>[];
  stmt.finalizeSync();

  console.log('[Sync] pushOutbox: pending items =', items.length);
  if (!items.length) {
    console.log('[Sync] pushOutbox: nothing to push, DONE');
    return;
  }

  const pushData: any = {
    status_changes: [],
    checkpoints: [],
    gps_waivers: [],
    photos: [],
    usages: [],
    measurements: [],
    checklist_items: [],
    reports: [],
  };

  for (const item of items) {
    const payload = JSON.parse(item.payload as string);
    switch (item.entity_type) {
      case 'status':
        pushData.status_changes.push(payload);
        break;
      case 'checkpoint':
        pushData.checkpoints.push(payload);
        break;
      case 'gps_waiver':
        pushData.gps_waivers.push(payload);
        break;
      case 'photo':
        pushData.photos.push(payload);
        break;
      case 'usage':
        pushData.usages.push(payload);
        break;
      case 'measurement':
        pushData.measurements.push(payload);
        break;
      case 'checklist':
        pushData.checklist_items.push(payload);
        break;
      case 'report':
        pushData.reports.push(payload);
        break;
    }
  }

  console.log('[Sync] pushOutbox: sending POST /api/mobile/sync', {
    status_changes: pushData.status_changes.length,
    checkpoints: pushData.checkpoints.length,
    gps_waivers: pushData.gps_waivers.length,
    photos: pushData.photos.length,
    usages: pushData.usages.length,
    measurements: pushData.measurements.length,
    checklist_items: pushData.checklist_items.length,
    reports: pushData.reports.length,
  });

  try {
    await apiRequest<any>('/api/mobile/sync', {
      method: 'POST',
      body: JSON.stringify(pushData),
    });
    console.log('[Sync] pushOutbox: POST OK, uploading photo binaries');

    // Upload binary photo files BEFORE marking outbox as sent
    const photoItems = items.filter((i) => i.entity_type === 'photo');
    if (photoItems.length > 0) {
      console.log('[Sync] pushOutbox: uploading', photoItems.length, 'photo files');
      for (const item of photoItems) {
        try {
          const uploaded = await uploadPhoto(item.entity_id as string);
          if (!uploaded) {
            throw new Error('photo binary upload returned false');
          }
          console.log('[Sync] pushOutbox: uploadPhoto OK', { entityId: item.entity_id });
        } catch (e) {
          console.error('[Sync] pushOutbox: uploadPhoto FAILED for', item.entity_id, e);
          // Mark all items as failed so they retry next time
          for (const failItem of items) {
            await db.runAsync(
              'UPDATE outbox SET status = ?, error = ?, retry_count = retry_count + 1 WHERE id = ?',
              'failed', String(e), failItem.id as string
            );
          }
          console.log('[Sync] pushOutbox: marked items as failed due to photo upload error');
          return;
        }
      }
    }

    // Only mark as sent after all photo binaries uploaded successfully
    for (const item of items) {
      await db.runAsync('UPDATE outbox SET status = ? WHERE id = ?', 'sent', item.id as string);
    }
    console.log('[Sync] pushOutbox: DONE');
  } catch (error) {
    console.error('[Sync] pushOutbox: POST FAILED', error);
    for (const item of items) {
      await db.runAsync(
        'UPDATE outbox SET status = ?, error = ?, retry_count = retry_count + 1 WHERE id = ?',
        'failed', String(error), item.id as string
      );
    }
    console.log('[Sync] pushOutbox: marked items as failed');
  }
}

export async function retryFailedItems(): Promise<{ retried: number; failed: number }> {
  console.log('[Sync] retryFailedItems called');
  const db = getDatabase();
  const stmt = db.prepareSync('SELECT * FROM outbox WHERE status = ?');
  const result = stmt.executeSync('failed');
  const failedItems = result.getAllSync() as Record<string, unknown>[];
  stmt.finalizeSync();

  if (!failedItems.length) {
    console.log('[Sync] retryFailedItems: no failed items');
    return { retried: 0, failed: 0 };
  }

  let retried = 0;
  let failed = 0;

  for (const item of failedItems) {
    try {
      await db.runAsync('UPDATE outbox SET status = ? WHERE id = ?', 'pending', item.id as string);
      retried++;
    } catch {
      failed++;
    }
  }

  console.log('[Sync] retryFailedItems: retried =', retried, 'failed =', failed);
  return { retried, failed };
}
