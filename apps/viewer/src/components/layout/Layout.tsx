import { Component, createSignal, For } from 'solid-js';
import { A, useLocation, useNavigate } from '@solidjs/router';
import type { RouteSectionProps } from '@solidjs/router';
import { setAuthHeader } from '@neshiman/api-client';
import UserBar from '../ui/UserBar';

const navItems = [
  { href: '/', label: 'Dashboard', icon: '◷' },
  { href: '/rooms', label: 'Rooms', icon: '☰' },
  { href: '/reservations', label: 'Reservations', icon: '☰' },
  { href: '/requests', label: 'Requests', icon: '☰' },
];

const Layout: Component<RouteSectionProps> = (props) => {
  const location = useLocation();
  const navigate = useNavigate();
  const [currentUserId, setCurrentUserId] = createSignal(
    localStorage.getItem('viewer_user_id') ?? ''
  );

  const handleSelectUser = (id: string) => {
    localStorage.setItem('viewer_user_id', id);
    setCurrentUserId(id);
    setAuthHeader(() => id);
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
            <UserBar currentUserId={currentUserId()} onSelectUser={handleSelectUser} />
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
