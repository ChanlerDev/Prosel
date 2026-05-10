import { request } from '@/lib/api/client';
import { readAuthState } from '@/lib/auth/store';
import type { ApiResponse, PaginationMeta } from '@/types/api';
import type { AdminCommentListParams, CommentListResult, CommentNode, CommentRefType, CommentStatus, SubmitCommentValues } from '@/types/comment';

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

async function listAdminWithMeta(path: string): Promise<CommentListResult> {
  const response = await request<ApiResponse<CommentNode[]>>(path, { cache: 'no-store', accessToken: accessToken(), rawResponse: true });
  return { comments: response.data ?? [], meta: response.meta ?? emptyMeta() };
}

export const commentsApi = {
  list: (refType: CommentRefType, refId: string) => request<CommentNode[]>(`/comments${queryString({ refType, refId })}`, { cache: 'no-store' }),
  submit: (body: SubmitCommentValues) => request<CommentNode>('/comments', { method: 'POST', body: JSON.stringify(body), cache: 'no-store' }),
  admin: {
    list: (params: AdminCommentListParams = {}) => listAdminWithMeta(`/admin/comments${queryString(params)}`),
    moderate: (id: string, status: CommentStatus) => request<{ ok: boolean }>(`/admin/comments/${id}/status`, { method: 'PATCH', body: JSON.stringify({ status }), cache: 'no-store', accessToken: accessToken() }),
    reply: (id: string, content: string) => request<CommentNode>(`/admin/comments/${id}/reply`, { method: 'POST', body: JSON.stringify({ content }), cache: 'no-store', accessToken: accessToken() }),
    delete: (id: string) => request<{ ok: boolean }>(`/admin/comments/${id}`, { method: 'DELETE', cache: 'no-store', accessToken: accessToken() }),
  },
};
