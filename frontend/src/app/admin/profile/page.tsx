import { ProfileForm } from '@/components/features/auth/profile-form';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AdminProfilePage() {
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">Profile</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Update the public admin profile used by content modules.</p>
      </div>
      <ProfileForm />
    </AdminShell>
  );
}
