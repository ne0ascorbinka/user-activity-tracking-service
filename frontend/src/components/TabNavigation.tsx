import React from 'react';
import { Radio, BarChart3 } from 'lucide-react';
import { ActiveTab } from '../types';

interface TabNavigationProps {
  activeTab: ActiveTab;
  onTabChange: (tab: ActiveTab) => void;
  eventCount: number;
  statCount: number;
}

export const TabNavigation: React.FC<TabNavigationProps> = ({
  activeTab,
  onTabChange,
  eventCount,
  statCount,
}) => {
  return (
    <div className="flex border-b border-slate-300">
      <button
        onClick={() => onTabChange('events')}
        className={`flex items-center space-x-2 px-5 py-3 font-mono text-sm font-bold tracking-wide border-t border-l border-r -mb-px cursor-pointer transition-colors ${
          activeTab === 'events'
            ? 'bg-[#0B192C] text-white border-[#0B192C]'
            : 'bg-white text-slate-600 border-transparent hover:text-[#0B192C] hover:bg-slate-50'
        }`}
      >
        <Radio className="w-4 h-4" />
        <span>LIVE EVENTS</span>
        <span
          className={`text-xs px-2 py-0.5 font-mono ${
            activeTab === 'events'
              ? 'bg-[#1E293B] text-slate-200 border border-slate-700'
              : 'bg-slate-100 text-slate-700 border border-slate-300'
          }`}
        >
          {eventCount}
        </span>
      </button>

      <button
        onClick={() => onTabChange('stats')}
        className={`flex items-center space-x-2 px-5 py-3 font-mono text-sm font-bold tracking-wide border-t border-l border-r -mb-px cursor-pointer transition-colors ${
          activeTab === 'stats'
            ? 'bg-[#0B192C] text-white border-[#0B192C]'
            : 'bg-white text-slate-600 border-transparent hover:text-[#0B192C] hover:bg-slate-50'
        }`}
      >
        <BarChart3 className="w-4 h-4" />
        <span>AGGREGATED STATS</span>
        <span
          className={`text-xs px-2 py-0.5 font-mono ${
            activeTab === 'stats'
              ? 'bg-[#1E293B] text-slate-200 border border-slate-700'
              : 'bg-slate-100 text-slate-700 border border-slate-300'
          }`}
        >
          {statCount}
        </span>
      </button>
    </div>
  );
};
