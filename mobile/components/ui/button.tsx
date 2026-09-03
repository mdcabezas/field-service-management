import React from 'react';
import { Pressable, Text, StyleSheet, ActivityIndicator, ViewStyle } from 'react-native';

interface ButtonProps {
  title: string;
  onPress: () => void;
  variant?: 'primary' | 'secondary' | 'danger';
  disabled?: boolean;
  loading?: boolean;
  style?: ViewStyle;
}

export const Button = React.memo(function Button({ title, onPress, variant = 'primary', disabled, loading, style }: ButtonProps) {
  const bgColor = variant === 'primary' ? '#000' : variant === 'danger' ? '#dc2626' : '#fff';
  const textColor = variant === 'primary' ? '#fff' : variant === 'danger' ? '#fff' : '#000';
  const borderColor = variant === 'secondary' ? '#000' : bgColor;

  return (
    <Pressable
      onPress={onPress}
      disabled={disabled || loading}
      hitSlop={{ top: 20, bottom: 20, left: 24, right: 24 }}
      style={[
        styles.button,
        { backgroundColor: bgColor, borderColor },
        disabled && styles.disabled,
        style,
      ]}
    >
      {loading ? (
        <ActivityIndicator color={textColor} size="small" />
      ) : (
        <Text style={[styles.text, { color: textColor }]}>{title}</Text>
      )}
    </Pressable>
  );
});

const styles = StyleSheet.create({
  button: {
    borderWidth: 2,
    borderRadius: 0,
    paddingVertical: 12,
    paddingHorizontal: 16,
    alignItems: 'center',
    justifyContent: 'center',
  },
  text: {
    fontFamily: 'monospace',
    fontSize: 14,
    fontWeight: 'bold',
    textTransform: 'uppercase',
  },
  disabled: {
    opacity: 0.5,
  },
});
