import { getDatabase } from './database';
import * as Crypto from 'expo-crypto';

export async function addToOutbox(entityType: string, entityId: string, operation: string, payload: any): Promise<void> {
  const db = getDatabase();
  await db.runAsync(
    'INSERT INTO outbox (id, entity_type, entity_id, operation, payload) VALUES (?, ?, ?, ?, ?)',
    Crypto.randomUUID(), entityType, entityId, operation, JSON.stringify(payload)
  );
}
