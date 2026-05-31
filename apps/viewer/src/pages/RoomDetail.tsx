import { Component, createMemo, createSignal, For, Show, onMount } from 'solid-js';
import { A, useParams } from '@solidjs/router';
import {
  useRoom,
  useSeatsByRoom,
  useReservationsByDate,
  useCreateReservation,
  useCancelReservation,
  useAdminCreateReservation,
  useCreateCrossTeamRequest,
  useTeams,
  useUsers,
  ApiError,
} from '@neshiman/api-client';
import type { definitions } from '@neshiman/api-types';

type Seat = definitions['dto.SeatResponse'];

const CELL = 64;

const TEAM_COLORS: { bg: string; cellBg: string; border: string; text: string }[] = [
  { bg: 'bg-blue-500', cellBg: 'bg-blue-100', border: 'border-blue-400', text: 'text-blue-800' },
  { bg: 'bg-green-500', cellBg: 'bg-green-100', border: 'border-green-400', text: 'text-green-800' },
  { bg: 'bg-purple-500', cellBg: 'bg-purple-100', border: 'border-purple-400', text: 'text-purple-800' },
  { bg: 'bg-amber-500', cellBg: 'bg-amber-100', border: 'border-amber-400', text: 'text-amber-800' },
  { bg: 'bg-pink-500', cellBg: 'bg-pink-100', border: 'border-pink-400', text: 'text-pink-800' },
  { bg: 'bg-cyan-500', cellBg: 'bg-cyan-100', border: 'border-cyan-400', text: 'text-cyan-800' },
  { bg: 'bg-indigo-500', cellBg: 'bg-indigo-100', border: 'border-indigo-400', text: 'text-indigo-800' },
  { bg: 'bg-rose-500', cellBg: 'bg-rose-100', border: 'border-rose-400', text: 'text-rose-800' },
];

const today = new Date().toISOString().slice(0, 10);

const gridCss = `
  .viewer-grid {
    background-image:
      linear-gradient(to right, rgba(0,0,0,0.06) 1px, transparent 1px),
      linear-gradient(to bottom, rgba(0,0,0,0.06) 1px, transparent 1px);
    background-size: ${CELL}px ${CELL}px;
  }
`;

