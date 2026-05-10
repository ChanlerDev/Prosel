import { Card } from '@/components/ui/card';
import type { ActivityLog } from '@/types/dashboard';

export function ActivityTimeline({ activities }: { activities: ActivityLog[] }) {
  if (activities.length === 0) {
    return <Card><p className="text-sm text-[var(--muted-foreground)]">No activity yet.</p></Card>;
  }
  return (
    <Card>
      <h2 className="font-semibold">Activity</h2>
      <div className="mt-4 grid gap-4">
        {activities.map((activity) => (
          <div className="border-l-2 border-[var(--border)] pl-4" key={activity.id}>
            <p className="text-sm font-medium">{activity.message || activity.action}</p>
            <p className="mt-1 text-xs text-[var(--muted-foreground)]">{activity.entityType} · {formatDate(activity.createdAt)}</p>
          </div>
        ))}
      </div>
    </Card>
  );
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat('en', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value));
}
