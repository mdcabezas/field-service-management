import React, { useEffect, useState } from 'react';
import { View, Text, StyleSheet, Alert } from 'react-native';
import { useRouter } from 'expo-router';
import { useAuthStore } from '../../stores/auth-store';
import { getDatabase } from '../../lib/database';
import { Header } from '../../components/ui/header';
import { Card } from '../../components/ui/card';
import { Badge } from '../../components/ui/badge';
import { Button } from '../../components/ui/button';
import { Input } from '../../components/ui/input';

export default function SettingsScreen() {
  const router = useRouter();
  const { user, logout } = useAuthStore();
  const [retentionDays, setRetentionDays] = useState('7');
  const [maxStorage, setMaxStorage] = useState('100');
  const [warningThreshold, setWarningThreshold] = useState('50');
  const [gpsThreshold, setGpsThreshold] = useState('50');
  const [email, setEmail] = useState(user?.email || '');
  const [password, setPassword] = useState('');

  useEffect(() => {
    loadConfig();
    if (user?.email) {
      setEmail(user.email);
    }
  }, [user]);

  const loadConfig = async () => {
    const db = getDatabase();
    const retentionStmt = db.prepareSync("SELECT value FROM config WHERE key = 'photo_retention_days'");
    const retentionResult = retentionStmt.executeSync();
    const retention = retentionResult.getFirstSync() as { value: string } | null;
    retentionStmt.finalizeSync();

    const maxStmt = db.prepareSync("SELECT value FROM config WHERE key = 'photo_max_storage_mb'");
    const maxResult = maxStmt.executeSync();
    const max = maxResult.getFirstSync() as { value: string } | null;
    maxStmt.finalizeSync();

    const warningStmt = db.prepareSync("SELECT value FROM config WHERE key = 'photo_warning_threshold_mb'");
    const warningResult = warningStmt.executeSync();
    const warning = warningResult.getFirstSync() as { value: string } | null;
    warningStmt.finalizeSync();

    const gpsStmt = db.prepareSync("SELECT value FROM config WHERE key = 'gps_accuracy_threshold'");
    const gpsResult = gpsStmt.executeSync();
    const gps = gpsResult.getFirstSync() as { value: string } | null;
    gpsStmt.finalizeSync();

    if (retention) setRetentionDays(retention.value);
    if (max) setMaxStorage(max.value);
    if (warning) setWarningThreshold(warning.value);
    if (gps) setGpsThreshold(gps.value);
  };

  const saveConfig = async () => {
    const db = getDatabase();
    await db.runAsync(
      "INSERT OR REPLACE INTO config (key, value) VALUES ('photo_retention_days', ?)",
      retentionDays
    );
    await db.runAsync(
      "INSERT OR REPLACE INTO config (key, value) VALUES ('photo_max_storage_mb', ?)",
      maxStorage
    );
    await db.runAsync(
      "INSERT OR REPLACE INTO config (key, value) VALUES ('photo_warning_threshold_mb', ?)",
      warningThreshold
    );
    await db.runAsync(
      "INSERT OR REPLACE INTO config (key, value) VALUES ('gps_accuracy_threshold', ?)",
      gpsThreshold
    );
    Alert.alert('Guardado', 'Configuración actualizada');
  };

  const handleReconnect = () => {
    if (!email.trim()) {
      Alert.alert('Error', 'Ingresa un email');
      return;
    }
    if (!password.trim()) {
      Alert.alert('Error', 'Ingresa la contraseña');
      return;
    }

    Alert.alert(
      'Reconectar',
      `Se cerrará la sesión actual y se conectará como ${email}. ¿Continuar?`,
      [
        { text: 'Cancelar', style: 'cancel' },
        {
          text: 'Reconectar',
          onPress: async () => {
            await logout();
            router.replace({
              pathname: '/login',
              params: { email, password },
            });
          },
        },
      ]
    );
  };

  const handleLogout = () => {
    Alert.alert('Cerrar sesión', '¿Estás seguro?', [
      { text: 'Cancelar', style: 'cancel' },
      {
        text: 'Cerrar sesión',
        style: 'destructive',
        onPress: async () => {
          await logout();
          router.replace('/login');
        },
      },
    ]);
  };

  return (
    <View style={styles.container}>
      <Header title="Configuración" />

      <View style={styles.content}>
        <Card title="Credenciales">
          <Input
            label="Email"
            value={email}
            onChangeText={setEmail}
            placeholder="usuario@ejemplo.com"
            keyboardType="email-address"
            autoCapitalize="none"
          />
          <Input
            label="Contraseña"
            value={password}
            onChangeText={setPassword}
            secureTextEntry
            placeholder="Contraseña"
          />
          <Button
            title="Guardar y reconectar"
            onPress={handleReconnect}
            variant="secondary"
            style={styles.reconnectButton}
          />
        </Card>

        <Card title="Sesión actual" style={styles.section}>
          <View style={styles.row}>
            <Text style={styles.label}>Email:</Text>
            <Text style={styles.value}>{user?.email || 'N/A'}</Text>
          </View>
          <View style={styles.row}>
            <Text style={styles.label}>Nombre:</Text>
            <Text style={styles.value}>{user?.name || 'N/A'}</Text>
          </View>
        </Card>

        <Card title="Fotos" style={styles.section}>
          <Input
            label="Días de retención"
            value={retentionDays}
            onChangeText={setRetentionDays}
            placeholder="7"
          />
          <Input
            label="Almacenamiento máximo (MB)"
            value={maxStorage}
            onChangeText={setMaxStorage}
            placeholder="100"
          />
          <Input
            label="Umbral de alerta (MB)"
            value={warningThreshold}
            onChangeText={setWarningThreshold}
            placeholder="50"
          />
          <Button title="Guardar configuración" onPress={saveConfig} />
        </Card>

        <Card title="GPS" style={styles.section}>
          <Input
            label="Umbral precisión (m)"
            value={gpsThreshold}
            onChangeText={setGpsThreshold}
            placeholder="50"
          />
          <Button title="Guardar configuración" onPress={saveConfig} />
        </Card>

        <Card title="Base de datos" style={styles.section}>
          <Button
            title="Limpiar fotos antiguas"
            onPress={() => Alert.alert('Limpiando...')}
            variant="secondary"
          />
        </Card>

        <Button
          title="Cerrar sesión"
          onPress={handleLogout}
          variant="danger"
          style={styles.logoutButton}
        />
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
  section: {
    marginTop: 16,
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
  reconnectButton: {
    marginTop: 8,
  },
  logoutButton: {
    marginTop: 24,
  },
});
