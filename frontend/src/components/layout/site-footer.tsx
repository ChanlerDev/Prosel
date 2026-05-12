import { RSSLink } from '@/components/features/subscribe/rss-link';
import { SubscribeForm } from '@/components/features/subscribe/subscribe-form';
import { SiteContainer } from '@/components/layout/site-container';

export function SiteFooter() {
  return (
    <footer className="border-t border-[var(--border)] py-8 text-sm text-[var(--muted-foreground)]">
      <SiteContainer>
        <div className="grid gap-5 md:grid-cols-[1fr_2fr] md:items-center">
          <div>
            <p>Powered by Prosel.</p>
            <p className="mt-2"><RSSLink /></p>
          </div>
          <SubscribeForm />
        </div>
      </SiteContainer>
    </footer>
  );
}
