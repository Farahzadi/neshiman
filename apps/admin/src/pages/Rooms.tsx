import { Component, createSignal, For } from 'solid-js';
import { A } from '@solidjs/router';
import { useRooms, useCreateRoom, useUpdateRoom, useDeleteRoom, ApiError } from '@neshiman/api-client';
import type { definitions } from '@neshiman/api-types';
import Modal from '../components/Modal';

type Room = definitions['dto.RoomResponse'];

const Rooms: Component = () => {
  const rooms = useRooms();
  const createRoom = useCreateRoom();
  const updateRoom = useUpdateRoom();
  const deleteRoom = useDeleteRoom();

  const [showForm, setShowForm] = createSignal(false);
  const [editingId, setEditingId] = createSignal<string | null>(null);
  const [name, setName] = createSignal('');
  const [gridWidth, setGridWidth] = createSignal(10);
  const [gridHeight, setGridHeight] = createSignal(8);
  const [error, setError] = createSignal('');

  const openCreate = () => {
    setEditingId(null);
    setName('');
    setGridWidth(10);
    setGridHeight(8);
    setError('');
    setShowForm(true);
  };

  const openEdit = (room: Room) => {
    setEditingId(room.id ?? null);
    setName(room.name ?? '');
    setGridWidth(room.grid_width ?? 10);
    setGridHeight(room.grid_height ?? 8);
    setError('');
    setShowForm(true);
  };

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError('');
    try {
      if (editingId()) {
        await updateRoom.mutateAsync({
          id: editingId()!,
          data: { name: name(), grid_width: gridWidth(), grid_height: gridHeight() },
        });
      } else {
        await createRoom.mutateAsync({ name: name(), grid_width: gridWidth(), grid_height: gridHeight() });
      }
      setShowForm(false);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Failed to save room');
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Delete this room and all its seats?')) return;
    try {
      await deleteRoom.mutateAsync(id);
    } catch (err) {
      alert(err instanceof ApiError ? err.message : 'Failed to delete room');
    }
  };

  return (
    <div class="p-6">
      <div class="flex items-center justify-between mb-6">
        <h1 class="text-2xl font-bold">Rooms</h1>
        <button onClick={openCreate} class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 text-sm">
          + New Room
        </button>
      </div>

      <div class="bg-white rounded-lg shadow overflow-hidden">
        <table class="w-full">
          <thead>
            <tr class="text-left text-sm text-gray-500 border-b">
              <th class="px-4 py-3 font-medium">Name</th>
              <th class="px-4 py-3 font-medium">Grid</th>
              <th class="px-4 py-3 font-medium">Created</th>
              <th class="px-4 py-3 font-medium">Actions</th>
            </tr>
          </thead>
          <tbody>
            <For each={rooms.data}>
              {(room) => (
                <tr class="border-b last:border-0 hover:bg-gray-50">
                  <td class="px-4 py-3 font-medium">{room.name}</td>
                  <td class="px-4 py-3 text-sm text-gray-600">{room.grid_width} × {room.grid_height}</td>
                  <td class="px-4 py-3 text-sm text-gray-500">
                    {room.created_at ? new Date(room.created_at).toLocaleDateString() : '-'}
                  </td>
                  <td class="px-4 py-3 text-sm space-x-3">
                    <A href={`/rooms/${room.id}/edit`} class="text-blue-600 hover:underline">
                      Edit Layout
                    </A>
                    <button onClick={() => openEdit(room)} class="text-gray-600 hover:underline">
                      Edit
                    </button>
                    <button onClick={() => handleDelete(room.id!)} class="text-red-600 hover:underline">
                      Delete
                    </button>
                  </td>
                </tr>
              )}
            </For>
            {rooms.data && rooms.data.length === 0 && (
              <tr>
                <td colspan="4" class="px-4 py-8 text-center text-gray-400">No rooms yet. Create one!</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <Modal open={showForm()} onClose={() => setShowForm(false)} title={editingId() ? 'Edit Room' : 'New Room'}>
        <form onSubmit={handleSubmit} class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
            <input
              value={name()}
              onInput={(e) => setName(e.currentTarget.value)}
              required
              class="w-full border rounded-md px-3 py-2 text-sm"
              placeholder="Conference Room A"
            />
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Width</label>
              <input
                type="number"
                min="1"
                max="50"
                value={gridWidth()}
                onInput={(e) => setGridWidth(Number(e.currentTarget.value))}
                required
                class="w-full border rounded-md px-3 py-2 text-sm"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Height</label>
              <input
                type="number"
                min="1"
                max="50"
                value={gridHeight()}
                onInput={(e) => setGridHeight(Number(e.currentTarget.value))}
                required
                class="w-full border rounded-md px-3 py-2 text-sm"
              />
            </div>
          </div>
          {error() && <p class="text-red-600 text-sm">{error()}</p>}
          <div class="flex justify-end gap-3 pt-2">
            <button type="button" onClick={() => setShowForm(false)} class="px-4 py-2 border rounded-md text-sm">
              Cancel
            </button>
            <button
              type="submit"
              disabled={createRoom.isPending || updateRoom.isPending}
              class="px-4 py-2 bg-blue-600 text-white rounded-md text-sm disabled:opacity-50"
            >
              {editingId() ? 'Update' : 'Create'}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  );
};

export default Rooms;
