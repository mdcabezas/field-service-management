import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, Alert } from 'react-native';
import { useRouter, useLocalSearchParams } from 'expo-router';
import { useAuthStore } from '../stores/auth-store';
import { Header } from '../components/ui/header';
import { Card } from '../components/ui/card';
import { Input } from '../components/ui/input';
import { Button } from '../components/ui/button';

export default function LoginScreen() {
  const router = useRouter();
  const params = useLocalSearchParams<{ email?: string; password?: string }>();
  const { login, isLoading, error } = useAuthStore();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  useEffect(() => {
    if (params.email) setEmail(params.email);
    if (params.password) setPassword(params.password);
  }, [params.email, params.password]);

  const handleLogin = async () => {
    if (!email || !password) {
      Alert.alert('Error', 'Ingresa email y contraseña');
      return;
    }

    try {
      await login(email, password);
      router.replace('/(tabs)');
    } catch (err) {
      Alert.alert('Error de login', String(err));
    }
  };

  return (
    <View style={styles.container}>
      <Header title="Localis Mobile" />

      <View style={styles.content}>
        <Card title="Iniciar sesión">
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

          {error && <Text style={styles.error}>{error}</Text>}

          <Button
            title={isLoading ? 'Ingresando...' : 'Iniciar sesión'}
            onPress={handleLogin}
            disabled={isLoading}
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
    justifyContent: 'center',
    padding: 16,
  },
  error: {
    fontFamily: 'monospace',
    fontSize: 12,
    color: '#dc2626',
    marginBottom: 12,
  },
  button: {
    marginTop: 8,
  },
});
