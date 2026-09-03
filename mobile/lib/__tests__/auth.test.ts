import { login, logout, getAccessToken } from '../auth';

jest.mock('expo-secure-store', () => ({
  getItemAsync: jest.fn().mockResolvedValue(null),
  setItemAsync: jest.fn().mockResolvedValue(undefined),
  deleteItemAsync: jest.fn().mockResolvedValue(undefined),
}));

const mockFetch = jest.fn();
global.fetch = mockFetch;

describe('auth', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should login successfully', async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: jest.fn().mockResolvedValue({
        access_token: 'test-token',
        refresh_token: 'test-refresh',
      }),
    });

    const result = await login('admin@localis.cl', 'password');
    expect(result.access_token).toBe('test-token');
  });

  it('should throw on login failure', async () => {
    mockFetch.mockResolvedValue({
      ok: false,
      json: jest.fn().mockResolvedValue({ error: 'Invalid credentials' }),
    });

    await expect(login('admin@localis.cl', 'wrong')).rejects.toThrow('Invalid credentials');
  });

  it('should logout', async () => {
    mockFetch.mockResolvedValue({ ok: true });
    await expect(logout()).resolves.not.toThrow();
  });

  it('should get access token', async () => {
    const token = await getAccessToken();
    expect(token).toBeNull();
  });

  it('should call fetch on login with email', async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: jest.fn().mockResolvedValue({
        access_token: 'token',
        refresh_token: 'refresh',
      }),
    });

    await login('admin@localis.cl', 'pass');
    expect(mockFetch).toHaveBeenCalled();
    const callBody = JSON.parse(mockFetch.mock.calls[0][1].body);
    expect(callBody.email).toBe('admin@localis.cl');
  });
});
