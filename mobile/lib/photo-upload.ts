import { Platform } from 'react-native';
import * as ImageManipulator from 'expo-image-manipulator';
import { Paths, Directory, File } from 'expo-file-system';
import * as Crypto from 'expo-crypto';
import { getDatabase } from './database';
import { apiRequest } from './auth';
import { addToOutbox } from './outbox';

interface PhotoResult {
  id: string;
  localPath: string;
  thumbnailPath: string;
}

async function ensureAccessibleUri(uri: string): Promise<string> {
  if (uri.startsWith('content://')) {
    console.log('[PhotoUpload] URI is content://, accessible directly');
    return uri;
  }
  if (uri.startsWith('file://')) {
    console.log('[PhotoUpload] URI is file://, copying to Documents for accessibility');
    const fileName = `camera_${Date.now()}_${Crypto.randomUUID().slice(0, 8)}.jpg`;
    const tempDir = new Directory(Paths.document.uri, 'temp_camera');
    if (!tempDir.exists) {
      tempDir.create();
      console.log('[PhotoUpload] Created temp_camera directory');
    }
    const destFile = new File(tempDir.uri, fileName);
    await new File(uri).copy(destFile);
    console.log('[PhotoUpload] Copied file:// to Documents:', destFile.uri);
    return destFile.uri;
  }
  console.warn('[PhotoUpload] Unknown URI scheme:', uri.slice(0, 20));
  return uri;
}

export async function captureAndCompressPhoto(
  uri: string,
  visitId: string,
  lat?: number,
  lng?: number
): Promise<PhotoResult> {
  console.log('[PhotoUpload] Starting captureAndCompressPhoto', { uri: uri.slice(0, 50), visitId, lat, lng });

  console.log('[PhotoUpload] Step 0: Ensuring URI accessibility');
  const accessibleUri = await ensureAccessibleUri(uri);
  console.log('[PhotoUpload] Step 0 OK: accessibleUri =', accessibleUri.slice(0, Math.min(50, accessibleUri.length)));

  console.log('[PhotoUpload] Step 1: Manipulating image (resize 640px)');
  let manipulated: ImageManipulator.ImageResult;
  try {
    manipulated = await ImageManipulator.manipulateAsync(
      accessibleUri,
      [{ resize: { width: 640 } }],
      { compress: 0.6, format: ImageManipulator.SaveFormat.JPEG }
    );
    console.log('[PhotoUpload] Step 1 OK: manipulated URI =', manipulated.uri);
  } catch (e) {
    console.error('[PhotoUpload] Step 1 FAILED: ImageManipulator.manipulateAsync (resize)', e);
    throw e;
  }

  console.log('[PhotoUpload] Step 2: Creating directories');
  const photoId = Crypto.randomUUID();
  const photosDir = new Directory(Paths.document.uri, 'photos');
  const thumbsDir = new Directory(Paths.document.uri, 'thumbs');

  try {
    if (!photosDir.exists) {
      console.log('[PhotoUpload] Creating photos directory');
      photosDir.create();
    }
    if (!thumbsDir.exists) {
      console.log('[PhotoUpload] Creating thumbs directory');
      thumbsDir.create();
    }
    console.log('[PhotoUpload] Step 2 OK: directories ready');
  } catch (e) {
    console.error('[PhotoUpload] Step 2 FAILED: Directory creation', e);
    throw e;
  }

  const localPath = `${photosDir.uri}/${photoId}.jpg`;
  const thumbnailPath = `${thumbsDir.uri}/${photoId}_thumb.jpg`;
  console.log('[PhotoUpload] Paths:', { localPath, thumbnailPath });

  console.log('[PhotoUpload] Step 3: Copying manipulated file to localPath');
  try {
    const sourceFile = new File(manipulated.uri);
    const destFile = new File(localPath);
    await sourceFile.copy(destFile);
    console.log('[PhotoUpload] Step 3 OK: file copied to', localPath);
  } catch (e) {
    console.error('[PhotoUpload] Step 3 FAILED: File copy', e);
    throw e;
  }

  console.log('[PhotoUpload] Step 4: Creating thumbnail (150px)');
  let thumbnail: ImageManipulator.ImageResult;
  try {
    thumbnail = await ImageManipulator.manipulateAsync(
      manipulated.uri,
      [{ resize: { width: 150 } }],
      { compress: 0.5, format: ImageManipulator.SaveFormat.JPEG }
    );
    console.log('[PhotoUpload] Step 4 OK: thumbnail URI =', thumbnail.uri);
  } catch (e) {
    console.error('[PhotoUpload] Step 4 FAILED: ImageManipulator.manipulateAsync (thumbnail)', e);
    throw e;
  }

  console.log('[PhotoUpload] Step 5: Copying thumbnail');
  try {
    const thumbSource = new File(thumbnail.uri);
    const thumbDest = new File(thumbnailPath);
    await thumbSource.copy(thumbDest);
    console.log('[PhotoUpload] Step 5 OK: thumbnail copied to', thumbnailPath);
  } catch (e) {
    console.error('[PhotoUpload] Step 5 FAILED: Thumbnail copy', e);
    throw e;
  }

  console.log('[PhotoUpload] Step 6: Inserting into database');
  try {
    const db = getDatabase();
    const photoFile = new File(localPath);
    const fileSize = photoFile.exists ? photoFile.size : 0;
    await db.runAsync(
      'INSERT INTO photos (id, visit_id, local_path, thumbnail_path, uploaded, file_size, lat, lng, timestamp) VALUES (?, ?, ?, ?, 0, ?, ?, ?, ?)',
      photoId, visitId, localPath, thumbnailPath, fileSize, lat || null, lng || null, new Date().toISOString()
    );
    console.log('[PhotoUpload] Step 6 OK: DB insert successful');
  } catch (e) {
    console.error('[PhotoUpload] Step 6 FAILED: DB insert', e);
    throw e;
  }

  console.log('[PhotoUpload] Step 7: Adding to outbox');
  try {
    await addToOutbox('photo', photoId, 'create', {
      visit_id: visitId,
      local_id: photoId,
      timestamp: new Date().toISOString(),
      lat,
      lng,
    });
    console.log('[PhotoUpload] Step 7 OK: Outbox entry added');
  } catch (e) {
    console.error('[PhotoUpload] Step 7 FAILED: addToOutbox', e);
    throw e;
  }

  console.log('[PhotoUpload] All steps completed successfully', { photoId, localPath, thumbnailPath });
  return { id: photoId, localPath, thumbnailPath };
}

