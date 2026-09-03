import React, { useEffect, useRef, useState } from 'react';
import { View, Text, StyleSheet, Alert, TouchableOpacity, ActivityIndicator } from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { CameraView as Camera, useCameraPermissions } from 'expo-camera';
import * as ImagePicker from 'expo-image-picker';
import { getDatabase } from '../../lib/database';
import { captureAndCompressPhoto } from '../../lib/photo-upload';
import { Header } from '../../components/ui/header';
import { Button } from '../../components/ui/button';

export default function CameraScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const [cameraType, setCameraType] = useState<'front' | 'back'>('back');
  const [isCapturing, setIsCapturing] = useState(false);
  const [hasPermission, setHasPermission] = useState(false);
  const cameraRef = useRef<any>(null);
  const [cameraPermission, requestCameraPermission] = useCameraPermissions();

  useEffect(() => {
    checkPermissions();
  }, []);

  const checkPermissions = async () => {
    const { status } = await requestCameraPermission();
    if (status === 'granted') {
      setHasPermission(true);
    }
  };

  const handleRequestCameraPermission = async () => {
    const { status } = await requestCameraPermission();
    if (status === 'granted') {
      // Force re-render
    }
  };

  const handleClose = () => {
    router.back();
  };

  if (!cameraPermission?.granted) {
    return (
      <View style={styles.container}>
        <Header title="Cámara" right={<TouchableOpacity onPress={handleClose}><Text style={styles.closeButtonText}>✕</Text></TouchableOpacity>} />
        <View style={styles.permissionContainer}>
          <Text style={styles.permissionText}>Permiso de cámara requerido</Text>
          <Button title="Solicitar permiso" onPress={handleRequestCameraPermission} variant="primary" />
          <Button title="Cancelar" onPress={handleClose} variant="secondary" />
        </View>
      </View>
    );
  }

  if (!hasPermission) {
    return (
      <View style={styles.container}>
        <Header title="Cámara" right={<TouchableOpacity onPress={handleClose}><Text style={styles.closeButtonText}>✕</Text></TouchableOpacity>} />
        <View style={styles.permissionContainer}>
          <ActivityIndicator size="large" color="#fff" />
          <Text style={styles.permissionText}>Inicializando cámara...</Text>
        </View>
      </View>
    );
  }

  const handleCapture = async () => {
    if (!cameraRef.current || !hasPermission) return;

    setIsCapturing(true);
    try {
      console.log('[CameraScreen] Taking picture with expo-camera');
      const photo = await cameraRef.current?.takePictureAsync({
        quality: 0.7,
        base64: false,
        exif: false,
      });
      console.log('[CameraScreen] Camera capture result:', { uri: photo?.uri?.slice(0, 50) });

      if (photo?.uri) {
        console.log('[CameraScreen] Calling captureAndCompressPhoto');
        await captureAndCompressPhoto(photo.uri, id!);
        console.log('[CameraScreen] captureAndCompressPhoto returned');
        Alert.alert('Éxito', 'Foto capturada correctamente', [
          {
            text: 'OK',
            onPress: () => router.back(),
          },
        ]);
      }
    } catch (error) {
      console.error('[CameraScreen] Capture ERROR:', error);
      Alert.alert('Error', 'Error al capturar foto');
    } finally {
      setIsCapturing(false);
    }
  };

  const handleFlipCamera = () => {
    setCameraType((prev) => (prev === 'back' ? 'front' : 'back'));
  };

  const handleGallery = async () => {
    const { status } = await ImagePicker.requestMediaLibraryPermissionsAsync();
    if (status !== 'granted') {
      Alert.alert('Permiso requerido', 'Se necesita acceso a la galería');
      return;
    }

    const result = await ImagePicker.launchImageLibraryAsync({
      mediaTypes: ['images'],
      quality: 0.7,
    });

    if (!result.canceled && result.assets[0]) {
      await captureAndCompressPhoto(result.assets[0].uri, id!);
      router.back();
    }
  };

  return (
    <View style={styles.container}>
      <Header
        title="Tomar foto"
        right={
          <TouchableOpacity onPress={handleClose} style={styles.closeButton}>
            <Text style={styles.closeButtonText}>✕</Text>
          </TouchableOpacity>
        }
      />

      <Camera
        ref={cameraRef}
        style={styles.cameraPreview}
        facing={cameraType}
      />

      <View style={styles.overlay}>
        <View style={styles.topOverlay}>
          <TouchableOpacity
            style={styles.controlButton}
            onPress={handleFlipCamera}
            disabled={isCapturing}
          >
            <Text style={styles.controlButtonText}>🔄</Text>
          </TouchableOpacity>
          <Text style={styles.cameraTitle}>Tomar foto</Text>
          <TouchableOpacity
            style={styles.controlButton}
            onPress={handleClose}
            disabled={isCapturing}
          >
            <Text style={styles.closeButtonText}>✕</Text>
          </TouchableOpacity>
        </View>

        <View style={styles.bottomOverlay}>
          <TouchableOpacity
            style={[styles.galleryButton, isCapturing && styles.buttonDisabled]}
            onPress={handleGallery}
            disabled={isCapturing}
          >
            <Text style={styles.galleryButtonText}>🖼 Galería</Text>
          </TouchableOpacity>

          <TouchableOpacity
            style={[styles.captureButton, isCapturing && styles.captureButtonDisabled]}
            onPress={handleCapture}
            disabled={isCapturing}
            activeOpacity={0.8}
          >
            <View style={[styles.captureInner, isCapturing && styles.captureInnerCapturing]} />
            {isCapturing && <ActivityIndicator size="small" color="#fff" style={styles.captureSpinner} />}
          </TouchableOpacity>
        </View>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#000',
  },
  cameraPreview: {
    flex: 1,
  },
  overlay: {
    flex: 1,
    backgroundColor: 'transparent',
  },
  topOverlay: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 16,
    paddingTop: 40,
  },
  bottomOverlay: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    alignItems: 'center',
    paddingHorizontal: 20,
    paddingBottom: 40,
  },
  cameraTitle: {
    color: '#fff',
    fontFamily: 'monospace',
    fontSize: 16,
    fontWeight: 'bold',
  },
  closeButton: {
    padding: 8,
  },
  closeButtonText: {
    color: '#fff',
    fontSize: 24,
  },
  controlButton: {
    width: 56,
    height: 56,
    borderRadius: 28,
    backgroundColor: 'rgba(255,255,255,0.2)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  controlButtonText: {
    fontSize: 24,
  },
  captureButton: {
    width: 80,
    height: 80,
    borderRadius: 40,
    borderWidth: 4,
    borderColor: '#fff',
    backgroundColor: 'rgba(255,255,255,0.2)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  captureButtonDisabled: {
    opacity: 0.5,
  },
  captureInner: {
    width: 60,
    height: 60,
    borderRadius: 30,
    backgroundColor: '#fff',
  },
  captureInnerCapturing: {
    backgroundColor: '#f59e0b',
  },
  captureSpinner: {
    position: 'absolute',
  },
  galleryButton: {
    paddingHorizontal: 24,
    paddingVertical: 12,
    borderRadius: 24,
    backgroundColor: 'rgba(255,255,255,0.2)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  galleryButtonText: {
    color: '#fff',
    fontFamily: 'monospace',
    fontSize: 14,
  },
  buttonDisabled: {
    opacity: 0.5,
  },
  permissionContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 24,
  },
  permissionText: {
    color: '#fff',
    fontFamily: 'monospace',
    fontSize: 16,
    marginTop: 16,
    textAlign: 'center',
  },
});