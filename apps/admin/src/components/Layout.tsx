import { Component, createMemo, For, createSignal } from 'solid-js';
import { A, useNavigate } from '@solidjs/router';
import type { RouteSectionProps } from '@solidjs/router';
import { setAuthHeader } from '@neshiman/api-client';

const icons: Record<string, () => any> = {
  Dashboard: () => (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 shrink-0">
      <rect x="3" y="3" width="7" height="7" />
      <rect x="14" y="3" width="7" height="7" />
      <rect x="14" y="14" width="7" height="7" />
      <rect x="3" y="14" width="7" height="7" />
    </svg>
  ),
  Rooms: () => (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 shrink-0">
      <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
      <polyline points="9 22 9 12 15 12 15 22" />
    </svg>
  ),
  Teams: () => (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 shrink-0">
      <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
      <circle cx="9" cy="7" r="4" />
      <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
      <path d="M16 3.13a4 4 0 0 1 0 7.75" />
    </svg>
  ),
  Users: () => (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 shrink-0">
      <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
      <circle cx="12" cy="7" r="4" />
    </svg>
  ),
  'Cross-Team Requests': () => (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 shrink-0">
      <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
    </svg>
  ),
};

const navItems = [
  { href: '/', label: 'Dashboard' },
  { href: '/rooms', label: 'Rooms' },
  { href: '/teams', label: 'Teams' },
  { href: '/users', label: 'Users' },
  { href: '/cross-team-requests', label: 'Cross-Team Requests' },
];

const CollapseIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 shrink-0">
    <polyline points="15 18 9 12 15 6" />
  </svg>
);

const ExpandIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 shrink-0">
    <polyline points="9 18 15 12 9 6" />
  </svg>
);

function getCurrentUser() {
  try { return JSON.parse(localStorage.getItem('neshiman_user') ?? 'null'); }
  catch { return null; }
}

const Layout: Component<RouteSectionProps> = (props) => {
  const navigate = useNavigate();
  const [collapsed, setCollapsed] = createSignal(false);
  const currentUser = createMemo(() => getCurrentUser());

  const handleLogout = () => {
    localStorage.removeItem('neshiman_token');
    localStorage.removeItem('neshiman_user');
    setAuthHeader(() => null);
    navigate('/login', { replace: true });
  };

  return (
    <div class="min-h-screen flex bg-gray-50">
      <nav
        class="bg-gray-900 text-white flex flex-col shrink-0 transition-[width] duration-200"
        classList={{ 'w-56': !collapsed(), 'w-16': collapsed() }}
      >
        <div class="p-4 border-b border-gray-700 flex items-center gap-3 overflow-hidden whitespace-nowrap">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6 shrink-0">
            <rect x="2" y="3" width="20" height="14" rx="2" ry="2" />
            <line x1="8" y1="21" x2="16" y2="21" />
            <line x1="12" y1="17" x2="12" y2="21" />
          </svg>
          <span classList={{ 'opacity-100': !collapsed(), 'opacity-0': collapsed() }}>
            <h1 class="text-lg font-bold">Neshiman</h1>
            <p class="text-xs text-gray-400">Admin Panel</p>
          </span>
        </div>
        <ul class="flex-1 p-3 space-y-1">
          <For each={navItems}>
            {(item) => {
              const Icon = icons[item.label];
              return (
                <li>
                  <A
                    href={item.href}
                    end={item.href === '/'}
                    activeClass="bg-gray-700 text-white"
                    inactiveClass="text-gray-300 hover:bg-gray-800 hover:text-white"
                    class="flex items-center gap-3 px-3 py-2 rounded-md text-sm transition-colors overflow-hidden whitespace-nowrap"
                  >
                    <Icon />
                    <span classList={{ 'opacity-100': !collapsed(), 'opacity-0': collapsed() }}>{item.label}</span>
                  </A>
                </li>
              );
            }}
          </For>
        </ul>
        <button
          onClick={() => setCollapsed((c) => !c)}
          class="flex items-center justify-center gap-3 p-4 border-t border-gray-700 text-gray-400 hover:text-white hover:bg-gray-800 transition-colors"
        >
          {collapsed() ? <ExpandIcon /> : <CollapseIcon />}
          <span class="text-sm" classList={{ 'opacity-100': !collapsed(), 'opacity-0': collapsed() }}>Collapse</span>
        </button>
      </nav>
      <main class="flex-1 flex flex-col overflow-hidden">
        <div class="flex items-center justify-end gap-3 px-6 py-2 bg-white border-b text-sm">
          <div class="flex items-center gap-2">
            <div class="w-7 h-7 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white text-xs font-bold shrink-0">
              {currentUser()?.name?.charAt(0)?.toUpperCase() ?? '?'}
            </div>
            <span class="font-medium text-gray-700">{currentUser()?.name ?? 'User'}</span>
            <span class="text-xs text-gray-400">{currentUser()?.role}</span>
          </div>
          <button
            onClick={handleLogout}
            class="text-xs text-gray-500 hover:text-red-600 transition-colors ml-2"
          >
            Logout
          </button>
        </div>
        <div class="flex flex-col flex-1 overflow-auto">
          {props.children}
        </div>
      </main>
    </div>
  );
};

export default Layout;
