import {
  canTransition,
  requiresChecklistCompletion,
  getValidTransitions,
  getStatusLabel,
  getStatusColor,
} from '../status-transition';

describe('status-transition utils', () => {
  describe('canTransition', () => {
    it('should allow assigned to en_route', () => {
      expect(canTransition('assigned', 'en_route')).toBe(true);
    });

    it('should allow assigned to cancelled', () => {
      expect(canTransition('assigned', 'cancelled')).toBe(true);
    });

    it('should not allow assigned to completed', () => {
      expect(canTransition('assigned', 'completed')).toBe(false);
    });

    it('should allow en_route to in_progress', () => {
      expect(canTransition('en_route', 'in_progress')).toBe(true);
    });

    it('should allow in_progress to completed', () => {
      expect(canTransition('in_progress', 'completed')).toBe(true);
    });

    it('should allow completed to assigned', () => {
      expect(canTransition('completed', 'assigned')).toBe(true);
    });

    it('should not allow completed to in_progress', () => {
      expect(canTransition('completed', 'in_progress')).toBe(false);
    });
  });

  describe('requiresChecklistCompletion', () => {
    it('should require checklist for completed', () => {
      expect(requiresChecklistCompletion('in_progress', 'completed')).toBe(true);
    });

    it('should not require checklist for en_route', () => {
      expect(requiresChecklistCompletion('assigned', 'en_route')).toBe(false);
    });
  });

  describe('getValidTransitions', () => {
    it('should return transitions for assigned', () => {
      const transitions = getValidTransitions('assigned');
      expect(transitions).toContain('en_route');
      expect(transitions).toContain('cancelled');
    });

    it('should return empty for unknown status', () => {
      const transitions = getValidTransitions('unknown' as any);
      expect(transitions).toEqual([]);
    });
  });

  describe('getStatusLabel', () => {
    it('should return Spanish labels', () => {
      expect(getStatusLabel('assigned')).toBe('Asignada');
      expect(getStatusLabel('en_route')).toBe('En camino');
      expect(getStatusLabel('in_progress')).toBe('En progreso');
      expect(getStatusLabel('completed')).toBe('Completada');
    });
  });

  describe('getStatusColor', () => {
    it('should return valid colors', () => {
      expect(getStatusColor('assigned')).toBe('#6b7280');
      expect(getStatusColor('en_route')).toBe('#2563eb');
      expect(getStatusColor('in_progress')).toBe('#f59e0b');
      expect(getStatusColor('completed')).toBe('#16a34a');
    });
  });
});