export async function uploadPhoto(photoId: string): Promise<boolean> {
  if (Platform.OS === 'web') {
    console.warn('[PhotoUpload] uploadPhoto not supported on web');
    return false;
  }
  console.log('[PhotoUpload] uploadPhoto called', { photoId });
  const db = getDatabase();
  const stmt = db.prepareSync('SELECT * FROM photos WHERE id = ?');
  const result = stmt.executeSync(photoId);
  const photo = result.getFirstSync() as Record<string, unknown> | null;
  stmt.finalizeSync();
  if (!photo) {
    console.error('[PhotoUpload] uploadPhoto: photo not found', { photoId });
    return false;
  }

  try {
    console.log('[PhotoUpload] uploadPhoto: preparing FormData');
    const formData = new FormData();
    formData.append('visit_id', photo.visit_id as string);
    const photoFile = new File(photo.local_path as string);
    formData.append('photo', photoFile, `${photoId}.jpg`);
    formData.append('timestamp', photo.timestamp as string);
    if (photo.lat) formData.append('lat', String(photo.lat));
    if (photo.lng) formData.append('lng', String(photo.lng));

    console.log('[PhotoUpload] uploadPhoto: calling apiRequest');
    const response = await apiRequest<any>('/api/mobile/photos', {
      method: 'POST',
      body: formData,
    });

    await db.runAsync('UPDATE photos SET uploaded = 1, server_id = ? WHERE id = ?', response.id, photoId);
    console.log('[PhotoUpload] uploadPhoto: success', { serverId: response.id });
    return true;
  } catch (error) {
    console.error('[PhotoUpload] uploadPhoto FAILED:', error);
    return false;
  }
}

export async function getPhotoStorageSize(): Promise<number> {
  const db = getDatabase();
  const stmt = db.prepareSync('SELECT COALESCE(SUM(file_size), 0) as total FROM photos');
  const result = stmt.executeSync();
  const row = result.getFirstSync() as { total: number } | null;
  stmt.finalizeSync();
  return row?.total || 0;
}

export async function cleanupOldPhotos(retentionDays: number, maxStorageMB: number): Promise<void> {
  const db = getDatabase();
  const cutoffDate = new Date();
  cutoffDate.setDate(cutoffDate.getDate() - retentionDays);

  await db.runAsync('DELETE FROM photos WHERE created_at < ? AND uploaded = 1', cutoffDate.toISOString());

  const storageSize = await getPhotoStorageSize();
  const maxBytes = maxStorageMB * 1024 * 1024;

  if (storageSize > maxBytes) {
    const stmt = db.prepareSync(
      'SELECT id, local_path, thumbnail_path FROM photos WHERE uploaded = 1 ORDER BY created_at ASC LIMIT 10'
    );
    const result = stmt.executeSync();
    const oldestPhotos = result.getAllSync() as Record<string, unknown>[];
    stmt.finalizeSync();

    for (const photo of oldestPhotos) {
      try {
        const localFile = new File(photo.local_path as string);
        if (localFile.exists) localFile.delete();
        const thumbFile = new File(photo.thumbnail_path as string);
        if (thumbFile.exists) thumbFile.delete();
      } catch {}
      await db.runAsync('DELETE FROM photos WHERE id = ?', photo.id as string);
    }
  }
}

export async function sharePhoto(photoId: string): Promise<void> {
  const db = getDatabase();
  const stmt = db.prepareSync('SELECT local_path FROM photos WHERE id = ?');
  const result = stmt.executeSync(photoId);
  const photo = result.getFirstSync() as { local_path: string } | null;
  stmt.finalizeSync();
  if (!photo) return;

  const Sharing = require('expo-sharing');
  if (await Sharing.isAvailableAsync()) {
    await Sharing.shareAsync(photo.local_path as string);
  }
}