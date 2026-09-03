import * as SQLite from 'expo-sqlite';

let db: SQLite.SQLiteDatabase | null = null;

export async function initDatabase(): Promise<SQLite.SQLiteDatabase> {
  if (!db) {
    db = await SQLite.openDatabaseAsync('localis-mobile');
    await initializeDatabase(db);
  }
  return db;
}

export function getDatabase(): SQLite.SQLiteDatabase {
  if (!db) {
    throw new Error('Database not initialized. Call initDatabase() first.');
  }
  return db;
}

async function initializeDatabase(database: SQLite.SQLiteDatabase): Promise<void> {
  const pragmas = [
    'PRAGMA journal_mode = MEMORY',
    'PRAGMA synchronous = NORMAL',
    'PRAGMA cache_size = 2000',
    'PRAGMA temp_store = MEMORY',
    'PRAGMA mmap_size = 67108864',
  ];

  for (const pragma of pragmas) {
    try {
      await database.execAsync(pragma);
    } catch (e) {
      console.warn('PRAGMA failed:', pragma, e);
    }
  }

  await database.execAsync(`
    CREATE TABLE IF NOT EXISTS materials (
      id TEXT PRIMARY KEY,
      name TEXT NOT NULL,
      sku TEXT,
      unit TEXT,
      unit_cost REAL
    )
  `);

  await database.execAsync(`
    CREATE TABLE IF NOT EXISTS tools (
      id TEXT PRIMARY KEY,
      name TEXT NOT NULL,
      sku TEXT,
      description TEXT
    )
  `);

  await database.execAsync(`
    CREATE TABLE IF NOT EXISTS epp (
      id TEXT PRIMARY KEY,
      name TEXT NOT NULL,
      sku TEXT,
      description TEXT
    )
  `);

  await database.execAsync(`
    CREATE TABLE IF NOT EXISTS checklist_templates (
      id TEXT PRIMARY KEY,
      name TEXT NOT NULL,
      description TEXT,
      visit_type TEXT,
      work_type TEXT,
      partner_id TEXT,
      stages_enabled INTEGER DEFAULT 0
    )
  `);

  await database.execAsync(`
    CREATE TABLE IF NOT EXISTS checklist_stages (
      id TEXT PRIMARY KEY,
      template_id TEXT NOT NULL,
      name TEXT NOT NULL,
      sort_order INTEGER DEFAULT 0,
      required INTEGER DEFAULT 0
    )
  `);

  await database.execAsync(`
    CREATE TABLE IF NOT EXISTS partners (
      id TEXT PRIMARY KEY,
      name TEXT NOT NULL
    )
  `);

  await database.execAsync(`
    CREATE TABLE IF NOT EXISTS partner_contacts (
      id TEXT PRIMARY KEY,
      partner_id TEXT NOT NULL,
      name TEXT NOT NULL,
      role TEXT,
      phone TEXT,
      email TEXT
    )
  `);

  await database.execAsync(`
    CREATE TABLE IF NOT EXISTS visits (
      id TEXT PRIMARY KEY,
      property_id TEXT,
      property_name TEXT,
      property_address TEXT,
      property_lat REAL,
      property_lng REAL,
      partner_id TEXT,
      partner_name TEXT,
      partner_order_id TEXT,
      partner_supervisor TEXT,
      type TEXT,
      status TEXT DEFAULT 'assigned',
      priority TEXT DEFAULT 'normal',
      scheduled_at TEXT,
      checklist_template_id TEXT,
      stages_enabled INTEGER DEFAULT 0,
      selected_stages TEXT,
      sla_response_hours INTEGER,
      sla_resolution_hours INTEGER,
      created_at TEXT DEFAULT (datetime('now')),
      updated_at TEXT DEFAULT (datetime('now'))
    )
  `);

  await database.execAsync(`
    CREATE TABLE IF NOT EXISTS checklist_items (
      id TEXT PRIMARY KEY,
      visit_id TEXT NOT NULL,
      stage_id TEXT,
      stage_name TEXT,
      type TEXT NOT NULL,
      ref_id TEXT NOT NULL,
      name TEXT NOT NULL,
      quantity REAL DEFAULT 1,
      unit TEXT,
      confirmed INTEGER DEFAULT 0,
      required INTEGER DEFAULT 0,
      created_at TEXT DEFAULT (datetime('now'))
    )
  `);

  await database.execAsync(`
    CREATE TABLE IF NOT EXISTS checkpoints (
      id TEXT PRIMARY KEY,
      visit_id TEXT NOT NULL,
      type TEXT NOT NULL,
      timestamp TEXT NOT NULL,
      lat REAL,
      lng REAL,
      accuracy REAL,
      distance_from_property REAL,
      synced INTEGER DEFAULT 0,
      created_at TEXT DEFAULT (datetime('now'))
    )
  `);

  await database.execAsync(`
    CREATE TABLE IF NOT EXISTS gps_waivers (
      id TEXT PRIMARY KEY,
      visit_id TEXT NOT NULL,
      type TEXT NOT NULL,
      timestamp TEXT NOT NULL,
      accepted_by TEXT,
      synced INTEGER DEFAULT 0,
      created_at TEXT DEFAULT (datetime('now'))
    )
  `);

  await database.execAsync(`
    CREATE TABLE IF NOT EXISTS photos (
      id TEXT PRIMARY KEY,
      visit_id TEXT NOT NULL,
      local_path TEXT NOT NULL,
      thumbnail_path TEXT,
      server_id TEXT,
      uploaded INTEGER DEFAULT 0,
      file_size INTEGER DEFAULT 0,
      lat REAL,
      lng REAL,
      timestamp TEXT,
      created_at TEXT DEFAULT (datetime('now'))
    )
  `);

  await database.execAsync(`
    CREATE TABLE IF NOT EXISTS measurements (
      id TEXT PRIMARY KEY,
      visit_id TEXT NOT NULL,
      key TEXT NOT NULL,
      value TEXT,
      unit TEXT,
      lat REAL,
      lng REAL,
      synced INTEGER DEFAULT 0,
      created_at TEXT DEFAULT (datetime('now'))
    )
  `);

  await database.execAsync(`
    CREATE TABLE IF NOT EXISTS usages (
      id TEXT PRIMARY KEY,
      visit_id TEXT NOT NULL,
      material_id TEXT NOT NULL,
      quantity REAL NOT NULL,
      notes TEXT,
      synced INTEGER DEFAULT 0,
      created_at TEXT DEFAULT (datetime('now'))
    )
  `);

  await database.execAsync(`
    CREATE TABLE IF NOT EXISTS audit_entries (
      id TEXT PRIMARY KEY,
      visit_id TEXT NOT NULL,
      timestamp TEXT NOT NULL,
      user_name TEXT,
      action TEXT,
      field_name TEXT,
      old_value TEXT,
      new_value TEXT,
      created_at TEXT DEFAULT (datetime('now'))
    )
  `);

  await database.execAsync(`
    CREATE TABLE IF NOT EXISTS outbox (
      id TEXT PRIMARY KEY,
      entity_type TEXT NOT NULL,
      entity_id TEXT NOT NULL,
      operation TEXT NOT NULL,
      payload TEXT NOT NULL,
      status TEXT DEFAULT 'pending',
      error TEXT,
      retry_count INTEGER DEFAULT 0,
      created_at TEXT DEFAULT (datetime('now'))
    )
  `);

  await database.execAsync(`
    CREATE TABLE IF NOT EXISTS config (
      key TEXT PRIMARY KEY,
      value TEXT NOT NULL
    )
  `);

  await database.execAsync(`
    INSERT OR IGNORE INTO config (key, value) VALUES ('gps_accuracy_threshold', '50')
  `);

  await database.execAsync(`
    CREATE TABLE IF NOT EXISTS sync_meta (
      key TEXT PRIMARY KEY,
      value TEXT NOT NULL
    )
  `);

  // Indexes
  await database.execAsync('CREATE INDEX IF NOT EXISTS idx_visits_status ON visits(status)');
  await database.execAsync('CREATE INDEX IF NOT EXISTS idx_visits_scheduled ON visits(scheduled_at)');
  await database.execAsync('CREATE INDEX IF NOT EXISTS idx_visits_partner ON visits(partner_id)');
  await database.execAsync('CREATE INDEX IF NOT EXISTS idx_checklist_items_visit ON checklist_items(visit_id)');
  await database.execAsync('CREATE INDEX IF NOT EXISTS idx_checklist_items_stage ON checklist_items(stage_id)');
  await database.execAsync('CREATE INDEX IF NOT EXISTS idx_checkpoints_visit ON checkpoints(visit_id)');
  await database.execAsync('CREATE INDEX IF NOT EXISTS idx_checkpoints_synced ON checkpoints(synced)');
  await database.execAsync('CREATE INDEX IF NOT EXISTS idx_photos_visit ON photos(visit_id)');
  await database.execAsync('CREATE INDEX IF NOT EXISTS idx_photos_uploaded ON photos(uploaded)');
  await database.execAsync('CREATE INDEX IF NOT EXISTS idx_outbox_status ON outbox(status)');
  await database.execAsync('CREATE INDEX IF NOT EXISTS idx_outbox_entity ON outbox(entity_type, entity_id)');
  await database.execAsync('CREATE INDEX IF NOT EXISTS idx_audit_visit ON audit_entries(visit_id)');
}

export async function closeDatabase(): Promise<void> {
  if (db) {
    await db.closeAsync();
    db = null;
  }
}
