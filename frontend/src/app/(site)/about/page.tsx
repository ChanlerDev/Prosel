import CustomPage from '@/app/(site)/pages/[slug]/page';

export { generateMetadata, revalidate } from '@/app/(site)/pages/[slug]/page';

export default function AboutPage() {
  return <CustomPage params={Promise.resolve({ slug: 'about' })} />;
}
