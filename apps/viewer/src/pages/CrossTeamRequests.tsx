import { Component, createMemo, createSignal, For, Show } from 'solid-js';
import {
  useAllCrossTeamRequests,
  usePendingByTeam,
  useApproveCrossTeamRequest,
  useRejectCrossTeamRequest,
  useUsers,
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
  const teamId = () => currentUser()?.team_id ?? '';
  const isSuperAdmin = () => currentUser()?.role === 'superadmin';
  const allUsers = useUsers(() => '');

  const [statusFilter, setStatusFilter] = createSignal('pending');

  const allRequests = useAllCrossTeamRequests();
  const pendingByTeam = usePendingByTeam(teamId);

  const filteredRequests = createMemo(() => {
    const all = allRequests.data ?? [];
    return all.filter((r) => r.status === statusFilter());
  });

  const approveMutation = useApproveCrossTeamRequest();
  const rejectMutation = useRejectCrossTeamRequest();

  const handleApprove = async (id: string) => {
    await approveMutation.mutateAsync(id);
  };

  const handleReject = async (id: string) => {
    await rejectMutation.mutateAsync(id);
  };

  const canManage = (req: { status?: string }) => {
    if (isSuperAdmin()) return req.status === 'pending';
    if (currentUser()?.role === 'team_admin') return req.status === 'pending';
    return false;
  };

  const isRequestForMyTeam = createMemo(() => {
    if (isSuperAdmin()) return () => true;
    const myTeamId = teamId();
    if (!myTeamId) return () => false;
    const pending = pendingByTeam.data ?? [];
    const ids = new Set(pending.map((r) => r.id));
    return (reqId: string | undefined) => reqId ? ids.has(reqId) : false;
  });

  const getName = (id: string | undefined) => {
    if (!id) return '-';
    return allUsers.data?.find((u) => u.id === id)?.name ?? id.slice(0, 8);
  };

  return (
    <div>
      <div class="mb-8">
        <h1 class="text-2xl font-bold text-gray-900">Cross-Team Requests</h1>
        <p class="text-sm text-gray-500 mt-1">
          {isSuperAdmin() ? 'All cross-team seat requests' : 'Requests for seats in your team'}
        </p>
      </div>

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
              <th class="px-5 py-3 font-medium">Requester</th>
              <th class="px-5 py-3 font-medium">Status</th>
              <th class="px-5 py-3 font-medium">Actions</th>
            </tr>
          </thead>
          <tbody>
            <For each={filteredRequests()}>
              {(req) => {
                const manage = canManage(req) && isRequestForMyTeam()(req.id);
                return (
                  <tr class="border-b last:border-0 hover:bg-gray-50 transition-colors">
                    <td class="px-5 py-3 text-gray-700">
                      {req.date ? new Date(req.date + 'T00:00:00').toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric' }) : '-'}
                    </td>
                    <td class="px-5 py-3 text-gray-700">{getName(req.requesting_user_id)}</td>
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
                    <td class="px-5 py-3">
                      <Show when={manage}>
                        <div class="flex items-center gap-2">
                          <button
                            onClick={() => handleApprove(req.id!)}
                            disabled={approveMutation.isPending || rejectMutation.isPending}
                            class="px-3 py-1 bg-green-600 text-white rounded-md text-xs hover:bg-green-700 disabled:opacity-40 transition-colors"
                          >
                            Approve
                          </button>
                          <button
                            onClick={() => handleReject(req.id!)}
                            disabled={approveMutation.isPending || rejectMutation.isPending}
                            class="px-3 py-1 bg-red-600 text-white rounded-md text-xs hover:bg-red-700 disabled:opacity-40 transition-colors"
                          >
                            Reject
                          </button>
                        </div>
                      </Show>
                    </td>
                  </tr>
                );
              }}
            </For>
            {filteredRequests().length === 0 && (
              <tr>
                <td colSpan="4" class="px-5 py-12 text-center text-gray-400">
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
    </div>
  );
};

export default CrossTeamRequests;
