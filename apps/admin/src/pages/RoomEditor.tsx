import { Component, createSignal, createMemo, createEffect, For, Show, onMount } from 'solid-js';
import { createStore } from 'solid-js/store';
import { A, useParams } from '@solidjs/router';
import {
  useRoom,
  useSeatsByRoom,
  useBulkSyncSeats,
  useTeams,
  ApiError,
} from '@neshiman/api-client';
import type { definitions } from '@neshiman/api-types';

type Seat = definitions['dto.SeatResponse'];

const CELL = 56;

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

const roomEditorStyle = `
  .room-editor-grid {
    background-image:
      linear-gradient(to right, rgba(0,0,0,0.08) 1px, transparent 1px),
      linear-gradient(to bottom, rgba(0,0,0,0.08) 1px, transparent 1px);
    background-size: ${CELL}px ${CELL}px;
  }
  .room-editor-cell {
    width: ${CELL - 8}px;
    height: ${CELL - 8}px;
    margin: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 8px;
    border: 2px solid transparent;
    font-size: 11px;
    font-weight: 700;
    cursor: pointer;
    user-select: none;
    box-shadow: 0 1px 3px rgba(0,0,0,0.08), 0 1px 2px rgba(0,0,0,0.04);
    transition: box-shadow 0.15s, transform 0.1s;
  }
  .room-editor-cell:hover {
    box-shadow: 0 0 0 2px rgba(59,130,246,0.5), 0 2px 6px rgba(0,0,0,0.1);
    transform: translateY(-1px);
    z-index: 5;
  }
  .room-editor-cell.selected {
    box-shadow: 0 0 0 3px #2563eb, 0 2px 8px rgba(37,99,235,0.2);
    z-index: 10;
  }
  .room-editor-cell.place-mode {
    cursor: crosshair;
  }
`;

