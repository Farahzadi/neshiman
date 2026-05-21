import { useQuery, useMutation, useQueryClient } from '@tanstack/solid-query';
import { apiFetch } from './client';
import { queryKeys } from './query-keys';
import type { definitions } from '@neshiman/api-types';

type Room = definitions['dto.RoomResponse'];
type CreateRoom = definitions['dto.CreateRoomRequest'];
type UpdateRoom = definitions['dto.UpdateRoomRequest'];

export function useRooms() {
  return useQuery(() => ({
    queryKey: queryKeys.rooms.all,
    queryFn: () => apiFetch<Room[]>('/v1/rooms'),
  }));
}

export function useRoom(id: () => string) {
  return useQuery(() => ({
    queryKey: queryKeys.rooms.detail(id()),
    queryFn: () => apiFetch<Room>(`/v1/rooms/${id()}`),
    enabled: !!id(),
  }));
}

export function useCreateRoom() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: (data: CreateRoom) =>
      apiFetch<Room>('/v1/rooms', { method: 'POST', body: JSON.stringify(data) }),
    onSuccess: () => qc.invalidateQueries({ queryKey: queryKeys.rooms.all }),
  }));
}

export function useUpdateRoom() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: ({ id, data }: { id: string; data: UpdateRoom }) =>
      apiFetch<Room>(`/v1/rooms/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    onSuccess: (_data, vars) => {
      qc.invalidateQueries({ queryKey: queryKeys.rooms.all });
      qc.invalidateQueries({ queryKey: queryKeys.rooms.detail(vars.id) });
    },
  }));
}

export function useDeleteRoom() {
  const qc = useQueryClient();
  return useMutation(() => ({
    mutationFn: (id: string) =>
      apiFetch<void>(`/v1/rooms/${id}`, { method: 'DELETE' }),
    onSuccess: () => qc.invalidateQueries({ queryKey: queryKeys.rooms.all }),
  }));
}
