// API base URL - can be configured via environment at build time
const API_BASE = 'http://localhost:8000/api'

async function request<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE}${path}`
  const response = await fetch(url, {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  })

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Unknown error' }))
    throw new Error(error.error || `HTTP ${response.status}`)
  }

  if (response.status === 204) {
    return undefined as T
  }

  return response.json()
}

// Team API
export interface Team {
  id: string
  name: string
  description: string
  leadRuntimeType: string
  leadRuntimeModel: string
  createdAt: string
  updatedAt: string
}

export interface CreateTeamInput {
  name: string
  description?: string
  leadRuntimeType?: string
  leadRuntimeModel?: string
}

export interface UpdateTeamInput {
  name?: string
  description?: string
  leadRuntimeType?: string
  leadRuntimeModel?: string
}

export const teamApi = {
  list: () => request<Team[]>('/teams'),
  get: (id: string) => request<Team>(`/teams/${id}`),
  create: (data: CreateTeamInput) => request<Team>('/teams', {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  update: (id: string, data: UpdateTeamInput) => request<Team>(`/teams/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  }),
  delete: (id: string) => request<void>(`/teams/${id}`, { method: 'DELETE' }),
}

// Project API
export interface Project {
  id: string
  teamId: string
  name: string
  path: string
  defaultBranch: string
  isActive: boolean
  createdAt: string
}

export interface CreateProjectInput {
  name: string
  path: string
  defaultBranch?: string
}

export interface UpdateProjectInput {
  name?: string
  path?: string
  defaultBranch?: string
  isActive?: boolean
}

export const projectApi = {
  list: (teamId: string) => request<Project[]>(`/teams/${teamId}/projects`),
  create: (teamId: string, data: CreateProjectInput) => request<Project>(`/teams/${teamId}/projects`, {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  update: (id: string, data: UpdateProjectInput) => request<Project>(`/teams/${id}/projects/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  }),
  delete: (teamId: string, id: string) => request<void>(`/teams/${teamId}/projects/${id}`, { method: 'DELETE' }),
}

// Department API
export interface Department {
  id: string
  teamId: string
  name: string
  description: string
  sortOrder: number
  createdAt: string
}

export interface CreateDepartmentInput {
  name: string
  description?: string
  sortOrder?: number
}

export interface UpdateDepartmentInput {
  name?: string
  description?: string
  sortOrder?: number
}

export const departmentApi = {
  list: (teamId: string) => request<Department[]>(`/teams/${teamId}/departments`),
  create: (teamId: string, data: CreateDepartmentInput) => request<Department>(`/teams/${teamId}/departments`, {
    method: 'POST',
    body: JSON.stringify(data),
  }),
}

// Task API
export interface Task {
  id: string
  teamId: string
  projectId: string | null
  assignedTo: string | null
  title: string
  description: string
  status: 'pending' | 'in_progress' | 'completed' | 'blocked'
  priority: 'low' | 'medium' | 'high'
  createdAt: string
  completedAt: string | null
}

export interface CreateTaskInput {
  title: string
  description?: string
  projectId?: string
  assignedTo?: string
  priority?: 'low' | 'medium' | 'high'
}

export interface UpdateTaskInput {
  title?: string
  description?: string
  status?: 'pending' | 'in_progress' | 'completed' | 'blocked'
  priority?: 'low' | 'medium' | 'high'
}

export const taskApi = {
  list: (teamId: string, status?: string) => {
    const query = status ? `?status=${status}` : ''
    return request<Task[]>(`/teams/${teamId}/tasks${query}`)
  },
  create: (teamId: string, data: CreateTaskInput) => request<Task>(`/teams/${teamId}/tasks`, {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  update: (teamId: string, id: string, data: UpdateTaskInput) => request<Task>(`/teams/${teamId}/tasks/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  }),
}
