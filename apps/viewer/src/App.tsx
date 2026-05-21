import { Component, lazy } from 'solid-js';
import { Router, Route } from '@solidjs/router';
import Providers from './providers';
import Dashboard from './pages/Dashboard';
import Reservations from './pages/Reservations';

const Devtools = lazy(() => import('./devtools'));

const App: Component = () => {
  return (
    <Providers>
      <Router>
        <Route path="/" component={Dashboard} />
        <Route path="/reservations" component={Reservations} />
      </Router>
      {import.meta.env.DEV && <Devtools />}
    </Providers>
  );
};

export default App;
