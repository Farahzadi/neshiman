import { useQuery, useMutation, useQueryClient } from '@tanstack/solid-query';
import { apiFetch } from './client';
import { queryKeys } from './query-keys';
import type { definitions } from '@neshiman/api-types';

type CrossTeamRequest = definitions['dto.CrossTeamRequestResponse'];
type CreateCrossTeamRequest = definitions['dto.CreateCrossTeamRequestRequest'];

export function useCrossTeamRequests(status: () => string) {
  return useQuery(() => ({
    queryKey: queryKeys.crossTeamRequests.all,
    queryFn: () => apiFetch<CrossTeamRequest[]>(`/v1/cross-team-requests?status=${status()}`),
    enabled: !!status(),
  }));
}

export function usePendingByTeam(teamId: () => string) {
  return useQuery(() => ({
    queryKey: queryKeys.crossTeamRequests.pendingByTeam(teamId()),
    queryFn: () =>
      apiFetch<CrossTeamRequest[]>(`/v1/cross-team-requests/pending-by-team?team_id=${teamId()}`),
    enabled: !!teamId(),
  }));
}

export function useCrossTeamRequest(id: () => string) {
  return useQuery(() => ({
    queryKey: queryKeys.crossTeamRequests.detail(id()),
    queryFn: () => apiFetch<CrossTeamRequest>(`/v1/cross-team-requests/${id()}`),
    enabled: !!id(),
  }));
}

export function useCreateCrossTeamRequest() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: (data: CreateCrossTeamRequest) =>
      apiFetch<CrossTeamRequest>('/v1/cross-team-requests', {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    onSuccess: () => qc.invalidateQueries({ queryKey: queryKeys.crossTeamRequests.all }),
  }));
}

export function useApproveCrossTeamRequest() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: (id: string) =>
      apiFetch<CrossTeamRequest>(`/v1/cross-team-requests/${id}/approve`, { method: 'PUT' }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.crossTeamRequests.all });
    },
  }));
}

export function useRejectCrossTeamRequest() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: (id: string) =>
      apiFetch<CrossTeamRequest>(`/v1/cross-team-requests/${id}/reject`, { method: 'PUT' }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.crossTeamRequests.all });
    },
  }));
}
