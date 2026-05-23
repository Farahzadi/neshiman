import { Component, createMemo, For } from 'solid-js';
import { A, useLocation, useNavigate } from '@solidjs/router';
import type { RouteSectionProps } from '@solidjs/router';
import { setAuthHeader } from '@neshiman/api-client';

const navItems = [
  { href: '/', label: 'Dashboard' },
  { href: '/rooms', label: 'Rooms' },
  { href: '/reservations', label: 'Reservations' },
  { href: '/requests', label: 'Requests' },
];

function getUser() {
  try {
    const raw = localStorage.getItem('neshiman_user');
    return raw ? JSON.parse(raw) : null;
  } catch { return null; }
}

const Layout: Component<RouteSectionProps> = (props) => {
  const location = useLocation();
  const navigate = useNavigate();

  const currentUser = createMemo(() => getUser());

  const handleLogout = () => {
    localStorage.removeItem('neshiman_token');
    localStorage.removeItem('neshiman_user');
    setAuthHeader(() => null);
    navigate('/login', { replace: true });
  };

  return (
    <div class="min-h-screen bg-gradient-to-br from-gray-50 to-blue-50">
      <header class="bg-white border-b border-gray-200 shadow-sm sticky top-0 z-30">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex items-center justify-between h-16">
            <div class="flex items-center gap-8">
              <A href="/" class="flex items-center gap-2.5 text-lg font-bold text-gray-900">
                <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-600 to-indigo-700 flex items-center justify-center text-white text-sm font-bold shadow-sm">
                  N
                </div>
                <span>Neshiman</span>
              </A>
              <nav class="hidden md:flex items-center gap-1">
                <For each={navItems}>
                  {(item) => {
                    const isActive = location.pathname === item.href ||
                      (item.href !== '/' && location.pathname.startsWith(item.href));
                    return (
                      <A
                        href={item.href}
                        class="px-3 py-2 text-sm font-medium rounded-lg transition-colors"
                        classList={{
                          'bg-blue-50 text-blue-700': isActive,
                          'text-gray-600 hover:text-gray-900 hover:bg-gray-50': !isActive,
                        }}
                      >
                        {item.label}
                      </A>
                    );
                  }}
                </For>
              </nav>
            </div>
            <div class="flex items-center gap-3">
              <div class="flex items-center gap-2 text-sm">
                <div class="w-7 h-7 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white text-xs font-bold shrink-0">
                  {currentUser()?.name?.charAt(0)?.toUpperCase() ?? '?'}
                </div>
                <span class="text-sm font-medium text-gray-700">{currentUser()?.name ?? 'User'}</span>
                <span class="text-xs text-gray-400 hidden sm:inline">{currentUser()?.role}</span>
              </div>
              <button
                onClick={handleLogout}
                class="text-xs text-gray-500 hover:text-red-600 transition-colors ml-2"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </header>

      <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {props.children}
      </main>
    </div>
  );
};

export default Layout;
