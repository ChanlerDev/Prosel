import { AdminShell } from '@/components/layout/admin-shell';
import { Card } from '@/components/ui/card';
import { api } from '@/lib/api/client';

export default async function AdminHomePage() {
  const health = await api.system.health().catch(() => null);

  return (
    <AdminShell>
      <div className="max-w-4xl">
        <h1 className="text-3xl font-semibold tracking-tight">Dashboard</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Admin dashboard module will extend this shell.</p>
        <div className="mt-6 grid gap-4 md:grid-cols-2">
          <Card>
            <h2 className="font-semibold">API health</h2>
            <p className="mt-2 text-sm text-[var(--muted-foreground)]">
              {health ? `${health.status} · db ${health.databaseOk ? 'ok' : 'down'} · redis ${health.redisOk ? 'ok' : 'down'}` : 'Backend unavailable'}
            </p>
          </Card>
          <Card>
            <h2 className="font-semibold">Next step</h2>
            <p className="mt-2 text-sm text-[var(--muted-foreground)]">Implement User Auth before protecting admin routes.</p>
          </Card>
        </div>
      </div>
    </AdminShell>
  );
}
