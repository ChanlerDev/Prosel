import type { PaginationMeta } from '@/types/api';

export type SubscriberStatus = 'pending' | 'active' | 'unsubscribed' | 'bounced';

export interface Subscriber {
  id: string;
  email: string;
  name?: string;
  status: SubscriberStatus;
  verifiedAt?: string;
  unsubscribedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface SubscribeRequest {
  email: string;
  name?: string;
}

export interface SubscriberListParams {
  page?: number;
  perPage?: number;
  status?: SubscriberStatus | '';
  search?: string;
}

export interface SubscriberListResult {
  subscribers: Subscriber[];
  meta: PaginationMeta;
}
