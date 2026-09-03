import React from 'react';
import { View, Text, StyleSheet } from 'react-native';

interface TimestampOverlayProps {
  timestamp: string;
  lat?: number | null;
  lng?: number | null;
}

export function TimestampOverlay({ timestamp, lat, lng }: TimestampOverlayProps) {
  const formatTimestamp = (ts: string) => {
    try {
      const date = new Date(ts);
      return date.toLocaleString('es-MX', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
      });
    } catch {
      return ts;
    }
  };

  const formatCoordinates = (latitude: number, longitude: number) => {
    return `${latitude.toFixed(6)}, ${longitude.toFixed(6)}`;
  };

  return (
    <View style={styles.container}>
      <Text style={styles.timestamp}>{formatTimestamp(timestamp)}</Text>
      {lat && lng && (
        <Text style={styles.coordinates}>{formatCoordinates(lat, lng)}</Text>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    position: 'absolute',
    bottom: 8,
    left: 8,
    backgroundColor: 'rgba(0, 0, 0, 0.7)',
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 4,
  },
  timestamp: {
    fontFamily: 'monospace',
    fontSize: 10,
    color: '#fff',
  },
  coordinates: {
    fontFamily: 'monospace',
    fontSize: 8,
    color: '#ccc',
    marginTop: 2,
  },
});