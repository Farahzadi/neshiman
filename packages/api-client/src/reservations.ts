import { useQuery, useMutation, useQueryClient } from '@tanstack/solid-query';
import { apiFetch } from './client';
import { queryKeys } from './query-keys';
import type { definitions } from '@neshiman/api-types';

type Reservation = definitions['dto.ReservationResponse'];
type CreateReservation = definitions['dto.CreateReservationRequest'];
type DeleteResponse = definitions['dto.DeleteResponse'];

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

type ReservationWithUser = definitions['dto.ReservationWithUserResponse'];
type AdminCreateReservation = definitions['dto.AdminCreateReservationRequest'];

export function useReservationsByRoomDate(roomId: () => string, date: () => string) {
  return useQuery(() => ({
    queryKey: ['reservations', 'by-room', roomId(), date()] as const,
    queryFn: () =>
      apiFetch<ReservationWithUser[]>(`/v1/reservations/by-room?room_id=${roomId()}&date=${date()}`),
    enabled: !!roomId() && !!date(),
  }));
}

export function useAdminCreateReservation() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: (data: AdminCreateReservation) =>
      apiFetch<Reservation>('/v1/reservations/admin', { method: 'POST', body: JSON.stringify(data) }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['reservations'] }),
  }));
}

export function useCancelReservation() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: (id: string) =>
      apiFetch<DeleteResponse>(`/v1/reservations/${id}`, { method: 'DELETE' }),
    onSuccess: () => qc.invalidateQueries({ queryKey: queryKeys.reservations.all }),
  }));
}
