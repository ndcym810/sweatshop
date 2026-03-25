import { create } from 'zustand'
import {
  teamApi,
  projectApi,
  taskApi,
  departmentApi,
  type Team,
  type Project,
  type Task,
  type Department,
  type CreateTeamInput
} from '../lib/api'

interface TeamState {
  teams: Team[]
  selectedTeamId: string | null
  projects: Project[]
  departments: Department[]
  tasks: Task[]
  isLoading: boolean
  error: string | null

  // Actions
  fetchTeams: () => Promise<void>
  setSelectedTeam: (id: string | null) => void
  addTeam: (input: CreateTeamInput) => Promise<void>
  clearError: () => void
  fetchProjects: (teamId: string) => Promise<void>
  fetchDepartments: (teamId: string) => Promise<void>
  fetchTasks: (teamId: string) => Promise<void>
}

export const useTeamStore = create<TeamState>((set, get) => ({
  teams: [],
  selectedTeamId: null,
  projects: [],
  departments: [],
  tasks: [],
  isLoading: false,
  error: null,

  fetchTeams: async () => {
    set({ isLoading: true, error: null })
    try {
      const teams = await teamApi.list()
      set({ teams, isLoading: false })
      // Auto-select first team if none selected
      if (teams.length > 0 && !get().selectedTeamId) {
        set({ selectedTeamId: teams[0].id })
        // Fetch related data
        get().fetchProjects(teams[0].id)
        get().fetchDepartments(teams[0].id)
        get().fetchTasks(teams[0].id)
      }
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false })
    }
  },

  setSelectedTeam: (id) => {
    set({ selectedTeamId: id })
    if (id) {
      get().fetchProjects(id)
      get().fetchDepartments(id)
      get().fetchTasks(id)
    }
  },

  addTeam: async (input) => {
    set({ isLoading: true, error: null })
    try {
      const team = await teamApi.create(input)
      set((state) => ({ teams: [...state.teams, team], isLoading: false }))
      // Auto-select the new team
      set({ selectedTeamId: team.id })
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false })
      throw error
    }
  },

  clearError: () => set({ error: null }),

  fetchProjects: async (teamId: string) => {
    try {
      const projects = await projectApi.list(teamId)
      set({ projects })
    } catch (error) {
      console.error('Failed to fetch projects:', error)
    }
  },

  fetchDepartments: async (teamId: string) => {
    try {
      const departments = await departmentApi.list(teamId)
      set({ departments })
    } catch (error) {
      console.error('Failed to fetch departments:', error)
    }
  },

  fetchTasks: async (teamId: string) => {
    try {
      const tasks = await taskApi.list(teamId)
      set({ tasks })
    } catch (error) {
      console.error('Failed to fetch tasks:', error)
    }
  },
}))