const RoomDetail: Component = () => {
  const params = useParams<{ id: string }>();
  const room = useRoom(() => params.id);
  const seats = useSeatsByRoom(() => params.id);
  const teams = useTeams();
  const allUsers = useUsers(() => '');

  const [selectedDate, setSelectedDate] = createSignal(today);
  const [showConfirm, setShowConfirm] = createSignal<Seat | null>(null);
  const [successMsg, setSuccessMsg] = createSignal('');
  const [errorMsg, setErrorMsg] = createSignal('');

  onMount(() => {
    const id = 'viewer-grid-style';
    if (!document.getElementById(id)) {
      const s = document.createElement('style');
      s.id = id;
      s.textContent = gridCss;
      document.head.appendChild(s);
    }
  });

  const reservations = useReservationsByDate(selectedDate);

  const createReservation = useCreateReservation();
  const cancelReservation = useCancelReservation();
  const adminCreateReservation = useAdminCreateReservation();
  const createCrossTeamRequest = useCreateCrossTeamRequest();

  const currentUserId = () => {
    try { return JSON.parse(localStorage.getItem('neshiman_user') ?? 'null')?.id ?? ''; }
    catch { return ''; }
  };
  const currentUser = createMemo(() => allUsers.data?.find((u) => u.id === currentUserId()));

  const isAdminUser = createMemo(() => {
    const role = currentUser()?.role;
    return role === 'superadmin' || role === 'team_admin';
  });

  const teamMembers = createMemo(() => {
    const teamId = currentUser()?.team_id;
    if (!teamId) return [];
    return allUsers.data?.filter((u) => u.team_id === teamId) ?? [];
  });

  const teamColorMap = createMemo(() => {
    const map = new Map<string, number>();
    teams.data?.forEach((t, i) => {
      if (t.id) map.set(t.id, i % TEAM_COLORS.length);
    });
    return map;
  });

  const reservedSeatIds = createMemo(() => {
    return new Set(reservations.data?.map((r) => r.seat_id) ?? []);
  });

  const myUserReservations = createMemo(() => {
    return reservations.data?.filter((r) => r.user_id === currentUserId()) ?? [];
  });

  const reservationBySeat = createMemo(() => {
    const map = new Map<string, definitions['dto.ReservationResponse']>();
    for (const r of reservations.data ?? []) {
      if (r.seat_id) map.set(r.seat_id, r);
    }
    return map;
  });

  const myReservedSeatIds = createMemo(() => {
    return new Set(myUserReservations().map((r) => r.seat_id));
  });

  const handleReserve = async (seatId: string) => {
    setErrorMsg('');
    setSuccessMsg('');
    try {
      await createReservation.mutateAsync({ date: selectedDate(), seat_id: seatId });
      setShowConfirm(null);
      setSuccessMsg('Reservation confirmed!');
      setTimeout(() => setSuccessMsg(''), 3000);
    } catch (err) {
      setErrorMsg(err instanceof ApiError ? err.message : 'Failed to reserve seat');
    }
  };

  const handleCrossTeamRequest = async (seatId: string) => {
    setErrorMsg('');
    setSuccessMsg('');
    try {
      await createCrossTeamRequest.mutateAsync({ date: selectedDate(), target_seat_id: seatId });
      setShowConfirm(null);
      setSuccessMsg('Cross-team request submitted!');
      setTimeout(() => setSuccessMsg(''), 3000);
    } catch (err) {
      setErrorMsg(err instanceof ApiError ? err.message : 'Failed to submit request');
    }
  };

  const sortedSeats = createMemo(() => {
    const arr = [...(seats.data ?? [])];
    arr.sort((a, b) => {
      if ((a.pos_y ?? 0) !== (b.pos_y ?? 0)) return (a.pos_y ?? 0) - (b.pos_y ?? 0);
      return (a.pos_x ?? 0) - (b.pos_x ?? 0);
    });
    return arr;
  });

  const isOwnTeam = (seat: Seat) => {
    return seat.team_id === currentUser()?.team_id;
  };

  return (
    <div>
      <div class="mb-6">
        <div class="flex items-center gap-2 text-sm text-gray-500 mb-1">
          <A href="/rooms" class="text-blue-600 hover:underline">Rooms</A>
          <span>/</span>
          <span class="text-gray-900 font-medium">{room.data?.name ?? 'Loading...'}</span>
        </div>
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">{room.data?.name ?? 'Loading...'}</h1>
            <p class="text-sm text-gray-500">{room.data?.grid_width} × {room.data?.grid_height} grid</p>
          </div>
          <div class="flex items-center gap-3">
            <div class="flex items-center gap-2 bg-white border border-gray-200 rounded-lg px-3 py-2 shadow-sm">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 text-gray-400"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" /><line x1="16" y1="2" x2="16" y2="6" /><line x1="8" y1="2" x2="8" y2="6" /><line x1="3" y1="10" x2="21" y2="10" /></svg>
              <input
                type="date"
                value={selectedDate()}
                onInput={(e) => { setSelectedDate(e.currentTarget.value); setShowConfirm(null); }}
                min={today}
                class="border-0 text-sm font-medium text-gray-700 focus:ring-0 p-0"
              />
            </div>
          </div>
        </div>
      </div>

      <Show when={successMsg()}>
        <div class="mb-4 px-4 py-3 bg-green-50 border border-green-200 rounded-lg text-sm text-green-700 flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 shrink-0"><polyline points="20 6 9 17 4 12" /></svg>
          {successMsg()}
        </div>
      </Show>

      <Show when={errorMsg()}>
        <div class="mb-4 px-4 py-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700 flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 shrink-0"><circle cx="12" cy="12" r="10" /><line x1="15" y1="9" x2="9" y2="15" /><line x1="9" y1="9" x2="15" y2="15" /></svg>
          {errorMsg()}
        </div>
      </Show>

      <Show when={!currentUserId()}>
        <div class="mb-4 px-4 py-3 bg-amber-50 border border-amber-200 rounded-lg text-sm text-amber-700">
          Please sign in to reserve seats.
        </div>
      </Show>

      <div class="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
        <div class="flex items-center gap-4 px-5 py-3 border-b border-gray-100 bg-gray-50/50 text-xs">
          <div class="flex items-center gap-1.5">
            <div class="w-3.5 h-3.5 rounded border border-green-400 bg-green-100" />
            <span class="text-gray-600">Available</span>
          </div>
          <div class="flex items-center gap-1.5">
            <div class="w-3.5 h-3.5 rounded border border-blue-400 bg-blue-200" />
            <span class="text-gray-600">My Reservation</span>
          </div>
          <div class="flex items-center gap-1.5">
            <div class="w-3.5 h-3.5 rounded border border-gray-300 bg-gray-100" />
            <span class="text-gray-600">Reserved</span>
          </div>
          <div class="flex items-center gap-1.5">
            <div class="w-3.5 h-3.5 rounded border border-purple-400 bg-purple-200" />
            <span class="text-gray-600">Permanent Seat</span>
          </div>
        </div>

        <div class="overflow-auto bg-white flex items-center justify-center p-8" style={{ 'min-height': '400px' }}>
          <Show when={room.data && seats.data} fallback={
            <div class="flex items-center gap-2 text-gray-400">
              <div class="w-5 h-5 border-2 border-blue-600 border-t-transparent rounded-full animate-spin" />
              <span class="text-sm">Loading seats...</span>
            </div>
          }>
            <div
              class="viewer-grid relative border border-gray-200 rounded-lg shadow-sm"
              style={{
                width: `${(room.data?.grid_width ?? 10) * CELL + 1}px`,
                height: `${(room.data?.grid_height ?? 8) * CELL + 1}px`,
              }}
            >
              <For each={sortedSeats()}>
                {(seat) => {
                  const colorIdx = teamColorMap().get(seat.team_id ?? '') ?? 0;
                  const colors = TEAM_COLORS[colorIdx];
                  const isReservedByMe = myReservedSeatIds().has(seat.id);
                  const isReserved = reservedSeatIds().has(seat.id);
                  const ownTeam = isOwnTeam(seat);

                  let stateClass: string;
                  const isPermanentAssigned = !!seat.assigned_user_id;
                  const isMyPermanentSeat = isPermanentAssigned && seat.assigned_user_id === currentUserId();

                  if (isReservedByMe) {
                    stateClass = 'bg-blue-200 border-blue-400 ring-2 ring-blue-300';
                  } else if (isReserved) {
                    stateClass = 'bg-gray-100 border-gray-300 opacity-60';
                  } else if (isMyPermanentSeat) {
                    stateClass = 'bg-purple-200 border-purple-400 ring-2 ring-purple-300';
                  } else if (isPermanentAssigned) {
                    stateClass = 'bg-purple-100 border-purple-300 opacity-60';
                  } else if (!ownTeam && currentUserId()) {
                    stateClass = `${colors.cellBg} ${colors.border} opacity-70`;
                  } else {
                    stateClass = `${colors.cellBg} ${colors.border} cursor-pointer hover:ring-2 hover:ring-green-400`;
                  }

                  return (
                    <div
                      class="absolute rounded-lg border-2 flex flex-col items-center justify-center transition-all select-none"
                      classList={{
                        'hover:shadow-md': !isReserved && !isReservedByMe,
                      }}
                      style={{
                        left: `${(seat.pos_x ?? 0) * CELL + 3}px`,
                        top: `${(seat.pos_y ?? 0) * CELL + 3}px`,
                        width: `${CELL - 6}px`,
                        height: `${CELL - 6}px`,
                      }}
                      onClick={() => {
                        if (!currentUserId()) return;
                        if (isReserved && !isReservedByMe) return;
                        if (isPermanentAssigned && !isMyPermanentSeat) return;
                        setShowConfirm(seat);
                      }}
                    >
                      <div class={`w-full h-full rounded-md flex items-center justify-center ${stateClass}`}>
                        <span class="text-[10px] font-bold leading-tight text-center px-0.5"
                          classList={{
                            'text-gray-700': !isReserved && !isPermanentAssigned,
                            'text-gray-400': !isReservedByMe && (isReserved || (isPermanentAssigned && !isMyPermanentSeat)),
                            'text-blue-800': isReservedByMe,
                            'text-purple-800': isMyPermanentSeat,
                          }}
                        >
                          {seat.label}
                          <Show when={isPermanentAssigned}>
                            <br /><span class="text-[7px] font-normal">{isMyPermanentSeat ? 'Your seat' : 'Assigned'}</span>
                          </Show>
                        </span>
                      </div>
                    </div>
                  );
                }}
              </For>

              <Show when={sortedSeats().length === 0}>
                <div class="absolute inset-0 flex items-center justify-center">
                  <p class="text-gray-300 text-sm">No seats in this room</p>
                </div>
              </Show>
            </div>
          </Show>
        </div>
      </div>

      <Show when={showConfirm()}>
        {(seat) => {
          const ownTeam = isOwnTeam(seat());
          const isReservedByMe = myReservedSeatIds().has(seat().id);
          const isAdminReserve = ownTeam && !isReservedByMe && isAdminUser();
          const [adminTargetUserId, setAdminTargetUserId] = createSignal('');
          return (
            <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/30" onClick={() => setShowConfirm(null)}>
              <div class="bg-white rounded-xl shadow-xl border border-gray-200 p-6 w-full max-w-sm mx-4" onClick={(e) => e.stopPropagation()}>
                <h3 class="text-lg font-bold text-gray-900 mb-2">
                  {isReservedByMe ? 'Cancel Reservation?' : isAdminReserve ? 'Reserve for Team Member' : ownTeam ? 'Reserve Seat' : 'Request Seat'}
                </h3>
                <div class="space-y-2 mb-5 text-sm text-gray-600">
                  <p><span class="font-medium text-gray-900">Seat:</span> {seat().label}</p>
                  <p><span class="font-medium text-gray-900">Date:</span> {selectedDate()}</p>
                  <p><span class="font-medium text-gray-900">Team:</span> {teams.data?.find((t) => t.id === seat().team_id)?.name ?? 'Unknown'}</p>
                  <Show when={isAdminReserve}>
                    <div>
                      <label class="block text-xs font-medium text-gray-500 mb-1">Reserve for</label>
                      <select
                        value={adminTargetUserId()}
                        onChange={(e) => setAdminTargetUserId(e.currentTarget.value)}
                        class="w-full border rounded-md px-2.5 py-1.5 text-sm"
                      >
                        <option value="">Select a team member...</option>
                        <For each={teamMembers()}>
                          {(u) => <option value={u.id!}>{u.name}</option>}
                        </For>
                      </select>
                    </div>
                  </Show>
                  <Show when={!ownTeam && !isReservedByMe}>
                    <div class="mt-3 px-3 py-2 bg-amber-50 border border-amber-200 rounded-lg text-xs text-amber-700">
                      This seat belongs to another team. A cross-team request will be sent for approval.
                    </div>
                  </Show>
                </div>
                <div class="flex gap-3 justify-end">
                  <button
                    onClick={() => setShowConfirm(null)}
                    class="px-4 py-2 border border-gray-200 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
                  >
                    Cancel
                  </button>
                  <button
                    onClick={() => {
                      if (isReservedByMe) {
                        const res = reservationBySeat().get(seat().id!);
                        cancelReservation.mutateAsync(res?.id ?? '').then(() => {
                          setShowConfirm(null);
                          setSuccessMsg('Reservation cancelled!');
                          setTimeout(() => setSuccessMsg(''), 3000);
                        }).catch((err) => {
                          setErrorMsg(err instanceof ApiError ? err.message : 'Failed to cancel');
                        });
                      } else if (isAdminReserve && adminTargetUserId()) {
                        adminCreateReservation.mutateAsync({ date: selectedDate(), seat_id: seat().id!, user_id: adminTargetUserId() }).then(() => {
                          setShowConfirm(null);
                          setSuccessMsg('Reservation created for team member!');
                          setTimeout(() => setSuccessMsg(''), 3000);
                        }).catch((err) => {
                          setErrorMsg(err instanceof ApiError ? err.message : 'Failed to reserve');
                        });
                      } else if (ownTeam) {
                        handleReserve(seat().id!);
                      } else {
                        handleCrossTeamRequest(seat().id!);
                      }
                    }}
                    disabled={createReservation.isPending || cancelReservation.isPending || createCrossTeamRequest.isPending || adminCreateReservation.isPending || (isAdminReserve && !adminTargetUserId())}
                    class="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors disabled:opacity-50 shadow-sm"
                  >
                    {createReservation.isPending || cancelReservation.isPending || createCrossTeamRequest.isPending || adminCreateReservation.isPending ? 'Processing...' : 'Confirm'}
                  </button>
                </div>
              </div>
            </div>
          );
        }}
      </Show>
    </div>
  );
};

export default RoomDetail;
