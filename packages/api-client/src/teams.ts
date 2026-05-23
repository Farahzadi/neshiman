import { useQuery, useMutation, useQueryClient } from '@tanstack/solid-query';
import { apiFetch } from './client';
import { queryKeys } from './query-keys';
import type { definitions } from '@neshiman/api-types';

type Team = definitions['dto.TeamResponse'];
type CreateTeam = definitions['dto.CreateTeamRequest'];
type DeleteResponse = definitions['dto.DeleteResponse'];

export function useTeams() {
  return useQuery(() => ({
    queryKey: queryKeys.teams.all,
    queryFn: () => apiFetch<Team[]>('/v1/teams'),
  }));
}

export function useTeam(id: () => string) {
  return useQuery(() => ({
    queryKey: queryKeys.teams.detail(id()),
    queryFn: () => apiFetch<Team>(`/v1/teams/${id()}`),
    enabled: !!id(),
  }));
}

export function useCreateTeam() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: (data: CreateTeam) =>
      apiFetch<Team>('/v1/teams', { method: 'POST', body: JSON.stringify(data) }),
    onSuccess: () => qc.invalidateQueries({ queryKey: queryKeys.teams.all }),
  }));
}

export function useDeleteTeam() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: (id: string) =>
      apiFetch<DeleteResponse>(`/v1/teams/${id}`, { method: 'DELETE' }),
    onSuccess: () => qc.invalidateQueries({ queryKey: queryKeys.teams.all }),
  }));
}
