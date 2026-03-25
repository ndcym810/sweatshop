import type { Team, Project, Department, Teammate, Task } from '../types/team'

export const mockTeams: Team[] = [
  {
    id: 'team-1',
    name: 'Team Alpha',
    description: 'Main development team',
    leadRuntimeType: 'claude-code',
    leadRuntimeModel: 'claude-opus-4-6',
    createdAt: '2026-03-20T10:00:00Z',
    updatedAt: '2026-03-20T10:00:00Z',
  },
  {
    id: 'team-2',
    name: 'Team Beta',
    description: 'Marketing and research',
    leadRuntimeType: 'claude-code',
    leadRuntimeModel: 'claude-sonnet-4-6',
    createdAt: '2026-03-21T10:00:00Z',
    updatedAt: '2026-03-21T10:00:00Z',
  },
]

export const mockProjects: Project[] = [
  {
    id: 'proj-1',
    teamId: 'team-1',
    name: 'E-commerce Platform',
    path: '/projects/ecommerce',
    defaultBranch: 'main',
    isActive: true,
  },
]

export const mockDepartments: Department[] = [
  { id: 'dept-1', teamId: 'team-1', name: 'Development', description: 'Software development', sortOrder: 0 },
  { id: 'dept-2', teamId: 'team-1', name: 'Marketing', description: 'Marketing and content', sortOrder: 1 },
  { id: 'dept-3', teamId: 'team-1', name: 'Deployment', description: 'DevOps and CI/CD', sortOrder: 2 },
]

export const mockTeammates: Teammate[] = [
  { id: 'tm-1', teamId: 'team-1', departmentId: 'dept-1', templateId: 'tpl-1', name: 'Frontend Dev 1', status: 'idle' },
  { id: 'tm-2', teamId: 'team-1', departmentId: 'dept-1', templateId: 'tpl-2', name: 'Backend Dev 1', status: 'running' },
  { id: 'tm-3', teamId: 'team-1', departmentId: 'dept-3', templateId: 'tpl-3', name: 'DevOps 1', status: 'idle' },
]

export const mockTasks: Task[] = [
  { id: 'task-1', teamId: 'team-1', projectId: 'proj-1', assignedTo: 'tm-2', title: 'Implement login API', description: 'Create authentication endpoints', status: 'in_progress', priority: 'high' },
  { id: 'task-2', teamId: 'team-1', projectId: 'proj-1', assignedTo: 'tm-1', title: 'Design checkout page', description: 'Create checkout UI', status: 'pending', priority: 'medium' },
  { id: 'task-3', teamId: 'team-1', projectId: 'proj-1', assignedTo: 'tm-3', title: 'Setup CI/CD', description: 'Configure pipeline', status: 'completed', priority: 'high' },
]
