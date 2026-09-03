import { useFonts } from 'expo-font';
import { Stack, Redirect } from 'expo-router';
import * as SplashScreen from 'expo-splash-screen';
import { useEffect, useState } from 'react';
import { View, Text, StyleSheet } from 'react-native';
import 'react-native-reanimated';
import { useAuthStore } from '../stores/auth-store';
import { initDatabase } from '../lib/database';

export {
  ErrorBoundary,
} from 'expo-router';

export const unstable_settings = {
  initialRouteName: '(tabs)',
};

SplashScreen.preventAutoHideAsync();

export default function RootLayout() {
  const [loaded, error] = useFonts({
    SpaceMono: require('../assets/fonts/SpaceMono-Regular.ttf'),
  });
  const [dbReady, setDbReady] = useState(false);
  const [dbError, setDbError] = useState<string | null>(null);

  useEffect(() => {
    if (error) throw error;
  }, [error]);

  useEffect(() => {
    initDatabase()
      .then(async () => {
        setDbReady(true);
        await useAuthStore.getState().checkAuth();
      })
      .catch((err) => {
        console.error('Database init failed:', err);
        setDbError(String(err));
        SplashScreen.hideAsync();
      });
  }, []);

  useEffect(() => {
    if (loaded && dbReady) {
      SplashScreen.hideAsync();
    }
  }, [loaded, dbReady]);

  if (dbError) {
    return (
      <View style={errorStyles.container}>
        <Text style={errorStyles.title}>Error de inicialización</Text>
        <Text style={errorStyles.message}>{dbError}</Text>
      </View>
    );
  }

  if (!loaded || !dbReady) {
    return null;
  }

  return <RootLayoutNav />;
}

function RootLayoutNav() {
  const { isAuthenticated, isLoading } = useAuthStore();

  if (isLoading) return null;

  return (
    <>
      {!isAuthenticated && <Redirect href="/login" />}
      {isAuthenticated && <Redirect href="/(tabs)" />}
      <Stack>
        <Stack.Screen name="login" options={{ headerShown: false }} />
        <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
        <Stack.Screen name="visit/[id]" options={{ headerShown: false, presentation: 'card' }} />
        <Stack.Screen name="visit/[id]/reopen" options={{ headerShown: false, presentation: 'modal' }} />
        <Stack.Screen name="visit/[id]/measurements" options={{ headerShown: false, presentation: 'card' }} />
        <Stack.Screen name="visit/[id]/usage" options={{ headerShown: false, presentation: 'card' }} />
        <Stack.Screen name="visit/[id]/audit" options={{ headerShown: false, presentation: 'card' }} />
        <Stack.Screen name="checklist/[id]" options={{ headerShown: false, presentation: 'card' }} />
        <Stack.Screen name="photo/[id]" options={{ headerShown: false, presentation: 'card' }} />
        <Stack.Screen name="gps/[id]" options={{ headerShown: false, presentation: 'card' }} />
      </Stack>
    </>
  );
}

const errorStyles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 24,
    backgroundColor: '#fff',
  },
  title: {
    fontFamily: 'monospace',
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 12,
  },
  message: {
    fontFamily: 'monospace',
    fontSize: 12,
    color: '#666',
    textAlign: 'center',
  },
});
