import { useState, useEffect, useCallback } from 'react';
import { Header } from './components/Header';
import { FilterBar } from './components/FilterBar';
import { TabNavigation } from './components/TabNavigation';
import { EventsTable } from './components/EventsTable';
import { StatsTable } from './components/StatsTable';
import { LoadingSkeleton } from './components/LoadingSkeleton';
import { ErrorBanner } from './components/ErrorBanner';
import { ActiveTab, EventItem, FilterQueryParams, FilterState, StatItem } from './types';
import { checkHealth, fetchEvents, fetchStats } from './services/api';
import { toRFC3339 } from './utils/formatters';

export function App() {
  const [activeTab, setActiveTab] = useState<ActiveTab>('events');
  const [events, setEvents] = useState<EventItem[]>([]);
  const [stats, setStats] = useState<StatItem[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [isHealthy, setIsHealthy] = useState<boolean | null>(null);

  const [activeFilters, setActiveFilters] = useState<FilterState>({
    userId: '',
    from: '',
    to: '',
  });

  // Check backend server health
  const verifyHealth = useCallback(async () => {
    try {
      await checkHealth();
      setIsHealthy(true);
    } catch {
      setIsHealthy(false);
    }
  }, []);

  // Convert string filter state to API query parameters
  const getFilterParams = useCallback((filters: FilterState): FilterQueryParams => {
    const params: FilterQueryParams = {};
    if (filters.userId) {
      const parsed = parseInt(filters.userId, 10);
      if (!isNaN(parsed) && parsed > 0) {
        params.user_id = parsed;
      }
    }
    if (filters.from) {
      const rfcFrom = toRFC3339(filters.from);
      if (rfcFrom) params.from = rfcFrom;
    }
    if (filters.to) {
      const rfcTo = toRFC3339(filters.to);
      if (rfcTo) params.to = rfcTo;
    }
    return params;
  }, []);

  // Load data for active tab (and both counts)
  const loadData = useCallback(async (filters: FilterState = activeFilters) => {
    setIsLoading(true);
    setError(null);

    const queryParams = getFilterParams(filters);

    try {
      const [eventsResult, statsResult] = await Promise.all([
        fetchEvents(queryParams),
        fetchStats(queryParams),
      ]);

      setEvents(eventsResult);
      setStats(statsResult);
      setIsHealthy(true);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'An unexpected error occurred while communicating with the server.';
      setError(msg);
      setIsHealthy(false);
    } finally {
      setIsLoading(false);
    }
  }, [activeFilters, getFilterParams]);

  // Initial load
  useEffect(() => {
    verifyHealth();
    loadData(activeFilters);
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const handleApplyFilters = (newFilters: FilterState) => {
    setActiveFilters(newFilters);
    loadData(newFilters);
  };

  const handleResetFilters = () => {
    const emptyFilters: FilterState = { userId: '', from: '', to: '' };
    setActiveFilters(emptyFilters);
    loadData(emptyFilters);
  };

  const isFiltered = Boolean(activeFilters.userId || activeFilters.from || activeFilters.to);

  return (
    <div className="min-h-screen bg-[#F8F7F4] flex flex-col font-sans">
      <Header
        isHealthy={isHealthy}
        onRefresh={() => loadData(activeFilters)}
        isLoading={isLoading}
      />

      <main className="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-6">
        {/* Top Filter Bar */}
        <FilterBar
          onApplyFilters={handleApplyFilters}
          onResetFilters={handleResetFilters}
          isLoading={isLoading}
        />

        {/* Error Banner */}
        {error && (
          <ErrorBanner
            message={error}
            onRetry={() => loadData(activeFilters)}
          />
        )}

        {/* Tab Navigation */}
        <div className="mt-2">
          <TabNavigation
            activeTab={activeTab}
            onTabChange={setActiveTab}
            eventCount={events.length}
            statCount={stats.length}
          />
        </div>

        {/* Main Content Area */}
        <div className="bg-white border-x border-b border-slate-300 p-4 min-h-[400px]">
          {isLoading ? (
            <LoadingSkeleton
              rowCount={8}
              colCount={activeTab === 'events' ? 5 : 5}
            />
          ) : activeTab === 'events' ? (
            <EventsTable
              events={events}
              isFiltered={isFiltered}
              onClearFilters={handleResetFilters}
            />
          ) : (
            <StatsTable
              stats={stats}
              isFiltered={isFiltered}
              onClearFilters={handleResetFilters}
            />
          )}
        </div>
      </main>

      {/* Sharp Footer */}
      <footer className="bg-white border-t border-slate-300 py-3 mt-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex flex-col sm:flex-row items-center justify-between text-[11px] font-mono text-slate-500 gap-2">
          <div>USER ACTIVITY TRACKING SERVICE &bull; CLIENT TELEMETRY</div>
          <div>POWERED BY GO + POSTGRESQL + REACT + TAILWIND</div>
        </div>
      </footer>
    </div>
  );
}

export default App;
