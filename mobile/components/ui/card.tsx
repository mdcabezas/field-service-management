import React, { ReactNode, isValidElement } from 'react';
import { View, Text, StyleSheet, ViewStyle } from 'react-native';

interface CardProps {
  title?: string;
  children: ReactNode;
  style?: ViewStyle;
}

function ensureTextChild(child: ReactNode): ReactNode {
  if (typeof child === 'string' || typeof child === 'number') {
    return <Text style={styles.autoText}>{String(child)}</Text>;
  }
  if (Array.isArray(child)) {
    return child.map(ensureTextChild);
  }
  if (isValidElement(child)) {
    return child;
  }
  return null;
}

export const Card = React.memo(function Card({ title, children, style }: CardProps) {
  return (
    <View style={[styles.card, style]}>
      {title && <Text style={styles.title}>{title}</Text>}
      {React.Children.map(children, ensureTextChild)}
    </View>
  );
});

const styles = StyleSheet.create({
  card: {
    borderWidth: 2,
    borderColor: '#000',
    borderRadius: 0,
    padding: 16,
    backgroundColor: '#fff',
  },
  title: {
    fontFamily: 'monospace',
    fontSize: 14,
    fontWeight: 'bold',
    textTransform: 'uppercase',
    marginBottom: 12,
  },
  autoText: {
    fontFamily: 'monospace',
    fontSize: 14,
  },
});