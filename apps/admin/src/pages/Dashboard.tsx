import { Component, For } from 'solid-js';
import { useRooms, useTeams } from '@neshiman/api-client';

const Dashboard: Component = () => {
  const rooms = useRooms();
  const teams = useTeams();

  const stats = () => [
    { label: 'Rooms', value: rooms.data?.length ?? '-', color: 'bg-blue-500' },
    { label: 'Teams', value: teams.data?.length ?? '-', color: 'bg-green-500' },
  ];

  return (
    <div>
      <h1 class="text-2xl font-bold mb-6">Dashboard</h1>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <For each={stats()}>
          {(stat) => (
            <div class="bg-white rounded-lg shadow p-5 border-l-4" style={{ 'border-left-color': stat.color.replace('bg-', '') === 'blue-500' ? '#3b82f6' : '#22c55e' }}>
              <p class="text-sm text-gray-500">{stat.label}</p>
              <p class="text-3xl font-bold mt-1">{stat.value}</p>
            </div>
          )}
        </For>
      </div>
    </div>
  );
};

export default Dashboard;
