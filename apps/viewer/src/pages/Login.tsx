import { Component, createSignal, Show } from 'solid-js';
import { useNavigate } from '@solidjs/router';
import { useLogin, setAuthHeader, ApiError } from '@neshiman/api-client';

const Login: Component = () => {
  const navigate = useNavigate();
  const login = useLogin();
  const [username, setUsername] = createSignal('');
  const [password, setPassword] = createSignal('');
  const [error, setError] = createSignal('');

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError('');
    try {
      const result = await login.mutateAsync({ username: username(), password: password() });
      localStorage.setItem('neshiman_token', result.token!);
      localStorage.setItem('neshiman_user', JSON.stringify(result.user));
      setAuthHeader(() => result.token!);
      navigate('/', { replace: true });
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Login failed');
    }
  };

  return (
    <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-50 to-blue-50">
      <div class="bg-white rounded-2xl shadow-xl border border-gray-200 p-8 w-full max-w-sm mx-4">
        <div class="flex items-center gap-2.5 mb-8 justify-center">
          <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-600 to-indigo-700 flex items-center justify-center text-white text-lg font-bold shadow-sm">
            N
          </div>
          <span class="text-xl font-bold text-gray-900">Neshiman</span>
        </div>

        <form onSubmit={handleSubmit} class="space-y-4">
          <div>
            <label for="username" class="block text-sm font-medium text-gray-700 mb-1">Username</label>
            <input
              id="username"
              type="text"
              value={username()}
              onInput={(e) => setUsername(e.currentTarget.value)}
              required
              placeholder="alice"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
            />
          </div>
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-1">Password</label>
            <input
              id="password"
              type="password"
              value={password()}
              onInput={(e) => setPassword(e.currentTarget.value)}
              required
              placeholder="Enter your password"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
            />
          </div>

          <Show when={error()}>
            <div class="px-3 py-2 bg-red-50 border border-red-200 rounded-lg text-xs text-red-700">
              {error()}
            </div>
          </Show>

          <button
            type="submit"
            disabled={login.isPending}
            class="w-full py-2.5 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors disabled:opacity-50 shadow-sm"
          >
            {login.isPending ? 'Signing in...' : 'Sign in'}
          </button>
        </form>

        <p class="mt-6 text-xs text-gray-400 text-center">
          Default seed: use any username with password "password"
        </p>
      </div>
    </div>
  );
};

export default Login;
