import type { Search, SearchFormData, Property } from '../types';

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = localStorage.getItem('token');
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...((options.headers as Record<string, string>) || {}),
  };
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const res = await fetch(path, { ...options, headers });

  if (res.status === 401) {
    localStorage.removeItem('token');
    window.location.href = '/login';
    throw new Error('Unauthorized');
  }

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error((body as { error?: string }).error || `Request failed: ${res.status}`);
  }

  if (res.status === 204) return undefined as T;
  return res.json();
}

export async function login(username: string, password: string): Promise<string> {
  const data = await request<{ token: string }>('/api/auth/login', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  });
  return data.token;
}

export async function getSearches(): Promise<Search[]> {
  return request<Search[]>('/api/searches');
}

export async function createSearch(data: SearchFormData): Promise<Search> {
  return request<Search>('/api/searches', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function updateSearch(id: string, data: SearchFormData): Promise<Search> {
  return request<Search>(`/api/searches/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export async function deleteSearch(id: string): Promise<void> {
  return request<void>(`/api/searches/${id}`, { method: 'DELETE' });
}

export async function getSearch(id: string): Promise<Search> {
  return request<Search>(`/api/searches/${id}`);
}

export async function getProperties(
  searchId: string,
  params: { limit?: number; offset?: number } = {}
): Promise<Property[]> {
  const query = new URLSearchParams();
  if (params.limit) query.set('limit', String(params.limit));
  if (params.offset) query.set('offset', String(params.offset));
  const qs = query.toString();
  return request<Property[]>(`/api/searches/${searchId}/properties${qs ? `?${qs}` : ''}`);
}

export async function triggerScrape(searchId: string): Promise<{ status: string }> {
  return request<{ status: string }>(`/api/searches/${searchId}/scrape`, { method: 'POST' });
}
