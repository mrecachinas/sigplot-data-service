import React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render, fireEvent } from '@testing-library/react';
import FileBrowser from '../components/FileBrowser';

describe('FileBrowser', () => {
  const defaultProps = {
    locations: ['TestDir', 'minio'],
    selectedLocation: '',
    onSelectLocation: vi.fn(),
    files: [],
    onSelectFile: vi.fn(),
    path: '',
    onGoBack: vi.fn(),
  };

  it('renders location list', () => {
    const { getByText } = render(<FileBrowser {...defaultProps} />);
    expect(getByText('TestDir')).toBeTruthy();
    expect(getByText('minio')).toBeTruthy();
  });

  it('calls onSelectLocation when location clicked', () => {
    const onSelect = vi.fn();
    const { getByText } = render(
      <FileBrowser {...defaultProps} onSelectLocation={onSelect} />
    );
    fireEvent.click(getByText('TestDir'));
    expect(onSelect).toHaveBeenCalledWith('TestDir');
  });

  it('shows file list when location selected', () => {
    const files = [
      { filename: 'data.tmp', type: 'file' },
      { filename: 'subdir', type: 'directory' },
    ];
    const { getByText } = render(
      <FileBrowser {...defaultProps} selectedLocation="TestDir" files={files} />
    );
    expect(getByText('data.tmp')).toBeTruthy();
    expect(getByText('subdir')).toBeTruthy();
  });

  it('calls onSelectFile when file clicked', () => {
    const onSelect = vi.fn();
    const file = { filename: 'data.tmp', type: 'file' };
    const { getByText } = render(
      <FileBrowser
        {...defaultProps}
        selectedLocation="TestDir"
        files={[file]}
        onSelectFile={onSelect}
      />
    );
    fireEvent.click(getByText('data.tmp'));
    expect(onSelect).toHaveBeenCalledWith(file);
  });

  it('shows back button when path is set', () => {
    const { getByText } = render(
      <FileBrowser {...defaultProps} selectedLocation="TestDir" path="subdir" />
    );
    expect(getByText('Back')).toBeTruthy();
  });

  it('hides back button when at root', () => {
    const { queryByText } = render(
      <FileBrowser {...defaultProps} selectedLocation="TestDir" path="" />
    );
    expect(queryByText('Back')).toBeNull();
  });
});
