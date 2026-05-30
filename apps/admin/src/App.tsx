import type { Component } from 'solid-js';
import { lazy } from 'solid-js';
import { Router, Route, Navigate } from '@solidjs/router';
import type { RouteSectionProps } from '@solidjs/router';
import Providers from './providers';
import Layout from './components/Layout';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Rooms from './pages/Rooms';
import RoomEditor from './pages/RoomEditor';
import Teams from './pages/Teams';
import Users from './pages/Users';
import CrossTeamRequests from './pages/CrossTeamRequests';
import Reservations from './pages/Reservations';
import Toast from './components/Toast';

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
        <Route path="/rooms/:id/edit" component={() => <ProtectedRoute><RoomEditor /></ProtectedRoute>} />
        <Route path="/teams" component={() => <ProtectedRoute><Teams /></ProtectedRoute>} />
        <Route path="/users" component={() => <ProtectedRoute><Users /></ProtectedRoute>} />
        <Route path="/cross-team-requests" component={() => <ProtectedRoute><CrossTeamRequests /></ProtectedRoute>} />
        <Route path="/reservations" component={() => <ProtectedRoute><Reservations /></ProtectedRoute>} />
      </Router>
      <Toast />
      {import.meta.env.DEV && <Devtools />}
    </Providers>
  );
};

export default App;
