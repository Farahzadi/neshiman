import { ParentComponent } from 'solid-js';
import { QueryClient, QueryClientProvider } from '@tanstack/solid-query';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      retry: 1,
    },
  },
});

const Providers: ParentComponent = (props) => (
  <QueryClientProvider client={queryClient}>
    {props.children}
  </QueryClientProvider>
);

export default Providers;
