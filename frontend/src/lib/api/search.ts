import { request } from '@/lib/api/client';
import { readAuthState } from '@/lib/auth/store';
import type { ApiResponse, PaginationMeta } from '@/types/api';
import type { SearchIndexStatus, SearchListResult, SearchParams, SearchResult } from '@/types/search';

function queryString(params: object) {
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== '') search.set(key, String(value));
  }
  const value = search.toString();
  return value ? `?${value}` : '';
}

function emptyMeta(): PaginationMeta {
  return { page: 1, perPage: 20, total: 0, totalPages: 0 };
}

function accessToken() {
  const auth = readAuthState();
  if (!auth.accessToken) throw new Error('Authentication required');
  return auth.accessToken;
}

async function searchWithMeta(params: SearchParams): Promise<SearchListResult> {
  const response = await request<ApiResponse<SearchResult[]>>(`/search${queryString(params)}`, { cache: 'no-store', rawResponse: true });
  return { results: response.data ?? [], meta: response.meta ?? emptyMeta() };
}

export const searchApi = {
  query: searchWithMeta,
  admin: {
    status: () => request<SearchIndexStatus>('/admin/search', { cache: 'no-store', accessToken: accessToken() }),
    rebuild: () => request<SearchIndexStatus>('/admin/search/rebuild', { method: 'POST', body: '{}', cache: 'no-store', accessToken: accessToken() }),
  },
};
