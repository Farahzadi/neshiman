import { Component, createEffect, createMemo, createSignal, For, Show } from 'solid-js';

import { useQueryClient } from '@tanstack/solid-query';
import {
  useRooms,
  useSeatsByRoom,
  useReservationsByRoomDate,
  useCreateReservation,
  useCancelReservation,
  useAdminCreateReservation,
  useCreateCrossTeamRequest,
  useTeams,
  useUsers,
  ApiError,
} from '@neshiman/api-client';
import type { definitions } from '@neshiman/api-types';
import { reconcile, unwrap } from 'solid-js/store';

type Seat = definitions['dto.SeatResponse'];


const DAY_INDEX: Record<string, number> = {
  sat: 6, sun: 0, mon: 1, tue: 2, wed: 3, thu: 4, fri: 5,
};

function getWeekStartDay(): string {
  try { return import.meta.env.VITE_WORK_WEEK_START ?? 'sat'; }
  catch { return 'sat'; }
}

function getWeekDays(): number {
  try { return Number(import.meta.env.VITE_WORK_WEEK_DAYS) || 5; }
  catch { return 5; }
}

// Return [startDate, workDayDates] for the week containing `date`
function getWorkWeek(date: Date, startDay: string, numDays: number): { startDate: Date; workDays: Date[] } {
  const dayIndex = DAY_INDEX[startDay] ?? 6;
  const currentDay = date.getDay();
  let diff = currentDay - dayIndex;
  if (diff < 0) diff += 7;

  const start = new Date(date);
  start.setDate(start.getDate() - diff);
  start.setHours(0, 0, 0, 0);

  const workDays: Date[] = [];
  for (let i = 0; i < numDays; i++) {
    const d = new Date(start);
    d.setDate(start.getDate() + i);
    workDays.push(d);
  }
  return { startDate: start, workDays };
}

function formatDate(d: Date): string {
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  return `${y}-${m}-${day}`;
}

function todayDate(): Date {
  const d = new Date();
  d.setHours(0, 0, 0, 0);
  return d;
}

function getCurrentUserId(): string {
  try { return JSON.parse(localStorage.getItem('neshiman_user') ?? 'null')?.id ?? ''; }
  catch { return ''; }
}


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

