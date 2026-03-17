import React, { useState, useEffect, useCallback } from 'react';
import FileBrowser from './components/FileBrowser';
import SigPlotViewer from './components/SigPlotViewer';
import { getLocations, getFiles, getFileUrl } from './api/sds';

export default function App() {
  const [locations, setLocations] = useState([]);
  const [selectedLocation, setSelectedLocation] = useState('');
  const [files, setFiles] = useState([]);
  const [path, setPath] = useState('');
  const [pathStack, setPathStack] = useState([]);
  const [rawHref, setRawHref] = useState(null);
  const [sdsHref, setSdsHref] = useState(null);

  // Fetch locations on mount and poll every 5s
  useEffect(() => {
    const fetchLocations = async () => {
      const locs = await getLocations();
      setLocations(locs);
    };
    fetchLocations();
    const interval = setInterval(fetchLocations, 5000);
    return () => clearInterval(interval);
  }, []);

  // Fetch files when location or path changes
  useEffect(() => {
    if (!selectedLocation) {
      setFiles([]);
      return;
    }
    const fetchFiles = async () => {
      const fileList = await getFiles(selectedLocation, path);
      setFiles(fileList);
    };
    fetchFiles();
  }, [selectedLocation, path]);

  const handleSelectLocation = useCallback((loc) => {
    setSelectedLocation(loc);
    setPath('');
    setPathStack([]);
    setRawHref(null);
    setSdsHref(null);
  }, []);

  const handleSelectFile = useCallback(
    (file) => {
      if (file.type === 'directory') {
        setPathStack((prev) => [...prev, path]);
        setPath((prev) => (prev ? `${prev}/${file.filename}` : file.filename));
      } else {
        const filePath = path ? `${path}/${file.filename}` : file.filename;
        setRawHref(getFileUrl(filePath, 'fs', selectedLocation));
        setSdsHref(getFileUrl(filePath, 'hdr', selectedLocation));
      }
    },
    [selectedLocation, path]
  );

  const handleGoBack = useCallback(() => {
    setPathStack((prev) => {
      const next = [...prev];
      const parentPath = next.pop();
      setPath(parentPath || '');
      return next;
    });
  }, []);

  return (
    <div className="app">
      <header>
        <h1>SigPlot Data Service</h1>
      </header>
      <main>
        <FileBrowser
          locations={locations}
          selectedLocation={selectedLocation}
          onSelectLocation={handleSelectLocation}
          files={files}
          onSelectFile={handleSelectFile}
          path={path}
          onGoBack={handleGoBack}
        />
        <SigPlotViewer rawHref={rawHref} sdsHref={sdsHref} />
      </main>
    </div>
  );
}
