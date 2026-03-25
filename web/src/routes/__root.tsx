/// <reference types="vite/client" />
import { useEffect } from 'react'
import type { ReactNode } from 'react'
import {
  Outlet,
  createRootRoute,
  HeadContent,
  Scripts,
} from '@tanstack/react-router'
import { MainLayout } from '../shared/components/layout/main-layout'
import { useTeamStore } from '../shared/stores'
import '../index.css'

export const Route = createRootRoute({
  head: () => ({
    meta: [
      {
        charSet: 'utf-8',
      },
      {
        name: 'viewport',
        content: 'width=device-width, initial-scale=1',
      },
      {
        title: 'Sweatshop',
      },
    ],
  }),
  component: RootComponent,
})

function RootComponent() {
  const fetchTeams = useTeamStore((state) => state.fetchTeams)

  useEffect(() => {
    fetchTeams()
  }, [fetchTeams])

  return (
    <RootDocument>
      <MainLayout>
        <Outlet />
      </MainLayout>
    </RootDocument>
  )
}

function RootDocument({ children }: Readonly<{ children: ReactNode }>) {
  return (
    <html lang="en" className="dark">
      <head>
        <HeadContent />
      </head>
      <body>
        {children}
        <Scripts />
      </body>
    </html>
  )
}
