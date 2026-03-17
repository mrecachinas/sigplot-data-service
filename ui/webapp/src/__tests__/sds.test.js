import { describe, it, expect, vi, beforeEach } from 'vitest';
import { getLocations, getFiles, getFileUrl } from '../api/sds';

describe('SDS API', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  describe('getLocations', () => {
    it('returns location names from API', async () => {
      global.fetch = vi.fn(() =>
        Promise.resolve({
          ok: true,
          json: () =>
            Promise.resolve([
              { location_name: 'TestDir' },
              { location_name: 'minio' },
            ]),
        })
      );

      const result = await getLocations();
      expect(result).toEqual(['TestDir', 'minio']);
      expect(fetch).toHaveBeenCalledWith('/sds/fs');
    });

    it('handles camelCase locationName', async () => {
      global.fetch = vi.fn(() =>
        Promise.resolve({
          ok: true,
          json: () => Promise.resolve([{ locationName: 'OldFormat' }]),
        })
      );

      const result = await getLocations();
      expect(result).toEqual(['OldFormat']);
    });

    it('returns empty array on error', async () => {
      global.fetch = vi.fn(() => Promise.reject(new Error('network')));
      const result = await getLocations();
      expect(result).toEqual([]);
    });
  });

  describe('getFiles', () => {
    it('fetches files for a location', async () => {
      const files = [
        { filename: 'test.tmp', type: 'file' },
        { filename: 'subdir', type: 'directory' },
      ];
      global.fetch = vi.fn(() =>
        Promise.resolve({ ok: true, json: () => Promise.resolve(files) })
      );

      const result = await getFiles('TestDir');
      expect(result).toEqual(files);
      expect(fetch).toHaveBeenCalledWith('/sds/fs/TestDir/');
    });

    it('includes path when provided', async () => {
      global.fetch = vi.fn(() =>
        Promise.resolve({ ok: true, json: () => Promise.resolve([]) })
      );

      await getFiles('TestDir', 'subdir/nested');
      expect(fetch).toHaveBeenCalledWith('/sds/fs/TestDir/subdir/nested');
    });

    it('returns empty array on error', async () => {
      global.fetch = vi.fn(() => Promise.reject(new Error('fail')));
      const result = await getFiles('TestDir');
      expect(result).toEqual([]);
    });
  });

  describe('getFileUrl', () => {
    it('builds correct URL for fs mode', () => {
      expect(getFileUrl('data.tmp', 'fs', 'TestDir')).toBe(
        '/sds/fs/TestDir/data.tmp'
      );
    });

    it('builds correct URL for hdr mode', () => {
      expect(getFileUrl('data.tmp', 'hdr', 'TestDir')).toBe(
        '/sds/hdr/TestDir/data.tmp'
      );
    });
  });
});
