import { create } from 'zustand';
import { getSyncState, fullSync, deltaSync, retryFailedItems } from '../lib/sync-engine';

interface SyncState {
  lastSync: string | null;
  isOnline: boolean;
  isSyncing: boolean;
  pendingItems: number;
  error: string | null;
  sync: () => Promise<void>;
  retry: () => Promise<void>;
  refreshState: () => Promise<void>;
}

export const useSyncStore = create<SyncState>((set, get) => ({
  lastSync: null,
  isOnline: true,
  isSyncing: false,
  pendingItems: 0,
  error: null,

  sync: async () => {
    if (get().isSyncing) {
      console.log('[SyncStore] sync: already syncing, skipping');
      return;
    }
    console.log('[SyncStore] sync: START, isSyncing = true');
    set({ isSyncing: true, error: null });
    try {
      const state = await getSyncState();
      console.log('[SyncStore] sync: lastSync =', state.lastSync, ', pendingItems =', state.pendingItems);
      if (state.lastSync) {
        console.log('[SyncStore] sync: calling deltaSync');
        await deltaSync();
      } else {
        console.log('[SyncStore] sync: calling fullSync');
        await fullSync();
      }
      const newState = await getSyncState();
      console.log('[SyncStore] sync: DONE, new lastSync =', newState.lastSync, ', pendingItems =', newState.pendingItems);
      set({
        lastSync: newState.lastSync,
        pendingItems: newState.pendingItems,
        isSyncing: false,
      });
    } catch (error) {
      console.error('[SyncStore] sync: FAILED', error);
      set({ error: String(error), isSyncing: false });
    }
  },

  retry: async () => {
    set({ isSyncing: true, error: null });
    try {
      await retryFailedItems();
      const state = await getSyncState();
      set({ pendingItems: state.pendingItems, isSyncing: false });
    } catch (error) {
      set({ error: String(error), isSyncing: false });
    }
  },

  refreshState: async () => {
    const state = await getSyncState();
    set({
      lastSync: state.lastSync,
      pendingItems: state.pendingItems,
    });
  },
}));
