import React, { useRef, useEffect } from 'react';
import { SigPlot, HrefLayer } from 'react-sigplot';

export default function SigPlotViewer({ rawHref, sdsHref }) {
  const rawPlotRef = useRef(null);
  const sdsPlotRef = useRef(null);
  const mimicSetup = useRef(false);

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

  useEffect(() => {
    mimicSetup.current = false;
  }, [rawHref, sdsHref]);

  return (
    <>
      <div className="plot-container">
        <h4>Raw File</h4>
        {rawHref ? (
          <SigPlot
            ref={rawPlotRef}
            height={350}
            width={550}
            options={{ all: true, expand: true, autol: 100, autohide_panbars: true, cmode: 6 }}
          >
            <HrefLayer href={rawHref} />
          </SigPlot>
        ) : (
          <p className="plot-placeholder">Select a file to view</p>
        )}
      </div>

      <div className="plot-container">
        <h4>SDS Tiled View</h4>
        {sdsHref ? (
          <SigPlot
            ref={sdsPlotRef}
            height={350}
            width={550}
            options={{ all: true, expand: true, autol: 100, autohide_panbars: true, cmode: 6 }}
          >
            <HrefLayer href={sdsHref} options={{ usetiles: true }} />
          </SigPlot>
        ) : (
          <p className="plot-placeholder">Select a file to view</p>
        )}
      </div>
    </>
  );
}
