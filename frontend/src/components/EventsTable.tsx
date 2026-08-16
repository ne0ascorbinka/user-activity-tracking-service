import React, { useState } from 'react';
import { EventItem, SortConfig } from '../types';
import { formatTimestamp, formatRelativeTime, formatMetadataPreview } from '../utils/formatters';
import { sortData } from '../utils/sorting';
import { MetadataViewer } from './MetadataViewer';
import { EmptyState } from './EmptyState';
import { Eye } from 'lucide-react';

interface EventsTableProps {
  events: EventItem[];
  isFiltered: boolean;
  onClearFilters: () => void;
}

export const EventsTable: React.FC<EventsTableProps> = ({
  events,
  isFiltered,
  onClearFilters,
}) => {
  const [sortConfig, setSortConfig] = useState<SortConfig<EventItem> | null>({
    key: 'id',
    direction: 'desc',
  });
  const [selectedMetadataEvent, setSelectedMetadataEvent] = useState<EventItem | null>(null);

  const handleSort = (key: keyof EventItem) => {
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

  const sortedEvents = sortData(events, sortConfig);

  const renderSortIndicator = (key: keyof EventItem) => {
    if (!sortConfig || sortConfig.key !== key) {
      return <span className="text-slate-500 opacity-40 ml-1 text-[10px]">▲▼</span>;
    }
    return (
      <span className="text-emerald-400 font-bold ml-1 text-xs">
        {sortConfig.direction === 'asc' ? '▲' : '▼'}
      </span>
    );
  };

  if (events.length === 0) {
    return (
      <EmptyState
        title="No Live Events Found"
        description={
          isFiltered
            ? 'No events match the specified filter criteria.'
            : 'No events have been ingested into the system yet.'
        }
        isFiltered={isFiltered}
        onClearFilters={onClearFilters}
      />
    );
  }

  return (
    <>
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
                onClick={() => handleSort('action')}
                className="py-3 px-4 cursor-pointer hover:bg-[#1E293B] transition-colors w-44"
              >
                <div className="flex items-center space-x-1">
                  <span>ACTION</span>
                  {renderSortIndicator('action')}
                </div>
              </th>

              <th className="py-3 px-4">
                <span>METADATA</span>
              </th>

              <th
                onClick={() => handleSort('created_at')}
                className="py-3 px-4 cursor-pointer hover:bg-[#1E293B] transition-colors w-56"
              >
                <div className="flex items-center space-x-1">
                  <span>CREATED AT</span>
                  {renderSortIndicator('created_at')}
                </div>
              </th>
            </tr>
          </thead>

          <tbody className="divide-y divide-slate-200 text-xs font-mono">
            {sortedEvents.map((evt) => {
              const metaInfo = formatMetadataPreview(evt.metadata);

              return (
                <tr
                  key={evt.id}
                  className="hover:bg-slate-50 transition-colors group"
                >
                  {/* ID */}
                  <td className="py-3 px-4 font-bold text-slate-700">
                    #{evt.id}
                  </td>

                  {/* USER ID */}
                  <td className="py-3 px-4">
                    <span className="inline-block px-2 py-0.5 bg-slate-100 border border-slate-300 font-bold text-[#0B192C]">
                      usr_{evt.user_id}
                    </span>
                  </td>

                  {/* ACTION */}
                  <td className="py-3 px-4">
                    <span className="inline-block px-2.5 py-0.5 bg-sky-50 text-sky-800 border border-sky-200 font-semibold uppercase tracking-wider text-[11px]">
                      {evt.action}
                    </span>
                  </td>

                  {/* METADATA */}
                  <td className="py-3 px-4">
                    <button
                      onClick={() => setSelectedMetadataEvent(evt)}
                      title="Inspect full JSON payload"
                      className="inline-flex items-center space-x-1.5 px-2 py-1 bg-slate-100 hover:bg-slate-200 border border-slate-300 text-slate-700 text-[11px] cursor-pointer group-hover:border-slate-400 transition-colors"
                    >
                      <Eye className="w-3 h-3 text-slate-500" />
                      <span className="font-mono text-slate-600 truncate max-w-xs">
                        {metaInfo.preview}
                      </span>
                    </button>
                  </td>

                  {/* CREATED AT */}
                  <td className="py-3 px-4 text-slate-600">
                    <div>{formatTimestamp(evt.created_at)}</div>
                    <div className="text-[10px] text-slate-400">
                      {formatRelativeTime(evt.created_at)}
                    </div>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      {selectedMetadataEvent && (
        <MetadataViewer
          eventId={selectedMetadataEvent.id}
          metadata={selectedMetadataEvent.metadata}
          onClose={() => setSelectedMetadataEvent(null)}
        />
      )}
    </>
  );
};
