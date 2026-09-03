jest.mock('expo-secure-store', () => ({
  getItemAsync: jest.fn(),
  setItemAsync: jest.fn().mockResolvedValue(undefined),
  deleteItemAsync: jest.fn().mockResolvedValue(undefined),
}));

describe('Auth Flow Integration', () => {
  beforeEach(() => {
    jest.resetModules();
    jest.clearAllMocks();
    jest.mock('expo-secure-store', () => ({
      getItemAsync: jest.fn(),
      setItemAsync: jest.fn().mockResolvedValue(undefined),
      deleteItemAsync: jest.fn().mockResolvedValue(undefined),
    }));
  });

  describe('Login Flow', () => {
    it('should store tokens on successful login', async () => {
      global.fetch = jest.fn().mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          access_token: 'access-123',
          refresh_token: 'refresh-123',
        }),
      });

      const { login } = require('../../lib/auth');
      const result = await login('admin@localis.cl', 'password');

      expect(result.access_token).toBe('access-123');
    });

    it('should throw on invalid credentials', async () => {
      global.fetch = jest.fn().mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Invalid credentials' }),
      });

      const { login } = require('../../lib/auth');
      await expect(login('admin@localis.cl', 'wrong')).rejects.toThrow('Invalid credentials');
    });
  });

  describe('Token Refresh Flow', () => {
    it('should refresh token successfully', async () => {
      const SecureStore = require('expo-secure-store');
      SecureStore.getItemAsync.mockResolvedValueOnce('old-access-token');
      SecureStore.getItemAsync.mockResolvedValueOnce('refresh-token');

      global.fetch = jest.fn().mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          access_token: 'new-access-token',
          refresh_token: 'new-refresh-token',
        }),
      });

      const { refreshToken } = require('../../lib/auth');
      const newToken = await refreshToken();

      expect(newToken).toBe('new-access-token');
    });

    it('should throw when no refresh token', async () => {
      const SecureStore = require('expo-secure-store');
      SecureStore.getItemAsync.mockResolvedValue(null);

      const { refreshToken } = require('../../lib/auth');
      await expect(refreshToken()).rejects.toThrow('No refresh token');
    });
  });

  describe('Logout Flow', () => {
    it('should clear tokens on logout', async () => {
      const SecureStore = require('expo-secure-store');
      SecureStore.getItemAsync.mockResolvedValueOnce('access-token');
      SecureStore.getItemAsync.mockResolvedValueOnce('refresh-token');

      global.fetch = jest.fn().mockResolvedValueOnce({
        ok: true,
      });

      const { logout } = require('../../lib/auth');
      await logout();

      expect(SecureStore.deleteItemAsync).toHaveBeenCalledWith('access_token');
      expect(SecureStore.deleteItemAsync).toHaveBeenCalledWith('refresh_token');
    });

    it('should clear tokens even if server logout fails', async () => {
      const SecureStore = require('expo-secure-store');
      SecureStore.getItemAsync.mockResolvedValueOnce('access-token');
      SecureStore.getItemAsync.mockResolvedValueOnce('refresh-token');

      global.fetch = jest.fn().mockRejectedValueOnce(new Error('Network error'));

      const { logout } = require('../../lib/auth');
      await logout();

      expect(SecureStore.deleteItemAsync).toHaveBeenCalledWith('access_token');
      expect(SecureStore.deleteItemAsync).toHaveBeenCalledWith('refresh_token');
    });
  });

  describe('Get Access Token', () => {
    it('should return token when available', async () => {
      const SecureStore = require('expo-secure-store');
      SecureStore.getItemAsync.mockResolvedValueOnce('my-token');
      SecureStore.getItemAsync.mockResolvedValueOnce('my-refresh-token');

      const { getAccessToken } = require('../../lib/auth');
      const token = await getAccessToken();

      expect(token).toBe('my-token');
    });

    it('should return null when no token', async () => {
      const SecureStore = require('expo-secure-store');
      SecureStore.getItemAsync.mockResolvedValue(null);

      const { getAccessToken } = require('../../lib/auth');
      const token = await getAccessToken();

      expect(token).toBeNull();
    });
  });
});
