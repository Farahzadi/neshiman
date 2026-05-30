import { useQuery, useMutation, useQueryClient } from '@tanstack/solid-query';
import { apiFetch } from './client';
import { queryKeys } from './query-keys';
import type { definitions } from '@neshiman/api-types';

type User = definitions['dto.UserResponse'];
type CreateUser = definitions['dto.CreateUserRequest'];
type UpdateWeeklyLimit = definitions['dto.UpdateWeeklyLimitRequest'];
type DeleteResponse = definitions['dto.DeleteResponse'];

export function useUsers(teamId: () => string) {
  return useQuery(() => ({
    queryKey: teamId() ? queryKeys.users.byTeam(teamId()) : queryKeys.users.all,
    queryFn: () => {
      const query = teamId() ? `?team_id=${teamId()}` : '';
      return apiFetch<User[]>(`/v1/users${query}`);
    },
  }));
}

export function useUser(id: () => string) {
  return useQuery(() => ({
    queryKey: queryKeys.users.detail(id()),
    queryFn: () => apiFetch<User>(`/v1/users/${id()}`),
    enabled: !!id(),
  }));
}

export function useUserByEmail(email: () => string) {
  return useQuery(() => ({
    queryKey: queryKeys.users.byEmail(email()),
    queryFn: () => apiFetch<User>(`/v1/users/by-email?email=${encodeURIComponent(email())}`),
    enabled: !!email(),
  }));
}

export function useCreateUser() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: (data: CreateUser) =>
      apiFetch<User>('/v1/users', { method: 'POST', body: JSON.stringify(data) }),
    onSuccess: () => qc.invalidateQueries({ queryKey: queryKeys.users.all }),
  }));
}

export function useDeleteUser() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: (id: string) =>
      apiFetch<DeleteResponse>(`/v1/users/${id}`, { method: 'DELETE' }),
    onSuccess: () => qc.invalidateQueries({ queryKey: queryKeys.users.all }),
  }));
}

export function useUpdateWeeklyLimit() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: ({ id, data }: { id: string; data: UpdateWeeklyLimit }) =>
      apiFetch<void>(`/v1/users/${id}/weekly-limit`, { method: 'PUT', body: JSON.stringify(data) }),
    onSuccess: (_data, vars) => {
      qc.invalidateQueries({ queryKey: queryKeys.users.detail(vars.id) });
      qc.invalidateQueries({ queryKey: queryKeys.users.all });
    },
  }));
}

type UpdateUserRole = definitions['dto.UpdateUserRoleRequest'];
type UpdateUserTeam = definitions['dto.UpdateUserTeamRequest'];

export function useUpdateUserRole() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: ({ id, data }: { id: string; data: UpdateUserRole }) =>
      apiFetch<User>(`/v1/users/${id}/role`, { method: 'PUT', body: JSON.stringify(data) }),
    onSuccess: (_data, vars) => {
      qc.invalidateQueries({ queryKey: queryKeys.users.detail(vars.id) });
      qc.invalidateQueries({ queryKey: queryKeys.users.all });
    },
  }));
}

export function useUpdateUserTeam() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: ({ id, data }: { id: string; data: UpdateUserTeam }) =>
      apiFetch<User>(`/v1/users/${id}/team`, { method: 'PUT', body: JSON.stringify(data) }),
    onSuccess: (_data, vars) => {
      qc.invalidateQueries({ queryKey: queryKeys.users.detail(vars.id) });
      qc.invalidateQueries({ queryKey: queryKeys.users.all });
    },
  }));
}

export function useSetPassword() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: ({ id, password }: { id: string; password: string }) =>
      apiFetch<void>(`/v1/users/${id}/password`, {
        method: 'PUT',
        body: JSON.stringify({ password }),
      }),
    onSuccess: (_data, vars) => {
      qc.invalidateQueries({ queryKey: queryKeys.users.detail(vars.id) });
    },
  }));
}
