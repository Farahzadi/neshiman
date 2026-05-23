import { Component, lazy } from 'solid-js';
import { Router, Route } from '@solidjs/router';
import Providers from './providers';
import Layout from './components/layout/Layout';
import Dashboard from './pages/Dashboard';
import Rooms from './pages/Rooms';
import RoomDetail from './pages/RoomDetail';
import Reservations from './pages/Reservations';
import CrossTeamRequests from './pages/CrossTeamRequests';

const Devtools = lazy(() => import('./devtools'));

const App: Component = () => {
  return (
    <Providers>
      <Router root={Layout}>
        <Route path="/" component={Dashboard} />
        <Route path="/rooms" component={Rooms} />
        <Route path="/rooms/:id" component={RoomDetail} />
        <Route path="/reservations" component={Reservations} />
        <Route path="/requests" component={CrossTeamRequests} />
      </Router>
      {import.meta.env.DEV && <Devtools />}
    </Providers>
  );
};

export default App;
