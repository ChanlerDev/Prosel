import { PasswordForm } from '@/components/features/auth/password-form';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AdminSecurityPage() {
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">Security</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Change the admin password. Existing sessions will be revoked.</p>
      </div>
      <PasswordForm />
    </AdminShell>
  );
}
