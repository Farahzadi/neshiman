import { Component } from 'solid-js';
import { A } from '@solidjs/router';

const Dashboard: Component = () => {
  return (
    <div class="p-6">
      <h1 class="text-2xl font-bold mb-6">Admin Dashboard</h1>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <A href="/rooms" class="p-4 border rounded-lg hover:shadow-md">
          <h2 class="font-semibold">Rooms</h2>
          <p class="text-sm text-gray-600">Manage rooms and grid sizes</p>
        </A>
        <A href="/teams" class="p-4 border rounded-lg hover:shadow-md">
          <h2 class="font-semibold">Teams</h2>
          <p class="text-sm text-gray-600">Manage teams and members</p>
        </A>
      </div>
    </div>
  );
};

export default Dashboard;
