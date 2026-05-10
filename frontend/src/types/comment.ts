import type { PaginationMeta } from '@/types/api';

export type CommentRefType = 'post' | 'note' | 'page';
export type CommentStatus = 'pending' | 'approved' | 'rejected' | 'spam';

export interface CommentNode {
  id: string;
  refType: CommentRefType;
  refId: string;
  parentId?: string;
  rootId?: string;
  authorName: string;
  authorEmail?: string;
  authorWebsite?: string;
  content: string;
  status: CommentStatus;
  isAdminReply: boolean;
  isPinned: boolean;
  replyCount: number;
  children?: CommentNode[];
  createdAt: string;
  updatedAt: string;
}

export interface SubmitCommentValues {
  refType: CommentRefType;
  refId: string;
  parentId?: string;
  authorName: string;
  authorEmail: string;
  authorWebsite?: string;
  content: string;
}

export interface AdminCommentListParams {
  page?: number;
  perPage?: number;
  status?: CommentStatus | '';
  refType?: CommentRefType | '';
  refId?: string;
  search?: string;
}

export interface CommentListResult {
  comments: CommentNode[];
  meta: PaginationMeta;
}
