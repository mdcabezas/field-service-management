import { useAuthStore } from '../../stores/auth-store';

describe('auth store', () => {
  beforeEach(() => {
    useAuthStore.setState({
      isAuthenticated: false,
      user: null,
      isLoading: false,
      error: null,
    });
  });

  it('should have initial state', () => {
    const state = useAuthStore.getState();

    expect(state.isAuthenticated).toBe(false);
    expect(state.user).toBeNull();
    expect(state.isLoading).toBe(false);
    expect(state.error).toBeNull();
  });

  it('should set loading state', () => {
    useAuthStore.setState({ isLoading: true });
    const state = useAuthStore.getState();

    expect(state.isLoading).toBe(true);
  });

  it('should set error', () => {
    useAuthStore.setState({ error: 'Test error' });
    const state = useAuthStore.getState();

    expect(state.error).toBe('Test error');
  });

  it('should clear error', () => {
    useAuthStore.setState({ error: 'Test error' });
    useAuthStore.setState({ error: null });
    const state = useAuthStore.getState();

    expect(state.error).toBeNull();
  });
});
