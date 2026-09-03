import React, { useCallback } from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { useFocusEffect } from 'expo-router';
import { useSyncStore } from '../../stores/sync-store';
import { Header } from '../../components/ui/header';
import { Card } from '../../components/ui/card';
import { Badge } from '../../components/ui/badge';
import { Button } from '../../components/ui/button';

export default function SyncScreen() {
  const { lastSync, isSyncing, pendingItems, error, sync, retry, refreshState } = useSyncStore();

  useFocusEffect(
    useCallback(() => {
      refreshState();
    }, [])
  );

  return (
    <View style={styles.container}>
      <Header title="Sincronización" />

      <View style={styles.content}>
        <Card title="Estado">
          <View style={styles.row}>
            <Text style={styles.label}>Estado:</Text>
            <Badge
              label={isSyncing ? 'Sincronizando' : 'Inactivo'}
              color={isSyncing ? '#f59e0b' : '#059669'}
            />
          </View>
          <View style={styles.row}>
            <Text style={styles.label}>Última sincronización:</Text>
            <Text style={styles.value}>
              {lastSync ? new Date(lastSync).toLocaleString() : 'Nunca'}
            </Text>
          </View>
          <View style={styles.row}>
            <Text style={styles.label}>Elementos pendientes:</Text>
            <Badge label={`${pendingItems}`} color={pendingItems > 0 ? '#f59e0b' : '#059669'} />
          </View>
          {error && (
            <View style={styles.row}>
              <Text style={styles.label}>Error:</Text>
              <Text style={styles.error}>{error}</Text>
            </View>
          )}
        </Card>

        <Card title="Acciones" style={styles.actions}>
          <Button
            title={isSyncing ? 'Sincronizando...' : 'Sincronizar ahora'}
            onPress={() => {
              console.log('[SyncScreen] Sync button pressed');
              sync();
            }}
            disabled={isSyncing}
            style={styles.button}
          />
          {pendingItems > 0 && (
            <Button
              title="Reintentar pendientes"
              onPress={() => {
                console.log('[SyncScreen] Retry button pressed');
                retry();
              }}
              variant="secondary"
              style={styles.button}
            />
          )}
        </Card>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  content: {
    flex: 1,
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
  },
  value: {
    fontFamily: 'monospace',
    fontSize: 12,
    color: '#666',
  },
  error: {
    fontFamily: 'monospace',
    fontSize: 12,
    color: '#dc2626',
  },
  actions: {
    marginTop: 16,
  },
  button: {
    marginTop: 8,
  },
});
