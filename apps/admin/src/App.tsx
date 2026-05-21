import { Component, lazy } from 'solid-js';
import { Router, Route } from '@solidjs/router';
import Providers from './providers';
import Dashboard from './pages/Dashboard';
import Rooms from './pages/Rooms';
import Seats from './pages/Seats';
import Teams from './pages/Teams';

const Devtools = lazy(() => import('./devtools'));

const App: Component = () => {
  return (
    <Providers>
      <Router>
        <Route path="/" component={Dashboard} />
        <Route path="/rooms" component={Rooms} />
        <Route path="/rooms/:id/seats" component={Seats} />
        <Route path="/teams" component={Teams} />
      </Router>
      {import.meta.env.DEV && <Devtools />}
    </Providers>
  );
};

export default App;
