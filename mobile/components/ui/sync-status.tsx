import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { useSyncStore } from '../../stores/sync-store';
import { Badge } from './badge';
import { Button } from './button';

interface SyncStatusProps {
  showRetry?: boolean;
}

export const SyncStatus = React.memo(function SyncStatus({ showRetry = true }: SyncStatusProps) {
  const { lastSync, pendingItems, isSyncing, error, retry } = useSyncStore();

  const formatLastSync = (time: string | null) => {
    if (!time) return 'Nunca';
    const date = new Date(time);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMin = Math.floor(diffMs / 60000);

    if (diffMin < 1) return 'Ahora mismo';
    if (diffMin < 60) return `Hace ${diffMin} min`;
    const diffHours = Math.floor(diffMin / 60);
    if (diffHours < 24) return `Hace ${diffHours}h`;
    return date.toLocaleDateString();
  };

  return (
    <View style={styles.container}>
      <View style={styles.row}>
        <Text style={styles.label}>Última sync:</Text>
        <Text style={styles.value}>{formatLastSync(lastSync)}</Text>
      </View>

      <View style={styles.row}>
        <Text style={styles.label}>Pendientes:</Text>
        <Badge
          label={`${pendingItems}`}
          color={pendingItems > 0 ? '#f59e0b' : '#16a34a'}
        />
      </View>

      {isSyncing && (
        <View style={styles.row}>
          <Text style={styles.label}>Estado:</Text>
          <Badge label="Sincronizando..." color="#2563eb" />
        </View>
      )}

      {error && (
        <View style={styles.errorContainer}>
          <Text style={styles.error}>{error}</Text>
          {showRetry && (
            <Button
              title="Reintentar"
              onPress={retry}
              variant="secondary"
              style={styles.retryButton}
            />
          )}
        </View>
      )}
    </View>
  );
});

const styles = StyleSheet.create({
  container: {
    padding: 16,
  },
  row: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 8,
    borderBottomWidth: 1,
    borderBottomColor: '#eee',
  },
  label: {
    fontFamily: 'monospace',
    fontSize: 12,
    textTransform: 'uppercase',
    color: '#666',
  },
  value: {
    fontFamily: 'monospace',
    fontSize: 12,
    fontWeight: 'bold',
  },
  errorContainer: {
    marginTop: 12,
    padding: 12,
    backgroundColor: '#fef2f2',
    borderWidth: 1,
    borderColor: '#fecaca',
  },
  error: {
    fontFamily: 'monospace',
    fontSize: 11,
    color: '#dc2626',
    marginBottom: 8,
  },
  retryButton: {
    alignSelf: 'flex-start',
  },
});
