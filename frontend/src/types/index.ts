export interface EventItem {
  id: number;
  user_id: number;
  action: string;
  metadata: Record<string, unknown> | null;
  created_at: string;
}

export interface StatItem {
  id: number;
  user_id: number;
  event_count: number;
  period_start: string;
  period_end: string;
  created_at: string;
  updated_at: string;
}

export interface FilterState {
  userId: string;
  from: string;
  to: string;
}

export interface FilterQueryParams {
  user_id?: number;
  from?: string;
  to?: string;
  limit?: number;
  offset?: number;
}

export type SortDirection = 'asc' | 'desc';

export interface SortConfig<T> {
  key: keyof T;
  direction: SortDirection;
}

export type ActiveTab = 'events' | 'stats';
