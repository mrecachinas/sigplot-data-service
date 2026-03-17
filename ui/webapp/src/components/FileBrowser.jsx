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
    <div className="file-browser">
      <h3>Locations</h3>
      <select
        value={selectedLocation}
        onChange={(e) => onSelectLocation(e.target.value)}
      >
        <option value="">-- Select Location --</option>
        {locations.map((loc) => (
          <option key={loc} value={loc}>
            {loc}
          </option>
        ))}
      </select>

      {selectedLocation && (
        <>
          <div className="path-bar">
            <h3>Files</h3>
            {path && (
              <button className="back-btn" onClick={onGoBack}>
                ← Back
              </button>
            )}
            {path && <span className="current-path">/{path}</span>}
          </div>
          <ul className="file-list">
            {files.map((file) => (
              <li
                key={file.filename}
                className={`file-item ${file.type}`}
                onClick={() => onSelectFile(file)}
              >
                <span className="file-icon">
                  {file.type === 'directory' ? '📁' : '📄'}
                </span>
                {file.filename}
              </li>
            ))}
            {files.length === 0 && (
              <li className="file-item empty">No files found</li>
            )}
          </ul>
        </>
      )}
    </div>
  );
}