const WeeklyCalendar: Component = () => {
  const rooms = useRooms();
  const teams = useTeams();
  const allUsers = useUsers(() => '');

  let savedRoomId = '';
  try { savedRoomId = localStorage.getItem('neshiman_calendar_room') ?? ''; } catch { savedRoomId = ''; }
  const [selectedRoomId, setSelectedRoomId] = createSignal(savedRoomId);

  const effectiveRoomId = createMemo(() => {
    const current = selectedRoomId();
    const roomsData = rooms.data;
    if (current && roomsData?.some((r) => r.id === current)) return current;
    return roomsData?.[0]?.id ?? '';
  });

  const handleRoomChange = (value: string) => {
    setSelectedRoomId(value);
    try { localStorage.setItem('neshiman_calendar_room', value); } catch { /* localStorage unavailable */ }
  };
  const [weekOffset, setWeekOffset] = createSignal(0);
  const [selectedModal, setSelectedModal] = createSignal<{ seat: Seat; date: Date; isReservedByMe: boolean; isAdminReserve?: boolean } | null>(null);
  const [selectedUserForReserve, setSelectedUserForReserve] = createSignal('');
  const [successMsg, setSuccessMsg] = createSignal('');
  const [errMsg, setErrMsg] = createSignal('');

  const currentUserId = () => getCurrentUserId();
  const currentUser = createMemo(() => allUsers.data?.find((u) => u.id === currentUserId()));

  const weekDays = createMemo(() => {
    const startDay = getWeekStartDay();
    const numDays = getWeekDays();
    const baseDate = new Date();
    baseDate.setDate(baseDate.getDate() + weekOffset() * 7);
    return getWorkWeek(baseDate, startDay, numDays);
  });

  const dateStrings = createMemo(() => weekDays().workDays.map(formatDate));

  const seats = useSeatsByRoom(effectiveRoomId);

  // We need per-date queries. Use the first date for now, extend later
  const day0Reservations = useReservationsByRoomDate(
    () => effectiveRoomId(),
    () => dateStrings()[0] ?? '',
  );

  const day1Reservations = useReservationsByRoomDate(
    () => effectiveRoomId(),
    () => dateStrings()[1] ?? '',
  );

  const day2Reservations = useReservationsByRoomDate(
    () => effectiveRoomId(),
    () => dateStrings()[2] ?? '',
  );

  const day3Reservations = useReservationsByRoomDate(
    () => effectiveRoomId(),
    () => dateStrings()[3] ?? '',
  );

  const day4Reservations = useReservationsByRoomDate(
    () => effectiveRoomId(),
    () => dateStrings()[4] ?? '',
  );

  const allReservations = createMemo(() => {
    const data = [
      ...(day0Reservations.data ?? []),
      ...(day1Reservations.data ?? []),
      ...(day2Reservations.data ?? []),
      ...(day3Reservations.data ?? []),
      ...(day4Reservations.data ?? []),
    ];
    return data;
  });

  const allReservationsReady = createMemo(() =>
    !day0Reservations.isLoading && !day1Reservations.isLoading &&
    !day2Reservations.isLoading && !day3Reservations.isLoading &&
    !day4Reservations.isLoading
  );

  const queryClient = useQueryClient();

  const createReservation = useCreateReservation();
  const cancelReservation = useCancelReservation();
  const adminCreateReservation = useAdminCreateReservation();
  const createCrossTeamRequest = useCreateCrossTeamRequest();

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

  const sortedSeats = createMemo(() => {
    const arr = [...(seats.data ?? [])];
    arr.sort((a, b) => {
      if ((a.pos_y ?? 0) !== (b.pos_y ?? 0)) return (a.pos_y ?? 0) - (b.pos_y ?? 0);
      return (a.pos_x ?? 0) - (b.pos_x ?? 0);
    });
    return arr;
  });

  type CalendarRow =
    | { kind: 'team-header'; teamName: string; seatCount: number; colorIdx: number }
    | { kind: 'seat'; seat: Seat };

  const calendarRows = createMemo(() => {
    const rows: CalendarRow[] = [];
    const byTeam = new Map<string, Seat[]>();
    for (const seat of sortedSeats()) {
      const tid = seat.team_id ?? '';
      if (!byTeam.has(tid)) byTeam.set(tid, []);
      byTeam.get(tid)!.push(seat);
    }
    for (const t of teams.data ?? []) {
      if (!t.id || !byTeam.has(t.id)) continue;
      const teamSeats = byTeam.get(t.id)!;
      rows.push({ kind: 'team-header', teamName: t.name ?? 'Unknown', seatCount: teamSeats.length, colorIdx: teamColorMap().get(t.id) ?? 0 });
      for (const seat of teamSeats) rows.push({ kind: 'seat', seat });
      byTeam.delete(t.id);
    }
    for (const [, teamSeats] of byTeam) {
      rows.push({ kind: 'team-header', teamName: 'Unassigned', seatCount: teamSeats.length, colorIdx: 0 });
      for (const seat of teamSeats) rows.push({ kind: 'seat', seat });
    }
    return rows;
  });

  const reservationMapByDateSeat = createMemo(() => {
    const map = new Map<string, Map<string, definitions['dto.ReservationWithUserResponse']>>();
    for (const r of allReservations()) {
      if (!r.seat_id || !r.date) continue;
      if (!map.has(r.date)) map.set(r.date, new Map());
      map.get(r.date)!.set(r.seat_id, r);
    }
    return map;
  });

  const getReservation = (dateStr: string, seatId: string) => {
    return reservationMapByDateSeat().get(dateStr)?.get(seatId);
  };

  const isOwnTeam = (seat: Seat) => {
    return seat.team_id === currentUser()?.team_id;
  };

  const isPastDate = (d: Date) => {
    const today = todayDate();
    return d < today;
  };

  const handleCellClick = (seat: Seat, date: Date, r: definitions['dto.ReservationWithUserResponse'] | undefined) => {
    if (!currentUserId()) return;
    if (isPastDate(date) && r?.user_id !== currentUserId()) return;

    const isMyPermanentSeat = !!seat.assigned_user_id && seat.assigned_user_id === currentUserId();

    if (isMyPermanentSeat) {
      return;
    }

    if (r && r.user_id === currentUserId()) {
      setSelectedModal({ seat, date, isReservedByMe: true });
    } else if (!r && isOwnTeam(seat) && isAdminUser()) {
      setSelectedModal({ seat, date, isReservedByMe: false, isAdminReserve: true });
      setSelectedUserForReserve('');
    } else if (!r && isOwnTeam(seat)) {
      setSelectedModal({ seat, date, isReservedByMe: false });
    } else if (!r) {
      setSelectedModal({ seat, date, isReservedByMe: false });
    }
  };

  const handleConfirm = async () => {
    const modal = selectedModal();
    if (!modal) return;

    const seatId = modal.seat.id!;
    const dateStr = formatDate(modal.date);
    const r = getReservation(dateStr, seatId);
    const roomId = effectiveRoomId();
    type RWU = definitions['dto.ReservationWithUserResponse'];

    try {
      if (modal.isReservedByMe) {
        await cancelReservation.mutateAsync(r!.id!);
        queryClient.setQueryData<RWU[]>(
          ['reservations', 'by-room', roomId, dateStr],
          (old) => (old ?? []).filter((res) => res.seat_id !== seatId),
        );
        setSuccessMsg('Reservation cancelled!');
      } else if (modal.isAdminReserve && selectedUserForReserve()) {
        const result = await adminCreateReservation.mutateAsync({ date: dateStr, seat_id: seatId, user_id: selectedUserForReserve() });
        const targetUser = allUsers.data?.find((u) => u.id === selectedUserForReserve());
        queryClient.setQueryData<RWU[]>(
          ['reservations', 'by-room', roomId, dateStr],
          (old) => [...(old ?? []), { ...result, user_name: targetUser?.name ?? '', seat_label: modal.seat.label ?? '' } as unknown as RWU]
        );
        setSuccessMsg('Reservation created for team member!');
      } else if (isOwnTeam(modal.seat)) {
        const result = await createReservation.mutateAsync({ date: dateStr, seat_id: seatId });
        const user = allUsers.data?.find((u) => u.id === currentUserId());
        queryClient.setQueryData<RWU[]>(
          ['reservations', 'by-room', roomId, dateStr],
          (old) => [...(old ?? []), { ...result, user_name: user?.name ?? '', seat_label: modal.seat.label ?? '' } as unknown as RWU]
        );
        setSuccessMsg('Reservation confirmed!');
      } else {
        await createCrossTeamRequest.mutateAsync({ date: dateStr, target_seat_id: seatId });
        setSuccessMsg('Cross-team request submitted!');
      }
      queryClient.refetchQueries()
      setSelectedModal(null);
      setSelectedUserForReserve('');
      setTimeout(() => setSuccessMsg(''), 3000);
    } catch (err) {
      setErrMsg(err instanceof ApiError ? err.message : 'Failed');
      setTimeout(() => setErrMsg(''), 5000);
    }
  };

  const navigateWeek = (direction: number) => {
    setWeekOffset((prev) => prev + direction);
  };

  const goToToday = () => {
    setWeekOffset(0);
  };

  const workDays = () => weekDays().workDays;
  const dateStrs = () => dateStrings();

  return (
    <div>
      <div class="flex items-center justify-between mb-6">
        <h1 class="text-2xl font-bold text-gray-900">Weekly Calendar</h1>
        <div class="flex items-center gap-3">
          <div class="flex items-center gap-2">
            <Show when={currentUser()}>
              <span class="text-sm text-gray-500">
                Limit: {allUsers.data?.find((u) => u.id === currentUserId())?.weekly_limit ?? '-'}/week
              </span>
            </Show>
          </div>
        </div>
      </div>

      <Show when={successMsg()}>
        <div class="mb-4 px-4 py-3 bg-blue-50 border border-blue-200 rounded-lg text-sm text-blue-700 flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 shrink-0"><circle cx="12" cy="12" r="10" /><line x1="12" y1="16" x2="12" y2="12" /><line x1="12" y1="8" x2="12.01" y2="8" /></svg>
          {successMsg()}
        </div>
      </Show>
      <Show when={errMsg()}>
        <div class="mb-4 px-4 py-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700 flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 shrink-0"><circle cx="12" cy="12" r="10" /><line x1="15" y1="9" x2="9" y2="15" /><line x1="9" y1="9" x2="15" y2="15" /></svg>
          {errMsg()}
        </div>
      </Show>

      <div class="flex flex-wrap gap-4 mb-6 items-end">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Room</label>
          <select
            value={effectiveRoomId()}
            onChange={(e) => handleRoomChange(e.currentTarget.value)}
            class="border rounded-md px-3 py-2 text-sm w-64"
          >
            <For each={rooms.data}>
              {(room) => <option value={room.id!}>{room.name} ({room.grid_width}×{room.grid_height})</option>}
            </For>
          </select>
        </div>
        <div class="flex items-center gap-2">
          <button
            onClick={() => navigateWeek(-1)}
            class="px-3 py-2 border rounded-md text-sm hover:bg-gray-50"
          >
            ← Prev
          </button>
          <button
            onClick={goToToday}
            class="px-3 py-2 bg-blue-600 text-white rounded-md text-sm hover:bg-blue-700"
          >
            Today
          </button>
          <button
            onClick={() => navigateWeek(1)}
            class="px-3 py-2 border rounded-md text-sm hover:bg-gray-50"
          >
            Next →
          </button>
        </div>
      </div>

      <Show
        when={effectiveRoomId()}
        fallback={
          <div class="bg-white rounded-xl border border-dashed border-gray-300 p-12 text-center text-gray-400">
            <p class="text-lg font-medium mb-1">Select a room</p>
            <p class="text-sm">Choose a room above to see the weekly seat calendar</p>
          </div>
        }
      >
        <Show when={allReservationsReady()}
          fallback={
            <div class="bg-white rounded-xl border border-gray-200 shadow-sm p-12 text-center text-gray-400">
              <div class="flex items-center justify-center gap-2">
                <div class="w-5 h-5 border-2 border-blue-600 border-t-transparent rounded-full animate-spin" />
                <span class="text-sm">Loading reservations...</span>
              </div>
            </div>
          }
        >
          <div class="bg-white rounded-xl border border-gray-200 shadow-sm overflow-auto">
            <div class="min-w-max">
              {/* Header row: day names */}
              <div class="flex border-b border-gray-200 bg-gray-50/80">
                <div class="w-24 shrink-0 px-3 py-3 text-xs font-medium text-gray-500 border-r border-gray-200">
                  Seat
                </div>
                <For each={workDays()}>
                  {(day) => {
                    const today = todayDate();
                    const isToday = day.getFullYear() === today.getFullYear() &&
                      day.getMonth() === today.getMonth() &&
                      day.getDate() === today.getDate();
                    return (
                      <div class="w-28 shrink-0 px-3 py-3 text-center border-r border-gray-200 last:border-r-0"
                        classList={{ 'bg-blue-50 border-l-2 border-l-blue-400': isToday }}
                      >
                        <div class={`text-xs font-medium ${isToday ? 'text-blue-600 font-semibold' : 'text-gray-500'}`}>
                          {['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'][day.getDay()]}
                        </div>
                        <div class={`text-sm font-bold ${isToday ? 'text-blue-600' : 'text-gray-800'}`}>
                          {isToday ? (
                            <span class="inline-flex items-center justify-center w-7 h-7 rounded-full bg-blue-600 text-white text-xs">
                              {day.getDate()}
                            </span>
                          ) : day.getDate()}
                        </div>
                      </div>
                    );
                  }}
                </For>
              </div>

              {/* Seat rows grouped by team */}
              <For each={calendarRows()}>
                {(row) => row.kind === 'team-header' ? (
                  <div class="flex items-center gap-2 px-3 py-1.5 bg-gray-50 border-b border-gray-200 text-xs font-semibold text-gray-500">
                    <div class="w-2.5 h-2.5 rounded-sm" classList={{ [TEAM_COLORS[row.colorIdx].bg]: true }} />
                    {row.teamName}
                    <span class="text-gray-400 font-normal">({row.seatCount} seats)</span>
                  </div>
                ) : (
                  <div class="flex border-b border-gray-100 last:border-b-0 hover:bg-gray-50/50">
                    <div class="w-24 shrink-0 px-3 py-3 text-sm font-medium text-gray-700 border-r border-gray-100 flex items-center gap-2">
                      <div class="w-2.5 h-2.5 rounded-sm" classList={{
                        [TEAM_COLORS[teamColorMap().get(row.seat.team_id ?? '') ?? 0].bg]: true,
                      }} />
                      {row.seat.label}
                    </div>
                    <For each={workDays()}>
                      {(day, i) => {
                        const ds = dateStrs()[i()];
                        const r = getReservation(ds, row.seat.id!);
                        const isMine = r?.user_id === currentUserId();
                        const past = isPastDate(day);
                        const ownTeam = isOwnTeam(row.seat);
                        const isPermanentAssigned = !!row.seat.assigned_user_id;
                        const isMyPermanentSeat = isPermanentAssigned && row.seat.assigned_user_id === currentUserId();
                        const today = todayDate();
                        const isToday = day.getFullYear() === today.getFullYear() &&
                          day.getMonth() === today.getMonth() &&
                          day.getDate() === today.getDate();

                        let cellClass = '';
                        if (isMyPermanentSeat) {
                          cellClass = 'bg-purple-100 border-purple-300 text-purple-700';
                        } else if (isPermanentAssigned) {
                          cellClass = 'bg-purple-50 border-purple-200 text-gray-400';
                        } else if (r && isMine) {
                          cellClass = 'bg-blue-100 border-blue-300';
                        } else if (r) {
                          cellClass = 'bg-gray-100 border-gray-200 text-gray-400';
                        } else if (past) {
                          cellClass = 'bg-gray-50 border-gray-100 text-gray-300';
                        } else if (ownTeam) {
                          cellClass = isToday ? 'bg-blue-50 border-blue-200 hover:bg-blue-100 cursor-pointer' : 'bg-green-50 border-green-200 hover:bg-green-100 cursor-pointer';
                        } else {
                          cellClass = isToday ? 'bg-blue-50 border-blue-200 hover:bg-blue-100 cursor-pointer' : 'bg-amber-50 border-amber-200 hover:bg-amber-100 cursor-pointer';
                        }

                        const showAction = !isPermanentAssigned || isMyPermanentSeat;

                        return (
                          <div
                            class={`w-28 shrink-0 px-3 py-3 text-xs text-center border-r border-gray-100 last:border-r-0 transition-colors ${cellClass}`}
                            classList={{
                              'border-l-2 border-l-blue-400': isToday,
                            }}
                            onClick={() => showAction ? handleCellClick(row.seat, day, r) : undefined}
                          >
                            <Show when={isMyPermanentSeat}>
                              <span class="text-purple-700 font-medium">✦ Permanent</span>
                            </Show>
                            <Show when={isPermanentAssigned && !isMyPermanentSeat}>
                              <span class="text-gray-400">{row.seat.assigned_user_name?.slice(0, 12) ?? 'Assigned'}</span>
                            </Show>
                            <Show when={!isPermanentAssigned && r && isMine}>
                              <span class="text-blue-700 font-medium">✓ You</span>
                            </Show>
                            <Show when={!isPermanentAssigned && r && !isMine}>
                              <span class="text-gray-400">{r?.user_name?.slice(0, 12) ?? 'Reserved'}</span>
                            </Show>
                            <Show when={!isPermanentAssigned && !r && !past}>
                              <span class={ownTeam ? 'text-green-600' : 'text-amber-600'}>{ownTeam ? 'Available' : 'Cross-team'}</span>
                            </Show>
                            <Show when={!isPermanentAssigned && !r && past}>
                              <span class="text-gray-300">—</span>
                            </Show>
                          </div>
                        );
                      }}
                    </For>
                  </div>
                )}
              </For>

              <Show when={calendarRows().length === 0}>
                <div class="p-8 text-center text-gray-400">No seats in this room</div>
              </Show>
            </div>
          </div>
        </Show>
      </Show>

      {/* Legend */}
      <div class="flex flex-wrap gap-4 mt-4 text-xs text-gray-500">
        <div class="flex items-center gap-1.5">
          <div class="w-3 h-3 rounded bg-green-50 border border-green-200" /> Available (own team)
        </div>
        <div class="flex items-center gap-1.5">
          <div class="w-3 h-3 rounded bg-amber-50 border border-amber-200" /> Cross-team request
        </div>
        <div class="flex items-center gap-1.5">
          <div class="w-3 h-3 rounded bg-blue-100 border border-blue-300" /> My reservation
        </div>
        <div class="flex items-center gap-1.5">
          <div class="w-3 h-3 rounded bg-gray-100 border border-gray-200" /> Reserved
        </div>
        <div class="flex items-center gap-1.5">
          <div class="w-3 h-3 rounded bg-purple-100 border border-purple-300" /> Permanent seat
        </div>
      </div>

      {/* Confirmation modal */}
      <Show when={selectedModal()}>
        {(modal) => (
          <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/30" onClick={() => setSelectedModal(null)}>
            <div class="bg-white rounded-xl shadow-xl border border-gray-200 p-6 w-full max-w-sm mx-4" onClick={(e) => e.stopPropagation()}>
              <h3 class="text-lg font-bold text-gray-900 mb-2">
                {modal().isReservedByMe ? 'Cancel Reservation?' : modal().isAdminReserve ? 'Reserve for Team Member' : isOwnTeam(modal().seat) ? 'Reserve Seat' : 'Request Seat'}
              </h3>
              <div class="space-y-2 mb-5 text-sm text-gray-600">
                <p><span class="font-medium text-gray-900">Seat:</span> {modal().seat.label}</p>
                <p><span class="font-medium text-gray-900">Date:</span> {formatDate(modal().date)}</p>
                <p><span class="font-medium text-gray-900">Team:</span> {teams.data?.find((t) => t.id === modal().seat.team_id)?.name ?? 'Unknown'}</p>
                <Show when={modal().isAdminReserve}>
                  <div>
                    <label class="block text-xs font-medium text-gray-500 mb-1">Reserve for</label>
                    <select
                      value={selectedUserForReserve()}
                      onChange={(e) => setSelectedUserForReserve(e.currentTarget.value)}
                      class="w-full border rounded-md px-2.5 py-1.5 text-sm"
                    >
                      <option value="">Select a team member...</option>
                      <For each={teamMembers()}>
                        {(u) => <option value={u.id!}>{u.name}</option>}
                      </For>
                    </select>
                  </div>
                </Show>
                <Show when={!isOwnTeam(modal().seat) && !modal().isReservedByMe && !modal().isAdminReserve}>
                  <div class="mt-3 px-3 py-2 bg-amber-50 border border-amber-200 rounded-lg text-xs text-amber-700">
                    This seat belongs to another team. A cross-team request will be sent for approval.
                  </div>
                </Show>
              </div>
              <div class="flex gap-3 justify-end">
                <button
                  onClick={() => setSelectedModal(null)}
                  class="px-4 py-2 border border-gray-200 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50"
                >
                  Cancel
                </button>
                <button
                  onClick={handleConfirm}
                  disabled={createReservation.isPending || cancelReservation.isPending || createCrossTeamRequest.isPending || adminCreateReservation.isPending || (modal().isAdminReserve && !selectedUserForReserve())}
                  class="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700 disabled:opacity-50"
                >
                  {createReservation.isPending || cancelReservation.isPending || createCrossTeamRequest.isPending || adminCreateReservation.isPending ? 'Processing...' : 'Confirm'}
                </button>
              </div>
            </div>
          </div>
        )}
      </Show>
    </div>
  );
};

export default WeeklyCalendar;
