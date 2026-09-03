import { useSyncStore } from '../../stores/sync-store';

describe('sync store', () => {
  beforeEach(() => {
    useSyncStore.setState({
      lastSync: null,
      isOnline: true,
      isSyncing: false,
      pendingItems: 0,
      error: null,
    });
  });

  it('should have initial state', () => {
    const state = useSyncStore.getState();

    expect(state.lastSync).toBeNull();
    expect(state.isOnline).toBe(true);
    expect(state.isSyncing).toBe(false);
    expect(state.pendingItems).toBe(0);
    expect(state.error).toBeNull();
  });

  it('should update pendingItems', () => {
    useSyncStore.setState({ pendingItems: 5 });
    const state = useSyncStore.getState();

    expect(state.pendingItems).toBe(5);
  });

  it('should set isSyncing', () => {
    useSyncStore.setState({ isSyncing: true });
    const state = useSyncStore.getState();

    expect(state.isSyncing).toBe(true);
  });

  it('should update lastSync', () => {
    const now = new Date().toISOString();
    useSyncStore.setState({ lastSync: now });
    const state = useSyncStore.getState();

    expect(state.lastSync).toBe(now);
  });

  it('should update isOnline', () => {
    useSyncStore.setState({ isOnline: false });
    const state = useSyncStore.getState();

    expect(state.isOnline).toBe(false);
  });
});
