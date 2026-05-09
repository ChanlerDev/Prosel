import { LoginForm } from '@/components/features/auth/login-form';
import { Card } from '@/components/ui/card';

export default function LoginPage() {
  return (
    <main className="flex min-h-screen items-center justify-center px-6">
      <Card className="w-full max-w-md">
        <p className="text-sm text-[var(--muted-foreground)]">Admin login</p>
        <h1 className="mt-2 text-2xl font-semibold">Sign in to Prosel</h1>
        <LoginForm />
      </Card>
    </main>
  );
}
