import { Component, createMemo, For } from 'solid-js';
import { useUsers } from '@neshiman/api-client';

interface UserBarProps {
  currentUserId: string;
  onSelectUser: (id: string) => void;
}

const UserBar: Component<UserBarProps> = (props) => {
  const users = useUsers(() => '');

  const currentUser = createMemo(() =>
    users.data?.find((u) => u.id === props.currentUserId)
  );

  return (
    <div class="flex items-center gap-3">
      <div class="flex items-center gap-2 text-sm">
        <div class="w-7 h-7 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white text-xs font-bold shrink-0">
          {currentUser()?.name?.charAt(0)?.toUpperCase() ?? '?'}
        </div>
        <select
          value={props.currentUserId}
          onChange={(e) => props.onSelectUser(e.currentTarget.value)}
          class="border-0 bg-transparent text-sm font-medium text-gray-700 cursor-pointer focus:ring-0 p-0 pr-6 appearance-none"
        >
          <For each={users.data}>
            {(u) => <option value={u.id!}>{u.name}</option>}
          </For>
        </select>
      </div>
      <span class="text-xs text-gray-400 hidden sm:inline">
        {currentUser()?.role === 'superadmin' ? 'Super Admin' : currentUser()?.role === 'team_admin' ? 'Team Admin' : 'Viewer'}
      </span>
    </div>
  );
};

export default UserBar;
