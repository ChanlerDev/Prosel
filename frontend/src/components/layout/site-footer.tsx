import { SiteContainer } from '@/components/layout/site-container';

export function SiteFooter() {
  return (
    <footer className="border-t border-[var(--border)] py-8 text-sm text-[var(--muted-foreground)]">
      <SiteContainer>Powered by Prosel.</SiteContainer>
    </footer>
  );
}
