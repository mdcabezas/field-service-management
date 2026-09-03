jest.mock('react-native', () => ({
  Platform: { OS: 'android' },
}));

jest.mock('expo-secure-store', () => ({
  getItemAsync: jest.fn().mockResolvedValue(null),
  setItemAsync: jest.fn().mockResolvedValue(undefined),
  deleteItemAsync: jest.fn().mockResolvedValue(undefined),
}));

jest.mock('expo-image-manipulator', () => ({
  manipulateAsync: jest.fn().mockResolvedValue({ uri: 'file:///mock/image.jpg', width: 640, height: 480 }),
  SaveFormat: { JPEG: 'jpeg' },
}));

jest.mock('expo-file-system', () => ({
  documentDirectory: 'file:///mock/documents/',
  cacheDirectory: 'file:///mock/cache/',
  copyAsync: jest.fn().mockResolvedValue(undefined),
  getInfoAsync: jest.fn().mockResolvedValue({ exists: true }),
  makeDirectoryAsync: jest.fn().mockResolvedValue(undefined),
  readDirectoryAsync: jest.fn().mockResolvedValue([]),
  deleteAsync: jest.fn().mockResolvedValue(undefined),
  Paths: {
    document: { uri: 'file:///mock/documents' },
  },
  Directory: jest.fn().mockImplementation(() => ({
    uri: 'file:///mock/photos/',
    exists: true,
    create: jest.fn(),
  })),
  File: jest.fn().mockImplementation((uri: string) => ({
    uri,
    name: uri.split('/').pop() || 'file.jpg',
    type: 'image/jpeg',
    exists: true,
    copy: jest.fn().mockResolvedValue(undefined),
    delete: jest.fn(),
    bytes: jest.fn().mockResolvedValue(new Uint8Array([0xff, 0xd8])),
    arrayBuffer: jest.fn().mockResolvedValue(new ArrayBuffer(2)),
  })),
}));

jest.mock('../../lib/database', () => ({
  getDatabase: jest.fn().mockReturnValue({
    prepareSync: jest.fn().mockReturnValue({
      executeSync: jest.fn().mockReturnValue({
        getFirstSync: jest.fn().mockReturnValue(null),
        getAllSync: jest.fn().mockReturnValue([]),
      }),
      finalizeSync: jest.fn(),
    }),
    runAsync: jest.fn().mockResolvedValue({}),
  }),
}));

jest.mock('../../lib/auth');

const mockApiRequest = require('../../lib/auth').apiRequest;

describe('Sync Flow Integration', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('Full Sync', () => {
    it('should fetch reference data and visits on full sync', async () => {
      mockApiRequest.mockResolvedValueOnce({
        reference_data: {
          materials: [{ id: '1', name: 'Material 1' }],
          tools: [],
          epp: [],
          templates: [],
          partners: [],
        },
        visits: [{ id: 'visit-1', status: 'assigned' }],
        checklist_items: [],
        photos: [],
        synced_at: '2024-01-01T00:00:00Z',
      });

      const { fullSync } = require('../../lib/sync-engine');
      await fullSync();

      expect(mockApiRequest).toHaveBeenCalled();
    });

    it('should handle empty response gracefully', async () => {
      mockApiRequest.mockResolvedValueOnce({
        reference_data: null,
        visits: [],
        checklist_items: [],
        photos: [],
        synced_at: '2024-01-01T00:00:00Z',
      });

      const { fullSync } = require('../../lib/sync-engine');
      await fullSync();

      expect(mockApiRequest).toHaveBeenCalled();
    });
  });

  describe('Delta Sync', () => {
    it('should fetch changes since last sync', async () => {
      mockApiRequest.mockResolvedValueOnce({
        visits: [{ id: 'visit-2', status: 'in_progress' }],
        checklist_items: [],
        photos: [],
        synced_at: '2024-01-02T00:00:00Z',
      });

      const { deltaSync } = require('../../lib/sync-engine');
      await deltaSync();

      expect(mockApiRequest).toHaveBeenCalled();
    });
  });
});