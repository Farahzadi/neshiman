import { Component, createSignal, For, Show, createMemo } from 'solid-js';
import { useParams, A } from '@solidjs/router';
import {
  useRoom,
  useSeatsByRoom,
  useCreateSeat,
  useDeleteSeat,
  useRotateSeat,
  useTeams,
  ApiError,
} from '@neshiman/api-client';
import type { definitions } from '@neshiman/api-types';
import Modal from '../components/Modal';

type Seat = definitions['dto.SeatResponse'];

const TEAM_COLORS: { bg: string; border: string; text: string }[] = [
  { bg: 'bg-blue-100', border: 'border-blue-400', text: 'text-blue-800' },
  { bg: 'bg-green-100', border: 'border-green-400', text: 'text-green-800' },
  { bg: 'bg-purple-100', border: 'border-purple-400', text: 'text-purple-800' },
  { bg: 'bg-amber-100', border: 'border-amber-400', text: 'text-amber-800' },
  { bg: 'bg-pink-100', border: 'border-pink-400', text: 'text-pink-800' },
  { bg: 'bg-indigo-100', border: 'border-indigo-400', text: 'text-indigo-800' },
  { bg: 'bg-cyan-100', border: 'border-cyan-400', text: 'text-cyan-800' },
  { bg: 'bg-rose-100', border: 'border-rose-400', text: 'text-rose-800' },
];

