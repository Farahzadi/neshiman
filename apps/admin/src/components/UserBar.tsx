import { Component, createSignal, createEffect, For, Show } from 'solid-js';
import { useUsers, setAuthHeader } from '@neshiman/api-client';

const STORAGE_KEY = 'neshiman_current_user';

function loadUserId(): string | null {
  try { return localStorage.getItem(STORAGE_KEY); } catch { return null; }
}

function saveUserId(id: string) {
  try { localStorage.setItem(STORAGE_KEY, id); } catch { /* noop */ }
}

const UserBar: Component = () => {
  const users = useUsers(() => '');
  const [userId, setUserId] = createSignal(loadUserId() ?? '');

  createEffect(() => {
    if (userId()) {
      setAuthHeader(() => userId());
    }
  });

  const handleSelect = (e: Event) => {
    const id = (e.currentTarget as HTMLSelectElement).value;
    setUserId(id);
    saveUserId(id);
  };

  const currentUser = () => users.data?.find((u) => u.id === userId());

  return (
    <div class="flex items-center justify-end gap-3 px-6 py-2 bg-white border-b">
      <span class="text-sm text-gray-500">Logged in as:</span>
      <select
        value={userId()}
        onChange={handleSelect}
        class="border rounded-md px-3 py-1.5 text-sm bg-white min-w-40"
      >
        <option value="">Select a user</option>
        <For each={users.data}>
          {(user) => (
            <option value={user.id!}>
              {user.name} ({user.role})
            </option>
          )}
        </For>
      </select>
      <Show when={currentUser()}>
        {(cu) => (
          <span class="text-xs text-gray-400">{cu().email}</span>
        )}
      </Show>
    </div>
  );
};

export default UserBar;
