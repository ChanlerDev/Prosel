import { request } from '@/lib/api/client';
import { readAuthState } from '@/lib/auth/store';
import type { ApiResponse, PaginationMeta } from '@/types/api';
import type { AdminPageListParams, Friend, FriendValues, Page, PageEditorValues, PageListParams, PageListResult } from '@/types/page';

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

async function listWithMeta(path: string, accessTokenValue?: string): Promise<PageListResult> {
  const response = await request<ApiResponse<Page[]>>(path, { cache: 'no-store', accessToken: accessTokenValue, rawResponse: true });
  return { pages: response.data ?? [], meta: response.meta ?? emptyMeta() };
}

export const pagesApi = {
  list: (params: PageListParams = {}) => listWithMeta(`/pages${queryString(params)}`),
  detail: (slug: string) => request<Page>(`/pages/${slug}`, { cache: 'no-store' }),
  friends: () => request<Friend[]>('/friends', { cache: 'no-store' }),
  admin: {
    list: (params: AdminPageListParams = {}) => listWithMeta(`/admin/pages${queryString(params)}`, accessToken()),
    detail: (id: string) => request<Page>(`/admin/pages/${id}`, { cache: 'no-store', accessToken: accessToken() }),
    create: (body: PageEditorValues) => request<Page>('/admin/pages', { method: 'POST', body: JSON.stringify(body), cache: 'no-store', accessToken: accessToken() }),
    update: (id: string, body: PageEditorValues) => request<Page>(`/admin/pages/${id}`, { method: 'PATCH', body: JSON.stringify(body), cache: 'no-store', accessToken: accessToken() }),
    delete: (id: string) => request<{ ok: boolean }>(`/admin/pages/${id}`, { method: 'DELETE', cache: 'no-store', accessToken: accessToken() }),
    friends: (status = '') => request<Friend[]>(`/admin/friends${queryString({ status })}`, { cache: 'no-store', accessToken: accessToken() }),
    createFriend: (body: FriendValues) => request<Friend>('/admin/friends', { method: 'POST', body: JSON.stringify(body), cache: 'no-store', accessToken: accessToken() }),
    updateFriend: (id: string, body: FriendValues) => request<Friend>(`/admin/friends/${id}`, { method: 'PATCH', body: JSON.stringify(body), cache: 'no-store', accessToken: accessToken() }),
    deleteFriend: (id: string) => request<{ ok: boolean }>(`/admin/friends/${id}`, { method: 'DELETE', cache: 'no-store', accessToken: accessToken() }),
  },
};
