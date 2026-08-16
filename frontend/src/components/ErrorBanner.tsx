import React from 'react';
import { AlertTriangle, RefreshCw } from 'lucide-react';

interface ErrorBannerProps {
  message: string;
  onRetry?: () => void;
}

export const ErrorBanner: React.FC<ErrorBannerProps> = ({ message, onRetry }) => {
  return (
    <div className="bg-rose-50 border-l-4 border-rose-600 p-4 mb-6 text-slate-800 shadow-sm border-y border-r border-rose-200">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
        <div className="flex items-start space-x-3">
          <AlertTriangle className="w-5 h-5 text-rose-600 shrink-0 mt-0.5" />
          <div>
            <h4 className="text-xs font-mono font-bold uppercase text-rose-900">
              Request Failed
            </h4>
            <p className="text-xs font-mono text-rose-700 mt-0.5 break-all">
              {message}
            </p>
          </div>
        </div>

        {onRetry && (
          <button
            onClick={onRetry}
            className="self-start sm:self-center flex items-center space-x-1.5 px-3 py-1.5 bg-rose-600 text-white hover:bg-rose-700 active:bg-rose-800 text-xs font-mono font-bold cursor-pointer transition-colors"
          >
            <RefreshCw className="w-3.5 h-3.5" />
            <span>RETRY</span>
          </button>
        )}
      </div>
    </div>
  );
};
