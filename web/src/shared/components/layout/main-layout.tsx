import { ReactNode } from 'react'
import { Header } from './header'
import { TeamsSidebar } from './teams-sidebar'
import { SecondarySidebar } from './secondary-sidebar'
import { FileBrowser } from './file-browser'
import { ClaudeTeamDashboard } from '../claude-teams/claude-team-dashboard'
import { useClaudeTeamStore } from '../../stores'

interface MainLayoutProps {
  children: ReactNode
}

export function MainLayout({ children }: MainLayoutProps) {
  const { selectedTeam: selectedClaudeTeam } = useClaudeTeamStore()

  return (
    <div className="flex h-screen flex-col">
      <Header />
      <div className="flex flex-1 overflow-hidden">
        <TeamsSidebar />
        <SecondarySidebar />
        <main className="flex-1 overflow-auto p-4">
          {selectedClaudeTeam ? <ClaudeTeamDashboard /> : children}
        </main>
        <FileBrowser />
      </div>
    </div>
  )
}
