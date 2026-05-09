import { Badge } from '@/components/ui/badge';
import type { PostStatus } from '@/types/post';

const labels: Record<PostStatus, string> = {
  draft: 'Draft',
  published: 'Published',
  archived: 'Archived',
};

export function PostStatusBadge({ status }: { status: PostStatus }) {
  return <Badge>{labels[status]}</Badge>;
}
