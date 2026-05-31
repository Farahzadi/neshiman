import { ParentComponent, onMount } from 'solid-js';
import { QueryClient, QueryClientProvider } from '@tanstack/solid-query';
import { setAuthHeader } from '@neshiman/api-client';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
    },
  },
});

const Providers: ParentComponent = (props) => {
  onMount(() => {
    const token = localStorage.getItem('neshiman_token');
    if (token) {
      setAuthHeader(() => token);
    }
  });

  return (
    <QueryClientProvider client={queryClient}>
      {props.children}
    </QueryClientProvider>
  );
};

export default Providers;
