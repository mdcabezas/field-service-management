import React from 'react';
import { View, Text, StyleSheet, ViewStyle } from 'react-native';

interface HeaderProps {
  title: string;
  right?: React.ReactNode;
  style?: ViewStyle;
}

export const Header = React.memo(function Header({ title, right, style }: HeaderProps) {
  return (
    <View style={[styles.header, style]}>
      <Text style={styles.title}>{title}</Text>
      {right}
    </View>
  );
});

const styles = StyleSheet.create({
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    borderBottomWidth: 2,
    borderBottomColor: '#000',
    paddingVertical: 12,
    paddingHorizontal: 16,
  },
  title: {
    fontFamily: 'monospace',
    fontSize: 18,
    fontWeight: 'bold',
    textTransform: 'uppercase',
  },
});
