import { create } from 'zustand'

type Tab = 'tasks' | 'departments' | 'lead'

interface UIState {
  sidebarCollapsed: boolean
  activeTab: Tab
  theme: 'light' | 'dark' | 'system'

  toggleSidebar: () => void
  setActiveTab: (tab: Tab) => void
  setTheme: (theme: 'light' | 'dark' | 'system') => void
}

export const useUIStore = create<UIState>((set) => ({
  sidebarCollapsed: false,
  activeTab: 'tasks',
  theme: 'dark',

  toggleSidebar: () =>
    set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed })),

  setActiveTab: (tab) => set({ activeTab: tab }),

  setTheme: (theme) => set({ theme }),
}))
