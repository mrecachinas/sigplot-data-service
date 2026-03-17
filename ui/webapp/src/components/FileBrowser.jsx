import React from 'react';

export default function FileBrowser({
  locations,
  selectedLocation,
  onSelectLocation,
  files,
  onSelectFile,
  path,
  onGoBack,
}) {
  return (
    <>
      <div className="location-panel">
        <h4>Choose Location:</h4>
        <ul className="list-group">
          {locations.map((loc) => (
            <button
              key={loc}
              className={`list-group-item${loc === selectedLocation ? ' active' : ''}`}
              onClick={() => onSelectLocation(loc)}
            >
              {loc}
            </button>
          ))}
          {locations.length === 0 && (
            <p className="empty-message">Loading locations...</p>
          )}
        </ul>
      </div>

      <div className="file-browser">
        <div className="file-browser-header">
          <h4>Choose File:</h4>
          {path && (
            <button className="back-btn" onClick={onGoBack}>
              Back
            </button>
          )}
        </div>
        {selectedLocation ? (
          <ul className="list-group">
            {files.map((file) => (
              <button
                key={file.filename}
                className="list-group-item"
                onClick={() => onSelectFile(file)}
              >
                <p className={`file-name ${file.type}`}>
                  {file.filename}
                </p>
              </button>
            ))}
            {files.length === 0 && (
              <p className="empty-message">No files found</p>
            )}
          </ul>
        ) : (
          <p className="empty-message">Select a location first</p>
        )}
      </div>
    </>
  );
}
