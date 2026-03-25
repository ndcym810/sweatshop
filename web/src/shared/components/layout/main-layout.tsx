import { ReactNode } from 'react'
import { Header } from './header'
import { TeamsSidebar } from './teams-sidebar'
import { SecondarySidebar } from './secondary-sidebar'
import { FileBrowser } from './file-browser'

interface MainLayoutProps {
  children: ReactNode
}

export function MainLayout({ children }: MainLayoutProps) {
  return (
    <div className="flex h-screen flex-col">
      <Header />
      <div className="flex flex-1 overflow-hidden">
        <TeamsSidebar />
        <SecondarySidebar />
        <main className="flex-1 overflow-auto p-4">{children}</main>
        <FileBrowser />
      </div>
    </div>
  )
}
