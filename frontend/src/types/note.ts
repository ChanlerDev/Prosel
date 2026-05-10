import type { PaginationMeta } from '@/types/api';

export type NoteStatus = 'draft' | 'published' | 'private' | 'archived';

export interface Note {
  id: string;
  authorId: string;
  title?: string;
  slug: string;
  contentMarkdown?: string;
  contentText: string;
  mood?: string;
  weather?: string;
  location?: string;
  status: NoteStatus;
  pinnedAt?: string;
  publishedAt?: string;
  viewCount: number;
  likeCount: number;
  commentCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface NoteListParams {
  page?: number;
  perPage?: number;
  search?: string;
}

export interface AdminNoteListParams extends NoteListParams {
  status?: NoteStatus | '';
}

export interface NoteEditorValues {
  title: string;
  slug: string;
  contentMarkdown: string;
  mood: string;
  weather: string;
  location: string;
  status: NoteStatus;
}

export interface NoteListResult {
  notes: Note[];
  meta: PaginationMeta;
}
