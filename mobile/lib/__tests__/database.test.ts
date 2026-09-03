import { getDatabase, closeDatabase, initDatabase } from '../database';

jest.mock('expo-sqlite', () => ({
  openDatabaseAsync: jest.fn().mockResolvedValue({
    execAsync: jest.fn(),
    prepareSync: jest.fn().mockReturnValue({
      executeSync: jest.fn().mockReturnValue({
        getFirstSync: jest.fn().mockReturnValue(null),
        getAllSync: jest.fn().mockReturnValue([]),
      }),
      finalizeSync: jest.fn(),
    }),
    runAsync: jest.fn(),
    closeAsync: jest.fn(),
  }),
}));

describe('database', () => {
  beforeEach(async () => {
    await closeDatabase();
  });

  it('should initialize database', async () => {
    const db = await initDatabase();
    expect(db).toBeDefined();
  });

  it('should return same instance on multiple calls', async () => {
    const db1 = await initDatabase();
    const db2 = getDatabase();
    expect(db1).toBe(db2);
  });

  it('should close database', async () => {
    await initDatabase();
    await closeDatabase();
    expect(() => getDatabase()).toThrow('Database not initialized');
  });

  it('should throw if not initialized', () => {
    expect(() => getDatabase()).toThrow('Database not initialized');
  });
});
