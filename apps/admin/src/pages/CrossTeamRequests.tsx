import { Component, createSignal, createMemo, For, Show } from 'solid-js';
import {
  useCrossTeamRequests,
  useApproveCrossTeamRequest,
  useRejectCrossTeamRequest,
  useUsers,
  useSeat,
  ApiError,
} from '@neshiman/api-client';
import type { definitions } from '@neshiman/api-types';

type CrossTeamRequest = definitions['dto.CrossTeamRequestResponse'];

const statusTabs = [
  { value: 'pending', label: 'Pending' },
  { value: 'approved', label: 'Approved' },
  { value: 'rejected', label: 'Rejected' },
] as const;

const SeatLabel: Component<{ seatId: string }> = (props) => {
  const seat = useSeat(() => props.seatId);
  return <span class="text-gray-500">{seat.data?.label ?? props.seatId.slice(0, 8)}</span>;
};

const CrossTeamRequests: Component = () => {
  const [statusFilter, setStatusFilter] = createSignal('pending');
  const requests = useCrossTeamRequests(statusFilter);
  const allUsers = useUsers(() => '');
  const approve = useApproveCrossTeamRequest();
  const reject = useRejectCrossTeamRequest();

  const userName = createMemo(() => {
    const map = new Map<string, string>();
    for (const u of allUsers.data ?? []) {
      if (u.id && u.name) map.set(u.id, u.name);
    }
    return map;
  });

  const handleApprove = async (id: string) => {
    try {
      await approve.mutateAsync(id);
    } catch (err) {
      alert(err instanceof ApiError ? err.message : 'Failed to approve request');
    }
  };

  const handleReject = async (id: string) => {
    try {
      await reject.mutateAsync(id);
    } catch (err) {
      alert(err instanceof ApiError ? err.message : 'Failed to reject request');
    }
  };

  return (
    <div>
      <h1 class="text-2xl font-bold mb-6">Cross-Team Requests</h1>

      <div class="flex gap-2 mb-4">
        <For each={statusTabs}>
          {(tab) => (
            <button
              onClick={() => setStatusFilter(tab.value)}
              class="px-4 py-2 text-sm rounded-md transition-colors capitalize"
              classList={{
                'bg-blue-600 text-white': statusFilter() === tab.value,
                'bg-gray-100 text-gray-600 hover:bg-gray-200': statusFilter() !== tab.value,
              }}
            >
              {tab.label}
            </button>
          )}
        </For>
      </div>

      <div class="bg-white rounded-lg shadow overflow-hidden">
        <table class="w-full text-sm">
          <thead>
            <tr class="text-left text-gray-500 border-b">
              <th class="px-4 py-3 font-medium">Date</th>
              <th class="px-4 py-3 font-medium">Requesting User</th>
              <th class="px-4 py-3 font-medium">Target Seat</th>
              <th class="px-4 py-3 font-medium">Status</th>
              <th class="px-4 py-3 font-medium">Actions</th>
            </tr>
          </thead>
          <tbody>
            <For each={requests.data}>
              {(req: CrossTeamRequest) => (
                <tr class="border-b last:border-0 hover:bg-gray-50">
                  <td class="px-4 py-3">{req.date}</td>
                  <td class="px-4 py-3 text-gray-700">{req.requesting_user_id ? userName().get(req.requesting_user_id) ?? req.requesting_user_id.slice(0, 12) : '-'}</td>
                  <td class="px-4 py-3">
                    <Show when={req.target_seat_id} fallback={<span class="text-gray-400">-</span>}>
                      <SeatLabel seatId={req.target_seat_id!} />
                    </Show>
                  </td>
                  <td class="px-4 py-3 capitalize">
                    <span
                      class="inline-block px-2 py-0.5 rounded-full text-xs font-medium"
                      classList={{
                        'bg-yellow-100 text-yellow-800': req.status === 'pending',
                        'bg-green-100 text-green-800': req.status === 'approved',
                        'bg-red-100 text-red-800': req.status === 'rejected',
                      }}
                    >
                      {req.status}
                    </span>
                  </td>
                  <td class="px-4 py-3">
                    <Show when={req.status === 'pending'}>
                      <div class="flex gap-2">
                        <button
                          onClick={() => handleApprove(req.id!)}
                          disabled={approve.isPending}
                          class="px-3 py-1 bg-green-600 text-white rounded text-xs hover:bg-green-700 disabled:opacity-50"
                        >
                          Approve
                        </button>
                        <button
                          onClick={() => handleReject(req.id!)}
                          disabled={reject.isPending}
                          class="px-3 py-1 bg-red-600 text-white rounded text-xs hover:bg-red-700 disabled:opacity-50"
                        >
                          Reject
                        </button>
                      </div>
                    </Show>
                  </td>
                </tr>
              )}
            </For>
            {requests.data && requests.data.length === 0 && (
              <tr>
                <td colspan="5" class="px-4 py-8 text-center text-gray-400">
                  No {statusFilter()} requests
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default CrossTeamRequests;
