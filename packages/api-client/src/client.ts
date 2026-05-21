type AuthHeaderFn = () => string | null;

let getAuthHeader: AuthHeaderFn | null = null;

/**
 * Configure a function that returns a Bearer token for authenticated requests.
 * Pass `null` (default) to send unauthenticated requests.
 */
export function setAuthHeader(fn: AuthHeaderFn) {
  getAuthHeader = fn;
}

const BASE_URL = import.meta.env.VITE_API_URL ?? '/api';

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

export async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  const token = getAuthHeader?.();
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const res = await fetch(`${BASE_URL}${path}`, {
    ...init,
    headers: { ...headers, ...init?.headers },
  });

  if (!res.ok) {
    const body = await res.text();
    throw new ApiError(res.status, body || res.statusText);
  }

  if (res.status === 204) return undefined as T;

  return res.json() as Promise<T>;
}
