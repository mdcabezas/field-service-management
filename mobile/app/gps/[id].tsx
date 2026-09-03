import React, { useEffect, useState } from 'react';
import { View, Text, FlatList, StyleSheet, Alert, ActivityIndicator } from 'react-native';
import { useLocalSearchParams } from 'expo-router';
import { getDatabase } from '../../lib/database';
import { captureGPS, captureGPSWithAccuracy, checkDistance, createCheckpoint, createGPSWaiver, MAX_ACCURACY } from '../../lib/gps';
import { Header } from '../../components/ui/header';
import { Card } from '../../components/ui/card';
import { Badge } from '../../components/ui/badge';
import { Button } from '../../components/ui/button';

interface Checkpoint {
  id: string;
  type: string;
  lat: number;
  lng: number;
  accuracy: number;
  distance_from_property: number | null;
  timestamp: string;
}

export default function GPSScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const [checkpoints, setCheckpoints] = useState<Checkpoint[]>([]);
  const [isCapturing, setIsCapturing] = useState(false);
  const [propertyLat, setPropertyLat] = useState<number | null>(null);
  const [propertyLng, setPropertyLng] = useState<number | null>(null);
  const [gpsThreshold, setGpsThreshold] = useState(50);
  const [isCapturingGPS, setIsCapturingGPS] = useState(false);
  const [currentAccuracy, setCurrentAccuracy] = useState<number | null>(null);

  useEffect(() => {
    loadCheckpoints();
    loadPropertyGPS();
    loadGpsThreshold();
  }, [id]);

  const loadCheckpoints = async () => {
    const db = getDatabase();
    const stmt = db.prepareSync(
      'SELECT * FROM checkpoints WHERE visit_id = ? ORDER BY timestamp DESC'
    );
    const result = stmt.executeSync(id);
    const rows = result.getAllSync() as Record<string, unknown>[];
    stmt.finalizeSync();
    setCheckpoints(
      rows.map((row) => ({
        id: row.id as string,
        type: row.type as string,
        lat: row.lat as number,
        lng: row.lng as number,
        accuracy: row.accuracy as number,
        distance_from_property: row.distance_from_property as number | null,
        timestamp: row.timestamp as string,
      }))
    );
  };

  const loadPropertyGPS = async () => {
    const db = getDatabase();
    const stmt = db.prepareSync('SELECT property_lat, property_lng FROM visits WHERE id = ?');
    const result = stmt.executeSync(id);
    const row = result.getFirstSync() as Record<string, unknown> | null;
    stmt.finalizeSync();
    if (row) {
      setPropertyLat(row.property_lat as number);
      setPropertyLng(row.property_lng as number);
    }
  };

  const loadGpsThreshold = async () => {
    const db = getDatabase();
    const stmt = db.prepareSync("SELECT value FROM config WHERE key = 'gps_accuracy_threshold'");
    const result = stmt.executeSync();
    const row = result.getFirstSync() as { value: string } | null;
    stmt.finalizeSync();
    if (row) {
      setGpsThreshold(parseInt(row.value, 10) || 50);
    }
  };

  const handleCapture = async (type: string) => {
    setIsCapturing(true);
    setIsCapturingGPS(true);
    setCurrentAccuracy(null);

    try {
      const gps = await captureGPSWithAccuracy(
        (accuracy) => setCurrentAccuracy(accuracy),
        30000,
        gpsThreshold
      );

      setIsCapturingGPS(false);
      setCurrentAccuracy(null);

      if (!gps) {
        Alert.alert(
          'GPS no disponible',
          'No se pudo obtener la ubicación. ¿Deseas aceptar una exención?',
          [
            { text: 'Cancelar', style: 'cancel' },
            {
              text: 'Aceptar exención',
              onPress: async () => {
                await createGPSWaiver(id!, type);
                loadCheckpoints();
              },
            },
          ]
        );
        return;
      }

      if (propertyLat && propertyLng) {
        const distance = checkDistance(gps.lat, gps.lng, propertyLat, propertyLng);
        if (!distance.isWithinThreshold) {
          Alert.alert(
            'Distancia excedida',
            `Estás a ${Math.round(distance.distance)}m de la propiedad (umbral: 500m)`,
            [
              { text: 'Cancelar', style: 'cancel' },
              {
                text: 'Continuar de todos modos',
                onPress: async () => {
                  await createCheckpoint(id!, type, gps, propertyLat, propertyLng);
                  loadCheckpoints();
                },
              },
            ]
          );
          return;
        }
      }

      await createCheckpoint(id!, type, gps, propertyLat || undefined, propertyLng || undefined);
      loadCheckpoints();
    } catch (error) {
      setIsCapturingGPS(false);
      setCurrentAccuracy(null);
      Alert.alert('Error', 'Error al capturar GPS');
    } finally {
      setIsCapturing(false);
    }
  };

  const handleRetryGPS = (type: string) => {
    handleCapture(type);
  };

  const renderCheckpoint = ({ item }: { item: Checkpoint }) => (
    <Card style={styles.checkpointCard}>
      <View style={styles.checkpointHeader}>
        <Badge
          label={item.type === 'arrival' ? 'Llegada' : 'Salida'}
          color={item.type === 'arrival' ? '#059669' : '#dc2626'}
        />
        <Text style={styles.checkpointTime}>
          {new Date(item.timestamp).toLocaleTimeString()}
        </Text>
      </View>
      <View style={styles.checkpointDetails}>
        <Text style={styles.detail}>Lat: {item.lat.toFixed(6)}</Text>
        <Text style={styles.detail}>Lng: {item.lng.toFixed(6)}</Text>
        <Text style={styles.detail}>Precisión: {Math.round(item.accuracy)}m</Text>
        {item.distance_from_property !== null && (
          <Text style={styles.detail}>
            Distancia: {Math.round(item.distance_from_property)}m
          </Text>
        )}
      </View>
    </Card>
  );

  return (
    <View style={styles.container}>
      <Header
        title="GPS"
        right={<Badge label={`${checkpoints.length}`} color="#000" />}
      />

      <View style={styles.actions}>
        <Button
          title={isCapturing ? 'Capturando...' : 'Registrar llegada'}
          onPress={() => handleCapture('arrival')}
          disabled={isCapturing}
          style={styles.actionButton}
        />
        <Button
          title={isCapturing ? 'Capturando...' : 'Registrar salida'}
          onPress={() => handleCapture('departure')}
          disabled={isCapturing}
          variant="secondary"
          style={styles.actionButton}
        />
      </View>

      {isCapturingGPS && (
        <Card style={styles.gpsProgressCard}>
          <View style={styles.gpsProgressContent}>
            <Text style={styles.gpsProgressTitle}>Obteniendo GPS...</Text>
            <View style={styles.gpsProgressAccuracy}>
              <Text style={styles.gpsAccuracyLabel}>
                Precisión actual:
                {' '}
                <Text style={styles.gpsAccuracyValue}>
                  {currentAccuracy !== null ? `${Math.round(currentAccuracy)}m` : 'Buscando...'}
                </Text>
              </Text>
              <Text style={styles.gpsRequiredLabel}>
                Requerido: ≤{gpsThreshold}m
              </Text>
            </View>
            <ActivityIndicator size="large" color="#000" style={styles.activityIndicator} />
          </View>
        </Card>
      )}

      {propertyLat && propertyLng && (
        <Card style={styles.propertyGPS}>
          <Text style={styles.propertyLabel}>GPS de la propiedad:</Text>
          <Text style={styles.propertyValue}>
            {propertyLat.toFixed(6)}, {propertyLng.toFixed(6)}
          </Text>
        </Card>
      )}

      <FlatList
        data={checkpoints}
        renderItem={renderCheckpoint}
        keyExtractor={(item) => item.id}
        contentContainerStyle={styles.list}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  actions: {
    flexDirection: 'row',
    padding: 16,
    gap: 8,
  },
  actionButton: {
    flex: 1,
  },
  propertyGPS: {
    marginHorizontal: 16,
    marginBottom: 16,
  },
  propertyLabel: {
    fontFamily: 'monospace',
    fontSize: 10,
    textTransform: 'uppercase',
    color: '#666',
  },
  propertyValue: {
    fontFamily: 'monospace',
    fontSize: 12,
  },
  list: {
    padding: 16,
  },
  checkpointCard: {
    marginBottom: 12,
  },
  checkpointHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  checkpointTime: {
    fontFamily: 'monospace',
    fontSize: 12,
  },
  checkpointDetails: {
    gap: 4,
  },
  detail: {
    fontFamily: 'monospace',
    fontSize: 11,
  },
  gpsProgressCard: {
    marginHorizontal: 16,
    marginBottom: 16,
    padding: 16,
  },
  gpsProgressContent: {
    alignItems: 'center',
  },
  gpsProgressTitle: {
    fontFamily: 'monospace',
    fontSize: 14,
    fontWeight: 'bold',
    marginBottom: 12,
  },
  gpsProgressAccuracy: {
    alignItems: 'center',
    marginBottom: 16,
  },
  gpsAccuracyLabel: {
    fontFamily: 'monospace',
    fontSize: 14,
    marginBottom: 4,
  },
  gpsAccuracyValue: {
    fontFamily: 'monospace',
    fontSize: 14,
    fontWeight: 'bold',
    color: '#dc2626',
  },
  gpsRequiredLabel: {
    fontFamily: 'monospace',
    fontSize: 12,
    color: '#666',
    marginTop: 4,
  },
  activityIndicator: {
    marginTop: 8,
  },
});