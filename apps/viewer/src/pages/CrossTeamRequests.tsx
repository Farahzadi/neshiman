import { Component, createMemo, createSignal, For, Show } from 'solid-js';
import {
  useCrossTeamRequests,
  useSeat,
} from '@neshiman/api-client';

const statusTabs = [
  { value: 'pending', label: 'Pending' },
  { value: 'approved', label: 'Approved' },
  { value: 'rejected', label: 'Rejected' },
] as const;

const CrossTeamRequests: Component = () => {
  const currentUser = createMemo(() => {
    try { return JSON.parse(localStorage.getItem('neshiman_user') ?? 'null'); }
    catch { return null; }
  });
  const userId = () => currentUser()?.id ?? '';

  const [statusFilter, setStatusFilter] = createSignal('pending');
  const requests = useCrossTeamRequests(statusFilter);

  return (
    <div>
      <div class="mb-8">
        <h1 class="text-2xl font-bold text-gray-900">Cross-Team Requests</h1>
        <p class="text-sm text-gray-500 mt-1">
          <Show when={currentUser()} fallback="Select a user from the top-right">
            {currentUser()?.name}'s cross-team seat requests
          </Show>
        </p>
      </div>

      <Show when={!!userId()} fallback={
        <div class="text-center py-20">
          <div class="w-14 h-14 rounded-2xl bg-gray-100 flex items-center justify-center mx-auto mb-3">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-7 h-7 text-gray-300"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12" /></svg>
          </div>
          <p class="text-gray-500">Please sign in to view requests.</p>
        </div>
      }>
        <div class="flex gap-2 mb-4">
          <For each={statusTabs}>
            {(tab) => (
              <button
                onClick={() => setStatusFilter(tab.value)}
                class="px-4 py-2 text-sm rounded-lg font-medium transition-colors"
                classList={{
                  'bg-blue-600 text-white shadow-sm': statusFilter() === tab.value,
                  'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50': statusFilter() !== tab.value,
                }}
              >
                {tab.label}
              </button>
            )}
          </For>
        </div>

        <div class="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
          <table class="w-full text-sm">
            <thead>
              <tr class="text-left text-gray-500 border-b bg-gray-50/50">
                <th class="px-5 py-3 font-medium">Date</th>
                <th class="px-5 py-3 font-medium">Seat</th>
                <th class="px-5 py-3 font-medium">Status</th>
              </tr>
            </thead>
            <tbody>
              <For each={requests.data}>
                {(req) => (
                  <tr class="border-b last:border-0 hover:bg-gray-50 transition-colors">
                    <td class="px-5 py-3 text-gray-700">
                      {req.date ? new Date(req.date + 'T00:00:00').toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric' }) : '-'}
                    </td>
                    <td class="px-5 py-3">
                      <SeatLabel seatId={req.target_seat_id} />
                    </td>
                    <td class="px-5 py-3">
                      <span
                        class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium"
                        classList={{
                          'bg-yellow-50 text-yellow-700 border border-yellow-200': req.status === 'pending',
                          'bg-green-50 text-green-700 border border-green-200': req.status === 'approved',
                          'bg-red-50 text-red-700 border border-red-200': req.status === 'rejected',
                        }}
                      >
                        <span
                          class="w-1.5 h-1.5 rounded-full"
                          classList={{
                            'bg-yellow-500': req.status === 'pending',
                            'bg-green-500': req.status === 'approved',
                            'bg-red-500': req.status === 'rejected',
                          }}
                        />
                        {req.status}
                      </span>
                    </td>
                  </tr>
                )}
              </For>
              {requests.data && requests.data.length === 0 && (
                <tr>
                  <td colSpan="3" class="px-5 py-12 text-center text-gray-400">
                    <div class="w-10 h-10 rounded-xl bg-gray-50 flex items-center justify-center mx-auto mb-2">
                      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-gray-300"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12" /></svg>
                    </div>
                    No {statusFilter()} requests
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </Show>
    </div>
  );
};

const SeatLabel: Component<{ seatId: string | undefined }> = (props) => {
  const seat = useSeat(() => props.seatId ?? '');
  return <span class="text-gray-700">{seat.data?.label ?? props.seatId?.slice(0, 8) ?? '-'}</span>;
};

export default CrossTeamRequests;
