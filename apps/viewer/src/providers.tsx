import { ParentComponent, onMount } from 'solid-js';
import { QueryClient, QueryClientProvider } from '@tanstack/solid-query';
import { setAuthHeader } from '@neshiman/api-client';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      retry: 1,
    },
  },
});

const Providers: ParentComponent = (props) => {
  onMount(() => {
    setAuthHeader(() => localStorage.getItem('viewer_user_id'));
  });

  return (
    <QueryClientProvider client={queryClient}>
      {props.children}
    </QueryClientProvider>
  );
};

export default Providers;
