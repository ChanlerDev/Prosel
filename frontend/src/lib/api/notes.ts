import { request } from '@/lib/api/client';
import { readAuthState } from '@/lib/auth/store';
import type { ApiResponse, PaginationMeta } from '@/types/api';
import type { AdminNoteListParams, Note, NoteEditorValues, NoteListParams, NoteListResult } from '@/types/note';

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

async function listWithMeta(path: string, accessTokenValue?: string): Promise<NoteListResult> {
  const response = await request<ApiResponse<Note[]>>(path, { cache: 'no-store', accessToken: accessTokenValue, rawResponse: true });
  return { notes: response.data ?? [], meta: response.meta ?? emptyMeta() };
}

export const notesApi = {
  list: (params: NoteListParams = {}) => listWithMeta(`/notes${queryString(params)}`),
  detail: (slug: string) => request<Note>(`/notes/${slug}`, { cache: 'no-store' }),
  admin: {
    list: (params: AdminNoteListParams = {}) => listWithMeta(`/admin/notes${queryString(params)}`, accessToken()),
    detail: (id: string) => request<Note>(`/admin/notes/${id}`, { cache: 'no-store', accessToken: accessToken() }),
    create: (body: NoteEditorValues) => request<Note>('/admin/notes', { method: 'POST', body: JSON.stringify(body), cache: 'no-store', accessToken: accessToken() }),
    update: (id: string, body: NoteEditorValues) => request<Note>(`/admin/notes/${id}`, { method: 'PATCH', body: JSON.stringify(body), cache: 'no-store', accessToken: accessToken() }),
    pin: (id: string, pinned: boolean) => request<{ ok: boolean }>(`/admin/notes/${id}/pin`, { method: 'POST', body: JSON.stringify({ pinned }), cache: 'no-store', accessToken: accessToken() }),
    delete: (id: string) => request<{ ok: boolean }>(`/admin/notes/${id}`, { method: 'DELETE', cache: 'no-store', accessToken: accessToken() }),
  },
};
