import { Component, createSignal, For, Show } from 'solid-js';
import {
  useTeams,
  useCreateTeam,
  useDeleteTeam,
  useUsers,
  useCreateUser,
  useDeleteUser,
  useUpdateWeeklyLimit,
  ApiError,
} from '@neshiman/api-client';
import Modal from '../components/Modal';

const Teams: Component = () => {
  const teams = useTeams();
  const createTeam = useCreateTeam();
  const deleteTeam = useDeleteTeam();

  const [showTeamForm, setShowTeamForm] = createSignal(false);
  const [teamName, setTeamName] = createSignal('');
  const [teamError, setTeamError] = createSignal('');

  const [expandedTeamId, setExpandedTeamId] = createSignal<string | null>(null);
  const teamUsers = useUsers(() => expandedTeamId() ?? '');

  const [showUserForm, setShowUserForm] = createSignal(false);
  const [userName, setUserName] = createSignal('');
  const [userEmail, setUserEmail] = createSignal('');
  const [userRole, setUserRole] = createSignal('viewer');
  const [userLimit, setUserLimit] = createSignal(2);
  const [userError, setUserError] = createSignal('');

  const [editingLimit, setEditingLimit] = createSignal<{ id: string; limit: number } | null>(null);

  const createUser = useCreateUser();
  const deleteUser = useDeleteUser();
  const updateLimit = useUpdateWeeklyLimit();

  const handleCreateTeam = async (e: Event) => {
    e.preventDefault();
    setTeamError('');
    try {
      await createTeam.mutateAsync({ name: teamName() });
      setShowTeamForm(false);
      setTeamName('');
    } catch (err) {
      setTeamError(err instanceof ApiError ? err.message : 'Failed to create team');
    }
  };

  const handleDeleteTeam = async (id: string) => {
    if (!confirm('Delete this team and all its users?')) return;
    try {
      await deleteTeam.mutateAsync(id);
      if (expandedTeamId() === id) setExpandedTeamId(null);
    } catch (err) {
      alert(err instanceof ApiError ? err.message : 'Failed to delete team');
    }
  };

  const handleCreateUser = async (e: Event) => {
    e.preventDefault();
    setUserError('');
    try {
      await createUser.mutateAsync({
        name: userName(),
        email: userEmail(),
        role: userRole(),
        team_id: expandedTeamId()!,
        weekly_limit: userLimit(),
      });
      setShowUserForm(false);
      setUserName('');
      setUserEmail('');
      setUserRole('viewer');
      setUserLimit(2);
    } catch (err) {
      setUserError(err instanceof ApiError ? err.message : 'Failed to create user');
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

  const handleUpdateLimit = async (id: string) => {
    const edit = editingLimit();
    if (!edit) return;
    try {
      await updateLimit.mutateAsync({ id, data: { weekly_limit: edit.limit } });
      setEditingLimit(null);
    } catch (err) {
      alert(err instanceof ApiError ? err.message : 'Failed to update limit');
    }
  };

  const toggleTeam = (teamId: string) => {
    setExpandedTeamId((prev) => (prev === teamId ? null : teamId));
    setEditingLimit(null);
  };

  return (
    <div class="p-6">
      <div class="flex items-center justify-between mb-6">
        <h1 class="text-2xl font-bold">Teams</h1>
        <button
          onClick={() => { setTeamName(''); setTeamError(''); setShowTeamForm(true); }}
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 text-sm"
        >
          + New Team
        </button>
      </div>

      <div class="space-y-3">
        <For each={teams.data}>
          {(team) => (
            <div class="bg-white rounded-lg shadow">
              <div
                onClick={() => toggleTeam(team.id!)}
                class="flex items-center justify-between px-4 py-3 cursor-pointer hover:bg-gray-50"
              >
                <span class="font-medium">{team.name}</span>
                <div class="flex items-center gap-3">
                  <span class="text-sm text-gray-400">{expandedTeamId() === team.id ? '▲' : '▼'}</span>
                  <button
                    onClick={(e) => { e.stopPropagation(); handleDeleteTeam(team.id!); }}
                    class="text-sm text-red-600 hover:underline"
                  >
                    Delete
                  </button>
                </div>
              </div>
              <Show when={expandedTeamId() === team.id}>
                <div class="border-t px-4 py-3">
                  <div class="flex items-center justify-between mb-3">
                    <h3 class="text-sm font-medium text-gray-500">Members</h3>
                    <button
                      onClick={() => { setShowUserForm(true); setUserError(''); }}
                      class="text-sm text-blue-600 hover:underline"
                    >
                      + Add Member
                    </button>
                  </div>
                  <table class="w-full text-sm">
                    <thead>
                      <tr class="text-left text-gray-400 border-b">
                        <th class="pb-2 font-medium">Name</th>
                        <th class="pb-2 font-medium">Email</th>
                        <th class="pb-2 font-medium">Role</th>
                        <th class="pb-2 font-medium">Weekly Limit</th>
                        <th class="pb-2 font-medium">Actions</th>
                      </tr>
                    </thead>
                    <tbody>
                      <For each={teamUsers.data}>
                        {(user) => (
                          <tr class="border-b last:border-0">
                            <td class="py-2">{user.name}</td>
                            <td class="py-2 text-gray-500">{user.email}</td>
                            <td class="py-2 capitalize">{user.role}</td>
                            <td class="py-2">
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
                                  <button
                                    onClick={() => handleUpdateLimit(user.id!)}
                                    class="text-green-600 hover:underline text-xs"
                                  >
                                    Save
                                  </button>
                                  <button
                                    onClick={() => setEditingLimit(null)}
                                    class="text-gray-400 hover:underline text-xs"
                                  >
                                    Cancel
                                  </button>
                                </div>
                              </Show>
                            </td>
                            <td class="py-2">
                              <button
                                onClick={() => handleDeleteUser(user.id!)}
                                class="text-red-600 hover:underline text-xs"
                              >
                                Remove
                              </button>
                            </td>
                          </tr>
                        )}
                      </For>
                      {teamUsers.data && teamUsers.data.length === 0 && (
                        <tr>
                          <td colspan="5" class="py-4 text-center text-gray-400">No members yet</td>
                        </tr>
                      )}
                    </tbody>
                  </table>
                </div>
              </Show>
            </div>
          )}
        </For>
        {teams.data && teams.data.length === 0 && (
          <p class="text-center text-gray-400 py-8">No teams yet. Create one!</p>
        )}
      </div>

      <Modal open={showTeamForm()} onClose={() => setShowTeamForm(false)} title="New Team">
        <form onSubmit={handleCreateTeam} class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
            <input
              value={teamName()}
              onInput={(e) => setTeamName(e.currentTarget.value)}
              required
              class="w-full border rounded-md px-3 py-2 text-sm"
              placeholder="Engineering"
            />
          </div>
          {teamError() && <p class="text-red-600 text-sm">{teamError()}</p>}
          <div class="flex justify-end gap-3 pt-2">
            <button type="button" onClick={() => setShowTeamForm(false)} class="px-4 py-2 border rounded-md text-sm">
              Cancel
            </button>
            <button type="submit" class="px-4 py-2 bg-blue-600 text-white rounded-md text-sm">
              Create
            </button>
          </div>
        </form>
      </Modal>

      <Modal open={showUserForm()} onClose={() => setShowUserForm(false)} title="Add Member">
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
          {userError() && <p class="text-red-600 text-sm">{userError()}</p>}
          <div class="flex justify-end gap-3 pt-2">
            <button type="button" onClick={() => setShowUserForm(false)} class="px-4 py-2 border rounded-md text-sm">
              Cancel
            </button>
            <button type="submit" class="px-4 py-2 bg-blue-600 text-white rounded-md text-sm">
              Add
            </button>
          </div>
        </form>
      </Modal>
    </div>
  );
};

export default Teams;
