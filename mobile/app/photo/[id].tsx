import React, { useEffect, useState, useCallback } from 'react';
import { View, Text, FlatList, Image, Pressable, StyleSheet, Alert } from 'react-native';
import { useLocalSearchParams, useRouter, useFocusEffect } from 'expo-router';
import * as ImagePicker from 'expo-image-picker';
import { getDatabase } from '../../lib/database';
import { captureAndCompressPhoto, sharePhoto } from '../../lib/photo-upload';
import { Header } from '../../components/ui/header';
import { Card } from '../../components/ui/card';
import { Badge } from '../../components/ui/badge';
import { Button } from '../../components/ui/button';
import { TimestampOverlay } from '../../components/ui/timestamp-overlay';

interface Photo {
  id: string;
  local_path: string;
  thumbnail_path: string;
  uploaded: boolean;
  lat: number | null;
  lng: number | null;
  created_at: string;
}

export default function PhotoScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const [photos, setPhotos] = useState<Photo[]>([]);

  useFocusEffect(
    useCallback(() => {
      console.log('[PhotoScreen] useFocusEffect fired, reloading photos');
      loadPhotos();
    }, [])
  );

  const loadPhotos = async () => {
    console.log('[PhotoScreen] loadPhotos called');
    const db = getDatabase();
    const stmt = db.prepareSync(
      'SELECT * FROM photos WHERE visit_id = ? ORDER BY created_at DESC'
    );
    const result = stmt.executeSync(id);
    const rows = result.getAllSync() as Record<string, unknown>[];
    stmt.finalizeSync();
    console.log('[PhotoScreen] loadPhotos found', rows.length, 'photos');
    setPhotos(
      rows.map((row) => ({
        id: row.id as string,
        local_path: row.local_path as string,
        thumbnail_path: row.thumbnail_path as string,
        uploaded: row.uploaded === 1,
        lat: row.lat as number | null,
        lng: row.lng as number | null,
        created_at: row.created_at as string,
      }))
    );
  };

  const handleCapture = async () => {
    router.push({ pathname: '/camera/[id]', params: { id } } as any);
  };

  const handlePickImage = async () => {
    console.log('[PhotoScreen] handlePickImage START');
    const { status } = await ImagePicker.requestMediaLibraryPermissionsAsync();
    if (status !== 'granted') {
      Alert.alert('Permiso requerido', 'Se necesita acceso a la galería');
      return;
    }

    console.log('[PhotoScreen] Launching image library');
    const result = await ImagePicker.launchImageLibraryAsync({
      mediaTypes: ['images'],
      quality: 0.7,
    });
    console.log('[PhotoScreen] Gallery result:', { canceled: result.canceled, assetsCount: result.assets?.length, uri: result.assets?.[0]?.uri?.slice(0, 50) });

    if (!result.canceled && result.assets[0]) {
      console.log('[PhotoScreen] Calling captureAndCompressPhoto (from gallery)');
      await captureAndCompressPhoto(result.assets[0].uri, id!);
      loadPhotos();
    }
  };

  const handleShare = async (photoId: string) => {
    await sharePhoto(photoId);
  };

  const handleLongPress = (photoId: string) => {
    Alert.alert('Compartir foto', '¿Deseas compartir esta foto?', [
      { text: 'Cancelar', style: 'cancel' },
      {
        text: 'Compartir',
        onPress: () => handleShare(photoId),
      },
    ]);
  };

  const handleDelete = (photoId: string) => {
    Alert.alert('Eliminar foto', '¿Estás seguro?', [
      { text: 'Cancelar', style: 'cancel' },
      {
        text: 'Eliminar',
        style: 'destructive',
        onPress: async () => {
          const db = getDatabase();
          await db.runAsync('DELETE FROM photos WHERE id = ?', photoId);
          loadPhotos();
        },
      },
    ]);
  };

  const renderPhoto = ({ item }: { item: Photo }) => (
    <Card style={styles.photoCard}>
      <Pressable onLongPress={() => handleLongPress(item.id)}>
        <View style={styles.imageContainer}>
          <Image source={{ uri: `file://${item.thumbnail_path}` }} style={styles.thumbnail} />
          <TimestampOverlay timestamp={item.created_at} lat={item.lat} lng={item.lng} />
        </View>
      </Pressable>
      <View style={styles.photoInfo}>
        <Text style={styles.photoDate}>
          {new Date(item.created_at).toLocaleString()}
        </Text>
        <Badge
          label={item.uploaded ? 'Subido' : 'Pendiente'}
          color={item.uploaded ? '#059669' : '#f59e0b'}
        />
      </View>
      <View style={styles.photoActions}>
        <Button
          title="Compartir"
          onPress={() => handleShare(item.id)}
          variant="secondary"
          style={styles.photoAction}
        />
        <Button
          title="Eliminar"
          onPress={() => handleDelete(item.id)}
          variant="danger"
          style={styles.photoAction}
        />
      </View>
    </Card>
  );

  return (
    <View style={styles.container}>
      <Header
        title="Fotos"
        right={<Badge label={`${photos.length}`} color="#000" />}
      />

      <View style={styles.actions}>
        <Button
          title="Tomar foto"
          onPress={handleCapture}
          style={styles.actionButton}
        />
        <Button
          title="Seleccionar de galería"
          onPress={handlePickImage}
          variant="secondary"
          style={styles.actionButton}
        />
      </View>

      <FlatList
        data={photos}
        renderItem={renderPhoto}
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
  list: {
    padding: 16,
  },
  photoCard: {
    marginBottom: 12,
  },
  imageContainer: {
    position: 'relative',
    marginBottom: 8,
  },
  thumbnail: {
    width: '100%',
    height: 200,
    resizeMode: 'cover',
  },
  photoInfo: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  photoDate: {
    fontFamily: 'monospace',
    fontSize: 10,
    color: '#666',
  },
  photoActions: {
    flexDirection: 'row',
    gap: 8,
  },
  photoAction: {
    flex: 1,
  },
});