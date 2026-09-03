import React, { useEffect, useState } from 'react';
import { View, Text, FlatList, StyleSheet, Alert, Pressable } from 'react-native';
import { useLocalSearchParams } from 'expo-router';
import { getDatabase } from '../../../lib/database';
import { addToOutbox } from '../../../lib/outbox';
import * as Crypto from 'expo-crypto';
import { Header } from '../../../components/ui/header';
import { Card } from '../../../components/ui/card';
import { Badge } from '../../../components/ui/badge';
import { Button } from '../../../components/ui/button';
import { Input } from '../../../components/ui/input';

interface Material {
  id: string;
  name: string;
  sku: string;
  unit: string;
  planned_quantity?: number;
}

interface Usage {
  id: string;
  material_id: string;
  material_name: string;
  quantity: number;
  notes: string;
  created_at: string;
}

interface MaterialSummary {
  material: Material;
  planned: number;
  used: number;
  difference: number;
}

export default function UsageScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const [materials, setMaterials] = useState<Material[]>([]);
  const [usages, setUsages] = useState<Usage[]>([]);
  const [selectedMaterial, setSelectedMaterial] = useState<Material | null>(null);
  const [quantity, setQuantity] = useState('');
  const [notes, setNotes] = useState('');
  const [isAdding, setIsAdding] = useState(false);

  useEffect(() => {
    loadMaterials();
    loadUsages();
  }, [id]);

  const loadMaterials = async () => {
    const db = getDatabase();
    const stmt = db.prepareSync('SELECT * FROM materials ORDER BY name');
    const result = stmt.executeSync();
    const rows = result.getAllSync() as Record<string, unknown>[];
    stmt.finalizeSync();
    setMaterials(
      rows.map((row) => ({
        id: row.id as string,
        name: row.name as string,
        sku: row.sku as string,
        unit: row.unit as string,
        planned_quantity: (row.planned_quantity as number) || 0,
      }))
    );
  };

  const loadUsages = async () => {
    const db = getDatabase();
    const stmt = db.prepareSync(
      `SELECT u.*, m.name as material_name
       FROM usages u
       LEFT JOIN materials m ON u.material_id = m.id
       WHERE u.visit_id = ?
       ORDER BY u.created_at DESC`
    );
    const result = stmt.executeSync(id);
    const rows = result.getAllSync() as Record<string, unknown>[];
    stmt.finalizeSync();
    setUsages(
      rows.map((row) => ({
        id: row.id as string,
        material_id: row.material_id as string,
        material_name: row.material_name as string,
        quantity: row.quantity as number,
        notes: row.notes as string,
        created_at: row.created_at as string,
      }))
    );
  };

  const getMaterialSummary = (): MaterialSummary[] => {
    const usageMap = new Map<string, number>();
    usages.forEach((u) => {
      const current = usageMap.get(u.material_id) || 0;
      usageMap.set(u.material_id, current + u.quantity);
    });

    return materials.map((m) => {
      const used = usageMap.get(m.id) || 0;
      const planned = m.planned_quantity || 0;
      return {
        material: m,
        planned,
        used,
        difference: used - planned,
      };
    });
  };

  const summary = getMaterialSummary();

  const handleAdd = async () => {
    console.log('[UsageScreen] handleAdd called', { selectedMaterial, quantity, id, isAdding });

    if (!id) {
      console.error('[UsageScreen] No visit ID available');
      Alert.alert('Error', 'ID de visita no válido. Vuelve a la lista de visitas.');
      return;
    }

    if (!selectedMaterial || !quantity.trim()) {
      console.log('[UsageScreen] Validation failed', { selectedMaterial: !!selectedMaterial, quantity: quantity.trim() });
      Alert.alert('Error', 'Selecciona material y cantidad');
      return;
    }

    setIsAdding(true);
    try {
      const usageId = Crypto.randomUUID();
      const db = getDatabase();

      console.log('[UsageScreen] Inserting usage', { usageId, visitId: id, materialId: selectedMaterial.id, quantity: parseFloat(quantity) });

      await db.runAsync(
        'INSERT INTO usages (id, visit_id, material_id, quantity, notes, synced) VALUES (?, ?, ?, ?, ?, 0)',
        usageId, id, selectedMaterial.id, parseFloat(quantity), notes
      );

      console.log('[UsageScreen] Adding to outbox');
      await addToOutbox('usage', usageId, 'create', {
        visit_id: id,
        material_id: selectedMaterial.id,
        quantity: parseFloat(quantity),
        notes,
        local_id: usageId,
      });

      console.log('[UsageScreen] Success, clearing form and reloading');
      setSelectedMaterial(null);
      setQuantity('');
      setNotes('');
      await loadUsages();
    } catch (error) {
      console.error('[UsageScreen] Add failed:', error);
      Alert.alert('Error', `Error al guardar: ${error instanceof Error ? error.message : String(error)}`);
    } finally {
      setIsAdding(false);
    }
  };

  const handleDelete = (usageId: string) => {
    Alert.alert('Eliminar', '¿Eliminar este uso?', [
      { text: 'Cancelar', style: 'cancel' },
      {
        text: 'Eliminar',
        style: 'destructive',
        onPress: async () => {
          const db = getDatabase();
          await db.runAsync('DELETE FROM usages WHERE id = ?', usageId);
          await addToOutbox('usage', usageId, 'delete', {
            visit_id: id,
            usage_id: usageId,
          });
          loadUsages();
        },
      },
    ]);
  };

  const renderUsage = ({ item }: { item: Usage }) => (
    <Card style={styles.usageCard}>
      <View style={styles.usageHeader}>
        <Text style={styles.usageName}>{item.material_name}</Text>
        <Button
          title="X"
          onPress={() => handleDelete(item.id)}
          variant="danger"
          style={styles.deleteButton}
        />
      </View>
      <Text style={styles.usageQuantity}>
        {item.quantity} {materials.find((m) => m.id === item.material_id)?.unit || ''}
      </Text>
      {item.notes && <Text style={styles.usageNotes}>{item.notes}</Text>}
      <Text style={styles.usageDate}>
        {new Date(item.created_at).toLocaleString()}
      </Text>
    </Card>
  );

  return (
    <View style={styles.container}>
      <Header
        title="Materiales"
        right={<Badge label={`${usages.length}`} color="#000" />}
      />

      {summary.some((s) => s.planned > 0) && (
        <Card title="Planificado vs Usado" style={styles.summaryCard}>
          {summary.filter((s) => s.planned > 0).map((s) => (
            <View key={s.material.id} style={styles.summaryRow}>
              <Text style={styles.summaryName}>{s.material.name}</Text>
              <View style={styles.summaryValues}>
                <Text style={styles.summaryPlanned}>
                  Plan: {s.planned} {s.material.unit}
                </Text>
                <Text style={styles.summaryUsed}>
                  Usado: {s.used} {s.material.unit}
                </Text>
                <Badge
                  label={s.difference > 0 ? `+${s.difference}` : `${s.difference}`}
                  color={s.difference > 0 ? '#dc2626' : s.difference < 0 ? '#059669' : '#6b7280'}
                />
              </View>
            </View>
          ))}
        </Card>
      )}

      <Card title="Agregar material" style={styles.addForm}>
        <Text style={styles.label}>Seleccionar material:</Text>
        {materials.length === 0 ? (
          <Text style={styles.noMaterials}>No hay materiales disponibles. Ejecuta una sincronización primero.</Text>
        ) : (
          <>
            <FlatList
              data={materials}
              horizontal
              showsHorizontalScrollIndicator={false}
              renderItem={({ item }) => (
                <Pressable
                  onPress={() => setSelectedMaterial(item)}
                  style={[
                    styles.materialChip,
                    selectedMaterial?.id === item.id && styles.materialChipSelected,
                  ]}
                  key={item.id}
                >
                  <Text style={[
                    styles.materialChipText,
                    selectedMaterial?.id === item.id && styles.materialChipTextSelected,
                  ]}>
                    {item.name}
                  </Text>
                </Pressable>
              )}
              keyExtractor={(item) => item.id}
              style={styles.materialList}
            />
          </>
        )}

        {selectedMaterial && (
          <Text style={styles.selectedMaterial}>
            Seleccionado: {selectedMaterial.name} ({selectedMaterial.sku})
          </Text>
        )}

        <View style={styles.valueRow}>
          <Input
            label="Cantidad"
            value={quantity}
            onChangeText={setQuantity}
            placeholder="0"
            style={styles.quantityInput}
          />
          <Input
            label="Notas"
            value={notes}
            onChangeText={setNotes}
            placeholder="Opcional"
            style={styles.notesInput}
          />
        </View>

        <Button
          title={isAdding ? 'Guardando...' : 'Agregar'}
          onPress={handleAdd}
          disabled={isAdding || !selectedMaterial}
        />
      </Card>

      <FlatList
        data={usages}
        renderItem={renderUsage}
        keyExtractor={(item) => item.id}
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
  summaryCard: {
    margin: 16,
    marginBottom: 0,
  },
  summaryRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 8,
    borderBottomWidth: 1,
    borderBottomColor: '#eee',
  },
  summaryName: {
    fontFamily: 'monospace',
    fontSize: 12,
    fontWeight: 'bold',
    flex: 1,
  },
  summaryValues: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  summaryPlanned: {
    fontFamily: 'monospace',
    fontSize: 10,
    color: '#666',
  },
  summaryUsed: {
    fontFamily: 'monospace',
    fontSize: 10,
    color: '#333',
  },
  addForm: {
    margin: 16,
  },
  label: {
    fontFamily: 'monospace',
    fontSize: 10,
    textTransform: 'uppercase',
    marginBottom: 8,
  },
  materialList: {
    marginBottom: 8,
  },
  materialChip: {
    marginRight: 8,
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderWidth: 2,
    borderColor: '#000',
    borderRadius: 0,
    backgroundColor: '#fff',
  },
  materialChipSelected: {
    backgroundColor: '#000',
    borderColor: '#000',
  },
  materialChipText: {
    fontFamily: 'monospace',
    fontSize: 12,
    fontWeight: 'bold',
    textTransform: 'uppercase',
    color: '#000',
  },
  materialChipTextSelected: {
    color: '#fff',
  },
  noMaterials: {
    fontFamily: 'monospace',
    fontSize: 12,
    color: '#f59e0b',
    marginTop: 8,
    marginBottom: 12,
  },
  selectedMaterial: {
    fontFamily: 'monospace',
    fontSize: 12,
    color: '#3b82f6',
    marginBottom: 12,
  },
  valueRow: {
    flexDirection: 'row',
    gap: 8,
  },
  quantityInput: {
    flex: 1,
  },
  notesInput: {
    flex: 2,
  },
  list: {
    padding: 16,
  },
  usageCard: {
    marginBottom: 12,
  },
  usageHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 4,
  },
  usageName: {
    fontFamily: 'monospace',
    fontSize: 14,
    fontWeight: 'bold',
  },
  deleteButton: {
    paddingHorizontal: 8,
    paddingVertical: 4,
  },
  usageQuantity: {
    fontFamily: 'monospace',
    fontSize: 16,
  },
  usageNotes: {
    fontFamily: 'monospace',
    fontSize: 11,
    color: '#666',
    marginTop: 4,
  },
  usageDate: {
    fontFamily: 'monospace',
    fontSize: 10,
    color: '#999',
    marginTop: 4,
  },
});
