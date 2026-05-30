import { Component, createMemo } from 'solid-js';
import { A } from '@solidjs/router';
import { useRooms, useTeams } from '@neshiman/api-client';

function getCurrentUser() {
  try {
    return JSON.parse(localStorage.getItem('neshiman_user') ?? 'null');
  } catch { return null; }
}

const Dashboard: Component = () => {
  const currentUser = createMemo(() => getCurrentUser());
  const rooms = useRooms();
  const teams = useTeams();

  return (
    <div>
      <div class="mb-8">
        <div class="flex items-center gap-3 mb-1">
          <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white text-sm font-bold shadow-sm">
            {currentUser()?.name?.charAt(0)?.toUpperCase() ?? '?'}
          </div>
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Welcome back, {currentUser()?.name ?? 'User'}</h1>
            <p class="text-sm text-gray-500">Here's your seat reservation overview</p>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
        <div class="bg-white rounded-xl border border-gray-200 p-5 shadow-sm hover:shadow-md transition-shadow">
          <div class="flex items-center gap-3">
            <div class="w-11 h-11 rounded-xl bg-blue-100 flex items-center justify-center text-blue-600 shrink-0">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" /><polyline points="9 22 9 12 15 12 15 22" /></svg>
            </div>
            <div>
              <p class="text-xs text-gray-500 font-medium uppercase tracking-wider">Rooms</p>
              <p class="text-2xl font-bold text-gray-900">{rooms.data?.length ?? '-'}</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-5 shadow-sm hover:shadow-md transition-shadow">
          <div class="flex items-center gap-3">
            <div class="w-11 h-11 rounded-xl bg-green-100 flex items-center justify-center text-green-600 shrink-0">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" /><circle cx="9" cy="7" r="4" /><path d="M23 21v-2a4 4 0 0 0-3-3.87" /><path d="M16 3.13a4 4 0 0 1 0 7.75" /></svg>
            </div>
            <div>
              <p class="text-xs text-gray-500 font-medium uppercase tracking-wider">Teams</p>
              <p class="text-2xl font-bold text-gray-900">{teams.data?.length ?? '-'}</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-5 shadow-sm hover:shadow-md transition-shadow">
          <div class="flex items-center gap-3">
            <div class="w-11 h-11 rounded-xl bg-amber-100 flex items-center justify-center text-amber-600 shrink-0">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" /><circle cx="12" cy="7" r="4" /></svg>
            </div>
            <div>
              <p class="text-xs text-gray-500 font-medium uppercase tracking-wider">Weekly Limit</p>
              <p class="text-2xl font-bold text-gray-900">{currentUser()?.weekly_limit ?? '-'}<span class="text-sm font-normal text-gray-400"> days</span></p>
            </div>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div class="bg-white rounded-xl border border-gray-200 shadow-sm">
          <div class="px-5 py-4 border-b border-gray-100">
            <h2 class="font-semibold text-gray-900">Quick Actions</h2>
          </div>
          <div class="p-5 space-y-2">
            <A href="/rooms" class="flex items-center gap-3 p-3 rounded-lg hover:bg-blue-50 transition-colors group">
              <div class="w-10 h-10 rounded-xl bg-blue-100 flex items-center justify-center text-blue-600 group-hover:bg-blue-200 transition-colors shrink-0">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" /><polyline points="9 22 9 12 15 12 15 22" /></svg>
              </div>
              <div class="flex-1">
                <p class="text-sm font-medium text-gray-900">Browse Rooms</p>
                <p class="text-xs text-gray-500">View seat maps and reserve a seat</p>
              </div>
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 text-gray-300 group-hover:text-gray-500"><polyline points="9 18 15 12 9 6" /></svg>
            </A>
            <A href="/reservations" class="flex items-center gap-3 p-3 rounded-lg hover:bg-green-50 transition-colors group">
              <div class="w-10 h-10 rounded-xl bg-green-100 flex items-center justify-center text-green-600 group-hover:bg-green-200 transition-colors shrink-0">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" /><line x1="16" y1="2" x2="16" y2="6" /><line x1="8" y1="2" x2="8" y2="6" /><line x1="3" y1="10" x2="21" y2="10" /></svg>
              </div>
              <div class="flex-1">
                <p class="text-sm font-medium text-gray-900">My Reservations</p>
                <p class="text-xs text-gray-500">View and manage your upcoming reservations</p>
              </div>
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 text-gray-300 group-hover:text-gray-500"><polyline points="9 18 15 12 9 6" /></svg>
            </A>
            <A href="/requests" class="flex items-center gap-3 p-3 rounded-lg hover:bg-purple-50 transition-colors group">
              <div class="w-10 h-10 rounded-xl bg-purple-100 flex items-center justify-center text-purple-600 group-hover:bg-purple-200 transition-colors shrink-0">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12" /></svg>
              </div>
              <div class="flex-1">
                <p class="text-sm font-medium text-gray-900">Cross-Team Requests</p>
                <p class="text-xs text-gray-500">Request seats in other team's areas</p>
              </div>
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 text-gray-300 group-hover:text-gray-500"><polyline points="9 18 15 12 9 6" /></svg>
            </A>
          </div>
        </div>

        <div class="bg-white rounded-xl border border-gray-200 shadow-sm">
          <div class="px-5 py-4 border-b border-gray-100">
            <h2 class="font-semibold text-gray-900">Your Profile</h2>
          </div>
          <div class="p-5 space-y-4">
            <div class="flex items-center justify-between">
              <span class="text-sm text-gray-500">Role</span>
              <span class="text-sm font-medium text-gray-900 capitalize">{currentUser()?.role ?? '-'}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-sm text-gray-500">Email</span>
              <span class="text-sm font-medium text-gray-900">{currentUser()?.email ?? '-'}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-sm text-gray-500">Team</span>
              <span class="text-sm font-medium text-gray-900">
                {teams.data?.find((t) => t.id === currentUser()?.team_id)?.name ?? '-'}
              </span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-sm text-gray-500">Weekly Limit</span>
              <span class="text-sm font-medium text-gray-900">{currentUser()?.weekly_limit ?? '-'} days</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
