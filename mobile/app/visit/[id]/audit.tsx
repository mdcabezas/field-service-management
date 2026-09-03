import React, { useEffect, useState } from 'react';
import { View, Text, FlatList, StyleSheet } from 'react-native';
import { useLocalSearchParams } from 'expo-router';
import { apiRequest } from '../../../lib/auth';
import { Header } from '../../../components/ui/header';
import { Card } from '../../../components/ui/card';
import { Badge } from '../../../components/ui/badge';
import { Button } from '../../../components/ui/button';

interface AuditEntry {
  id: string;
  timestamp: string;
  user_name: string;
  action: string;
  field_name: string;
  old_value: string;
  new_value: string;
}

export default function AuditScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const [entries, setEntries] = useState<AuditEntry[]>([]);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    loadAudit();
  }, [id, page]);

  const loadAudit = async () => {
    setIsLoading(true);
    try {
      const response = await apiRequest<{
        entries: AuditEntry[];
        total: number;
        page: number;
        limit: number;
      }>(`/api/visits/${id}/audit?page=${page}&limit=20`);

      setEntries(response.entries);
      setTotal(response.total);
    } catch (error) {
      console.error('Failed to load audit:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const getActionColor = (action: string) => {
    switch (action) {
      case 'create':
        return '#16a34a';
      case 'update':
        return '#2563eb';
      case 'delete':
        return '#dc2626';
      default:
        return '#666';
    }
  };

  const renderEntry = ({ item }: { item: AuditEntry }) => (
    <Card style={styles.entryCard}>
      <View style={styles.entryHeader}>
        <Badge label={item.action.toUpperCase()} color={getActionColor(item.action)} />
        <Text style={styles.entryTime}>
          {new Date(item.timestamp).toLocaleString()}
        </Text>
      </View>
      <Text style={styles.entryUser}>{item.user_name}</Text>
      {item.field_name && (
        <View style={styles.fieldRow}>
          <Text style={styles.fieldName}>{item.field_name}:</Text>
          {item.old_value && (
            <Text style={styles.oldValue}>{item.old_value}</Text>
          )}
          {item.old_value && item.new_value && (
            <Text style={styles.arrow}> → </Text>
          )}
          {item.new_value && (
            <Text style={styles.newValue}>{item.new_value}</Text>
          )}
        </View>
      )}
    </Card>
  );

  return (
    <View style={styles.container}>
      <Header
        title="Historial"
        right={<Badge label={`${total}`} color="#000" />}
      />

      <FlatList
        data={entries}
        renderItem={renderEntry}
        keyExtractor={(item) => item.id}
        contentContainerStyle={styles.list}
        onEndReached={() => {
          if (entries.length < total) {
            setPage(page + 1);
          }
        }}
        onEndReachedThreshold={0.5}
        ListFooterComponent={
          isLoading ? (
            <Text style={styles.loading}>Cargando...</Text>
          ) : null
        }
        windowSize={5}
        maxToRenderPerBatch={10}
        updateCellsBatchingPeriod={50}
        initialNumToRender={10}
        removeClippedSubviews={true}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  list: {
    padding: 16,
  },
  entryCard: {
    marginBottom: 12,
  },
  entryHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  entryTime: {
    fontFamily: 'monospace',
    fontSize: 10,
    color: '#666',
  },
  entryUser: {
    fontFamily: 'monospace',
    fontSize: 12,
    marginBottom: 4,
  },
  fieldRow: {
    flexDirection: 'row',
    alignItems: 'center',
    flexWrap: 'wrap',
  },
  fieldName: {
    fontFamily: 'monospace',
    fontSize: 11,
    fontWeight: 'bold',
    textTransform: 'uppercase',
  },
  oldValue: {
    fontFamily: 'monospace',
    fontSize: 11,
    color: '#dc2626',
    textDecorationLine: 'line-through',
  },
  arrow: {
    fontFamily: 'monospace',
    fontSize: 11,
    color: '#666',
  },
  newValue: {
    fontFamily: 'monospace',
    fontSize: 11,
    color: '#16a34a',
  },
  loading: {
    fontFamily: 'monospace',
    fontSize: 12,
    color: '#666',
    textAlign: 'center',
    padding: 16,
  },
});
