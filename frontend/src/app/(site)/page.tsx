import { Badge } from '@/components/ui/badge';
import { Card } from '@/components/ui/card';
import { SiteContainer } from '@/components/layout/site-container';
import { SiteFooter } from '@/components/layout/site-footer';
import { SiteHeader } from '@/components/layout/site-header';
import { api } from '@/lib/api/client';

export const revalidate = 60;

export default async function HomePage() {
  const settings = await api.settings.public().catch(() => null);
  const siteName = settings?.site_name ?? 'Prosel';
  const description = settings?.site_description ?? 'A personal blog powered by Prosel';

  return (
    <>
      <SiteHeader />
      <main className="py-16">
        <SiteContainer>
          <div className="max-w-3xl">
            <Badge>Project setup</Badge>
            <h1 className="mt-6 text-5xl font-semibold tracking-tight">{siteName}</h1>
            <p className="mt-5 text-lg leading-8 text-[var(--muted-foreground)]">{description}</p>
          </div>
          <div className="mt-10 grid gap-4 md:grid-cols-3">
            <Card>
              <h2 className="font-semibold">Public site</h2>
              <p className="mt-2 text-sm text-[var(--muted-foreground)]">Ready for posts, taxonomy, notes, and pages.</p>
            </Card>
            <Card>
              <h2 className="font-semibold">Admin shell</h2>
              <p className="mt-2 text-sm text-[var(--muted-foreground)]">Prepared for authenticated content management.</p>
            </Card>
            <Card>
              <h2 className="font-semibold">API contract</h2>
              <p className="mt-2 text-sm text-[var(--muted-foreground)]">Connected to backend health and public settings.</p>
            </Card>
          </div>
        </SiteContainer>
      </main>
      <SiteFooter />
    </>
  );
}
