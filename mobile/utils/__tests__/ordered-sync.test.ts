import { orderOutboxItems, validateOrder } from '../ordered-sync';

type EntityType = 'checkpoint' | 'photo' | 'usage' | 'measurement' | 'report' | 'status';

interface OutboxItem {
  id: string;
  entity_type: EntityType;
  entity_id: string;
  operation: string;
  payload: string;
  created_at: string;
}

describe('ordered-sync utils', () => {
  it('should sort outbox items by entity type priority', () => {
    const items: OutboxItem[] = [
      { id: '1', entity_type: 'status', entity_id: 'e1', operation: 'create', payload: '{}', created_at: '2026-01-01T00:00:00Z' },
      { id: '2', entity_type: 'checkpoint', entity_id: 'e2', operation: 'create', payload: '{}', created_at: '2026-01-01T00:00:00Z' },
      { id: '3', entity_type: 'photo', entity_id: 'e3', operation: 'create', payload: '{}', created_at: '2026-01-01T00:00:00Z' },
    ];

    const sorted = orderOutboxItems(items);

    expect(sorted[0].entity_type).toBe('checkpoint');
    expect(sorted[1].entity_type).toBe('photo');
    expect(sorted[2].entity_type).toBe('status');
  });

  it('should sort by created_at within same entity type', () => {
    const items: OutboxItem[] = [
      { id: '3', entity_type: 'checkpoint', entity_id: 'e3', operation: 'create', payload: '{}', created_at: '2026-01-03T00:00:00Z' },
      { id: '1', entity_type: 'checkpoint', entity_id: 'e1', operation: 'create', payload: '{}', created_at: '2026-01-01T00:00:00Z' },
      { id: '2', entity_type: 'checkpoint', entity_id: 'e2', operation: 'create', payload: '{}', created_at: '2026-01-02T00:00:00Z' },
    ];

    const sorted = orderOutboxItems(items);

    expect(sorted[0].id).toBe('1');
    expect(sorted[1].id).toBe('2');
    expect(sorted[2].id).toBe('3');
  });

  it('should validate correct order', () => {
    const items: OutboxItem[] = [
      { id: '1', entity_type: 'checkpoint', entity_id: 'e1', operation: 'create', payload: '{}', created_at: '2026-01-01T00:00:00Z' },
      { id: '2', entity_type: 'photo', entity_id: 'e2', operation: 'create', payload: '{}', created_at: '2026-01-01T00:00:00Z' },
      { id: '3', entity_type: 'status', entity_id: 'e3', operation: 'create', payload: '{}', created_at: '2026-01-01T00:00:00Z' },
    ];

    expect(validateOrder(items)).toBe(true);
  });

  it('should detect incorrect order', () => {
    const items: OutboxItem[] = [
      { id: '1', entity_type: 'status', entity_id: 'e1', operation: 'create', payload: '{}', created_at: '2026-01-01T00:00:00Z' },
      { id: '2', entity_type: 'checkpoint', entity_id: 'e2', operation: 'create', payload: '{}', created_at: '2026-01-01T00:00:00Z' },
    ];

    expect(validateOrder(items)).toBe(false);
  });
});
