import { Component, For, Show } from 'solid-js';
import { A } from '@solidjs/router';
import { useRooms } from '@neshiman/api-client';

const Rooms: Component = () => {
  const rooms = useRooms();

  return (
    <div>
      <div class="mb-8">
        <h1 class="text-2xl font-bold text-gray-900">Rooms</h1>
        <p class="text-sm text-gray-500 mt-1">Select a room to view its seat map and reserve a seat</p>
      </div>

      <Show when={rooms.data} fallback={
        <div class="flex items-center justify-center py-20">
          <div class="w-8 h-8 border-2 border-blue-600 border-t-transparent rounded-full animate-spin" />
        </div>
      }>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <For each={rooms.data}>
            {(room) => (
              <A
                href={`/rooms/${room.id}`}
                class="block bg-white rounded-xl border border-gray-200 shadow-sm hover:shadow-md hover:border-blue-200 transition-all group"
              >
                <div class="p-5">
                  <div class="flex items-center gap-3 mb-4">
                    <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white shadow-sm group-hover:shadow-md transition-shadow">
                      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5">
                        <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" /><polyline points="9 22 9 12 15 12 15 22" />
                      </svg>
                    </div>
                    <div class="flex-1 min-w-0">
                      <h2 class="font-semibold text-gray-900 truncate">{room.name}</h2>
                      <p class="text-xs text-gray-500">{room.grid_width} × {room.grid_height} grid</p>
                    </div>
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 text-gray-300 group-hover:text-gray-500 transition-colors"><polyline points="9 18 15 12 9 6" /></svg>
                  </div>
                  <div class="flex items-center gap-4 text-xs text-gray-400">
                    <span class="flex items-center gap-1">
                      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3.5 h-3.5"><rect x="3" y="3" width="18" height="18" rx="2" ry="2" /><line x1="3" y1="9" x2="21" y2="9" /><line x1="9" y1="21" x2="9" y2="9" /></svg>
                      {room.grid_width} × {room.grid_height}
                    </span>
                    <Show when={room.created_at}>
                      <span class="flex items-center gap-1">
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3.5 h-3.5"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" /><line x1="16" y1="2" x2="16" y2="6" /><line x1="8" y1="2" x2="8" y2="6" /><line x1="3" y1="10" x2="21" y2="10" /></svg>
                        {new Date(room.created_at!).toLocaleDateString()}
                      </span>
                    </Show>
                  </div>
                </div>
                <div class="px-5 py-3 border-t border-gray-100 bg-gray-50/50 rounded-b-xl">
                  <div class="flex items-center gap-2 text-xs text-blue-600 font-medium">
                    View Seat Map
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3.5 h-3.5"><line x1="5" y1="12" x2="19" y2="12" /><polyline points="12 5 19 12 12 19" /></svg>
                  </div>
                </div>
              </A>
            )}
          </For>
        </div>

        <Show when={rooms.data && rooms.data.length === 0}>
          <div class="text-center py-20">
            <div class="w-16 h-16 rounded-2xl bg-gray-100 flex items-center justify-center mx-auto mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-8 h-8 text-gray-300"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" /><polyline points="9 22 9 12 15 12 15 22" /></svg>
            </div>
            <p class="text-gray-500">No rooms available yet.</p>
          </div>
        </Show>
      </Show>
    </div>
  );
};

export default Rooms;
