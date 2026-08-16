import { SortConfig } from '../types';

export function sortData<T>(data: T[], sortConfig: SortConfig<T> | null): T[] {
  if (!sortConfig) return data;

  const { key, direction } = sortConfig;

  return [...data].sort((a, b) => {
    const valA = a[key];
    const valB = b[key];

    if (valA === valB) return 0;
    if (valA === undefined || valA === null) return 1;
    if (valB === undefined || valB === null) return -1;

    // Handle number comparisons
    if (typeof valA === 'number' && typeof valB === 'number') {
      return direction === 'asc' ? valA - valB : valB - valA;
    }

    // Handle string / date comparisons
    const strA = String(valA).toLowerCase();
    const strB = String(valB).toLowerCase();

    if (strA < strB) {
      return direction === 'asc' ? -1 : 1;
    }
    if (strA > strB) {
      return direction === 'asc' ? 1 : -1;
    }
    return 0;
  });
}
