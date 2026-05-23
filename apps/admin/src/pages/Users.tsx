import { Component, createSignal, createMemo, For, Show } from 'solid-js';
import {
  useTeams,
  useUsers,
  useCreateUser,
  useDeleteUser,
  useUpdateWeeklyLimit,
  ApiError,
} from '@neshiman/api-client';
import type { definitions } from '@neshiman/api-types';
import Modal from '../components/Modal';

type User = definitions['dto.UserResponse'];

const Users: Component = () => {
  const teams = useTeams();

  const [selectedTeamId, setSelectedTeamId] = createSignal('');
  const users = useUsers(selectedTeamId);

  const createUser = useCreateUser();
  const deleteUser = useDeleteUser();
  const updateLimit = useUpdateWeeklyLimit();

  const [showForm, setShowForm] = createSignal(false);
  const [userName, setUserName] = createSignal('');
  const [userEmail, setUserEmail] = createSignal('');
  const [userRole, setUserRole] = createSignal('viewer');
  const [userLimit, setUserLimit] = createSignal(2);
  const [error, setError] = createSignal('');

  const [editingLimit, setEditingLimit] = createSignal<{ id: string; limit: number } | null>(null);

  const teamName = createMemo(() => {
    const map = new Map<string, string>();
    for (const t of teams.data ?? []) {
      if (t.id && t.name) map.set(t.id, t.name);
    }
    return map;
  });

  const handleCreateUser = async (e: Event) => {
    e.preventDefault();
    setError('');
    try {
      await createUser.mutateAsync({
        name: userName(),
        email: userEmail(),
        role: userRole(),
        team_id: selectedTeamId(),
        weekly_limit: userLimit(),
      });
      setShowForm(false);
      setUserName('');
      setUserEmail('');
      setUserRole('viewer');
      setUserLimit(2);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Failed to create user');
    }
  };

  const handleDeleteUser = async (id: string) => {
    if (!confirm('Delete this user?')) return;
    try {
      await deleteUser.mutateAsync(id);
    } catch (err) {
      alert(err instanceof ApiError ? err.message : 'Failed to delete user');
    }
  };

  const handleSaveLimit = async (id: string) => {
    const edit = editingLimit();
    if (!edit) return;
    try {
      await updateLimit.mutateAsync({ id, data: { weekly_limit: edit.limit } });
      setEditingLimit(null);
    } catch (err) {
      alert(err instanceof ApiError ? err.message : 'Failed to update limit');
    }
  };

  return (
    <div class="p-6">
      <div class="flex items-center justify-between mb-6">
        <h1 class="text-2xl font-bold">Users</h1>
        <button
          onClick={() => { setShowForm(true); setError(''); }}
          disabled={!selectedTeamId()}
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 text-sm disabled:opacity-50"
        >
          + New User
        </button>
      </div>

      <div class="mb-4">
        <label class="block text-sm font-medium text-gray-700 mb-1">Filter by Team</label>
        <select
          value={selectedTeamId()}
          onChange={(e) => setSelectedTeamId(e.currentTarget.value)}
          class="border rounded-md px-3 py-2 text-sm w-64"
        >
          <option value="">All users</option>
          <For each={teams.data}>
            {(team) => <option value={team.id!}>{team.name}</option>}
          </For>
        </select>
      </div>

      <Show when={true}>
        <div class="bg-white rounded-lg shadow overflow-hidden">
          <table class="w-full text-sm">
            <thead>
              <tr class="text-left text-gray-500 border-b">
                <th class="px-4 py-3 font-medium">Name</th>
                <th class="px-4 py-3 font-medium">Email</th>
                <th class="px-4 py-3 font-medium">Team</th>
                <th class="px-4 py-3 font-medium">Role</th>
                <th class="px-4 py-3 font-medium">Weekly Limit</th>
                <th class="px-4 py-3 font-medium">Actions</th>
              </tr>
            </thead>
            <tbody>
              <For each={users.data}>
                {(user: User) => (
                  <tr class="border-b last:border-0 hover:bg-gray-50">
                    <td class="px-4 py-3">{user.name}</td>
                    <td class="px-4 py-3 text-gray-500">{user.email}</td>
                    <td class="px-4 py-3 text-gray-600">{user.team_id ? teamName().get(user.team_id) ?? '-' : '-'}</td>
                    <td class="px-4 py-3 capitalize">{user.role}</td>
                    <td class="px-4 py-3">
                      <Show
                        when={editingLimit()?.id === user.id}
                        fallback={
                          <button
                            onClick={() => setEditingLimit({ id: user.id!, limit: user.weekly_limit ?? 2 })}
                            class="text-blue-600 hover:underline"
                          >
                            {user.weekly_limit ?? '-'}/week
                          </button>
                        }
                      >
                        <div class="flex items-center gap-1">
                          <input
                            type="number"
                            min="0"
                            max="7"
                            value={editingLimit()?.limit ?? 2}
                            onInput={(e) => setEditingLimit({ id: user.id!, limit: Number(e.currentTarget.value) })}
                            class="w-16 border rounded px-2 py-1 text-sm"
                          />
                          <button onClick={() => handleSaveLimit(user.id!)} class="text-green-600 hover:underline text-xs">
                            Save
                          </button>
                          <button onClick={() => setEditingLimit(null)} class="text-gray-400 hover:underline text-xs">
                            Cancel
                          </button>
                        </div>
                      </Show>
                    </td>
                    <td class="px-4 py-3">
                      <button onClick={() => handleDeleteUser(user.id!)} class="text-red-600 hover:underline text-xs">
                        Delete
                      </button>
                    </td>
                  </tr>
                )}
              </For>
              {users.data && users.data.length === 0 && (
                <tr>
                  <td colspan="6" class="px-4 py-8 text-center text-gray-400">No users in this team</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </Show>

      <Modal open={showForm()} onClose={() => setShowForm(false)} title="New User">
        <form onSubmit={handleCreateUser} class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
            <input
              value={userName()}
              onInput={(e) => setUserName(e.currentTarget.value)}
              required
              class="w-full border rounded-md px-3 py-2 text-sm"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
            <input
              type="email"
              value={userEmail()}
              onInput={(e) => setUserEmail(e.currentTarget.value)}
              required
              class="w-full border rounded-md px-3 py-2 text-sm"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Role</label>
            <select
              value={userRole()}
              onChange={(e) => setUserRole(e.currentTarget.value)}
              class="w-full border rounded-md px-3 py-2 text-sm"
            >
              <option value="viewer">Viewer</option>
              <option value="team_admin">Team Admin</option>
              <option value="superadmin">Superadmin</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Weekly Limit (days)</label>
            <input
              type="number"
              min="0"
              max="7"
              value={userLimit()}
              onInput={(e) => setUserLimit(Number(e.currentTarget.value))}
              class="w-full border rounded-md px-3 py-2 text-sm"
            />
          </div>
          {error() && <p class="text-red-600 text-sm">{error()}</p>}
          <div class="flex justify-end gap-3 pt-2">
            <button type="button" onClick={() => setShowForm(false)} class="px-4 py-2 border rounded-md text-sm">
              Cancel
            </button>
            <button type="submit" disabled={createUser.isPending} class="px-4 py-2 bg-blue-600 text-white rounded-md text-sm">
              Create
            </button>
          </div>
        </form>
      </Modal>
    </div>
  );
};

export default Users;
