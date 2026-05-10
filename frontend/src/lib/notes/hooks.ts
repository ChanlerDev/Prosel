'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';

import { notesApi } from '@/lib/api/notes';
import type { AdminNoteListParams, NoteEditorValues } from '@/types/note';

export const noteKeys = {
  list: (params: Record<string, unknown>) => ['notes', params] as const,
  detail: (slug: string) => ['notes', slug] as const,
  adminList: (params: AdminNoteListParams) => ['admin', 'notes', params] as const,
  adminDetail: (id: string) => ['admin', 'notes', id] as const,
};

export function useAdminNotes(params: AdminNoteListParams) {
  return useQuery({
    queryKey: noteKeys.adminList(params),
    queryFn: () => notesApi.admin.list(params),
  });
}

export function useAdminNote(id: string) {
  return useQuery({
    queryKey: noteKeys.adminDetail(id),
    queryFn: () => notesApi.admin.detail(id),
    enabled: Boolean(id),
  });
}

export function useCreateNote() {
  const router = useRouter();
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (body: NoteEditorValues) => notesApi.admin.create(body),
    onSuccess: async (note) => {
      await queryClient.invalidateQueries({ queryKey: ['admin', 'notes'] });
      await queryClient.invalidateQueries({ queryKey: ['notes'] });
      router.replace(`/admin/notes/${note.id}/edit`);
    },
  });
}

export function useUpdateNote(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (body: NoteEditorValues) => notesApi.admin.update(id, body),
    onSuccess: async (note) => {
      await queryClient.invalidateQueries({ queryKey: ['admin', 'notes'] });
      await queryClient.invalidateQueries({ queryKey: noteKeys.adminDetail(note.id) });
      await queryClient.invalidateQueries({ queryKey: ['notes'] });
    },
  });
}

export function usePinNote() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, pinned }: { id: string; pinned: boolean }) => notesApi.admin.pin(id, pinned),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['admin', 'notes'] });
      await queryClient.invalidateQueries({ queryKey: ['notes'] });
    },
  });
}

export function useDeleteNote() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => notesApi.admin.delete(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['admin', 'notes'] });
      await queryClient.invalidateQueries({ queryKey: ['notes'] });
    },
  });
}
