export interface DashboardStats {
  totalPosts: number;
  publishedPosts: number;
  draftPosts: number;
  pendingComments: number;
  todayViews: number;
  totalViews: number;
  categories: number;
  tags: number;
  topics: number;
}

export interface PostSummary {
  id: string;
  title: string;
  slug: string;
  status: string;
  viewCount: number;
  publishedAt?: string;
  updatedAt: string;
}

export interface ActivityLog {
  id: string;
  actorId?: string;
  action: string;
  entityType?: string;
  entityId?: string;
  message?: string;
  createdAt: string;
}

export interface DashboardOverview {
  stats: DashboardStats;
  recentPosts: PostSummary[];
  activities: ActivityLog[];
}
