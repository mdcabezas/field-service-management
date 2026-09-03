type EntityType = 'checkpoint' | 'photo' | 'usage' | 'measurement' | 'report' | 'status';

interface OutboxItem {
  id: string;
  entity_type: EntityType;
  entity_id: string;
  operation: string;
  payload: string;
  created_at: string;
}

const DEPENDENCY_ORDER: EntityType[] = [
  'checkpoint',
  'photo',
  'usage',
  'measurement',
  'report',
  'status',
];

export function orderOutboxItems(items: OutboxItem[]): OutboxItem[] {
  return [...items].sort((a, b) => {
    const orderA = DEPENDENCY_ORDER.indexOf(a.entity_type);
    const orderB = DEPENDENCY_ORDER.indexOf(b.entity_type);
    if (orderA !== orderB) {
      return orderA - orderB;
    }
    return new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
  });
}

export function validateOrder(items: OutboxItem[]): boolean {
  let lastOrder = -1;
  for (const item of items) {
    const currentOrder = DEPENDENCY_ORDER.indexOf(item.entity_type);
    if (currentOrder < lastOrder) {
      return false;
    }
    lastOrder = currentOrder;
  }
  return true;
}
