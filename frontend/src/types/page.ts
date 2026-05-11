import type { PaginationMeta } from '@/types/api';

export type PageTemplate = 'default' | 'about' | 'friends' | 'projects';
export type PageStatus = 'draft' | 'published' | 'archived';
export type FriendStatus = 'active' | 'pending' | 'hidden';

export interface Page {
  id: string;
  authorId: string;
  title: string;
  slug: string;
  subtitle?: string;
  contentMarkdown?: string;
  contentText: string;
  template: PageTemplate;
  status: PageStatus;
  sortOrder: number;
  seoTitle?: string;
  seoDescription?: string;
  viewCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface Friend {
  id: string;
  name: string;
  url: string;
  avatarUrl?: string;
  description?: string;
  status: FriendStatus;
  sortOrder: number;
  createdAt: string;
  updatedAt: string;
}

export interface PageListParams {
  page?: number;
  perPage?: number;
  search?: string;
}

export interface AdminPageListParams extends PageListParams {
  status?: PageStatus | '';
}

export interface PageEditorValues {
  title: string;
  slug: string;
  subtitle: string;
  contentMarkdown: string;
  template: PageTemplate;
  status: PageStatus;
  sortOrder: number;
  seoTitle: string;
  seoDescription: string;
}

export interface FriendValues {
  name: string;
  url: string;
  avatarUrl: string;
  description: string;
  status: FriendStatus;
  sortOrder: number;
}

export interface PageListResult {
  pages: Page[];
  meta: PaginationMeta;
}
