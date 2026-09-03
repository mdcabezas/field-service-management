import React from 'react';
import { View, Text, StyleSheet } from 'react-native';

interface BadgeProps {
  label: string;
  color?: string;
}

export const Badge = React.memo(function Badge({ label, color = '#000' }: BadgeProps) {
  return (
    <View style={[styles.badge, { backgroundColor: color }]}>
      <Text style={styles.text}>{label}</Text>
    </View>
  );
});

const styles = StyleSheet.create({
  badge: {
    borderWidth: 2,
    borderColor: '#000',
    borderRadius: 0,
    paddingHorizontal: 8,
    paddingVertical: 4,
    alignSelf: 'flex-start',
  },
  text: {
    fontFamily: 'monospace',
    fontSize: 10,
    fontWeight: 'bold',
    textTransform: 'uppercase',
    color: '#fff',
  },
});
