import React, { useEffect, useState } from 'react';
import { View, Text, ScrollView, Pressable, StyleSheet, Alert } from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { useVisitStore } from '../../stores/visit-store';
import { getStatusColor } from '../../utils/status-transition';
import { Header } from '../../components/ui/header';
import { Card } from '../../components/ui/card';
import { Badge } from '../../components/ui/badge';
import { Button } from '../../components/ui/button';
import { StatusTransitionComponent } from '../../components/ui/status-transition';

export default function VisitDetailScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const { currentVisit, loadVisitDetail, updateVisitStatus } = useVisitStore();
  const [showHistory, setShowHistory] = useState(false);
  const [showContacts, setShowContacts] = useState(false);

  useEffect(() => {
    if (id) loadVisitDetail(id);
  }, [id]);

  if (!currentVisit) {
    return (
      <View style={styles.container}>
        <Header title="Cargando..." />
        <Text style={styles.loading}>Cargando detalles...</Text>
      </View>
    );
  }

  const handleStatusChange = (newStatus: string) => {
    Alert.alert(
      'Confirmar cambio de estado',
      `¿Cambiar estado a "${newStatus}"?`,
      [
        { text: 'Cancelar', style: 'cancel' },
        {
          text: 'Confirmar',
          onPress: async () => {
            try {
              await updateVisitStatus(currentVisit.id, newStatus);
            } catch (error) {
              Alert.alert('Error', 'No se pudo cambiar el estado');
            }
          },
        },
      ]
    );
  };

  return (
    <View style={styles.container}>
      <Header
        title={currentVisit.property_name || 'Visita'}
        right={
          <Badge
            label={currentVisit.status}
            color={getStatusColor(currentVisit.status as any)}
          />
        }
      />

      <ScrollView style={styles.content}>
        <Card title="Propiedad">
          <Text style={styles.label}>Nombre:</Text>
          <Text style={styles.value}>{String(currentVisit.property_name || 'Sin nombre')}</Text>
          <Text style={styles.label}>Dirección:</Text>
          <Text style={styles.value}>{String(currentVisit.property_address || 'Sin dirección')}</Text>
          {currentVisit.property_lat && currentVisit.property_lng && (
            <View>
              <Text style={styles.label}>GPS:</Text>
              <Text style={styles.value}>
                {String(currentVisit.property_lat.toFixed(6))}, {String(currentVisit.property_lng.toFixed(6))}
              </Text>
            </View>
          )}
        </Card>

        {currentVisit.partner_name && (
          <Card title="Partner">
            <Text style={styles.label}>Nombre:</Text>
            <Text style={styles.value}>{String(currentVisit.partner_name)}</Text>
            {currentVisit.partner_order_id && (
              <View>
                <Text style={styles.label}>Orden:</Text>
                <Text style={styles.value}>{String(currentVisit.partner_order_id)}</Text>
              </View>
            )}
            {currentVisit.partner_supervisor && (
              <View>
                <Text style={styles.label}>Supervisor:</Text>
                <Text style={styles.value}>{String(currentVisit.partner_supervisor)}</Text>
              </View>
            )}
            <Pressable onPress={() => setShowContacts(!showContacts)}>
              <Text style={styles.link}>
                {showContacts ? 'Ocultar contactos' : 'Ver contactos'}
              </Text>
            </Pressable>
          </Card>
        )}

        <Card title="Información">
          <Text style={styles.label}>Tipo:</Text>
          <Text style={styles.value}>{String(currentVisit.type || 'Sin tipo')}</Text>
          <Text style={styles.label}>Prioridad:</Text>
          <Text style={styles.value}>{String(currentVisit.priority)}</Text>
          <Text style={styles.label}>Programada:</Text>
          <Text style={styles.value}>
            {currentVisit.scheduled_at
              ? String(new Date(currentVisit.scheduled_at).toLocaleString())
              : 'Sin fecha'}
          </Text>
          {currentVisit.sla_response_hours && (
            <View>
              <Text style={styles.label}>SLA Respuesta:</Text>
              <Text style={styles.value}>{String(currentVisit.sla_response_hours)}h</Text>
            </View>
          )}
          {currentVisit.sla_resolution_hours && (
            <View>
              <Text style={styles.label}>SLA Resolución:</Text>
              <Text style={styles.value}>{String(currentVisit.sla_resolution_hours)}h</Text>
            </View>
          )}
        </Card>

        <Card title="Estado">
          <StatusTransitionComponent
            currentStatus={currentVisit.status}
            onTransition={handleStatusChange}
          />
        </Card>

        <Card title="Acciones">
          <View style={styles.actions}>
            <Button
              title="Checklist"
              onPress={() => router.push(`/checklist/${id}`)}
              variant="secondary"
              style={styles.actionButton}
            />
            <Button
              title="Fotos"
              onPress={() => router.push(`/photo/${id}`)}
              variant="secondary"
              style={styles.actionButton}
            />
            <Button
              title="GPS"
              onPress={() => router.push(`/gps/${id}`)}
              variant="secondary"
              style={styles.actionButton}
            />
            <Button
              title="Mediciones"
              onPress={() => router.push(`/visit/${id}/measurements`)}
              variant="secondary"
              style={styles.actionButton}
            />
            <Button
              title="Materiales"
              onPress={() => router.push(`/visit/${id}/usage`)}
              variant="secondary"
              style={styles.actionButton}
            />
          </View>
        </Card>

        <Pressable onPress={() => setShowHistory(!showHistory)}>
          <Card title="Historial">
            <Text style={styles.link}>
              {showHistory ? 'Ocultar historial' : 'Ver historial de la propiedad'}
            </Text>
          </Card>
        </Pressable>
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  loading: {
    padding: 16,
    fontFamily: 'monospace',
    textAlign: 'center',
  },
  content: {
    flex: 1,
    padding: 16,
  },
  label: {
    fontFamily: 'monospace',
    fontSize: 10,
    textTransform: 'uppercase',
    color: '#666',
    marginTop: 8,
  },
  value: {
    fontFamily: 'monospace',
    fontSize: 14,
    marginBottom: 4,
  },
  link: {
    fontFamily: 'monospace',
    fontSize: 12,
    color: '#3b82f6',
    marginTop: 8,
  },
  actions: {
    gap: 8,
  },
  actionButton: {
    marginTop: 0,
  },
});
