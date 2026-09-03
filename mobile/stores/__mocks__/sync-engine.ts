export async function getSyncState() {
  return { lastSync: null, pendingItems: 0 };
}

export async function fullSync() {}

export async function deltaSync() {}

export async function retryFailedItems() {}
