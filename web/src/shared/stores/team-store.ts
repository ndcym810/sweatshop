import { create } from 'zustand'
import type { Team, Project, Department, Teammate, Task } from '../types/team'
import { mockTeams, mockProjects, mockDepartments, mockTeammates, mockTasks } from '../lib/mock-data'

interface TeamState {
  teams: Team[]
  selectedTeamId: string | null
  projects: Project[]
  departments: Department[]
  teammates: Teammate[]
  tasks: Task[]

  setSelectedTeam: (id: string | null) => void
  addTeam: (team: Team) => void
  updateTeam: (id: string, updates: Partial<Team>) => void
  deleteTeam: (id: string) => void
}

export const useTeamStore = create<TeamState>((set) => ({
  // Initialize with mock data
  teams: mockTeams,
  selectedTeamId: 'team-1',
  projects: mockProjects,
  departments: mockDepartments,
  teammates: mockTeammates,
  tasks: mockTasks,

  setSelectedTeam: (id) => set({ selectedTeamId: id }),

  addTeam: (team) => set((state) => ({ teams: [...state.teams, team] })),

  updateTeam: (id, updates) =>
    set((state) => ({
      teams: state.teams.map((t) => (t.id === id ? { ...t, ...updates } : t)),
    })),

  deleteTeam: (id) =>
    set((state) => ({
      teams: state.teams.filter((t) => t.id !== id),
      selectedTeamId: state.selectedTeamId === id ? null : state.selectedTeamId,
    })),
}))
