import React, { useState } from 'react';
import { View, Text, StyleSheet, Alert } from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { apiRequest } from '../../../lib/auth';
import { getDatabase } from '../../../lib/database';
import { addToOutbox } from '../../../lib/outbox';
import { useVisitStore } from '../../../stores/visit-store';
import { Header } from '../../../components/ui/header';
import { Card } from '../../../components/ui/card';
import { Input } from '../../../components/ui/input';
import { Button } from '../../../components/ui/button';

export default function ReopenScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const { loadVisitDetail } = useVisitStore();
  const [reason, setReason] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleReopen = async () => {
    if (!reason.trim()) {
      Alert.alert('Error', 'Ingresa un motivo para reabrir');
      return;
    }

    setIsSubmitting(true);
    try {
      await apiRequest(`/api/visits/${id}/reopen`, {
        method: 'POST',
        body: JSON.stringify({ reason }),
      });

      const db = getDatabase();
      await db.runAsync(
        "UPDATE visits SET status = 'assigned', updated_at = datetime('now') WHERE id = ?",
        id
      );

      await addToOutbox('status', id as string, 'update', {
        visit_id: id,
        status: 'assigned',
        reason,
        timestamp: new Date().toISOString(),
      });

      await loadVisitDetail(id as string);

      Alert.alert('Éxito', 'Visita reabierta', [
        { text: 'OK', onPress: () => router.back() },
      ]);
    } catch (error) {
      Alert.alert('Error', 'No se pudo reabrir la visita');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <View style={styles.container}>
      <Header title="Reabrir visita" />

      <View style={styles.content}>
        <Card title="Motivo de reapertura">
          <Text style={styles.description}>
            Ingresa el motivo por el cual deseas reabrir esta visita. La visita debe haber sido completada en las últimas 24 horas.
          </Text>
          <Input
            label="Motivo"
            value={reason}
            onChangeText={setReason}
            placeholder="Ej: Se olvidó documentar un paso"
            multiline
            numberOfLines={4}
          />
          <Button
            title={isSubmitting ? 'Reabriendo...' : 'Reabrir visita'}
            onPress={handleReopen}
            disabled={isSubmitting}
            style={styles.button}
          />
          <Button
            title="Cancelar"
            onPress={() => router.back()}
            variant="secondary"
            style={styles.button}
          />
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
  description: {
    fontFamily: 'monospace',
    fontSize: 12,
    color: '#666',
    marginBottom: 16,
  },
  button: {
    marginTop: 8,
  },
});
