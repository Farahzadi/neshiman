import { useMutation } from '@tanstack/solid-query';
import { apiFetch } from './client';
import type { definitions } from '@neshiman/api-types';

type LoginRequest = definitions['dto.LoginRequest'];
type LoginResponse = definitions['dto.LoginResponse'];

export function useLogin() {
  return useMutation(() => ({
    mutationFn: (data: LoginRequest) =>
      apiFetch<LoginResponse>('/v1/auth/login', { method: 'POST', body: JSON.stringify(data) }),
  }));
}
