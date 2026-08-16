/**
 * Formats an ISO 8601 string to a crisp readable timestamp
 * e.g., "2026-08-16 11:28:15 UTC"
 */
export function formatTimestamp(isoString?: string | null): string {
  if (!isoString) return '—';
  try {
    const date = new Date(isoString);
    if (isNaN(date.getTime())) return isoString;

    const pad = (n: number) => n.toString().padStart(2, '0');
    const y = date.getUTCFullYear();
    const m = pad(date.getUTCMonth() + 1);
    const d = pad(date.getUTCDate());
    const hh = pad(date.getUTCHours());
    const mm = pad(date.getUTCMinutes());
    const ss = pad(date.getUTCSeconds());

    return `${y}-${m}-${d} ${hh}:${mm}:${ss} UTC`;
  } catch {
    return isoString;
  }
}

/**
 * Formats a timestamp to a relative time string (e.g. "5m ago", "2h ago")
 */
export function formatRelativeTime(isoString?: string | null): string {
  if (!isoString) return '';
  try {
    const date = new Date(isoString);
    const diff = Math.floor((Date.now() - date.getTime()) / 1000);

    if (diff < 5) return 'just now';
    if (diff < 60) return `${diff}s ago`;
    if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
    if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
    return `${Math.floor(diff / 86400)}d ago`;
  } catch {
    return '';
  }
}

/**
 * Converts a datetime-local input string (YYYY-MM-DDTHH:mm) or date string into RFC3339 UTC string
 */
export function toRFC3339(localDateTimeStr: string): string | undefined {
  if (!localDateTimeStr || !localDateTimeStr.trim()) return undefined;
  try {
    const date = new Date(localDateTimeStr);
    if (isNaN(date.getTime())) return undefined;
    return date.toISOString();
  } catch {
    return undefined;
  }
}

/**
 * Formats JSON metadata for concise inline preview
 */
export function formatMetadataPreview(metadata: unknown): { preview: string; count: number } {
  if (!metadata || typeof metadata !== 'object') {
    return { preview: '{}', count: 0 };
  }
  const keys = Object.keys(metadata);
  if (keys.length === 0) {
    return { preview: '{}', count: 0 };
  }

  const entries = Object.entries(metadata as Record<string, unknown>).slice(0, 2);
  const parts = entries.map(([k, v]) => {
    const valStr = typeof v === 'object' ? '{...}' : JSON.stringify(v);
    return `${k}: ${valStr}`;
  });

  const preview = `{ ${parts.join(', ')}${keys.length > 2 ? ', …' : ''} }`;
  return { preview, count: keys.length };
}
