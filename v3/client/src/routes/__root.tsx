import { createRootRoute, Link, Outlet } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/router-devtools'
import {
  useQuery,
  useQueryClient,
  QueryClient,
  QueryClientProvider,
} from '@tanstack/react-query'
import MainBody, { NavigationBar } from '../components/main_body_wrapper'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'

const queryClient = new QueryClient()

export const Route = createRootRoute({
  component: () => (
  <QueryClientProvider client={queryClient} >
      <MainBody>
      <NavigationBar/>
      <hr />
      <Outlet />
      <TanStackRouterDevtools />
    </MainBody>
    <ReactQueryDevtools initialIsOpen />
  </QueryClientProvider>
  ),
})