const Seats: Component = () => {
  const params = useParams<{ id: string }>();
  const room = useRoom(() => params.id);
  const seats = useSeatsByRoom(() => params.id);
  const teams = useTeams();
  const createSeat = useCreateSeat();
  const deleteSeat = useDeleteSeat();
  const rotateSeat = useRotateSeat();

  const [showForm, setShowForm] = createSignal(false);
  const [selectedPos, setSelectedPos] = createSignal<{ x: number; y: number } | null>(null);
  const [selectedSeat, setSelectedSeat] = createSignal<Seat | null>(null);
  const [newLabel, setNewLabel] = createSignal('');
  const [newTeamId, setNewTeamId] = createSignal('');
  const [newRotation, setNewRotation] = createSignal(0);
  const [error, setError] = createSignal('');

  const gridWidth = () => room.data?.grid_width ?? 0;
  const gridHeight = () => room.data?.grid_height ?? 0;

  const seatMap = createMemo(() => {
    const map = new Map<string, Seat>();
    for (const seat of seats.data ?? []) {
      map.set(`${seat.pos_x},${seat.pos_y}`, seat);
    }
    return map;
  });

  const teamColor = (teamId: string | undefined) => {
    if (!teamId) return { bg: 'bg-gray-100', border: 'border-gray-300', text: 'text-gray-700' };
    const idx = teams.data?.findIndex((t) => t.id === teamId) ?? -1;
    return TEAM_COLORS[idx >= 0 ? idx % TEAM_COLORS.length : 0];
  };

  const teamName = (teamId: string | undefined) => {
    if (!teamId) return 'No team';
    return teams.data?.find((t) => t.id === teamId)?.name ?? 'Unknown';
  };

  const handleCellClick = (x: number, y: number) => {
    const existing = seatMap().get(`${x},${y}`);
    if (existing) {
      setSelectedSeat(existing);
      setSelectedPos(null);
      setShowForm(false);
    } else {
      setSelectedPos({ x, y });
      setSelectedSeat(null);
      setNewLabel('');
      setNewTeamId('');
      setNewRotation(0);
      setError('');
      setShowForm(true);
    }
  };

  const handleCreateSeat = async (e: Event) => {
    e.preventDefault();
    const pos = selectedPos();
    if (!pos) return;
    setError('');
    try {
      await createSeat.mutateAsync({
        label: newLabel(),
        pos_x: pos.x,
        pos_y: pos.y,
        room_id: params.id,
        team_id: newTeamId() || undefined,
        rotation: newRotation(),
      });
      setShowForm(false);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Failed to create seat');
    }
  };

  const handleDeleteSeat = async (id: string) => {
    try {
      await deleteSeat.mutateAsync(id);
      setSelectedSeat(null);
    } catch (err) {
      alert(err instanceof ApiError ? err.message : 'Failed to delete seat');
    }
  };

  const handleRotateSeat = async (id: string) => {
    const seat = seats.data?.find((s) => s.id === id);
    const currentRotation = seat?.rotation ?? 0;
    try {
      await rotateSeat.mutateAsync({ id, data: { rotation: (currentRotation + 90) % 360 } });
    } catch (err) {
      alert(err instanceof ApiError ? err.message : 'Failed to rotate seat');
    }
  };

  const cells = createMemo(() => {
    const result: { x: number; y: number; seat: Seat | null }[] = [];
    for (let y = 0; y < gridHeight(); y++) {
      for (let x = 0; x < gridWidth(); x++) {
        result.push({ x, y, seat: seatMap().get(`${x},${y}`) ?? null });
      }
    }
    return result;
  });

  const loading = () => room.isLoading || seats.isLoading;

  return (
    <div>
      <div class="flex items-center gap-3 mb-6">
        <A href="/rooms" class="text-blue-600 hover:underline text-sm">Rooms</A>
        <span class="text-gray-400">/</span>
        <h1 class="text-2xl font-bold">{room.data?.name ?? 'Loading...'}</h1>
      </div>
      <p class="text-sm text-gray-500 mb-4">
        {gridWidth()} × {gridHeight()} grid &middot; {seats.data?.length ?? 0} seats
      </p>

      <Show when={!loading()} fallback={<p class="text-gray-400">Loading...</p>}>
        <div
          class="grid gap-1.5"
          style={{
            'grid-template-columns': `repeat(${gridWidth()}, 56px)`,
            'grid-template-rows': `repeat(${gridHeight()}, 56px)`,
          }}
        >
          <For each={cells()}>
            {(cell) => {
              const colors = cell.seat ? teamColor(cell.seat.team_id) : null;
              return (
                <div
                  onClick={() => handleCellClick(cell.x, cell.y)}
                  class="flex items-center justify-center rounded-md border-2 text-xs font-medium cursor-pointer select-none transition-colors"
                  classList={{
                    'hover:bg-gray-100 hover:border-gray-400': !cell.seat,
                    'bg-gray-50 border-gray-200': !cell.seat,
                    [`${colors?.bg} ${colors?.border} ${colors?.text}`]: !!cell.seat,
                  }}
                  title={cell.seat ? `${cell.seat.label} (${teamName(cell.seat.team_id)})` : `(${cell.x}, ${cell.y})`}
                >
                  <Show when={cell.seat} fallback={<span class="text-gray-300">+</span>}>
                    <span>{cell.seat!.label}</span>
                  </Show>
                </div>
              );
            }}
          </For>
        </div>
      </Show>

      <Show when={selectedSeat()}>
        {(seat) => (
          <div class="mt-6 p-4 border rounded-lg bg-white shadow-sm">
            <h3 class="font-semibold mb-2">Seat: {seat().label}</h3>
            <p class="text-sm text-gray-600 mb-3">
              Position ({seat().pos_x}, {seat().pos_y}) &middot;
              Rotation {seat().rotation}° &middot;
              Team: {teamName(seat().team_id)}
            </p>
            <div class="flex gap-2">
              <button
                onClick={() => handleRotateSeat(seat().id!)}
                disabled={rotateSeat.isPending}
                class="px-3 py-1.5 border rounded-md text-sm hover:bg-gray-50"
              >
                Rotate +90°
              </button>
              <button
                onClick={() => handleDeleteSeat(seat().id!)}
                disabled={deleteSeat.isPending}
                class="px-3 py-1.5 bg-red-600 text-white rounded-md text-sm hover:bg-red-700"
              >
                Delete
              </button>
            </div>
          </div>
        )}
      </Show>

      <Modal open={showForm()} onClose={() => setShowForm(false)} title="Add Seat">
        <form onSubmit={handleCreateSeat} class="space-y-4">
          <p class="text-sm text-gray-500">
            Position: ({selectedPos()?.x}, {selectedPos()?.y})
          </p>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Label</label>
            <input
              value={newLabel()}
              onInput={(e) => setNewLabel(e.currentTarget.value)}
              required
              class="w-full border rounded-md px-3 py-2 text-sm"
              placeholder="A1"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Team</label>
            <select
              value={newTeamId()}
              onChange={(e) => setNewTeamId(e.currentTarget.value)}
              class="w-full border rounded-md px-3 py-2 text-sm"
            >
              <option value="">No team</option>
              <For each={teams.data}>
                {(team) => <option value={team.id!}>{team.name}</option>}
              </For>
            </select>
          </div>
          {error() && <p class="text-red-600 text-sm">{error()}</p>}
          <div class="flex justify-end gap-3 pt-2">
            <button type="button" onClick={() => setShowForm(false)} class="px-4 py-2 border rounded-md text-sm">
              Cancel
            </button>
            <button
              type="submit"
              disabled={createSeat.isPending}
              class="px-4 py-2 bg-blue-600 text-white rounded-md text-sm disabled:opacity-50"
            >
              Create
            </button>
          </div>
        </form>
      </Modal>
    </div>
  );
};

export default Seats;
