import { useQuery, useMutation, useQueryClient } from '@tanstack/solid-query';
import { apiFetch } from './client';
import { queryKeys } from './query-keys';
import type { definitions } from '@neshiman/api-types';

type Reservation = definitions['dto.ReservationResponse'];
type CreateReservation = definitions['dto.CreateReservationRequest'];

export function useReservationsByDate(date: () => string) {
  return useQuery(() => ({
    queryKey: queryKeys.reservations.byDate(date()),
    queryFn: () => apiFetch<Reservation[]>(`/v1/reservations?date=${date()}`),
    enabled: !!date(),
  }));
}

export function useReservationsByUserDate(userId: () => string, date: () => string) {
  return useQuery(() => ({
    queryKey: queryKeys.reservations.byUserDate(userId(), date()),
    queryFn: () => apiFetch<Reservation[]>(`/v1/reservations?user_id=${userId()}&date=${date()}`),
    enabled: !!userId() && !!date(),
  }));
}

export function useCreateReservation() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: (data: CreateReservation) =>
      apiFetch<Reservation>('/v1/reservations', { method: 'POST', body: JSON.stringify(data) }),
    onSuccess: () => qc.invalidateQueries({ queryKey: queryKeys.reservations.all }),
  }));
}
