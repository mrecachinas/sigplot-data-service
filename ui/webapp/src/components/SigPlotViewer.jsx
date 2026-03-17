import React, { useRef, useEffect, useCallback } from 'react';
import { Plot } from 'sigplot';

function SigPlotPanel({ href, options, plotRef, title }) {
  const containerRef = useRef(null);

  useEffect(() => {
    if (!containerRef.current) return;
    const plot = new Plot(containerRef.current, {
      all: true,
      expand: true,
      autol: 100,
      autohide_panbars: true,
      cmode: 6,
      ...options,
    });
    plotRef.current = plot;
    return () => {
      plot.dispose && plot.dispose();
      plotRef.current = null;
    };
  }, []);

  useEffect(() => {
    const plot = plotRef.current;
    if (!plot || !href) return;
    // Remove existing overlays
    plot.deoverlay();
    plot.overlay_href(href);
  }, [href]);

  return (
    <div className="plot-container">
      <h4>{title}</h4>
      {href ? (
        <div ref={containerRef} style={{ width: 550, height: 350 }} />
      ) : (
        <p className="plot-placeholder">Select a file to view</p>
      )}
    </div>
  );
}

export default function SigPlotViewer({ rawHref, sdsHref }) {
  const rawPlotRef = useRef(null);
  const sdsPlotRef = useRef(null);

  return (
    <>
      <SigPlotPanel
        href={rawHref}
        plotRef={rawPlotRef}
        title="Raw File"
      />
      <SigPlotPanel
        href={sdsHref}
        plotRef={sdsPlotRef}
        title="SDS View"
      />
    </>
  );
}
