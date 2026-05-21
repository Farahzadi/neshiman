import { Component } from 'solid-js';
import { A } from '@solidjs/router';

const Dashboard: Component = () => {
  return (
    <div class="p-6">
      <h1 class="text-2xl font-bold mb-6">Seat Reservation</h1>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <A href="/reservations" class="p-4 border rounded-lg hover:shadow-md">
          <h2 class="font-semibold">My Reservations</h2>
          <p class="text-sm text-gray-600">View and manage your seat reservations</p>
        </A>
      </div>
    </div>
  );
};

export default Dashboard;
