import * as SecureStore from 'expo-secure-store';

const API_URL = process.env.EXPO_PUBLIC_API_URL || 'http://192.168.1.49:8081';
const API_TIMEOUT = 30000;

const pendingRequests = new Map<string, Promise<any>>();

function getRequestKey(path: string, options: RequestInit = {}): string {
  const method = options.method || 'GET';
  const body = options.body ? JSON.stringify(options.body) : '';
  return `${method}:${path}:${body}`;
}

interface AuthTokens {
  access_token: string;
  refresh_token: string;
}

interface LoginResponse {
  access_token: string;
  refresh_token: string;
}

export async function login(email: string, password: string): Promise<LoginResponse> {
  const response = await fetch(`${API_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Login failed');
  }

  const data = await response.json();
  await saveTokens({ access_token: data.access_token, refresh_token: data.refresh_token });
  return data;
}

export async function refreshToken(): Promise<string> {
  const tokens = await getTokens();
  if (!tokens?.refresh_token) {
    throw new Error('No refresh token');
  }

  const response = await fetch(`${API_URL}/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: tokens.refresh_token }),
  });

  if (!response.ok) {
    throw new Error('Token refresh failed');
  }

  const data = await response.json();
  await saveTokens({ access_token: data.access_token, refresh_token: data.refresh_token });
  return data.access_token;
}

export async function logout(): Promise<void> {
  const tokens = await getTokens();
  if (tokens?.access_token) {
    try {
      await fetch(`${API_URL}/auth/logout`, {
        method: 'POST',
        headers: { Authorization: `Bearer ${tokens.access_token}` },
      });
    } catch {
      // Ignore logout errors
    }
  }
  await clearTokens();
}

export async function getAccessToken(): Promise<string | null> {
  const tokens = await getTokens();
  return tokens?.access_token || null;
}

export async function getMe(): Promise<{ id: string; email: string; name: string; role: string }> {
  const token = await getAccessToken();
  if (!token) throw new Error('Not authenticated');

  const response = await fetch(`${API_URL}/auth/me`, {
    headers: { Authorization: `Bearer ${token}` },
  });

  if (!response.ok) throw new Error('Failed to get user info');
  return response.json();
}

async function getTokens(): Promise<AuthTokens | null> {
  const accessToken = await SecureStore.getItemAsync('access_token');
  const refreshToken = await SecureStore.getItemAsync('refresh_token');
  if (accessToken && refreshToken) {
    return { access_token: accessToken, refresh_token: refreshToken };
  }
  return null;
}

async function saveTokens(tokens: AuthTokens): Promise<void> {
  await SecureStore.setItemAsync('access_token', tokens.access_token);
  await SecureStore.setItemAsync('refresh_token', tokens.refresh_token);
}

async function clearTokens(): Promise<void> {
  await SecureStore.deleteItemAsync('access_token');
  await SecureStore.deleteItemAsync('refresh_token');
}

export async function apiRequest<T>(path: string, options: RequestInit = {}): Promise<T> {
  const key = getRequestKey(path, options);
  
  // Deduplicate GET requests
  if ((options.method || 'GET') === 'GET') {
    const existing = pendingRequests.get(key);
    if (existing) {
      return existing as Promise<T>;
    }
  }

  const token = await getAccessToken();
  if (!token) {
    throw new Error('Not authenticated');
  }

  const promise = (async () => {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), API_TIMEOUT);
    try {
      const isFormData = options.body instanceof FormData;
      const response = await fetch(`${API_URL}${path}`, {
        ...options,
        signal: controller.signal,
        headers: {
          ...(!isFormData ? { 'Content-Type': 'application/json' } : {}),
          Authorization: `Bearer ${token}`,
          ...options.headers,
        },
      });
      return response;
    } finally {
      clearTimeout(timeoutId);
    }
  })();

  if ((options.method || 'GET') === 'GET') {
    pendingRequests.set(key, promise);
  }

  try {
    const response = await promise;

    if (response.status === 401) {
      try {
        const newToken = await refreshToken();
        const retryController = new AbortController();
        const retryTimeoutId = setTimeout(() => retryController.abort(), API_TIMEOUT);
        let retryResponse: Response;
        try {
          const isFormData = options.body instanceof FormData;
          retryResponse = await fetch(`${API_URL}${path}`, {
            ...options,
            signal: retryController.signal,
            headers: {
              ...(!isFormData ? { 'Content-Type': 'application/json' } : {}),
              Authorization: `Bearer ${newToken}`,
              ...options.headers,
            },
          });
        } finally {
          clearTimeout(retryTimeoutId);
        }
        if (!retryResponse.ok) {
          throw new Error(`API error: ${retryResponse.status}`);
        }
        return retryResponse.json();
      } catch {
        await logout();
        throw new Error('Session expired');
      }
    }

    if (!response.ok) {
      throw new Error(`API error: ${response.status}`);
    }

    return response.json();
  } finally {
    pendingRequests.delete(key);
  }
}
