import { EventItem, FilterQueryParams, StatItem } from '../types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

function buildQueryString(params: FilterQueryParams): string {
  const query = new URLSearchParams();
  if (params.user_id !== undefined && params.user_id > 0) {
    query.set('user_id', params.user_id.toString());
  }
  if (params.from) {
    query.set('from', params.from);
  }
  if (params.to) {
    query.set('to', params.to);
  }
  if (params.limit) {
    query.set('limit', params.limit.toString());
  }
  if (params.offset !== undefined && params.offset > 0) {
    query.set('offset', params.offset.toString());
  }

  const qs = query.toString();
  return qs ? `?${qs}` : '';
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    let errorMessage = `HTTP ${response.status} ${response.statusText}`;
    try {
      const errorData = await response.json();
      if (errorData && typeof errorData.error === 'string') {
        errorMessage = errorData.error;
      }
    } catch {
      // Use fallback error message
    }
    throw new Error(errorMessage);
  }
  return response.json();
}

export async function fetchEvents(params: FilterQueryParams = {}): Promise<EventItem[]> {
  const url = `${API_BASE_URL}/api/v1/events${buildQueryString(params)}`;
  const response = await fetch(url, {
    headers: {
      Accept: 'application/json',
    },
  });
  const data = await handleResponse<EventItem[] | null>(response);
  return data ?? [];
}

export async function fetchStats(params: FilterQueryParams = {}): Promise<StatItem[]> {
  const url = `${API_BASE_URL}/api/v1/stats${buildQueryString(params)}`;
  const response = await fetch(url, {
    headers: {
      Accept: 'application/json',
    },
  });
  const data = await handleResponse<StatItem[] | null>(response);
  return data ?? [];
}

export async function checkHealth(): Promise<{ status: string }> {
  const url = `${API_BASE_URL}/health`;
  const response = await fetch(url);
  return handleResponse<{ status: string }>(response);
}
