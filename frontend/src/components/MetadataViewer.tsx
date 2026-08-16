import React, { useState } from 'react';
import { Copy, Check, X, Code2 } from 'lucide-react';

interface MetadataViewerProps {
  metadata: unknown;
  eventId: number;
  onClose: () => void;
}

export const MetadataViewer: React.FC<MetadataViewerProps> = ({
  metadata,
  eventId,
  onClose,
}) => {
  const [copied, setCopied] = useState(false);

  const formattedJson = JSON.stringify(metadata ?? {}, null, 2);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(formattedJson);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      // Ignore clipboard write failures
    }
  };

  return (
    <div className="fixed inset-0 z-50 bg-black/60 flex items-center justify-center p-4">
      <div className="bg-white border-2 border-[#0B192C] w-full max-w-xl shadow-2xl animate-in fade-in zoom-in-95 duration-150">
        {/* Header */}
        <div className="bg-[#0B192C] text-white px-4 py-3 flex items-center justify-between">
          <div className="flex items-center space-x-2 text-xs font-mono font-bold tracking-wider uppercase">
            <Code2 className="w-4 h-4 text-slate-300" />
            <span>Event #{eventId} &bull; Metadata Payload</span>
          </div>
          <button
            onClick={onClose}
            className="text-slate-300 hover:text-white cursor-pointer p-1"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Content */}
        <div className="p-4 bg-[#F8F7F4]">
          <div className="relative">
            <button
              onClick={handleCopy}
              className="absolute top-3 right-3 flex items-center space-x-1 px-2.5 py-1 bg-white border border-slate-300 text-xs font-mono text-slate-700 hover:bg-slate-50 cursor-pointer shadow-sm"
            >
              {copied ? (
                <>
                  <Check className="w-3.5 h-3.5 text-emerald-600" />
                  <span className="text-emerald-700 font-bold">COPIED</span>
                </>
              ) : (
                <>
                  <Copy className="w-3.5 h-3.5" />
                  <span>COPY JSON</span>
                </>
              )}
            </button>
            <pre className="p-4 bg-[#0B192C] text-emerald-400 font-mono text-xs overflow-x-auto max-h-96 leading-relaxed select-all">
              {formattedJson}
            </pre>
          </div>
        </div>

        {/* Footer */}
        <div className="bg-white border-t border-slate-300 px-4 py-3 flex justify-end">
          <button
            onClick={onClose}
            className="px-4 py-1.5 bg-[#0B192C] text-white text-xs font-mono font-bold cursor-pointer hover:bg-[#1E293B]"
          >
            CLOSE
          </button>
        </div>
      </div>
    </div>
  );
};
