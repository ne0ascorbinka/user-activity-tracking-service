import React from 'react';

interface LoadingSkeletonProps {
  rowCount?: number;
  colCount?: number;
}

export const LoadingSkeleton: React.FC<LoadingSkeletonProps> = ({
  rowCount = 8,
  colCount = 5,
}) => {
  return (
    <div className="w-full border border-slate-300 bg-white overflow-hidden animate-pulse">
      <div className="bg-[#0B192C] h-10 w-full" />
      <div className="divide-y divide-slate-200">
        {Array.from({ length: rowCount }).map((_, rIndex) => (
          <div
            key={rIndex}
            className="flex items-center px-4 py-3 space-x-4 bg-white"
          >
            {Array.from({ length: colCount }).map((_, cIndex) => (
              <div
                key={cIndex}
                className={`h-4 bg-slate-200 ${
                  cIndex === 0
                    ? 'w-12'
                    : cIndex === 1
                    ? 'w-20'
                    : cIndex === 2
                    ? 'w-32'
                    : cIndex === 3
                    ? 'flex-1'
                    : 'w-40'
                }`}
              />
            ))}
          </div>
        ))}
      </div>
    </div>
  );
};
