import { Card } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';

export function LoadingState() {
  return (
    <Card className="space-y-3">
      <Skeleton className="h-4 w-24" />
      <Skeleton className="h-8 w-full" />
    </Card>
  );
}

export function EmptyState({ title, description }: { title: string; description: string }) {
  return (
    <Card>
      <h2 className="font-semibold">{title}</h2>
      <p className="mt-2 text-sm text-[var(--muted-foreground)]">{description}</p>
    </Card>
  );
}

export function ApiErrorState({ message }: { message: string }) {
  return (
    <Card className="border-red-300 bg-red-50 text-red-950">
      <h2 className="font-semibold">Something went wrong</h2>
      <p className="mt-2 text-sm">{message}</p>
    </Card>
  );
}
