import React, { useState } from 'react';
import { Filter, RotateCcw, Search } from 'lucide-react';
import { FilterState } from '../types';

interface FilterBarProps {
  onApplyFilters: (filters: FilterState) => void;
  onResetFilters: () => void;
  isLoading: boolean;
}

export const FilterBar: React.FC<FilterBarProps> = ({
  onApplyFilters,
  onResetFilters,
  isLoading,
}) => {
  const [userId, setUserId] = useState('');
  const [from, setFrom] = useState('');
  const [to, setTo] = useState('');

  const handleApply = (e?: React.FormEvent) => {
    if (e) e.preventDefault();
    onApplyFilters({
      userId: userId.trim(),
      from: from.trim(),
      to: to.trim(),
    });
  };

  const handleReset = () => {
    setUserId('');
    setFrom('');
    setTo('');
    onResetFilters();
  };

  const isFiltering = Boolean(userId || from || to);

  return (
    <div className="bg-white border border-slate-300 p-4 shadow-sm mb-6">
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center space-x-2 text-xs font-mono font-bold tracking-wider text-[#0B192C] uppercase">
          <Filter className="w-4 h-4 text-slate-500" />
          <span>Query Filters</span>
        </div>
        {isFiltering && (
          <span className="text-[11px] font-mono text-emerald-700 bg-emerald-50 border border-emerald-300 px-2 py-0.5">
            FILTER ACTIVE
          </span>
        )}
      </div>

      <form onSubmit={handleApply} className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 items-end">
        {/* User ID input */}
        <div>
          <label className="block text-xs font-mono text-slate-600 uppercase mb-1">
            User ID (Integer)
          </label>
          <input
            type="number"
            min="1"
            placeholder="e.g. 42"
            value={userId}
            onChange={(e) => setUserId(e.target.value)}
            className="w-full bg-[#F8F7F4] border border-slate-300 px-3 py-2 text-xs font-mono text-[#0B192C] focus:outline-none focus:border-[#0B192C] focus:bg-white transition-colors"
          />
        </div>

        {/* From date */}
        <div>
          <label className="block text-xs font-mono text-slate-600 uppercase mb-1">
            From Timestamp (UTC / ISO)
          </label>
          <input
            type="datetime-local"
            value={from}
            onChange={(e) => setFrom(e.target.value)}
            className="w-full bg-[#F8F7F4] border border-slate-300 px-3 py-2 text-xs font-mono text-[#0B192C] focus:outline-none focus:border-[#0B192C] focus:bg-white transition-colors"
          />
        </div>

        {/* To date */}
        <div>
          <label className="block text-xs font-mono text-slate-600 uppercase mb-1">
            To Timestamp (UTC / ISO)
          </label>
          <input
            type="datetime-local"
            value={to}
            onChange={(e) => setTo(e.target.value)}
            className="w-full bg-[#F8F7F4] border border-slate-300 px-3 py-2 text-xs font-mono text-[#0B192C] focus:outline-none focus:border-[#0B192C] focus:bg-white transition-colors"
          />
        </div>

        {/* Actions */}
        <div className="flex items-center space-x-2">
          <button
            type="submit"
            disabled={isLoading}
            className="flex-1 flex items-center justify-center space-x-1.5 bg-[#0B192C] text-white hover:bg-[#1E293B] active:bg-black px-4 py-2 text-xs font-mono font-bold tracking-wider cursor-pointer transition-colors disabled:opacity-50"
          >
            <Search className="w-3.5 h-3.5" />
            <span>APPLY</span>
          </button>

          <button
            type="button"
            onClick={handleReset}
            disabled={isLoading && !isFiltering}
            title="Reset Filters"
            className="flex items-center justify-center space-x-1.5 bg-slate-100 text-slate-700 hover:bg-slate-200 active:bg-slate-300 border border-slate-300 px-3 py-2 text-xs font-mono font-bold cursor-pointer transition-colors"
          >
            <RotateCcw className="w-3.5 h-3.5" />
            <span>RESET</span>
          </button>
        </div>
      </form>
    </div>
  );
};
