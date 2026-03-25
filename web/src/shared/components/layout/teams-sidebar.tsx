import { Plus, ChevronLeft, ChevronRight, Loader2 } from 'lucide-react'
import { Button } from '../ui/button'
import { useTeamStore, useUIStore } from '../../stores'
import { cn } from '../../lib/utils'

export function TeamsSidebar() {
  const { teams, selectedTeamId, setSelectedTeam, isLoading } = useTeamStore()
  const { sidebarCollapsed, toggleSidebar } = useUIStore()

  return (
    <div
      className={cn(
        'flex flex-col border-r transition-all duration-300',
        sidebarCollapsed ? 'w-12' : 'w-56'
      )}
    >
      <div className="flex items-center justify-between p-2">
        {!sidebarCollapsed && (
          <span className="text-sm font-medium">Teams</span>
        )}
        <Button
          variant="ghost"
          size="icon"
          onClick={toggleSidebar}
          className="h-8 w-8"
        >
          {sidebarCollapsed ? (
            <ChevronRight className="h-4 w-4" />
          ) : (
            <ChevronLeft className="h-4 w-4" />
          )}
        </Button>
      </div>

      {!sidebarCollapsed && (
        <>
          <div className="p-2">
            <Button variant="outline" size="sm" className="w-full">
              <Plus className="mr-2 h-4 w-4" />
              New Team
            </Button>
          </div>

          <div className="flex-1 overflow-auto p-2">
            {isLoading ? (
              <div className="flex items-center justify-center py-4">
                <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
              </div>
            ) : (
              <>
                {teams.map((team) => (
                  <button
                    key={team.id}
                    onClick={() => setSelectedTeam(team.id)}
                    className={cn(
                      'flex w-full items-center rounded-md px-2 py-1.5 text-sm transition-colors cursor-pointer',
                      selectedTeamId === team.id
                        ? 'bg-accent text-accent-foreground'
                        : 'hover:bg-accent'
                    )}
                  >
                    <span className="truncate">{team.name}</span>
                  </button>
                ))}
                {teams.length === 0 && (
                  <p className="px-2 py-4 text-center text-sm text-muted-foreground">
                    No teams yet
                  </p>
                )}
              </>
            )}
          </div>
        </>
      )}
    </div>
  )
}
