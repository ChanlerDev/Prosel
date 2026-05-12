import { request } from '@/lib/api/client';
import { readAuthState } from '@/lib/auth/store';
import type { ApiResponse, PaginationMeta } from '@/types/api';
import type { SubscribeRequest, Subscriber, SubscriberListParams, SubscriberListResult } from '@/types/subscribe';

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

async function listWithMeta(path: string): Promise<SubscriberListResult> {
  const response = await request<ApiResponse<Subscriber[]>>(path, { cache: 'no-store', accessToken: accessToken(), rawResponse: true });
  return { subscribers: response.data ?? [], meta: response.meta ?? emptyMeta() };
}

export const subscribeApi = {
  create: (body: SubscribeRequest) => request<Subscriber>('/subscribe', { method: 'POST', body: JSON.stringify(body), cache: 'no-store' }),
  verify: (token: string) => request<{ ok: boolean }>(`/subscribe/verify${queryString({ token })}`, { cache: 'no-store' }),
  unsubscribe: (token: string) => request<{ ok: boolean }>(`/subscribe/unsubscribe${queryString({ token })}`, { cache: 'no-store' }),
  admin: {
    list: (params: SubscriberListParams = {}) => listWithMeta(`/admin/subscribers${queryString(params)}`),
    notifyPost: (postId: string) => request<{ ok: boolean }>('/admin/subscribers/notify-post', { method: 'POST', body: JSON.stringify({ postId }), cache: 'no-store', accessToken: accessToken() }),
  },
};
