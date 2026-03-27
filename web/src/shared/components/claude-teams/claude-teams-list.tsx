import { useEffect } from 'react'
import { Users, Loader2, AlertCircle } from 'lucide-react'
import { useClaudeTeamStore } from '../../stores/claude-team-store'
import { cn } from '../../lib/utils'

export function ClaudeTeamsList() {
  const { teams, selectedTeam, isLoading, error, fetchTeams, selectTeam } = useClaudeTeamStore()

  useEffect(() => {
    fetchTeams()
  }, [fetchTeams])

  if (isLoading && teams.length === 0) {
    return (
      <div className="flex items-center justify-center py-8">
        <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex items-center gap-2 p-4 text-sm text-destructive">
        <AlertCircle className="h-4 w-4" />
        <span>{error}</span>
      </div>
    )
  }

  if (teams.length === 0) {
    return (
      <div className="p-4 text-center text-sm text-muted-foreground">
        <Users className="mx-auto mb-2 h-8 w-8 opacity-50" />
        <p>No Claude Code teams found</p>
        <p className="mt-1 text-xs">
          Spawn a team using <code className="rounded bg-muted px-1">claude</code> CLI
        </p>
      </div>
    )
  }

  return (
    <div className="space-y-1 p-2">
      <h3 className="mb-2 px-2 text-xs font-medium uppercase text-muted-foreground">
        Claude Teams
      </h3>
      {teams.map((teamName) => (
        <button
          key={teamName}
          onClick={() => selectTeam(teamName)}
          className={cn(
            'flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors',
            selectedTeam?.name === teamName
              ? 'bg-accent text-accent-foreground'
              : 'hover:bg-accent'
          )}
        >
          <Users className="h-4 w-4" />
          <span className="truncate">{teamName}</span>
        </button>
      ))}
    </div>
  )
}
