import { useQuery, useMutation, useQueryClient } from '@tanstack/solid-query';
import { apiFetch } from './client';
import { queryKeys } from './query-keys';
import type { definitions } from '@neshiman/api-types';

type Seat = definitions['dto.SeatResponse'];
type CreateSeat = definitions['dto.CreateSeatRequest'];
type MoveSeat = definitions['dto.MoveSeatRequest'];
type BulkSyncSeatsRequest = definitions['dto.BulkSyncSeatsRequest'];
type BulkSyncSeatsResponse = definitions['dto.BulkSyncSeatsResponse'];

export function useSeatsByRoom(roomId: () => string) {
  return useQuery(() => ({
    queryKey: queryKeys.rooms.seats(roomId()),
    queryFn: () => apiFetch<Seat[]>(`/v1/seats?room_id=${roomId()}`),
    enabled: !!roomId(),
  }));
}

export function useSeat(id: () => string) {
  return useQuery(() => ({
    queryKey: queryKeys.seats.detail(id()),
    queryFn: () => apiFetch<Seat>(`/v1/seats/${id()}`),
    enabled: !!id(),
  }));
}

export function useCreateSeat() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: (data: CreateSeat) =>
      apiFetch<Seat>('/v1/seats', { method: 'POST', body: JSON.stringify(data) }),
    onSuccess: (seat) => {
      qc.invalidateQueries({ queryKey: queryKeys.rooms.seats(seat.room_id!) });
      qc.invalidateQueries({ queryKey: queryKeys.seats.all });
    },
  }));
}

export function useDeleteSeat() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: (id: string) =>
      apiFetch<void>(`/v1/seats/${id}`, { method: 'DELETE' }),
    onSuccess: () => qc.invalidateQueries({ queryKey: queryKeys.seats.all }),
  }));
}

export function useMoveSeat() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: ({ id, data }: { id: string; data: MoveSeat }) =>
      apiFetch<Seat>(`/v1/seats/${id}/move`, { method: 'PUT', body: JSON.stringify(data) }),
    onSuccess: (seat) => {
      qc.invalidateQueries({ queryKey: queryKeys.seats.detail(seat.id!) });
      qc.invalidateQueries({ queryKey: queryKeys.rooms.seats(seat.room_id!) });
    },
  }));
}

export function useBulkSyncSeats() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: ({ roomId, data }: { roomId: string; data: BulkSyncSeatsRequest }) =>
      apiFetch<BulkSyncSeatsResponse>(`/v1/rooms/${roomId}/seats`, { method: 'PUT', body: JSON.stringify(data) }),
    onSuccess: (result, vars) => {
      qc.setQueryData(queryKeys.rooms.seats(vars.roomId), result.seats);
      qc.invalidateQueries({ queryKey: queryKeys.seats.all });
    },
  }));
}
