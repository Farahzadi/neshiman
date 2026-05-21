import { Component, For } from 'solid-js';
import { A, useLocation } from '@solidjs/router';
import type { RouteSectionProps } from '@solidjs/router';
import UserBar from './UserBar';

const navItems = [
  { href: '/', label: 'Dashboard' },
  { href: '/rooms', label: 'Rooms' },
  { href: '/teams', label: 'Teams' },
  { href: '/users', label: 'Users' },
  { href: '/cross-team-requests', label: 'Cross-Team Requests' },
];

const Layout: Component<RouteSectionProps> = (props) => {
  const location = useLocation();

  return (
    <div class="min-h-screen flex bg-gray-50">
      <nav class="w-56 bg-gray-900 text-white flex flex-col shrink-0">
        <div class="p-4 border-b border-gray-700">
          <h1 class="text-lg font-bold">Neshiman</h1>
          <p class="text-xs text-gray-400">Admin Panel</p>
        </div>
        <ul class="flex-1 p-3 space-y-1">
          <For each={navItems}>
            {(item) => (
              <li>
                <A
                  href={item.href}
                  class="block px-3 py-2 rounded-md text-sm transition-colors"
                  classList={{
                    'bg-gray-700 text-white': location.pathname === item.href || (item.href !== '/' && location.pathname.startsWith(item.href)),
                    'text-gray-300 hover:bg-gray-800 hover:text-white': !(location.pathname === item.href || (item.href !== '/' && location.pathname.startsWith(item.href))),
                  }}
                >
                  {item.label}
                </A>
              </li>
            )}
          </For>
        </ul>
      </nav>
      <main class="flex-1 flex flex-col overflow-auto">
        <UserBar />
        <div class="p-6 flex-1">
          {props.children}
        </div>
      </main>
    </div>
  );
};

export default Layout;
