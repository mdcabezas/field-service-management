import * as Location from 'expo-location';
import * as Crypto from 'expo-crypto';
import { getDatabase } from './database';
import { addToOutbox } from './outbox';

export const MAX_ACCURACY = 50;
const DEFAULT_THRESHOLD = 500;
const EARTH_RADIUS = 6371000;

export interface GPSResult {
  lat: number;
  lng: number;
  accuracy: number;
}

export interface DistanceResult {
  distance: number;
  isWithinThreshold: boolean;
}

export function isAccuracyAcceptable(gps: GPSResult, maxAccuracy: number = MAX_ACCURACY): boolean {
  return gps.accuracy === 0 || gps.accuracy <= maxAccuracy;
}

export async function captureGPS(timeoutMs: number = 15000): Promise<GPSResult | null> {
  try {
    const { status } = await Location.requestForegroundPermissionsAsync();
    if (status !== 'granted') {
      return null;
    }

    const location = await Location.getCurrentPositionAsync({
      accuracy: Location.Accuracy.Balanced,
    });

    return {
      lat: location.coords.latitude,
      lng: location.coords.longitude,
      accuracy: location.coords.accuracy || 0,
    };
  } catch (error) {
    console.error('GPS capture failed:', error);
    return null;
  }
}

export async function captureGPSWithAccuracy(
  onProgress: (accuracy: number | null) => void,
  timeoutMs: number = 30000,
  maxAccuracy: number = MAX_ACCURACY
): Promise<GPSResult | null> {
  const { status } = await Location.requestForegroundPermissionsAsync();
  if (status !== 'granted') {
    return null;
  }

  return new Promise((resolve) => {
    let resolved = false;

    let subscription: Location.LocationSubscription | null = null;

    Location.watchPositionAsync(
      {
        accuracy: Location.Accuracy.High,
        timeInterval: 1000,
        distanceInterval: 0,
      },
      (location) => {
        if (resolved) return;

        const accuracy = location.coords.accuracy || 0;
        onProgress(accuracy);

        if (accuracy === 0 || accuracy <= maxAccuracy) {
          resolved = true;
          if (subscription) {
            subscription.remove();
          }
          resolve({
            lat: location.coords.latitude,
            lng: location.coords.longitude,
            accuracy,
          });
        }
      }
    ).then((sub) => {
      subscription = sub;
    });

    setTimeout(() => {
      if (resolved) return;
      resolved = true;
      if (subscription) {
        subscription.remove();
      }

      Location.getCurrentPositionAsync({
        accuracy: Location.Accuracy.Balanced,
      })
        .then((loc) => {
          resolve({
            lat: loc.coords.latitude,
            lng: loc.coords.longitude,
            accuracy: loc.coords.accuracy || 0,
          });
        })
        .catch(() => resolve(null));
    }, timeoutMs);
  });
}

export function calculateDistance(lat1: number, lng1: number, lat2: number, lng2: number): number {
  const dLat = toRad(lat2 - lat1);
  const dLng = toRad(lng2 - lng1);
  const a =
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos(toRad(lat1)) * Math.cos(toRad(lat2)) * Math.sin(dLng / 2) * Math.sin(dLng / 2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  return EARTH_RADIUS * c;
}

function toRad(deg: number): number {
  return (deg * Math.PI) / 180;
}

export function checkDistance(
  techLat: number,
  techLng: number,
  propertyLat: number,
  propertyLng: number,
  threshold: number = DEFAULT_THRESHOLD
): DistanceResult {
  const distance = calculateDistance(techLat, techLng, propertyLat, propertyLng);
  return {
    distance,
    isWithinThreshold: distance <= threshold,
  };
}

export async function createCheckpoint(
  visitId: string,
  type: string,
  gps: GPSResult,
  propertyLat?: number,
  propertyLng?: number
): Promise<string> {
  const db = getDatabase();
  const checkpointId = Crypto.randomUUID();

  let distanceFromProperty: number | null = null;
  if (propertyLat && propertyLng) {
    distanceFromProperty = calculateDistance(gps.lat, gps.lng, propertyLat, propertyLng);
  }

  await db.runAsync(
    'INSERT INTO checkpoints (id, visit_id, type, timestamp, lat, lng, accuracy, distance_from_property, synced) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0)',
    checkpointId, visitId, type, new Date().toISOString(), gps.lat, gps.lng, gps.accuracy, distanceFromProperty
  );

  await addToOutbox('checkpoint', checkpointId, 'create', {
    visit_id: visitId,
    type,
    timestamp: new Date().toISOString(),
    lat: gps.lat,
    lng: gps.lng,
    accuracy: gps.accuracy,
    local_id: checkpointId,
  });

  return checkpointId;
}

export async function createGPSWaiver(visitId: string, type: string): Promise<string> {
  const db = getDatabase();
  const waiverId = Crypto.randomUUID();

  await db.runAsync(
    'INSERT INTO gps_waivers (id, visit_id, type, timestamp, accepted_by, synced) VALUES (?, ?, ?, ?, ?, 0)',
    waiverId, visitId, type, new Date().toISOString(), 'Technician'
  );

  await addToOutbox('gps_waiver', waiverId, 'create', {
    visit_id: visitId,
    type,
    timestamp: new Date().toISOString(),
    local_id: waiverId,
  });

  return waiverId;
}