import React, { useEffect, useState } from 'react';
import { View, Text, FlatList, Pressable, StyleSheet } from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { getDatabase } from '../../lib/database';
import { addToOutbox } from '../../lib/outbox';
import { Header } from '../../components/ui/header';
import { Card } from '../../components/ui/card';
import { Badge } from '../../components/ui/badge';

interface ChecklistItem {
  id: string;
  stage_id: string;
  stage_name: string;
  type: string;
  name: string;
  quantity: number;
  unit: string;
  confirmed: boolean;
  required: boolean;
}

interface StageGroup {
  name: string;
  items: ChecklistItem[];
  confirmedCount: number;
  totalCount: number;
}

export default function ChecklistScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const [items, setItems] = useState<ChecklistItem[]>([]);
  const [expandedStages, setExpandedStages] = useState<Set<string>>(new Set());

  useEffect(() => {
    loadChecklist();
  }, [id]);

  const loadChecklist = async () => {
    const db = getDatabase();
    const stmt = db.prepareSync(
      'SELECT * FROM checklist_items WHERE visit_id = ? ORDER BY stage_name, type, name'
    );
    const result = stmt.executeSync(id);
    const rows = result.getAllSync() as Record<string, unknown>[];
    stmt.finalizeSync();
    setItems(
      rows.map((row) => ({
        id: row.id as string,
        stage_id: row.stage_id as string,
        stage_name: row.stage_name as string,
        type: row.type as string,
        name: row.name as string,
        quantity: row.quantity as number,
        unit: row.unit as string,
        confirmed: row.confirmed === 1,
        required: row.required === 1,
      }))
    );
  };

  const toggleItem = async (itemId: string) => {
    const db = getDatabase();
    const item = items.find((i) => i.id === itemId);
    if (!item) return;

    const newConfirmed = !item.confirmed;
    await db.runAsync('UPDATE checklist_items SET confirmed = ? WHERE id = ?', newConfirmed ? 1 : 0, itemId);

    await addToOutbox('checklist', itemId, 'update', {
      visit_id: id,
      item_id: itemId,
      confirmed: newConfirmed,
      timestamp: new Date().toISOString(),
    });

    setItems((prev) =>
      prev.map((i) => (i.id === itemId ? { ...i, confirmed: newConfirmed } : i))
    );
  };

  const toggleStage = (stageName: string) => {
    setExpandedStages((prev) => {
      const next = new Set(prev);
      if (next.has(stageName)) {
        next.delete(stageName);
      } else {
        next.add(stageName);
      }
      return next;
    });
  };

  const stageGroups: StageGroup[] = [];
  const stageMap = new Map<string, ChecklistItem[]>();

  for (const item of items) {
    const stageName = item.stage_name || 'Sin etapa';
    if (!stageMap.has(stageName)) {
      stageMap.set(stageName, []);
    }
    stageMap.get(stageName)!.push(item);
  }

  for (const [name, stageItems] of stageMap) {
    const confirmedCount = stageItems.filter((i) => i.confirmed).length;
    stageGroups.push({
      name,
      items: stageItems,
      confirmedCount,
      totalCount: stageItems.length,
    });
  }

  const renderStage = (group: StageGroup) => {
    const isExpanded = expandedStages.has(group.name);
    const allConfirmed = group.confirmedCount === group.totalCount;

    return (
      <View key={group.name} style={styles.stageContainer}>
        <Pressable onPress={() => toggleStage(group.name)}>
          <Card style={styles.stageHeader}>
            <View style={styles.stageHeaderContent}>
              <Text style={styles.stageName}>{group.name}</Text>
              <View style={styles.stageBadges}>
                <Badge
                  label={`${group.confirmedCount}/${group.totalCount}`}
                  color={allConfirmed ? '#059669' : '#f59e0b'}
                />
                {group.items.some((i) => i.required && !i.confirmed) && (
                  <Badge label="Requerido" color="#dc2626" />
                )}
              </View>
            </View>
            <Text style={styles.expandIcon}>{isExpanded ? '−' : '+'}</Text>
          </Card>
        </Pressable>

        {isExpanded && (
          <View style={styles.itemsList}>
            {group.items.map((item) => (
              <Pressable
                key={item.id}
                onPress={() => toggleItem(item.id)}
                style={[styles.itemRow, item.confirmed && styles.itemConfirmed]}
              >
                <View style={[styles.checkbox, item.confirmed && styles.checkboxChecked]}>
                  {item.confirmed && <Text style={styles.checkmark}>✓</Text>}
                </View>
                <View style={styles.itemInfo}>
                  <Text style={styles.itemName}>{item.name}</Text>
                  <Text style={styles.itemDetails}>
                    {item.type} • {item.quantity} {item.unit}
                  </Text>
                </View>
                {item.required && !item.confirmed && (
                  <Badge label="Req" color="#dc2626" />
                )}
              </Pressable>
            ))}
          </View>
        )}
      </View>
    );
  };

  return (
    <View style={styles.container}>
      <Header
        title="Checklist"
        right={
          <Badge
            label={`${items.filter((i) => i.confirmed).length}/${items.length}`}
            color={
              items.length > 0 && items.every((i) => i.confirmed)
                ? '#059669'
                : '#f59e0b'
            }
          />
        }
      />

      <FlatList
        data={stageGroups}
        renderItem={({ item }) => renderStage(item)}
        keyExtractor={(item) => item.name}
        contentContainerStyle={styles.list}
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
  stageContainer: {
    marginBottom: 16,
  },
  stageHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  stageHeaderContent: {
    flex: 1,
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  stageName: {
    fontFamily: 'monospace',
    fontSize: 14,
    fontWeight: 'bold',
    textTransform: 'uppercase',
  },
  stageBadges: {
    flexDirection: 'row',
    gap: 4,
  },
  expandIcon: {
    fontFamily: 'monospace',
    fontSize: 18,
    fontWeight: 'bold',
    marginLeft: 8,
  },
  itemsList: {
    borderLeftWidth: 2,
    borderLeftColor: '#000',
    marginLeft: 8,
  },
  itemRow: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 12,
    borderBottomWidth: 1,
    borderBottomColor: '#eee',
  },
  itemConfirmed: {
    backgroundColor: '#f0fdf4',
  },
  checkbox: {
    width: 24,
    height: 24,
    borderWidth: 2,
    borderColor: '#000',
    borderRadius: 0,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 12,
  },
  checkboxChecked: {
    backgroundColor: '#059669',
    borderColor: '#059669',
  },
  checkmark: {
    color: '#fff',
    fontFamily: 'monospace',
    fontWeight: 'bold',
  },
  itemInfo: {
    flex: 1,
  },
  itemName: {
    fontFamily: 'monospace',
    fontSize: 12,
    fontWeight: 'bold',
  },
  itemDetails: {
    fontFamily: 'monospace',
    fontSize: 10,
    color: '#666',
    textTransform: 'uppercase',
  },
});
