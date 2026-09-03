import React, { useEffect } from 'react';
import { View, FlatList, Pressable, StyleSheet, ActivityIndicator, Text } from 'react-native';
import { useRouter } from 'expo-router';
import { useVisitStore } from '../../stores/visit-store';
import { useSyncStore } from '../../stores/sync-store';
import { Header } from '../../components/ui/header';
import { Card } from '../../components/ui/card';
import { Badge } from '../../components/ui/badge';
import { Button } from '../../components/ui/button';

const STATUS_COLORS: Record<string, string> = {
  assigned: '#6b7280',
  scheduled: '#3b82f6',
  en_route: '#f59e0b',
  in_progress: '#10b981',
  completed: '#059669',
  cancelled: '#dc2626',
};

export default function VisitListScreen() {
  const router = useRouter();
  const { visits, isLoading, filter, setFilter, loadVisits } = useVisitStore();
  const { sync, isSyncing, lastSync, pendingItems } = useSyncStore();

  useEffect(() => {
    loadVisits();
  }, []);

  const handleRefresh = () => {
    sync().then(() => loadVisits());
  };

  const renderVisit = ({ item }: { item: any }) => (
    <Pressable onPress={() => router.push(`/visit/${item.id}`)}>
      <Card style={styles.visitCard}>
        <View style={styles.visitHeader}>
          <Text style={styles.propertyName}>{item.property_name || 'Sin nombre'}</Text>
          <Badge label={item.status} color={STATUS_COLORS[item.status] || '#6b7280'} />
        </View>
        <Text style={styles.address}>{item.property_address || 'Sin dirección'}</Text>
        {item.partner_name ? (
          <Text style={styles.partner}>Partner: {item.partner_name}</Text>
        ) : null}
        <View style={styles.visitFooter}>
          <Text style={styles.time}>{item.scheduled_at ? new Date(item.scheduled_at).toLocaleTimeString() : '--:--'}</Text>
          <Text style={styles.priority}>{item.priority}</Text>
        </View>
      </Card>
    </Pressable>
  );

  return (
    <View style={styles.container}>
      <Header
        title="Visitas"
        right={
          <View style={styles.headerRight}>
            <Text style={styles.syncInfo}>
              {isSyncing ? 'Syncing...' : lastSync ? 'Synced' : 'Not synced'}
            </Text>
            {pendingItems > 0 && (
              <Badge label={`${pendingItems} pending`} color="#f59e0b" />
            )}
          </View>
        }
      />

      <View style={styles.filters}>
        {['all', 'assigned', 'in_progress', 'completed'].map((f) => (
          <Pressable
            key={f}
            onPress={() => setFilter(f)}
            style={[styles.filterButton, filter === f && styles.filterActive]}
          >
            <Text style={[styles.filterText, filter === f && styles.filterTextActive]}>
              {f === 'all' ? 'Todos' : f}
            </Text>
          </Pressable>
        ))}
      </View>

      <Button
        title={isSyncing ? 'Sincronizando...' : 'Sincronizar'}
        onPress={handleRefresh}
        disabled={isSyncing}
        style={styles.syncButton}
      />

      {isLoading ? (
        <ActivityIndicator size="large" style={styles.loading} />
      ) : (
        <FlatList
          data={visits}
          renderItem={renderVisit}
          keyExtractor={(item) => item.id}
          contentContainerStyle={styles.list}
          windowSize={5}
          maxToRenderPerBatch={10}
          updateCellsBatchingPeriod={50}
          initialNumToRender={10}
          removeClippedSubviews={true}
        />
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  headerRight: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  syncInfo: {
    fontFamily: 'monospace',
    fontSize: 10,
    textTransform: 'uppercase',
  },
  filters: {
    flexDirection: 'row',
    borderBottomWidth: 2,
    borderBottomColor: '#000',
  },
  filterButton: {
    flex: 1,
    paddingVertical: 12,
    alignItems: 'center',
    borderRightWidth: 1,
    borderRightColor: '#000',
  },
  filterActive: {
    backgroundColor: '#000',
  },
  filterText: {
    fontFamily: 'monospace',
    fontSize: 10,
    textTransform: 'uppercase',
  },
  filterTextActive: {
    color: '#fff',
  },
  syncButton: {
    margin: 16,
  },
  loading: {
    flex: 1,
    justifyContent: 'center',
  },
  list: {
    padding: 16,
  },
  visitCard: {
    marginBottom: 12,
  },
  visitHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  propertyName: {
    fontFamily: 'monospace',
    fontSize: 14,
    fontWeight: 'bold',
    textTransform: 'uppercase',
    flex: 1,
  },
  address: {
    fontFamily: 'monospace',
    fontSize: 12,
    color: '#666',
    marginBottom: 4,
  },
  partner: {
    fontFamily: 'monospace',
    fontSize: 10,
    color: '#3b82f6',
    marginBottom: 4,
  },
  visitFooter: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginTop: 8,
  },
  time: {
    fontFamily: 'monospace',
    fontSize: 12,
  },
  priority: {
    fontFamily: 'monospace',
    fontSize: 10,
    textTransform: 'uppercase',
  },
});
