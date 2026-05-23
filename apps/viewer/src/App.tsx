import type { Component, ParentProps } from 'solid-js';
import { lazy } from 'solid-js';
import { Router, Route, Navigate } from '@solidjs/router';
import type { RouteSectionProps } from '@solidjs/router';
import Providers from './providers';
import Layout from './components/layout/Layout';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Rooms from './pages/Rooms';
import RoomDetail from './pages/RoomDetail';
import Reservations from './pages/Reservations';
import CrossTeamRequests from './pages/CrossTeamRequests';

const Devtools = lazy(() => import('./devtools'));

function ProtectedRoute(props: { children: any }) {
  const token = localStorage.getItem('neshiman_token');
  if (!token) {
    return <Navigate href="/login" />;
  }
  return props.children;
}

const RootLayout: Component<RouteSectionProps> = (props) => {
  if (props.location.pathname === '/login') {
    return <>{props.children}</>;
  }
  return <Layout {...props} />;
};

const App: Component = () => {
  return (
    <Providers>
      <Router root={RootLayout}>
        <Route path="/login" component={Login} />
        <Route path="/" component={() => <ProtectedRoute><Dashboard /></ProtectedRoute>} />
        <Route path="/rooms" component={() => <ProtectedRoute><Rooms /></ProtectedRoute>} />
        <Route path="/rooms/:id" component={() => <ProtectedRoute><RoomDetail /></ProtectedRoute>} />
        <Route path="/reservations" component={() => <ProtectedRoute><Reservations /></ProtectedRoute>} />
        <Route path="/requests" component={() => <ProtectedRoute><CrossTeamRequests /></ProtectedRoute>} />
      </Router>
      {import.meta.env.DEV && <Devtools />}
    </Providers>
  );
};

export default App;
