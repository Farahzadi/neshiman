import { Component } from 'solid-js';
import { Router, Routes, Route } from '@solidjs/router';
import Dashboard from './pages/Dashboard';
import Reservations from './pages/Reservations';

const App: Component = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" component={Dashboard} />
        <Route path="/reservations" component={Reservations} />
      </Routes>
    </Router>
  );
};

export default App;
