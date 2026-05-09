import type { PaginationMeta } from '@/types/api';

export type PostStatus = 'draft' | 'published' | 'archived';

export interface Post {
  id: string;
  authorId: string;
  categoryId?: string;
  tagIds?: string[];
  title: string;
  slug: string;
  excerpt?: string;
  contentMarkdown?: string;
  contentText?: string;
  coverImage?: string;
  status: PostStatus;
  featured: boolean;
  pinnedAt?: string;
  publishedAt?: string;
  seoTitle?: string;
  seoDescription?: string;
  viewCount: number;
  likeCount: number;
  commentCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface PostListParams {
  page?: number;
  perPage?: number;
  category?: string;
  featured?: boolean;
}

export interface AdminPostListParams {
  page?: number;
  perPage?: number;
  status?: PostStatus | '';
  search?: string;
}

export interface PostEditorValues {
  title: string;
  slug: string;
  excerpt: string;
  contentMarkdown: string;
  coverImage: string;
  categoryId: string;
  tagIds: string[];
  featured: boolean;
  seoTitle: string;
  seoDescription: string;
}

export interface PostListResult {
  posts: Post[];
  meta: PaginationMeta;
}
