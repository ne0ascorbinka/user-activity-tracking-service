import React, { useState } from 'react';
import { StatItem, SortConfig } from '../types';
import { formatTimestamp, formatRelativeTime } from '../utils/formatters';
import { sortData } from '../utils/sorting';
import { EmptyState } from './EmptyState';
import { Layers } from 'lucide-react';

interface StatsTableProps {
  stats: StatItem[];
  isFiltered: boolean;
  onClearFilters: () => void;
}

export const StatsTable: React.FC<StatsTableProps> = ({
  stats,
  isFiltered,
  onClearFilters,
}) => {
  const [sortConfig, setSortConfig] = useState<SortConfig<StatItem> | null>({
    key: 'period_start',
    direction: 'desc',
  });

  const handleSort = (key: keyof StatItem) => {
    setSortConfig((prev) => {
      if (prev && prev.key === key) {
        return {
          key,
          direction: prev.direction === 'asc' ? 'desc' : 'asc',
        };
      }
      return { key, direction: 'asc' };
    });
  };

  const sortedStats = sortData(stats, sortConfig);

  const renderSortIndicator = (key: keyof StatItem) => {
    if (!sortConfig || sortConfig.key !== key) {
      return <span className="text-slate-500 opacity-40 ml-1 text-[10px]">▲▼</span>;
    }
    return (
      <span className="text-emerald-400 font-bold ml-1 text-xs">
        {sortConfig.direction === 'asc' ? '▲' : '▼'}
      </span>
    );
  };

  if (stats.length === 0) {
    return (
      <EmptyState
        title="No Aggregated Stats Found"
        description={
          isFiltered
            ? 'No statistics match the specified filter criteria.'
            : 'Background aggregation worker runs every 4 hours. No periodic stats exist yet.'
        }
        isFiltered={isFiltered}
        onClearFilters={onClearFilters}
      />
    );
  }

  return (
    <div className="w-full bg-white border border-slate-300 overflow-x-auto shadow-sm">
      <table className="w-full text-left border-collapse">
        <thead>
          <tr className="bg-[#0B192C] text-white text-xs font-mono tracking-wider select-none border-b border-[#1E293B]">
            <th
              onClick={() => handleSort('id')}
              className="py-3 px-4 cursor-pointer hover:bg-[#1E293B] transition-colors w-20"
            >
              <div className="flex items-center space-x-1">
                <span>ID</span>
                {renderSortIndicator('id')}
              </div>
            </th>

            <th
              onClick={() => handleSort('user_id')}
              className="py-3 px-4 cursor-pointer hover:bg-[#1E293B] transition-colors w-28"
            >
              <div className="flex items-center space-x-1">
                <span>USER ID</span>
                {renderSortIndicator('user_id')}
              </div>
            </th>

            <th
              onClick={() => handleSort('event_count')}
              className="py-3 px-4 cursor-pointer hover:bg-[#1E293B] transition-colors w-36"
            >
              <div className="flex items-center space-x-1">
                <span>EVENT COUNT</span>
                {renderSortIndicator('event_count')}
              </div>
            </th>

            <th
              onClick={() => handleSort('period_start')}
              className="py-3 px-4 cursor-pointer hover:bg-[#1E293B] transition-colors"
            >
              <div className="flex items-center space-x-1">
                <span>4-HOUR WINDOW (START &rarr; END)</span>
                {renderSortIndicator('period_start')}
              </div>
            </th>

            <th
              onClick={() => handleSort('updated_at')}
              className="py-3 px-4 cursor-pointer hover:bg-[#1E293B] transition-colors w-52"
            >
              <div className="flex items-center space-x-1">
                <span>LAST AGGREGATED</span>
                {renderSortIndicator('updated_at')}
              </div>
            </th>
          </tr>
        </thead>

        <tbody className="divide-y divide-slate-200 text-xs font-mono">
          {sortedStats.map((stat) => (
            <tr
              key={stat.id}
              className="hover:bg-slate-50 transition-colors group"
            >
              {/* ID */}
              <td className="py-3 px-4 font-bold text-slate-700">
                #{stat.id}
              </td>

              {/* USER ID */}
              <td className="py-3 px-4">
                <span className="inline-block px-2 py-0.5 bg-slate-100 border border-slate-300 font-bold text-[#0B192C]">
                  usr_{stat.user_id}
                </span>
              </td>

              {/* EVENT COUNT */}
              <td className="py-3 px-4">
                <div className="flex items-center space-x-2">
                  <span className="inline-flex items-center space-x-1 px-2.5 py-0.5 bg-emerald-50 text-emerald-800 border border-emerald-300 font-bold text-xs">
                    <Layers className="w-3 h-3 text-emerald-600" />
                    <span>{stat.event_count.toLocaleString()}</span>
                  </span>
                  <span className="text-[10px] text-slate-400">events</span>
                </div>
              </td>

              {/* PERIOD WINDOW */}
              <td className="py-3 px-4 text-slate-700">
                <div className="flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-2">
                  <span className="bg-slate-100 px-2 py-0.5 border border-slate-200">
                    {formatTimestamp(stat.period_start)}
                  </span>
                  <span className="text-slate-400">&rarr;</span>
                  <span className="bg-slate-100 px-2 py-0.5 border border-slate-200">
                    {formatTimestamp(stat.period_end)}
                  </span>
                </div>
              </td>

              {/* UPDATED AT */}
              <td className="py-3 px-4 text-slate-600">
                <div>{formatTimestamp(stat.updated_at || stat.created_at)}</div>
                <div className="text-[10px] text-slate-400">
                  {formatRelativeTime(stat.updated_at || stat.created_at)}
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
