import React, { useEffect, useState } from 'react';
import { View, Text, FlatList, StyleSheet, Alert } from 'react-native';
import { useLocalSearchParams } from 'expo-router';
import { getDatabase } from '../../../lib/database';
import { captureGPS } from '../../../lib/gps';
import { addToOutbox } from '../../../lib/outbox';
import * as Crypto from 'expo-crypto';
import { Header } from '../../../components/ui/header';
import { Card } from '../../../components/ui/card';
import { Badge } from '../../../components/ui/badge';
import { Button } from '../../../components/ui/button';
import { Input } from '../../../components/ui/input';

interface Measurement {
  id: string;
  key: string;
  value: string;
  unit: string;
  lat: number | null;
  lng: number | null;
  created_at: string;
}

export default function MeasurementScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const [measurements, setMeasurements] = useState<Measurement[]>([]);
  const [key, setKey] = useState('');
  const [value, setValue] = useState('');
  const [unit, setUnit] = useState('');
  const [isAdding, setIsAdding] = useState(false);

  useEffect(() => {
    loadMeasurements();
  }, [id]);

  const loadMeasurements = async () => {
    const db = getDatabase();
    const stmt = db.prepareSync(
      'SELECT * FROM measurements WHERE visit_id = ? ORDER BY created_at DESC'
    );
    const result = stmt.executeSync(id);
    const rows = result.getAllSync() as Record<string, unknown>[];
    stmt.finalizeSync();
    setMeasurements(
      rows.map((row) => ({
        id: row.id as string,
        key: row.key as string,
        value: row.value as string,
        unit: row.unit as string,
        lat: row.lat as number | null,
        lng: row.lng as number | null,
        created_at: row.created_at as string,
      }))
    );
  };

  const handleAdd = async () => {
    if (!key.trim() || !value.trim()) {
      Alert.alert('Error', 'Ingresa clave y valor');
      return;
    }

    setIsAdding(true);
    try {
      const gps = await captureGPS();
      const measurementId = Crypto.randomUUID();

      const db = getDatabase();
      await db.runAsync(
        'INSERT INTO measurements (id, visit_id, key, value, unit, lat, lng, synced) VALUES (?, ?, ?, ?, ?, ?, ?, 0)',
        measurementId, id, key, value, unit, gps?.lat || null, gps?.lng || null
      );

      await addToOutbox('measurement', measurementId, 'create', {
        visit_id: id,
        key,
        value,
        unit,
        lat: gps?.lat,
        lng: gps?.lng,
        local_id: measurementId,
      });

      setKey('');
      setValue('');
      setUnit('');
      loadMeasurements();
    } catch (error) {
      Alert.alert('Error', 'Error al guardar');
    } finally {
      setIsAdding(false);
    }
  };

  const handleDelete = (measurementId: string) => {
    Alert.alert('Eliminar', '¿Eliminar esta medición?', [
      { text: 'Cancelar', style: 'cancel' },
      {
        text: 'Eliminar',
        style: 'destructive',
        onPress: async () => {
          const db = getDatabase();
          await db.runAsync('DELETE FROM measurements WHERE id = ?', measurementId);
          await addToOutbox('measurement', measurementId, 'delete', {
            visit_id: id,
            measurement_id: measurementId,
          });
          loadMeasurements();
        },
      },
    ]);
  };

  const renderMeasurement = ({ item }: { item: Measurement }) => (
    <Card style={styles.measurementCard}>
      <View style={styles.measurementHeader}>
        <Text style={styles.measurementKey}>{item.key}</Text>
        <Button
          title="X"
          onPress={() => handleDelete(item.id)}
          variant="danger"
          style={styles.deleteButton}
        />
      </View>
      <Text style={styles.measurementValue}>
        {item.value} {item.unit}
      </Text>
      {item.lat && item.lng && (
        <Text style={styles.measurementGPS}>
          GPS: {item.lat.toFixed(6)}, {item.lng.toFixed(6)}
        </Text>
      )}
      <Text style={styles.measurementDate}>
        {new Date(item.created_at).toLocaleString()}
      </Text>
    </Card>
  );

  return (
    <View style={styles.container}>
      <Header
        title="Mediciones"
        right={<Badge label={`${measurements.length}`} color="#000" />}
      />

      <Card title="Agregar medición" style={styles.addForm}>
        <Input
          label="Clave"
          value={key}
          onChangeText={setKey}
          placeholder="Ej: Voltaje"
        />
        <View style={styles.valueRow}>
          <Input
            label="Valor"
            value={value}
            onChangeText={setValue}
            placeholder="Ej: 220"
            style={styles.valueInput}
          />
          <Input
            label="Unidad"
            value={unit}
            onChangeText={setUnit}
            placeholder="Ej: V"
            style={styles.unitInput}
          />
        </View>
        <Button
          title={isAdding ? 'Guardando...' : 'Agregar'}
          onPress={handleAdd}
          disabled={isAdding}
        />
      </Card>

      <FlatList
        data={measurements}
        renderItem={renderMeasurement}
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
  addForm: {
    margin: 16,
  },
  valueRow: {
    flexDirection: 'row',
    gap: 8,
  },
  valueInput: {
    flex: 2,
  },
  unitInput: {
    flex: 1,
  },
  list: {
    padding: 16,
  },
  measurementCard: {
    marginBottom: 12,
  },
  measurementHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 4,
  },
  measurementKey: {
    fontFamily: 'monospace',
    fontSize: 14,
    fontWeight: 'bold',
    textTransform: 'uppercase',
  },
  deleteButton: {
    paddingHorizontal: 8,
    paddingVertical: 4,
  },
  measurementValue: {
    fontFamily: 'monospace',
    fontSize: 18,
    marginBottom: 4,
  },
  measurementGPS: {
    fontFamily: 'monospace',
    fontSize: 10,
    color: '#666',
  },
  measurementDate: {
    fontFamily: 'monospace',
    fontSize: 10,
    color: '#999',
    marginTop: 4,
  },
});
