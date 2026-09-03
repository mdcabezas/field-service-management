import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { Button } from './button';
import { Badge } from './badge';
import {
  canTransition,
  requiresChecklistCompletion,
  getValidTransitions,
  getStatusLabel,
  getStatusColor,
} from '../../utils/status-transition';

interface BlockingStage {
  name: string;
  requiredItems: number;
  completedItems: number;
}

interface StatusTransitionProps {
  currentStatus: string;
  onTransition: (newStatus: string) => void;
  checklistComplete?: boolean;
  blockingStages?: BlockingStage[];
}

export const StatusTransitionComponent = React.memo(function StatusTransitionComponent({
  currentStatus,
  onTransition,
  checklistComplete = true,
  blockingStages = [],
}: StatusTransitionProps) {
  const validTransitions = getValidTransitions(currentStatus as any);
  const showBlocking = blockingStages.length > 0 && !checklistComplete;

  if (validTransitions.length === 0) {
    return (
      <View style={styles.container}>
        <Badge
          label={getStatusLabel(currentStatus as any)}
          color={getStatusColor(currentStatus as any)}
        />
        <Text style={styles.noTransitions}>Sin acciones disponibles</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <View style={styles.currentStatus}>
        <Text style={styles.label}>Estado actual:</Text>
        <Badge
          label={getStatusLabel(currentStatus as any)}
          color={getStatusColor(currentStatus as any)}
        />
      </View>

      {showBlocking && (
        <View style={styles.blockingContainer}>
          <Text style={styles.blockingTitle}>Etapas incompletas:</Text>
          {blockingStages.map((stage, index) => (
            <View key={index} style={styles.stageRow}>
              <Text style={styles.stageName}>{stage.name}</Text>
              <Text style={styles.stageProgress}>
                {stage.completedItems}/{stage.requiredItems}
              </Text>
            </View>
          ))}
        </View>
      )}

      <View style={styles.actions}>
        {validTransitions.map((status) => {
          const needsChecklist = requiresChecklistCompletion(currentStatus as any, status);
          const disabled = needsChecklist && !checklistComplete;

          return (
            <Button
              key={status}
              title={getStatusLabel(status)}
              onPress={() => onTransition(status)}
              variant={status === 'completed' ? 'primary' : 'secondary'}
              disabled={disabled}
              style={styles.transitionButton}
            />
          );
        })}
      </View>
    </View>
  );
});

const styles = StyleSheet.create({
  container: {
    padding: 16,
  },
  currentStatus: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    marginBottom: 16,
  },
  label: {
    fontFamily: 'monospace',
    fontSize: 12,
    textTransform: 'uppercase',
  },
  actions: {
    gap: 8,
  },
  transitionButton: {
    marginBottom: 8,
  },
  noTransitions: {
    fontFamily: 'monospace',
    fontSize: 12,
    color: '#666',
    marginTop: 8,
  },
  blockingContainer: {
    marginBottom: 16,
    padding: 12,
    backgroundColor: '#fef3c7',
    borderWidth: 1,
    borderColor: '#fcd34d',
  },
  blockingTitle: {
    fontFamily: 'monospace',
    fontSize: 12,
    fontWeight: 'bold',
    textTransform: 'uppercase',
    marginBottom: 8,
  },
  stageRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingVertical: 4,
  },
  stageName: {
    fontFamily: 'monospace',
    fontSize: 11,
  },
  stageProgress: {
    fontFamily: 'monospace',
    fontSize: 11,
    fontWeight: 'bold',
  },
});
