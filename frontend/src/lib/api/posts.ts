import { readAuthState } from '@/lib/auth/store';
import { request } from '@/lib/api/client';
import type { AdminPostListParams, Post, PostEditorValues, PostListParams, PostListResult } from '@/types/post';
import type { ApiResponse, PaginationMeta } from '@/types/api';

function queryString(params: object) {
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== '') {
      search.set(key, String(value));
    }
  }
  const value = search.toString();
  return value ? `?${value}` : '';
}

async function listWithMeta(path: string, accessToken?: string): Promise<PostListResult> {
  const response = await request<ApiResponse<Post[]>>(path, { cache: 'no-store', accessToken, rawResponse: true });
  return { posts: response.data ?? [], meta: response.meta ?? emptyMeta() };
}

function emptyMeta(): PaginationMeta {
  return { page: 1, perPage: 20, total: 0, totalPages: 0 };
}

function accessToken() {
  const auth = readAuthState();
  if (!auth.accessToken) {
    throw new Error('Authentication required');
  }
  return auth.accessToken;
}

export const postsApi = {
  list: (params: PostListParams = {}) => listWithMeta(`/posts${queryString(params)}`),
  detail: (slug: string) => request<Post>(`/posts/${slug}`, { cache: 'no-store' }),
  admin: {
    list: (params: AdminPostListParams = {}) => listWithMeta(`/admin/posts${queryString(params)}`, accessToken()),
    detail: (id: string) => request<Post>(`/admin/posts/${id}`, { cache: 'no-store', accessToken: accessToken() }),
    create: (body: PostEditorValues) => request<Post>('/admin/posts', { method: 'POST', body: JSON.stringify(body), cache: 'no-store', accessToken: accessToken() }),
    update: (id: string, body: PostEditorValues) => request<Post>(`/admin/posts/${id}`, { method: 'PATCH', body: JSON.stringify(body), cache: 'no-store', accessToken: accessToken() }),
    delete: (id: string) => request<{ ok: boolean }>(`/admin/posts/${id}`, { method: 'DELETE', cache: 'no-store', accessToken: accessToken() }),
    publish: (id: string) => request<Post>(`/admin/posts/${id}/publish`, { method: 'POST', body: '{}', cache: 'no-store', accessToken: accessToken() }),
    unpublish: (id: string) => request<Post>(`/admin/posts/${id}/unpublish`, { method: 'POST', body: '{}', cache: 'no-store', accessToken: accessToken() }),
  },
};
