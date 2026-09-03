import { create } from 'zustand';
import { getDatabase } from '../lib/database';
import { addToOutbox } from '../lib/outbox';

interface Visit {
  id: string;
  property_id: string;
  property_name: string;
  property_address: string;
  property_lat: number;
  property_lng: number;
  partner_id: string;
  partner_name: string;
  partner_order_id: string;
  partner_supervisor: string;
  type: string;
  status: string;
  priority: string;
  scheduled_at: string;
  checklist_template_id: string;
  stages_enabled: boolean;
  selected_stages: string[];
  sla_response_hours: number;
  sla_resolution_hours: number;
}

interface VisitState {
  visits: Visit[];
  currentVisit: Visit | null;
  filter: string;
  isLoading: boolean;
  setFilter: (filter: string) => void;
  loadVisits: () => Promise<void>;
  selectVisit: (visit: Visit | null) => void;
  loadVisitDetail: (visitId: string) => Promise<void>;
  updateVisitStatus: (visitId: string, newStatus: string) => Promise<void>;
}

export const useVisitStore = create<VisitState>((set, get) => ({
  visits: [],
  currentVisit: null,
  filter: 'all',
  isLoading: false,

  setFilter: (filter: string) => {
    set({ filter });
    get().loadVisits();
  },

  loadVisits: async () => {
    set({ isLoading: true });
    const db = getDatabase();
    const { filter } = get();

    let query = 'SELECT * FROM visits';
    const params: any[] = [];

    if (filter !== 'all') {
      query += ' WHERE status = ?';
      params.push(filter);
    }

    query += ' ORDER BY scheduled_at DESC LIMIT 50';

    const statement = db.prepareSync(query);
    const result = statement.executeSync(...params);
    const rows = result.getAllSync() as Record<string, unknown>[];
    statement.finalizeSync();

    const visits = rows.map((row) => ({
      id: String(row.id ?? ''),
      property_id: String(row.property_id ?? ''),
      property_name: String(row.property_name ?? ''),
      property_address: String(row.property_address ?? ''),
      property_lat: Number(row.property_lat ?? 0),
      property_lng: Number(row.property_lng ?? 0),
      partner_id: String(row.partner_id ?? ''),
      partner_name: String(row.partner_name ?? ''),
      partner_order_id: String(row.partner_order_id ?? ''),
      partner_supervisor: String(row.partner_supervisor ?? ''),
      type: String(row.type ?? ''),
      status: String(row.status ?? ''),
      priority: String(row.priority ?? ''),
      scheduled_at: String(row.scheduled_at ?? ''),
      checklist_template_id: String(row.checklist_template_id ?? ''),
      stages_enabled: row.stages_enabled === 1,
      selected_stages: JSON.parse(String(row.selected_stages ?? '[]')),
      sla_response_hours: Number(row.sla_response_hours ?? 0),
      sla_resolution_hours: Number(row.sla_resolution_hours ?? 0),
    })) as Visit[];

    set({ visits, isLoading: false });
  },

  selectVisit: (visit: Visit | null) => {
    set({ currentVisit: visit });
  },

  loadVisitDetail: async (visitId: string) => {
    const db = getDatabase();
    const statement = db.prepareSync('SELECT * FROM visits WHERE id = ?');
    const result = statement.executeSync(visitId);
    const row = result.getFirstSync() as Record<string, unknown> | null;
    statement.finalizeSync();

    if (row) {
      const visit: Visit = {
        id: String(row.id ?? ''),
        property_id: String(row.property_id ?? ''),
        property_name: String(row.property_name ?? ''),
        property_address: String(row.property_address ?? ''),
        property_lat: Number(row.property_lat ?? 0),
        property_lng: Number(row.property_lng ?? 0),
        partner_id: String(row.partner_id ?? ''),
        partner_name: String(row.partner_name ?? ''),
        partner_order_id: String(row.partner_order_id ?? ''),
        partner_supervisor: String(row.partner_supervisor ?? ''),
        type: String(row.type ?? ''),
        status: String(row.status ?? ''),
        priority: String(row.priority ?? ''),
        scheduled_at: String(row.scheduled_at ?? ''),
        checklist_template_id: String(row.checklist_template_id ?? ''),
        stages_enabled: row.stages_enabled === 1,
        selected_stages: JSON.parse(String(row.selected_stages ?? '[]')),
        sla_response_hours: Number(row.sla_response_hours ?? 0),
        sla_resolution_hours: Number(row.sla_resolution_hours ?? 0),
      };
      set({ currentVisit: visit });
    }
  },

  updateVisitStatus: async (visitId: string, newStatus: string) => {
    const db = getDatabase();
    const now = new Date().toISOString();

    await db.runAsync(
      'UPDATE visits SET status = ?, updated_at = ? WHERE id = ?',
      newStatus, now, visitId
    );

    await addToOutbox('status', visitId, 'update', {
      visit_id: visitId,
      status: newStatus,
      timestamp: now,
    });

    const { currentVisit, visits } = get();
    if (currentVisit?.id === visitId) {
      set({ currentVisit: { ...currentVisit, status: newStatus } });
    }
    set({
      visits: visits.map((v) =>
        v.id === visitId ? { ...v, status: newStatus } : v
      ),
    });
  },
}));
