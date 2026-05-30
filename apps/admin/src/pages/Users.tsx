import { Component, createSignal, createMemo, For, Show } from 'solid-js';
import {
  useTeams,
  useUsers,
  useCreateUser,
  useDeleteUser,
  useUpdateWeeklyLimit,
  useUpdateUserRole,
  useUpdateUserTeam,
  useSetPassword,
  ApiError,
} from '@neshiman/api-client';
import type { definitions } from '@neshiman/api-types';
import Modal from '../components/Modal';
import ConfirmModal from '../components/ConfirmModal';
import { showToast } from '../stores/toast';

type User = definitions['dto.UserResponse'];

function getCurrentUser() {
  try { return JSON.parse(localStorage.getItem('neshiman_user') ?? 'null'); }
  catch { return null; }
}

const Users: Component = () => {
  const teams = useTeams();

  const [selectedTeamId, setSelectedTeamId] = createSignal('');
  const users = useUsers(selectedTeamId);

  const createUser = useCreateUser();
  const deleteUser = useDeleteUser();
  const updateLimit = useUpdateWeeklyLimit();
  const updateRole = useUpdateUserRole();
  const updateTeam = useUpdateUserTeam();
  const setPassword = useSetPassword();

  const [showForm, setShowForm] = createSignal(false);
  const [userName, setUserName] = createSignal('');
  const [userEmail, setUserEmail] = createSignal('');
  const [userRole, setUserRole] = createSignal('viewer');
  const [userLimit, setUserLimit] = createSignal(2);
  const [error, setError] = createSignal('');

  const [editingLimit, setEditingLimit] = createSignal<{ id: string; limit: number } | null>(null);
  const [deleteTarget, setDeleteTarget] = createSignal<string | null>(null);

  const [editingUser, setEditingUser] = createSignal<User | null>(null);
  const [editRole, setEditRole] = createSignal('');
  const [editPassword, setEditPassword] = createSignal('');
  const [editTeamId, setEditTeamId] = createSignal('');

  const currentUser = createMemo(() => getCurrentUser());
  const isSuperAdmin = createMemo(() => currentUser()?.role === 'superadmin');
  const isTeamAdmin = createMemo(() => currentUser()?.role === 'team_admin');
  const canEdit = createMemo(() => isSuperAdmin() || isTeamAdmin());

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
    setDeleteTarget(id);
  };

  const confirmDeleteUser = async () => {
    const id = deleteTarget();
    if (!id) return;
    try {
      await deleteUser.mutateAsync(id);
      setDeleteTarget(null);
      showToast('User deleted.', 'success');
    } catch (err) {
      showToast(err instanceof ApiError ? err.message : 'Failed to delete user', 'error');
    }
  };

  const handleSaveLimit = async (id: string) => {
    const edit = editingLimit();
    if (!edit) return;
    try {
      await updateLimit.mutateAsync({ id, data: { weekly_limit: edit.limit } });
      setEditingLimit(null);
      showToast('Weekly limit updated.', 'success');
    } catch (err) {
      showToast(err instanceof ApiError ? err.message : 'Failed to update limit', 'error');
    }
  };

  const openEditUser = (user: User) => {
    setEditingUser(user);
    setEditRole(user.role ?? 'viewer');
    setEditPassword('');
    setEditTeamId(user.team_id ?? '');
  };

  const handleEditUser = async (e: Event) => {
    e.preventDefault();
    const user = editingUser();
    if (!user?.id) return;

    try {
      // Update role if changed
      if (editRole() !== (user.role ?? 'viewer')) {
        await updateRole.mutateAsync({ id: user.id, data: { role: editRole() } });
      }
      // Update team if changed and superadmin
      if (isSuperAdmin() && editTeamId() !== (user.team_id ?? '')) {
        await updateTeam.mutateAsync({ id: user.id, data: { team_id: editTeamId() as string | undefined } });
      }
      // Set password if provided
      if (editPassword()) {
        await setPassword.mutateAsync({ id: user.id, password: editPassword() });
      }
      setEditingUser(null);
      showToast('User updated.', 'success');
    } catch (err) {
      showToast(err instanceof ApiError ? err.message : 'Failed to update user', 'error');
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
                            max="5"
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
                    <td class="px-4 py-3 flex gap-2">
                      <Show when={canEdit()}>
                        <button onClick={() => openEditUser(user)} class="text-blue-600 hover:underline text-xs">
                          Edit
                        </button>
                      </Show>
                      <Show when={user.role !== 'superadmin'}>
                        <button onClick={() => handleDeleteUser(user.id!)} class="text-red-600 hover:underline text-xs">
                          Delete
                        </button>
                      </Show>
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
            <label class="block text-sm font-medium text-gray-700 mb-1">Email (optional)</label>
            <input
              type="text"
              value={userEmail()}
              onInput={(e) => setUserEmail(e.currentTarget.value)}
              placeholder="user@example.com"
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
              max="5"
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

      <Modal open={editingUser() !== null} onClose={() => setEditingUser(null)} title={`Edit User - ${editingUser()?.name ?? ''}`}>
        <form onSubmit={handleEditUser} class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Role</label>
            <select
              value={editRole()}
              onChange={(e) => setEditRole(e.currentTarget.value)}
              class="w-full border rounded-md px-3 py-2 text-sm"
            >
              <option value="viewer">Viewer</option>
              <option value="team_admin">Team Admin</option>
              <option value="superadmin">Superadmin</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">New Password</label>
            <input
              type="password"
              value={editPassword()}
              onInput={(e) => setEditPassword(e.currentTarget.value)}
              placeholder="Leave blank to keep current"
              class="w-full border rounded-md px-3 py-2 text-sm"
            />
          </div>
          <Show when={isSuperAdmin()}>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Team</label>
              <select
                value={editTeamId()}
                onChange={(e) => setEditTeamId(e.currentTarget.value)}
                class="w-full border rounded-md px-3 py-2 text-sm"
              >
                <option value="">No team</option>
                <For each={teams.data}>
                  {(team) => <option value={team.id!}>{team.name}</option>}
                </For>
              </select>
            </div>
          </Show>
          <div class="flex justify-end gap-3 pt-2">
            <button type="button" onClick={() => setEditingUser(null)} class="px-4 py-2 border rounded-md text-sm">
              Cancel
            </button>
            <button type="submit" disabled={updateRole.isPending || setPassword.isPending} class="px-4 py-2 bg-blue-600 text-white rounded-md text-sm">
              Save
            </button>
          </div>
        </form>
      </Modal>

      <ConfirmModal
        open={deleteTarget() !== null}
        title="Delete User"
        description="Are you sure you want to delete this user?"
        onConfirm={confirmDeleteUser}
        onCancel={() => setDeleteTarget(null)}
      />
    </div>
  );
};

export default Users;