import { request } from '@/lib/api/client';
import { readAuthState } from '@/lib/auth/store';
import type { ApiResponse, PaginationMeta } from '@/types/api';
import type { Post } from '@/types/post';
import type { CategoryNode, CategoryValues, Tag, TagValues, TaxonomyPostsResult, Topic, TopicValues } from '@/types/taxonomy';

function queryString(params: object) {
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== '') search.set(key, String(value));
  }
  const value = search.toString();
  return value ? `?${value}` : '';
}

function accessToken() {
  const auth = readAuthState();
  if (!auth.accessToken) throw new Error('Authentication required');
  return auth.accessToken;
}

function emptyMeta(): PaginationMeta {
  return { page: 1, perPage: 20, total: 0, totalPages: 0 };
}

async function postList(path: string): Promise<TaxonomyPostsResult> {
  const response = await request<ApiResponse<Post[]>>(path, { cache: 'no-store', rawResponse: true });
  return { posts: response.data ?? [], meta: response.meta ?? emptyMeta() };
}

export const taxonomyApi = {
  categories: () => request<CategoryNode[]>('/categories'),
  categoryPosts: (slug: string, params: { page?: number; perPage?: number } = {}) => postList(`/categories/${slug}/posts${queryString(params)}`),
  tags: () => request<Tag[]>('/tags'),
  tagPosts: (slug: string, params: { page?: number; perPage?: number } = {}) => postList(`/tags/${slug}/posts${queryString(params)}`),
  topics: () => request<Topic[]>('/topics'),
  topic: (slug: string) => request<Topic>(`/topics/${slug}`, { cache: 'no-store' }),
  admin: {
    createCategory: (body: CategoryValues) => request<CategoryNode>('/admin/categories', { method: 'POST', body: JSON.stringify(body), cache: 'no-store', accessToken: accessToken() }),
    updateCategory: (id: string, body: CategoryValues) => request<CategoryNode>(`/admin/categories/${id}`, { method: 'PATCH', body: JSON.stringify(body), cache: 'no-store', accessToken: accessToken() }),
    deleteCategory: (id: string) => request<{ ok: boolean }>(`/admin/categories/${id}`, { method: 'DELETE', cache: 'no-store', accessToken: accessToken() }),
    createTag: (body: TagValues) => request<Tag>('/admin/tags', { method: 'POST', body: JSON.stringify(body), cache: 'no-store', accessToken: accessToken() }),
    updateTag: (id: string, body: TagValues) => request<Tag>(`/admin/tags/${id}`, { method: 'PATCH', body: JSON.stringify(body), cache: 'no-store', accessToken: accessToken() }),
    deleteTag: (id: string) => request<{ ok: boolean }>(`/admin/tags/${id}`, { method: 'DELETE', cache: 'no-store', accessToken: accessToken() }),
    createTopic: (body: TopicValues) => request<Topic>('/admin/topics', { method: 'POST', body: JSON.stringify(body), cache: 'no-store', accessToken: accessToken() }),
    updateTopic: (id: string, body: TopicValues) => request<Topic>(`/admin/topics/${id}`, { method: 'PATCH', body: JSON.stringify(body), cache: 'no-store', accessToken: accessToken() }),
    deleteTopic: (id: string) => request<{ ok: boolean }>(`/admin/topics/${id}`, { method: 'DELETE', cache: 'no-store', accessToken: accessToken() }),
  },
};
