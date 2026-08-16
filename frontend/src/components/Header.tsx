import React from 'react';
import { Activity, RefreshCw } from 'lucide-react';

interface HeaderProps {
  isHealthy: boolean | null;
  onRefresh: () => void;
  isLoading: boolean;
}

export const Header: React.FC<HeaderProps> = ({
  isHealthy,
  onRefresh,
  isLoading,
}) => {
  return (
    <header className="w-full bg-[#0B192C] text-white border-b border-[#1E293B]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-white text-[#0B192C]">
              <Activity className="w-6 h-6" />
            </div>
            <div>
              <div className="flex items-center space-x-3">
                <h1 className="text-xl font-bold tracking-tight uppercase font-mono">
                  Activity Tracker
                </h1>
                <span className="text-xs bg-[#1E293B] text-slate-300 px-2 py-0.5 font-mono border border-slate-700">
                  v1.0
                </span>
              </div>
              <p className="text-xs text-slate-400 font-mono mt-0.5">
                USER TELEMETRY &bull; 4-HOUR PERIODIC AGGREGATIONS
              </p>
            </div>
          </div>

          <div className="flex items-center space-x-3 self-end sm:self-auto">
            <div className="flex items-center space-x-2 px-3 py-1.5 bg-[#1E293B] border border-slate-700 text-xs font-mono">
              <span
                className={`w-2 h-2 inline-block ${
                  isHealthy === true
                    ? 'bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.8)]'
                    : isHealthy === false
                    ? 'bg-rose-500 shadow-[0_0_8px_rgba(244,63,94,0.8)]'
                    : 'bg-amber-400'
                }`}
              />
              <span className="text-slate-300">
                {isHealthy === true
                  ? 'API ONLINE'
                  : isHealthy === false
                  ? 'API OFFLINE'
                  : 'CHECKING...'}
              </span>
            </div>

            <button
              onClick={onRefresh}
              disabled={isLoading}
              title="Refresh Data"
              className="flex items-center space-x-1.5 px-3 py-1.5 bg-white text-[#0B192C] hover:bg-slate-100 active:bg-slate-200 transition-colors text-xs font-mono font-bold cursor-pointer disabled:opacity-50"
            >
              <RefreshCw className={`w-3.5 h-3.5 ${isLoading ? 'animate-spin' : ''}`} />
              <span>REFRESH</span>
            </button>
          </div>
        </div>
      </div>
    </header>
  );
};
