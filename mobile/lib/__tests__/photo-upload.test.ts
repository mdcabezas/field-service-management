import {
  captureAndCompressPhoto,
  uploadPhoto,
  getPhotoStorageSize,
  cleanupOldPhotos,
  sharePhoto,
} from '../photo-upload';

jest.mock('react-native', () => ({
  Platform: { OS: 'android' },
}));

jest.mock('expo-image-manipulator', () => ({
  manipulateAsync: jest.fn().mockResolvedValue({
    uri: 'file:///mock/path.jpg',
    width: 800,
    height: 600,
  }),
  SaveFormat: {
    JPEG: 'jpeg',
  },
}));

jest.mock('expo-file-system', () => ({
  Paths: {
    document: { uri: 'file:///mock/documents' },
  },
  Directory: jest.fn().mockImplementation((baseUri, name) => ({
    uri: `${baseUri}/${name}`,
    exists: true,
    create: jest.fn(),
  })),
  File: jest.fn().mockImplementation((uri) => ({
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

jest.mock('../database', () => ({
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

jest.mock('../auth', () => ({
  apiRequest: jest.fn().mockResolvedValue({ id: 'server-photo-id' }),
}));

jest.mock('../sync-engine', () => ({
  addToOutbox: jest.fn().mockResolvedValue(undefined),
}));

describe('photo-upload', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should capture and compress photo', async () => {
    const result = await captureAndCompressPhoto(
      'file:///original/photo.jpg',
      'visit-123'
    );

    expect(result).toHaveProperty('id');
    expect(result).toHaveProperty('localPath');
    expect(result).toHaveProperty('thumbnailPath');
  });

  it('should upload photo successfully', async () => {
    const success = await uploadPhoto('photo-123');

    expect(typeof success).toBe('boolean');
  });

  it('should get photo storage size', async () => {
    const size = await getPhotoStorageSize();

    expect(typeof size).toBe('number');
    expect(size).toBeGreaterThanOrEqual(0);
  });

  it('should cleanup old photos', async () => {
    await expect(
      cleanupOldPhotos(7, 100)
    ).resolves.not.toThrow();
  });

  it('should share photo', async () => {
    await expect(
      sharePhoto('photo-123')
    ).resolves.not.toThrow();
  });

  it('should handle missing photo on upload', async () => {
    const success = await uploadPhoto('nonexistent-photo');

    expect(success).toBe(false);
  });

  it('should handle missing photo on share', async () => {
    await expect(
      sharePhoto('nonexistent-photo')
    ).resolves.not.toThrow();
  });

  it('should handle missing photo path on share', async () => {
    await expect(
      sharePhoto('photo-without-path')
    ).resolves.not.toThrow();
  });
});
