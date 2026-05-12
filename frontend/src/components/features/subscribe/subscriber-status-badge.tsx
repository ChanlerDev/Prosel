import { Badge } from '@/components/ui/badge';
import type { SubscriberStatus } from '@/types/subscribe';

export function SubscriberStatusBadge({ status }: { status: SubscriberStatus }) {
  const color = status === 'active' ? 'bg-emerald-100 text-emerald-800' : status === 'pending' ? 'bg-amber-100 text-amber-900' : status === 'unsubscribed' ? 'bg-zinc-100 text-zinc-700' : 'bg-red-100 text-red-800';
  return <Badge className={color}>{status}</Badge>;
}