const RoomEditor: Component = () => {
  const params = useParams<{ id: string }>();
  const room = useRoom(() => params.id);
  const serverSeats = useSeatsByRoom(() => params.id);
  const teams = useTeams();
  const bulkSync = useBulkSyncSeats();

  const [localSeats, setLocalSeats] = createStore<Seat[]>([]);
  const [originalJson, setOriginalJson] = createSignal('');
  const [selectedIdx, setSelectedIdx] = createSignal<number | null>(null);
  const [mode, setMode] = createSignal<'select' | 'place'>('select');
  const [placeTeamId, setPlaceTeamId] = createSignal('');
  const [saving, setSaving] = createSignal(false);
  const [error, setError] = createSignal('');
  const [statusMsg, setStatusMsg] = createSignal('');
  const [initialized, setInitialized] = createSignal(false);

  let gridEl!: HTMLDivElement;

  onMount(() => {
    applyStyleTag();
  });

  const applyStyleTag = () => {
    const id = 'room-editor-style';
    if (!document.getElementById(id)) {
      const s = document.createElement('style');
      s.id = id;
      s.textContent = roomEditorStyle;
      document.head.appendChild(s);
    }
  };

  createEffect(() => {
    const data = serverSeats.data;
    if (!data || initialized()) return;
    const sorted = [...data].sort((a: Seat, b: Seat) => {
      if ((a.pos_y ?? 0) !== (b.pos_y ?? 0)) return (a.pos_y ?? 0) - (b.pos_y ?? 0);
      return (a.pos_x ?? 0) - (b.pos_x ?? 0);
    });
    setLocalSeats(sorted as Seat[]);
    setOriginalJson(JSON.stringify(sorted));
    setSelectedIdx(null);
    setInitialized(true);
  });

  const isDirty = createMemo(() => {
    const current = JSON.stringify([...localSeats]);
    return current !== originalJson();
  });

  const teamColor = (seat: Seat) => {
    const idx = teams.data?.findIndex((t) => t.id === seat.team_id) ?? -1;
    return TEAM_COLORS[idx >= 0 ? idx % TEAM_COLORS.length : 0];
  };

  const incLabel = () => {
    const max = localSeats.reduce((m: number, s: Seat) => {
      const n = parseInt((s.label || '').replace(/\D/g, '') || '0', 10);
      return n > m ? n : m;
    }, 0);
    return `S${max + 1}`;
  };

  const getGridPos = (clientX: number, clientY: number) => {
    const rect = gridEl.getBoundingClientRect();
    const sx = clientX - rect.left;
    const sy = clientY - rect.top;
    return {
      x: Math.floor(sx / CELL),
      y: Math.floor(sy / CELL),
    };
  };

  const findSeatAt = (gx: number, gy: number) => {
    for (let i = 0; i < localSeats.length; i++) {
      if (localSeats[i].pos_x === gx && localSeats[i].pos_y === gy) return i;
    }
    return -1;
  };

  const handleClick = (e: MouseEvent) => {
    if (e.button !== 0) return;
    const gp = getGridPos(e.clientX, e.clientY);
    const hit = findSeatAt(gp.x, gp.y);
    if (mode() === 'select') {
      setSelectedIdx(hit >= 0 ? hit : null);
    } else {
      if (hit >= 0) {
        setSelectedIdx(hit);
      } else {
        placeSeat(gp.x, gp.y);
      }
    }
  };

  const placeSeat = (gx: number, gy: number) => {
    const gw = room.data?.grid_width ?? 0;
    const gh = room.data?.grid_height ?? 0;
    if (gx < 0 || gx >= gw || gy < 0 || gy >= gh) return;
    if (findSeatAt(gx, gy) >= 0) return;
    if (!placeTeamId()) return;
    const label = incLabel();
    const idx = localSeats.length;
    const newSeat: Seat = {
      room_id: params.id,
      team_id: placeTeamId(),
      label,
      pos_x: gx,
      pos_y: gy,
      rotation: 0,
    };
    setLocalSeats(idx, newSeat);
    setSelectedIdx(idx);
  };

  const handleDelete = () => {
    const idx = selectedIdx();
    if (idx === null || idx < 0 || idx >= localSeats.length) return;
    setLocalSeats(localSeats.filter((_: Seat, i: number) => i !== idx));
    setSelectedIdx(null);
  };

  const handleRotate = () => {
    const idx = selectedIdx();
    if (idx === null || idx < 0 || idx >= localSeats.length) return;
    const current = localSeats[idx].rotation ?? 0;
    const newRotation = (current + 90) % 360;
    setLocalSeats(idx, 'rotation', newRotation);
  };

  const handleSave = async () => {
    setSaving(true);
    setError('');
    setStatusMsg('');
    try {
      const seats = localSeats.map((s) => ({
        id: s.id ?? undefined,
        team_id: s.team_id!,
        label: s.label ?? '',
        pos_x: s.pos_x ?? 0,
        pos_y: s.pos_y ?? 0,
        rotation: s.rotation ?? 0,
      }));
      const result = await bulkSync.mutateAsync({
        roomId: params.id,
        data: { seats },
      });
      const updated = [...(result.seats ?? [])].sort((a: Seat, b: Seat) => {
        if ((a.pos_y ?? 0) !== (b.pos_y ?? 0)) return (a.pos_y ?? 0) - (b.pos_y ?? 0);
        return (a.pos_x ?? 0) - (b.pos_x ?? 0);
      });
      setLocalSeats(updated as Seat[]);
      setOriginalJson(JSON.stringify(updated));
      setSelectedIdx(null);
      setStatusMsg('Saved');
      setTimeout(() => setStatusMsg(''), 2000);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Failed to save');
    } finally {
      setSaving(false);
    }
  };

  const selectedSeat = () => {
    const idx = selectedIdx();
    if (idx === null || idx < 0 || idx >= localSeats.length) return null;
    return localSeats[idx];
  };

  const updateSelected = (field: string, value: unknown) => {
    const idx = selectedIdx();
    if (idx === null || idx < 0 || idx >= localSeats.length) return;
    setLocalSeats(idx, field as never, value as never);
  };

  return (
    <div class="flex flex-col flex-1 min-h-0">
      <div class="flex items-center justify-between px-6 py-3 border-b bg-white shrink-0 shadow-sm">
        <div class="flex items-center gap-3">
          <A href="/rooms" class="text-blue-600 hover:underline text-sm">Rooms</A>
          <span class="text-gray-400">/</span>
          <h1 class="text-lg font-semibold">{room.data?.name ?? 'Loading...'}</h1>
          <span class="text-xs text-gray-400">
            {room.data?.grid_width} × {room.data?.grid_height}
          </span>
        </div>
        <div class="flex items-center gap-4">
          <div class="flex items-center gap-1 bg-gray-100 rounded-lg p-0.5">
            <button
              onClick={() => setMode('select')}
              class={`px-3 py-1.5 text-xs rounded-md font-medium transition-colors ${mode() === 'select' ? 'bg-white shadow-sm text-gray-900' : 'text-gray-500 hover:text-gray-700'}`}
            >
              Select
            </button>
            <button
              onClick={() => setMode('place')}
              class={`px-3 py-1.5 text-xs rounded-md font-medium transition-colors ${mode() === 'place' ? 'bg-white shadow-sm text-gray-900' : 'text-gray-500 hover:text-gray-700'}`}
            >
              Place
            </button>
          </div>

          <Show when={mode() === 'place'}>
            <select
              value={placeTeamId()}
              onChange={(e) => setPlaceTeamId(e.currentTarget.value)}
              class="border rounded-md px-2 py-1.5 text-xs"
            >
              <option value="">Select team</option>
              <For each={teams.data}>
                {(t) => <option value={t.id!}>{t.name}</option>}
              </For>
            </select>
          </Show>

          <div class="flex items-center gap-2">
            <Show when={statusMsg()}>
              <span class="text-xs text-green-600">{statusMsg()}</span>
            </Show>
            <button
              onClick={handleSave}
              disabled={!isDirty() || saving()}
              class="px-4 py-1.5 bg-blue-600 text-white rounded-md text-sm disabled:opacity-40 disabled:cursor-not-allowed hover:bg-blue-700 transition-colors"
            >
              {saving() ? 'Saving...' : 'Save'}
            </button>
          </div>
        </div>
      </div>

      {error() && (
        <div class="px-6 py-2 bg-red-50 border-b text-red-700 text-sm">{error()}</div>
      )}

      <div class="flex flex-1 min-h-0">
        <div
          class="flex-1 overflow-auto bg-white flex items-center justify-center"
        >
          <div
            ref={gridEl!}
            class="room-editor-grid relative border border-gray-200 rounded-md shadow-sm bg-white"
            onClick={handleClick}
            style={{
              cursor: mode() === 'place' ? 'crosshair' : 'default',
              width: `${(room.data?.grid_width ?? 10) * CELL + 1}px`,
              height: `${(room.data?.grid_height ?? 8) * CELL + 1}px`,
            }}
          >
            <For each={localSeats}>
              {(seat, idx) => {
                const colors = teamColor(seat);
                const isSelected = selectedIdx() === idx();
                return (
                  <div
                    class={`room-editor-cell absolute ${colors.bg} ${colors.border}`}
                    classList={{
                      selected: isSelected,
                      'place-mode': mode() === 'place',
                    }}
                    style={{
                      left: `${(seat.pos_x ?? 0) * CELL}px`,
                      top: `${(seat.pos_y ?? 0) * CELL}px`,
                      transform: `rotate(${seat.rotation ?? 0}deg)`,
                      'z-index': isSelected ? 20 : 1,
                    }}
                  >
                    <span class={`text-center leading-tight ${colors.text}`}>
                      {seat.label}
                    </span>
                  </div>
                );
              }}
            </For>

            <Show when={localSeats.length === 0 && !serverSeats.isLoading}>
              <div class="absolute inset-0 flex items-center justify-center pointer-events-none">
                <p class="text-gray-300 text-sm">{mode() === 'place' ? 'Click to place seats' : 'No seats yet. Switch to Place mode.'}</p>
              </div>
            </Show>
          </div>
        </div>

        <Show when={selectedSeat()}>
          {(seat) => (
            <div class="w-64 border-l bg-white p-4 shrink-0 overflow-y-auto shadow-lg z-10">
              <h3 class="font-semibold text-sm mb-4">Seat Properties</h3>

              <div class="space-y-3">
                <div>
                  <label class="block text-xs font-medium text-gray-500 mb-1">Label</label>
                  <input
                    value={seat().label ?? ''}
                    onInput={(e) => updateSelected('label', e.currentTarget.value)}
                    class="w-full border rounded-md px-2.5 py-1.5 text-sm"
                  />
                </div>

                <div>
                  <label class="block text-xs font-medium text-gray-500 mb-1">Team</label>
                  <select
                    value={seat().team_id ?? ''}
                    onChange={(e) => updateSelected('team_id', e.currentTarget.value)}
                    class="w-full border rounded-md px-2.5 py-1.5 text-sm"
                  >
                    <For each={teams.data}>
                      {(t) => <option value={t.id!}>{t.name}</option>}
                    </For>
                  </select>
                </div>

                <div>
                  <label class="block text-xs font-medium text-gray-500 mb-1">Position</label>
                  <p class="text-sm text-gray-700">({seat().pos_x}, {seat().pos_y})</p>
                </div>

                <div>
                  <label class="block text-xs font-medium text-gray-500 mb-1">Rotation</label>
                  <div class="flex gap-1.5">
                    <button
                      onClick={() => updateSelected('rotation', 0)}
                      class={`px-2.5 py-1 text-xs rounded border ${seat().rotation === 0 || !seat().rotation ? 'bg-blue-100 border-blue-400 text-blue-800' : 'hover:bg-gray-50'}`}
                    >0°</button>
                    <button
                      onClick={() => updateSelected('rotation', 90)}
                      class={`px-2.5 py-1 text-xs rounded border ${seat().rotation === 90 ? 'bg-blue-100 border-blue-400 text-blue-800' : 'hover:bg-gray-50'}`}
                    >90°</button>
                    <button
                      onClick={() => updateSelected('rotation', 180)}
                      class={`px-2.5 py-1 text-xs rounded border ${seat().rotation === 180 ? 'bg-blue-100 border-blue-400 text-blue-800' : 'hover:bg-gray-50'}`}
                    >180°</button>
                    <button
                      onClick={() => updateSelected('rotation', 270)}
                      class={`px-2.5 py-1 text-xs rounded border ${seat().rotation === 270 ? 'bg-blue-100 border-blue-400 text-blue-800' : 'hover:bg-gray-50'}`}
                    >270°</button>
                    <button
                      onClick={handleRotate}
                      class="px-2.5 py-1 text-xs rounded border hover:bg-gray-50"
                      title="Rotate +90°"
                    >↻</button>
                  </div>
                </div>

                <hr class="my-3" />

                <button
                  onClick={handleDelete}
                  class="w-full px-3 py-2 bg-red-600 text-white rounded-md text-sm hover:bg-red-700 transition-colors"
                >
                  Delete Seat
                </button>
              </div>
            </div>
          )}
        </Show>
      </div>
    </div>
  );
};

export default RoomEditor;
