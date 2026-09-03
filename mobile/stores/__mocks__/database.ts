export function getDatabase() {
  return {
    prepareSync: jest.fn().mockReturnValue({
      executeSync: jest.fn().mockReturnValue({
        getFirstSync: jest.fn().mockReturnValue(null),
        getAllSync: jest.fn().mockReturnValue([]),
      }),
      finalizeSync: jest.fn(),
    }),
    runAsync: jest.fn(),
  };
}

export function initDatabase() {
  return Promise.resolve(getDatabase());
}
