type VisitStatus =
  | 'assigned'
  | 'en_route'
  | 'in_progress'
  | 'completed'
  | 'cancelled'
  | 'rescheduled';

interface StatusTransition {
  from: VisitStatus;
  to: VisitStatus;
  requiresChecklist?: boolean;
}

const VALID_TRANSITIONS: Record<VisitStatus, VisitStatus[]> = {
  assigned: ['en_route', 'cancelled'],
  en_route: ['in_progress', 'cancelled'],
  in_progress: ['completed', 'cancelled'],
  completed: ['assigned'],
  cancelled: ['assigned'],
  rescheduled: ['assigned'],
};

const CHECKLIST_REQUIRED_TRANSITIONS: VisitStatus[] = ['completed'];

export function canTransition(from: VisitStatus, to: VisitStatus): boolean {
  const allowed = VALID_TRANSITIONS[from];
  return allowed ? allowed.includes(to) : false;
}

export function requiresChecklistCompletion(from: VisitStatus, to: VisitStatus): boolean {
  return CHECKLIST_REQUIRED_TRANSITIONS.includes(to);
}

export function getValidTransitions(currentStatus: VisitStatus): VisitStatus[] {
  return VALID_TRANSITIONS[currentStatus] || [];
}

export function getStatusLabel(status: VisitStatus): string {
  const labels: Record<VisitStatus, string> = {
    assigned: 'Asignada',
    en_route: 'En camino',
    in_progress: 'En progreso',
    completed: 'Completada',
    cancelled: 'Cancelada',
    rescheduled: 'Reprogramada',
  };
  return labels[status] || status;
}

export function getStatusColor(status: VisitStatus): string {
  const colors: Record<VisitStatus, string> = {
    assigned: '#6b7280',
    en_route: '#2563eb',
    in_progress: '#f59e0b',
    completed: '#16a34a',
    cancelled: '#dc2626',
    rescheduled: '#8b5cf6',
  };
  return colors[status] || '#6b7280';
}
