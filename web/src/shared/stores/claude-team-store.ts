import { create } from 'zustand'
import {
  claudeTeamApi,
  type ClaudeTeam,
  type ClaudeMember,
  type InboxMessage,
  type SendMessageInput,
} from '../lib/api'

interface ClaudeTeamState {
  teams: string[]
  selectedTeam: ClaudeTeam | null
  selectedMember: ClaudeMember | null
  inbox: InboxMessage[]
  isLoading: boolean
  error: string | null
  wsConnected: boolean

  // Actions
  fetchTeams: () => Promise<void>
  selectTeam: (name: string) => Promise<void>
  clearSelection: () => void
  selectMember: (member: ClaudeMember | null) => void
  fetchInbox: (teamName: string, agentName: string) => Promise<void>
  refreshCurrentInbox: () => Promise<void>
  sendMessage: (input: SendMessageInput) => Promise<void>
  markRead: (timestamp: string) => Promise<void>
  setWsConnected: (connected: boolean) => void
  clearError: () => void
}

export const useClaudeTeamStore = create<ClaudeTeamState>((set, get) => ({
  teams: [],
  selectedTeam: null,
  selectedMember: null,
  inbox: [],
  isLoading: false,
  error: null,
  wsConnected: false,

  fetchTeams: async () => {
    set({ isLoading: true, error: null })
    try {
      const teams = await claudeTeamApi.list()
      set({ teams, isLoading: false })
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false })
    }
  },

  selectTeam: async (name: string) => {
    set({ isLoading: true, error: null })
    try {
      const team = await claudeTeamApi.get(name)
      set({ selectedTeam: team, selectedMember: null, inbox: [], isLoading: false })
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false })
    }
  },

  clearSelection: () => {
    set({ selectedTeam: null, selectedMember: null, inbox: [] })
  },

  selectMember: (member) => {
    set({ selectedMember: member })
    if (member && get().selectedTeam) {
      get().fetchInbox(get().selectedTeam!.name, member.name)
    }
  },

  fetchInbox: async (teamName: string, agentName: string) => {
    try {
      const inbox = await claudeTeamApi.getInbox(teamName, agentName)
      set({ inbox })
    } catch (error) {
      console.error('Failed to fetch inbox:', error)
    }
  },

  refreshCurrentInbox: async () => {
    const { selectedTeam, selectedMember } = get()
    if (selectedTeam && selectedMember) {
      await get().fetchInbox(selectedTeam.name, selectedMember.name)
    }
  },

  sendMessage: async (input: SendMessageInput) => {
    const { selectedTeam } = get()
    if (!selectedTeam) return

    set({ isLoading: true, error: null })
    try {
      await claudeTeamApi.sendMessage(selectedTeam.name, input)
      set({ isLoading: false })
      // Refresh inbox if we're viewing the recipient
      if (get().selectedMember?.name === input.to) {
        get().fetchInbox(selectedTeam.name, input.to)
      }
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false })
      throw error
    }
  },

  markRead: async (timestamp: string) => {
    const { selectedTeam, selectedMember } = get()
    if (!selectedTeam || !selectedMember) return

    try {
      await claudeTeamApi.markRead(selectedTeam.name, selectedMember.name, timestamp)
      // Update local state
      set((state) => ({
        inbox: state.inbox.map((msg) =>
          msg.timestamp === timestamp ? { ...msg, read: true } : msg
        ),
      }))
    } catch (error) {
      console.error('Failed to mark message as read:', error)
    }
  },

  setWsConnected: (connected) => set({ wsConnected: connected }),
  clearError: () => set({ error: null }),
}))
