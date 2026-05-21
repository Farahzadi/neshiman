import { Component } from 'solid-js';
import { Router, Route } from '@solidjs/router';
import Dashboard from './pages/Dashboard';
import Reservations from './pages/Reservations';

const App: Component = () => {
  return (
    <Router>
      <Route path="/" component={Dashboard} />
      <Route path="/reservations" component={Reservations} />
    </Router>
  );
};

export default App;
