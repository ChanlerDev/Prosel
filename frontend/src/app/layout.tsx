import type { Metadata } from 'next';
import type { ReactNode } from 'react';

import { AnalyticsTracker } from '@/components/features/analytics/analytics-tracker';
import { QueryProvider } from '@/lib/query/provider';

import './globals.css';

export const metadata: Metadata = {
  title: 'Prosel',
  description: 'A personal blog powered by Prosel',
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body>
        <QueryProvider>
          <AnalyticsTracker />
          {children}
        </QueryProvider>
      </body>
    </html>
  );
}
