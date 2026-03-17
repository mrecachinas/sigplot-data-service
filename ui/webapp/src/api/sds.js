const SDS_BASE = '/sds';

export async function getLocations() {
  try {
    const res = await fetch(`${SDS_BASE}/fs`);
    if (!res.ok) return [];
    const data = await res.json();
    return data.map((loc) => loc.location_name || loc.locationName);
  } catch {
    return [];
  }
}

export async function getFiles(location, path = '') {
  try {
    const url = path
      ? `${SDS_BASE}/fs/${location}/${path}`
      : `${SDS_BASE}/fs/${location}/`;
    const res = await fetch(url);
    if (!res.ok) return [];
    return await res.json();
  } catch {
    return [];
  }
}

export function getFileUrl(file, mode, location) {
  return `${SDS_BASE}/${mode}/${location}/${file}`;
}
