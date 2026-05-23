import { Component, createMemo, createSignal, For, Show } from 'solid-js';
import { A } from '@solidjs/router';
import { useReservationsByUserDate, useCancelReservation, useUsers, ApiError } from '@neshiman/api-client';

const today = () => new Date().toISOString().slice(0, 10);

const formatDate = (s: string) => {
  const d = new Date(s + 'T00:00:00');
  return d.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric' });
};

const Reservations: Component = () => {
  const userId = () => localStorage.getItem('viewer_user_id') ?? '';
  const allUsers = useUsers(() => '');
  const currentUser = createMemo(() => allUsers.data?.find((u) => u.id === userId()));

  const [selectedDate, setSelectedDate] = createSignal(today());
  const reservations = useReservationsByUserDate(userId, selectedDate);
  const cancelReservation = useCancelReservation();
  const [errorMsg, setErrorMsg] = createSignal('');

  const weeklyLimit = createMemo(() => currentUser()?.weekly_limit ?? 2);

  return (
    <div>
      <div class="mb-8">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">My Reservations</h1>
            <p class="text-sm text-gray-500 mt-1">
              <Show when={currentUser()} fallback="Select a user from the top-right">
                {currentUser()?.name}'s reservations
              </Show>
            </p>
          </div>
          <div class="flex items-center gap-3">
            <div class="text-sm text-gray-500">
              Limit: <span class="font-semibold text-gray-900">{weeklyLimit()}</span> days/week
            </div>
            <div class="flex items-center gap-2 bg-white border border-gray-200 rounded-lg px-3 py-2 shadow-sm">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 text-gray-400"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" /><line x1="16" y1="2" x2="16" y2="6" /><line x1="8" y1="2" x2="8" y2="6" /><line x1="3" y1="10" x2="21" y2="10" /></svg>
              <input
                type="date"
                value={selectedDate()}
                onInput={(e) => setSelectedDate(e.currentTarget.value)}
                class="border-0 text-sm font-medium text-gray-700 focus:ring-0 p-0"
              />
            </div>
          </div>
        </div>
        <Show when={reservations.data && reservations.data.length >= weeklyLimit()}>
          <div class="mt-3 px-4 py-2 bg-amber-50 border border-amber-200 rounded-lg text-xs text-amber-700 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3.5 h-3.5 shrink-0"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z" /><line x1="12" y1="9" x2="12" y2="13" /><line x1="12" y1="17" x2="12.01" y2="17" /></svg>
            You have reached your weekly limit of {weeklyLimit()} days.
          </div>
        </Show>
        <Show when={errorMsg()}>
          <div class="mt-3 px-4 py-2 bg-red-50 border border-red-200 rounded-lg text-xs text-red-700">
            {errorMsg()}
          </div>
        </Show>
      </div>

      <Show when={!!userId()} fallback={
        <div class="text-center py-20">
          <div class="w-14 h-14 rounded-2xl bg-gray-100 flex items-center justify-center mx-auto mb-3">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-7 h-7 text-gray-300"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" /><line x1="16" y1="2" x2="16" y2="6" /><line x1="8" y1="2" x2="8" y2="6" /><line x1="3" y1="10" x2="21" y2="10" /></svg>
          </div>
          <p class="text-gray-500">Select a user from the top-right dropdown.</p>
        </div>
      }>
        <div class="bg-white rounded-xl border border-gray-200 shadow-sm">
          <div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between">
            <h2 class="font-semibold text-gray-900">{formatDate(selectedDate())}</h2>
            <span class="text-xs text-gray-500">{reservations.data?.length ?? 0} reservation(s)</span>
          </div>

          <Show when={reservations.data} fallback={
            <div class="flex items-center justify-center py-12">
              <div class="w-6 h-6 border-2 border-blue-600 border-t-transparent rounded-full animate-spin" />
            </div>
          }>
            <div class="divide-y divide-gray-100">
              <For each={reservations.data}>
                {(r) => (
                  <div class="px-5 py-4 flex items-center justify-between hover:bg-gray-50 transition-colors">
                    <div class="flex items-center gap-3">
                      <div class="w-9 h-9 rounded-lg bg-blue-100 flex items-center justify-center text-blue-600 text-xs font-bold shrink-0">
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4"><rect x="2" y="3" width="20" height="14" rx="2" ry="2" /><line x1="8" y1="21" x2="16" y2="21" /><line x1="12" y1="17" x2="12" y2="21" /></svg>
                      </div>
                      <div>
                        <p class="text-sm font-medium text-gray-900">Seat {r.seat_id?.slice(0, 8)}</p>
                        <p class="text-xs text-gray-500">ID: {r.id?.slice(0, 8)}</p>
                      </div>
                    </div>
                    <button
                      onClick={() => {
                        cancelReservation.mutateAsync(r.id!, {
                          onSuccess: () => setErrorMsg(''),
                          onError: (err) => setErrorMsg(err instanceof ApiError ? err.message : 'Failed to cancel'),
                        });
                      }}
                      disabled={cancelReservation.isPending}
                      class="px-3 py-1.5 text-xs font-medium text-red-600 hover:bg-red-50 border border-red-200 rounded-lg transition-colors disabled:opacity-50"
                    >
                      Cancel
                    </button>
                  </div>
                )}
              </For>
              <Show when={reservations.data && reservations.data.length === 0}>
                <div class="text-center py-12">
                  <div class="w-12 h-12 rounded-xl bg-gray-50 flex items-center justify-center mx-auto mb-3">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6 text-gray-300"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" /><line x1="16" y1="2" x2="16" y2="6" /><line x1="8" y1="2" x2="8" y2="6" /><line x1="3" y1="10" x2="21" y2="10" /></svg>
                  </div>
                  <p class="text-gray-500 font-medium">No reservations on this date</p>
                  <p class="text-xs text-gray-400 mt-1">
                    <A href="/rooms" class="text-blue-600 hover:underline">Browse rooms</A> to reserve a seat
                  </p>
                </div>
              </Show>
            </div>
          </Show>
        </div>
      </Show>
    </div>
  );
};

export default Reservations;
