import React, { useRef, useEffect, useMemo, useState } from 'react';
import { Plot } from 'sigplot';
import { AlertIcon } from './icons';

const PLOT_OPTIONS = {
  all: true,
  expand: true,
  autol: 100,
  autohide_panbars: true,
};

function removeLayer(plot, layer) {
  if (!plot || layer == null) return;
  if (Array.isArray(layer)) {
    layer.forEach((entry) => removeLayer(plot, entry));
    return;
  }
  if (plot.remove_layer) plot.remove_layer(layer);
}

function hasSize(el) {
  return el.clientWidth > 0 && el.clientHeight > 0;
}

function SigPlotPanel({ href, layerOptions, title, fileName }) {
  const containerRef = useRef(null);
  const latestHrefRef = useRef(null);
  const [plot, setPlot] = useState(null);
  const [status, setStatus] = useState('idle');

  // Create the plot once its container has a real size (SigPlot sizes its
  // canvases from the container on construction), then keep it in sync with
  // layout changes via checkresize().
  useEffect(() => {
    const el = containerRef.current;
    if (!el) return undefined;
    const canObserve = typeof ResizeObserver !== 'undefined';
    let instance = null;

    const sync = () => {
      if (instance) {
        if (hasSize(el)) instance.checkresize();
        return;
      }
      if (canObserve && !hasSize(el)) return;
      instance = new Plot(el, PLOT_OPTIONS);
      setPlot(instance);
    };

    sync();
    const observer = canObserve ? new ResizeObserver(sync) : null;
    if (observer) observer.observe(el);

    return () => {
      if (observer) observer.disconnect();
      latestHrefRef.current = null;
      if (instance) {
        if (instance.deoverlay) instance.deoverlay();
        if (instance.disable_listeners) instance.disable_listeners();
      }
      el.replaceChildren();
    };
  }, []);

  useEffect(() => {
    if (!plot) return undefined;
    latestHrefRef.current = href;
    plot.deoverlay();
    if (!href) {
      setStatus('idle');
      return undefined;
    }
    setStatus('loading');
    const layer = plot.overlay_href(
      href,
      {
        onload: () => {
          if (latestHrefRef.current !== href) {
            removeLayer(plot, layer);
            return;
          }
          setStatus('ready');
        },
        onerror: () => {
          if (latestHrefRef.current === href) setStatus('error');
        },
      },
      layerOptions && { ...layerOptions }
    );
    return () => {
      if (latestHrefRef.current === href) latestHrefRef.current = null;
    };
  }, [plot, href, layerOptions]);

  return (
    <section className="plot-card" aria-label={title} aria-busy={status === 'loading'}>
      <header className="plot-card-header">
        <h2 className="plot-card-title">{title}</h2>
        {fileName && (
          <span className="plot-card-file" title={fileName}>
            {fileName}
          </span>
        )}
        {status === 'loading' && (
          <span className="plot-card-status">
            <span className="spinner" aria-hidden="true" />
            Loading
          </span>
        )}
      </header>
      <div className="plot-card-body">
        <div ref={containerRef} className="plot-canvas" data-testid="plot-canvas" />
        {status === 'error' && (
          <div className="plot-card-error" role="alert">
            <AlertIcon size={20} />
            <span>Couldn’t load this file.</span>
          </div>
        )}
      </div>
    </section>
  );
}

export default function SigPlotViewer({ rawHref, sdsHref, fileName }) {
  const sdsLayerOptions = useMemo(() => ({ layerType: 'SDS' }), []);

  return (
    <>
      <SigPlotPanel href={rawHref} title="Raw file" fileName={fileName} />
      <SigPlotPanel
        href={sdsHref}
        layerOptions={sdsLayerOptions}
        title="SDS tiled"
        fileName={fileName}
      />
    </>
  );
}
