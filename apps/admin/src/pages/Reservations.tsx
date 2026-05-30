import { Component, createSignal, createMemo, For, Show } from 'solid-js';
import {
  useRooms,
  useUsers,
  useSeatsByRoom,
  useReservationsByRoomDate,
  useAdminCreateReservation,
  useCancelReservation,
  ApiError,
} from '@neshiman/api-client';
import { showToast } from '../stores/toast';

function getCurrentUser() {
  try { return JSON.parse(localStorage.getItem('neshiman_user') ?? 'null'); }
  catch { return null; }
}

function todayStr() {
  const d = new Date();
  return d.toISOString().slice(0, 10);
}

const Reservations: Component = () => {
  const rooms = useRooms();
  const [selectedRoomId, setSelectedRoomId] = createSignal('');
  const [date, setDate] = createSignal(todayStr());

  const reservations = useReservationsByRoomDate(selectedRoomId, date);
  const seats = useSeatsByRoom(selectedRoomId);
  const users = useUsers(() => '');
  const cancelReservation = useCancelReservation();
  const adminCreate = useAdminCreateReservation();

  const [showReserveForm, setShowReserveForm] = createSignal(false);
  const [reserveUserId, setReserveUserId] = createSignal('');
  const [reserveSeatId, setReserveSeatId] = createSignal('');

  const currentUser = createMemo(() => getCurrentUser());
  const isSuperAdmin = createMemo(() => currentUser()?.role === 'superadmin');

  const userNameMap = createMemo(() => {
    const map = new Map<string, string>();
    for (const u of users.data ?? []) {
      if (u.id && u.name) map.set(u.id, u.name);
    }
    return map;
  });

  const teamUserIds = createMemo(() => {
    if (!currentUser()?.team_id) return new Set<string>();
    const set = new Set<string>();
    for (const u of users.data ?? []) {
      if (u.team_id === currentUser()?.team_id && u.id) set.add(u.id);
    }
    return set;
  });

  const handleCancel = async (id: string | undefined) => {
    if (!id) return;
    try {
      await cancelReservation.mutateAsync(id);
      showToast('Reservation cancelled.', 'success');
    } catch (err) {
      showToast(err instanceof ApiError ? err.message : 'Failed to cancel', 'error');
    }
  };

  const handleAdminReserve = async (e: Event) => {
    e.preventDefault();
    if (!reserveUserId() || !reserveSeatId()) return;
    try {
      await adminCreate.mutateAsync({
        user_id: reserveUserId(),
        seat_id: reserveSeatId(),
        date: date(),
      });
      setShowReserveForm(false);
      setReserveUserId('');
      setReserveSeatId('');
      showToast('Reservation created.', 'success');
    } catch (err) {
      showToast(err instanceof ApiError ? err.message : 'Failed to create reservation', 'error');
    }
  };

  const availableSeats = createMemo(() => {
    const reservedSeatIds = new Set((reservations.data ?? []).map((r) => r.seat_id));
    return (seats.data ?? []).filter((s) => s.id && !reservedSeatIds.has(s.id));
  });

  return (
    <div class="p-6">
      <h1 class="text-2xl font-bold mb-6">Reservation Management</h1>

      <div class="flex gap-4 mb-6 items-end">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Room</label>
          <select
            value={selectedRoomId()}
            onChange={(e) => setSelectedRoomId(e.currentTarget.value)}
            class="border rounded-md px-3 py-2 text-sm w-64"
          >
            <option value="">Select a room</option>
            <For each={rooms.data}>
              {(room) => <option value={room.id!}>{room.name}</option>}
            </For>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Date</label>
          <input
            type="date"
            value={date()}
            onInput={(e) => setDate(e.currentTarget.value)}
            class="border rounded-md px-3 py-2 text-sm"
          />
        </div>
        <Show when={selectedRoomId()}>
          <button
            onClick={() => setShowReserveForm(true)}
            class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 text-sm"
          >
            Reserve for User
          </button>
        </Show>
      </div>

      <Show
        when={selectedRoomId()}
        fallback={<p class="text-gray-400 text-center py-8">Select a room and date to view reservations</p>}
      >
        <div class="bg-white rounded-lg shadow overflow-hidden">
          <table class="w-full text-sm">
            <thead>
              <tr class="text-left text-gray-500 border-b">
                <th class="px-4 py-3 font-medium">User</th>
                <th class="px-4 py-3 font-medium">Seat</th>
                <th class="px-4 py-3 font-medium">Date</th>
                <th class="px-4 py-3 font-medium">Actions</th>
              </tr>
            </thead>
            <tbody>
              <For each={reservations.data}>
                {(r) => (
                  <tr class="border-b last:border-0 hover:bg-gray-50">
                    <td class="px-4 py-3">{r.user_id ? (userNameMap().get(r.user_id) ?? r.user_id.slice(0, 8)) : '-'}</td>
                    <td class="px-4 py-3">{r.seat_label}</td>
                    <td class="px-4 py-3">{r.date}</td>
                    <td class="px-4 py-3">
                      <button
                        onClick={() => handleCancel(r.id)}
                        disabled={cancelReservation.isPending}
                        class="text-red-600 hover:underline text-xs"
                      >
                        Cancel
                      </button>
                    </td>
                  </tr>
                )}
              </For>
              {reservations.data && reservations.data.length === 0 && (
                <tr>
                  <td colspan="4" class="px-4 py-8 text-center text-gray-400">No reservations for this room and date</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </Show>

      <Show when={showReserveForm()}>
        <div class="fixed inset-0 bg-black bg-opacity-30 flex items-center justify-center z-50">
          <div class="bg-white rounded-lg shadow-xl w-full max-w-md">
            <form onSubmit={handleAdminReserve} class="p-6 space-y-4">
              <h3 class="text-lg font-bold">Reserve for User</h3>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">User</label>
                <select
                  value={reserveUserId()}
                  onChange={(e) => setReserveUserId(e.currentTarget.value)}
                  required
                  class="w-full border rounded-md px-3 py-2 text-sm"
                >
                  <option value="">Select a user</option>
                  <For each={users.data}>
                    {(u) => {
                    const show = isSuperAdmin() || (u.id != null && teamUserIds().has(u.id));
                    return show ? <option value={u.id!}>{u.name} ({u.role})</option> : null;
                    }}
                  </For>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Seat</label>
                <select
                  value={reserveSeatId()}
                  onChange={(e) => setReserveSeatId(e.currentTarget.value)}
                  required
                  class="w-full border rounded-md px-3 py-2 text-sm"
                >
                  <option value="">Select a seat</option>
                  <For each={availableSeats()}>
                    {(s) => <option value={s.id!}>{s.label}</option>}
                  </For>
                </select>
              </div>
              <div class="flex justify-end gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => { setShowReserveForm(false); setReserveUserId(''); setReserveSeatId(''); }}
                  class="px-4 py-2 border rounded-md text-sm"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={adminCreate.isPending}
                  class="px-4 py-2 bg-blue-600 text-white rounded-md text-sm"
                >
                  Reserve
                </button>
              </div>
            </form>
          </div>
        </div>
      </Show>
    </div>
  );
};

export default Reservations;