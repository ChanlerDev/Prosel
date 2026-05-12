import type { PaginationMeta } from '@/types/api';

export type FileStatus = 'attached' | 'orphan' | 'deleted';

export interface FileAsset {
  id: string;
  uploaderId?: string;
  originalName: string;
  fileName: string;
  storageType: 'local' | 's3';
  objectKey: string;
  publicUrl: string;
  mimeType: string;
  byteSize: number;
  width?: number;
  height?: number;
  refType?: string;
  refId?: string;
  status: FileStatus;
  createdAt: string;
  updatedAt: string;
}

export interface AdminFileListParams {
  page?: number;
  perPage?: number;
  type?: string;
  status?: FileStatus | '';
  search?: string;
}

export interface FileListResult {
  files: FileAsset[];
  meta: PaginationMeta;
}
