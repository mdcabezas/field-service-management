jest.mock('react-native', () => ({
  Platform: { OS: 'android' },
}));

jest.mock('expo-image-manipulator', () => ({
  manipulateAsync: jest.fn().mockResolvedValue({ uri: 'file:///tmp/manipulated.jpg' }),
}));

jest.mock('expo-file-system', () => ({
  Paths: {
    document: { uri: 'file:///documents/' },
  },
  Directory: jest.fn().mockImplementation(() => ({
    uri: 'file:///documents/photos/',
    create: jest.fn(),
    contains: jest.fn().mockResolvedValue(true),
    list: jest.fn().mockResolvedValue([]),
  })),
  File: jest.fn().mockImplementation((path: string) => ({
    uri: path,
    name: path.split('/').pop() || 'file.jpg',
    type: 'image/jpeg',
    exists: jest.fn().mockResolvedValue(true),
    copy: jest.fn(),
    delete: jest.fn(),
    base64: jest.fn().mockResolvedValue('base64data'),
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

jest.mock('../../lib/auth', () => ({
  apiRequest: jest.fn().mockResolvedValue({}),
}));

describe('Photo Flow Integration', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('Cleanup Old Photos', () => {
    it('should call cleanup without errors', async () => {
      const { cleanupOldPhotos } = require('../../lib/photo-upload');
      await expect(cleanupOldPhotos(7, 100)).resolves.not.toThrow();
    });
  });

  describe('Storage Stats', () => {
    it('should calculate storage size', async () => {
      const { getPhotoStorageSize } = require('../../lib/photo-upload');
      const size = await getPhotoStorageSize();

      expect(typeof size).toBe('number');
      expect(size).toBeGreaterThanOrEqual(0);
    });
  });
});