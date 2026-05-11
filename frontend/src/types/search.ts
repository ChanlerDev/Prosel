import type { PaginationMeta } from '@/types/api';

export type SearchRefType = 'post' | 'note' | 'page';

export interface SearchResult {
  refType: SearchRefType;
  refId: string;
  title: string;
  slug?: string;
  excerpt?: string;
  rank: number;
}

export interface SearchParams {
  q: string;
  type?: SearchRefType | '';
  page?: number;
  perPage?: number;
}

export interface SearchListResult {
  results: SearchResult[];
  meta: PaginationMeta;
}

export interface SearchIndexStatus {
  total: number;
  posts: number;
  notes: number;
  pages: number;
  updatedAt?: string;
}
