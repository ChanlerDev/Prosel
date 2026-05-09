import type { PostListResult } from '@/types/post';

export interface CategoryNode {
  id: string;
  parentId?: string;
  name: string;
  slug: string;
  description?: string;
  sortOrder: number;
  postCount: number;
  children: CategoryNode[];
  createdAt: string;
  updatedAt: string;
}

export interface Tag {
  id: string;
  name: string;
  slug: string;
  color?: string;
  description?: string;
  postCount?: number;
  createdAt: string;
  updatedAt: string;
}

export interface TopicItem {
  refType: 'post' | 'note';
  refId: string;
  title?: string;
  slug?: string;
  sortOrder: number;
}

export interface Topic {
  id: string;
  name: string;
  slug: string;
  description?: string;
  coverImage?: string;
  sortOrder: number;
  items?: TopicItem[];
  createdAt: string;
  updatedAt: string;
}

export interface CategoryValues {
  parentId: string;
  name: string;
  slug: string;
  description: string;
  sortOrder: number;
}

export interface TagValues {
  name: string;
  slug: string;
  color: string;
  description: string;
}

export interface TopicValues {
  name: string;
  slug: string;
  description: string;
  coverImage: string;
  sortOrder: number;
  items: TopicItem[];
}

export type TaxonomyPostsResult = PostListResult;
