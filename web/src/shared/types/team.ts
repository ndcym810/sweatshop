export interface Team {
  id: string
  name: string
  description: string
  leadRuntimeType: string
  leadRuntimeModel: string
  createdAt: string
  updatedAt: string
}

export interface Project {
  id: string
  teamId: string
  name: string
  path: string
  defaultBranch: string
  isActive: boolean
}

export interface Department {
  id: string
  teamId: string
  name: string
  description: string
  sortOrder: number
}

export interface Teammate {
  id: string
  teamId: string
  departmentId: string
  templateId: string
  name: string
  status: 'idle' | 'running' | 'stopped' | 'error'
}

export interface Task {
  id: string
  teamId: string
  projectId: string | null
  assignedTo: string | null
  title: string
  description: string
  status: 'pending' | 'in_progress' | 'completed' | 'blocked'
  priority: 'low' | 'medium' | 'high'
}
