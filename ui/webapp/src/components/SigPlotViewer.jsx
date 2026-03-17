import React, { useRef, useEffect } from 'react';
import { SigPlot, HrefLayer } from 'react-sigplot';

export default function SigPlotViewer({ rawHref, sdsHref }) {
  const rawPlotRef = useRef(null);
  const sdsPlotRef = useRef(null);
  const mimicSetup = useRef(false);

  // Set up bidirectional mimic after both plots mount
  useEffect(() => {
    if (mimicSetup.current) return;

    const rawPlot = rawPlotRef.current && rawPlotRef.current.plot;
    const sdsPlot = sdsPlotRef.current && sdsPlotRef.current.plot;

    if (rawPlot && sdsPlot) {
      const mask = { xzoom: true, yzoom: true, xpan: true, ypan: true, unzoom: true };
      try {
        rawPlot.mimic(sdsPlot, mask);
        sdsPlot.mimic(rawPlot, mask);
        mimicSetup.current = true;
      } catch {
        // mimic may not be available in all sigplot versions
      }
    }
  });

  // Reset mimic when files change
  useEffect(() => {
    mimicSetup.current = false;
  }, [rawHref, sdsHref]);

  if (!rawHref && !sdsHref) {
    return (
      <div className="plot-area">
        <p className="plot-placeholder">Select a file to view</p>
      </div>
    );
  }

  return (
    <div className="plot-area">
      <div className="plot-container">
        <h3>Raw File</h3>
        <SigPlot
          ref={rawPlotRef}
          height={400}
          width={600}
          options={{ all: true, expand: true, autol: 100, autohide_panbars: true }}
        >
          {rawHref && <HrefLayer href={rawHref} />}
        </SigPlot>
      </div>

      <div className="plot-container">
        <h3>SDS Tiled View</h3>
        <SigPlot
          ref={sdsPlotRef}
          height={400}
          width={600}
          options={{ all: true, expand: true, autol: 100, autohide_panbars: true, cmode: 6 }}
        >
          {sdsHref && <HrefLayer href={sdsHref} options={{ usetiles: true }} />}
        </SigPlot>
      </div>
    </div>
  );
}
