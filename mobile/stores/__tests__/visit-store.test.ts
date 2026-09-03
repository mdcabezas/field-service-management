import { useVisitStore } from '../../stores/visit-store';

describe('visit store', () => {
  beforeEach(() => {
    useVisitStore.setState({
      visits: [],
      currentVisit: null,
      filter: 'all',
      isLoading: false,
    });
  });

  it('should have initial state', () => {
    const state = useVisitStore.getState();

    expect(state.visits).toEqual([]);
    expect(state.currentVisit).toBeNull();
    expect(state.filter).toBe('all');
    expect(state.isLoading).toBe(false);
  });

  it('should set filter', () => {
    useVisitStore.getState().setFilter('in_progress');
    const state = useVisitStore.getState();

    expect(state.filter).toBe('in_progress');
  });

  it('should set currentVisit', () => {
    const mockVisit = {
      id: '123',
      property_id: '456',
      property_name: 'Test Property',
      property_address: '123 Main St',
      property_lat: 19.4326,
      property_lng: -99.1332,
      partner_id: '789',
      partner_name: 'Test Partner',
      partner_order_id: 'ORD-001',
      partner_supervisor: 'Supervisor',
      type: 'maintenance',
      status: 'assigned',
      priority: 'normal',
      scheduled_at: new Date().toISOString(),
      checklist_template_id: 'tpl-1',
      stages_enabled: true,
      selected_stages: ['electrical', 'plumbing'],
      sla_response_hours: 4,
      sla_resolution_hours: 24,
    };

    useVisitStore.getState().selectVisit(mockVisit);
    const state = useVisitStore.getState();

    expect(state.currentVisit).toEqual(mockVisit);
  });

  it('should clear currentVisit', () => {
    useVisitStore.getState().selectVisit(null);
    const state = useVisitStore.getState();

    expect(state.currentVisit).toBeNull();
  });
});
