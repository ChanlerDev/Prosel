const publicApiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? 'http://localhost:8080/api/v1';
const serverApiBaseUrl = process.env.API_BASE_URL ?? publicApiBaseUrl;

export const env = {
  apiBaseUrl: typeof window === 'undefined' ? serverApiBaseUrl : publicApiBaseUrl,
};
