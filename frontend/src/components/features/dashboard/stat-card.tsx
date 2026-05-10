import { Card } from '@/components/ui/card';

export function StatCard({ title, value, detail }: { title: string; value: number | string; detail?: string }) {
  return (
    <Card>
      <p className="text-sm text-[var(--muted-foreground)]">{title}</p>
      <p className="mt-3 text-3xl font-semibold tracking-tight">{value}</p>
      {detail ? <p className="mt-2 text-xs text-[var(--muted-foreground)]">{detail}</p> : null}
    </Card>
  );
}
