import { Component } from 'solid-js';
import { Router, Routes, Route } from '@solidjs/router';
import Dashboard from './pages/Dashboard';
import Rooms from './pages/Rooms';
import Seats from './pages/Seats';
import Teams from './pages/Teams';

const App: Component = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" component={Dashboard} />
        <Route path="/rooms" component={Rooms} />
        <Route path="/rooms/:id/seats" component={Seats} />
        <Route path="/teams" component={Teams} />
      </Routes>
    </Router>
  );
};

export default App;
