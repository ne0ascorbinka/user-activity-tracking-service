import React from 'react';
import { Inbox, RotateCcw } from 'lucide-react';

interface EmptyStateProps {
  title?: string;
  description?: string;
  isFiltered?: boolean;
  onClearFilters?: () => void;
}

export const EmptyState: React.FC<EmptyStateProps> = ({
  title = 'No Records Found',
  description = 'There are currently no records available to display.',
  isFiltered = false,
  onClearFilters,
}) => {
  return (
    <div className="w-full bg-white border border-slate-300 p-12 text-center my-4">
      <div className="inline-flex p-3 bg-slate-100 border border-slate-300 text-slate-500 mb-3">
        <Inbox className="w-8 h-8" />
      </div>
      <h3 className="text-sm font-mono font-bold text-[#0B192C] uppercase tracking-wider">
        {title}
      </h3>
      <p className="text-xs text-slate-500 font-mono mt-1 max-w-sm mx-auto">
        {description}
      </p>

      {isFiltered && onClearFilters && (
        <div className="mt-4">
          <button
            onClick={onClearFilters}
            className="inline-flex items-center space-x-1.5 px-3 py-1.5 bg-[#0B192C] text-white hover:bg-[#1E293B] text-xs font-mono font-bold cursor-pointer"
          >
            <RotateCcw className="w-3.5 h-3.5" />
            <span>RESET FILTERS</span>
          </button>
        </div>
      )}
    </div>
  );
};
