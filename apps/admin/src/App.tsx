import { Component, lazy } from 'solid-js';
import { Router, Route } from '@solidjs/router';
import Providers from './providers';
import Layout from './components/Layout';
import Dashboard from './pages/Dashboard';
import Rooms from './pages/Rooms';
import Seats from './pages/Seats';
import Teams from './pages/Teams';
import Users from './pages/Users';
import CrossTeamRequests from './pages/CrossTeamRequests';

const Devtools = lazy(() => import('./devtools'));

const App: Component = () => {
  return (
    <Providers>
      <Router root={Layout}>
        <Route path="/" component={Dashboard} />
        <Route path="/rooms" component={Rooms} />
        <Route path="/rooms/:id/seats" component={Seats} />
        <Route path="/teams" component={Teams} />
        <Route path="/users" component={Users} />
        <Route path="/cross-team-requests" component={CrossTeamRequests} />
      </Router>
      {import.meta.env.DEV && <Devtools />}
    </Providers>
  );
};

export default App;
