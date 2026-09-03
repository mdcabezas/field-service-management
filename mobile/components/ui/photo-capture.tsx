import React, { useState, useEffect, useCallback } from 'react';
import { View, Text, Image, StyleSheet, Alert } from 'react-native';
import { Button } from './button';
import { captureAndCompressPhoto } from '../../lib/photo-upload';

interface PhotoCaptureProps {
  visitId: string;
  onPhotoCaptured?: (photoId: string) => void;
  maxPhotos?: number;
  currentPhotoCount?: number;
}

export function PhotoCapture({
  visitId,
  onPhotoCaptured,
  maxPhotos = 12,
  currentPhotoCount = 0,
}: PhotoCaptureProps) {
  const [isCapturing, setIsCapturing] = useState(false);
  const [lastPhoto, setLastPhoto] = useState<string | null>(null);
  const [ImagePicker, setImagePicker] = useState<any>(null);

  useEffect(() => {
    import('expo-image-picker').then(setImagePicker);
  }, []);

  const takePhoto = useCallback(async () => {
    if (!ImagePicker) return;
    if (currentPhotoCount >= maxPhotos) {
      Alert.alert('Límite alcanzado', `Máximo ${maxPhotos} fotos por visita`);
      return;
    }

    const permission = await ImagePicker.requestCameraPermissionsAsync();
    if (!permission.granted) {
      Alert.alert('Permiso requerido', 'Se necesita acceso a la cámara');
      return;
    }

    setIsCapturing(true);
    try {
      const result = await ImagePicker.launchCameraAsync({
        mediaTypes: ['images'],
        quality: 0.7,
        allowsEditing: false,
      });

      if (!result.canceled && result.assets[0]) {
        const photo = await captureAndCompressPhoto(
          result.assets[0].uri,
          visitId
        );
        setLastPhoto(photo.localPath);
        onPhotoCaptured?.(photo.id);
      }
    } catch (error) {
      Alert.alert('Error', 'No se pudo capturar la foto');
    } finally {
      setIsCapturing(false);
    }
  }, [ImagePicker, currentPhotoCount, maxPhotos, visitId, onPhotoCaptured]);

  const pickFromGallery = useCallback(async () => {
    if (!ImagePicker) return;
    if (currentPhotoCount >= maxPhotos) {
      Alert.alert('Límite alcanzado', `Máximo ${maxPhotos} fotos por visita`);
      return;
    }

    const permission = await ImagePicker.requestMediaLibraryPermissionsAsync();
    if (!permission.granted) {
      Alert.alert('Permiso requerido', 'Se necesita acceso a la galería');
      return;
    }

    setIsCapturing(true);
    try {
      const result = await ImagePicker.launchImageLibraryAsync({
        mediaTypes: ['images'],
        quality: 0.7,
        allowsMultipleSelection: false,
      });

      if (!result.canceled && result.assets[0]) {
        const photo = await captureAndCompressPhoto(
          result.assets[0].uri,
          visitId
        );
        setLastPhoto(photo.localPath);
        onPhotoCaptured?.(photo.id);
      }
    } catch (error) {
      Alert.alert('Error', 'No se pudo cargar la foto');
    } finally {
      setIsCapturing(false);
    }
  }, [ImagePicker, currentPhotoCount, maxPhotos, visitId, onPhotoCaptured]);

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Fotos</Text>
        <Text style={styles.count}>
          {currentPhotoCount}/{maxPhotos}
        </Text>
      </View>

      <View style={styles.actions}>
        <Button
          title={isCapturing ? 'Capturando...' : 'Tomar foto'}
          onPress={takePhoto}
          disabled={isCapturing || currentPhotoCount >= maxPhotos || !ImagePicker}
          style={styles.button}
        />
        <Button
          title="Galería"
          onPress={pickFromGallery}
          disabled={isCapturing || currentPhotoCount >= maxPhotos || !ImagePicker}
          variant="secondary"
          style={styles.button}
        />
      </View>

      {lastPhoto && (
        <View style={styles.preview}>
          <Text style={styles.previewLabel}>Última foto:</Text>
          <Image source={{ uri: lastPhoto }} style={styles.previewImage} />
        </View>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    padding: 16,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  title: {
    fontFamily: 'monospace',
    fontSize: 14,
    fontWeight: 'bold',
    textTransform: 'uppercase',
  },
  count: {
    fontFamily: 'monospace',
    fontSize: 12,
    color: '#666',
  },
  actions: {
    flexDirection: 'row',
    gap: 8,
  },
  button: {
    flex: 1,
  },
  preview: {
    marginTop: 16,
  },
  previewLabel: {
    fontFamily: 'monospace',
    fontSize: 11,
    textTransform: 'uppercase',
    marginBottom: 8,
  },
  previewImage: {
    width: '100%',
    height: 200,
    resizeMode: 'cover',
    borderWidth: 2,
    borderColor: '#000',
  },
});