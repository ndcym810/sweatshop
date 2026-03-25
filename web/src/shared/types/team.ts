// Re-export types from api.ts for backwards compatibility
export type {
  Team,
  Project,
  Task,
  Department,
  CreateTeamInput,
  UpdateTeamInput,
  CreateProjectInput,
  UpdateProjectInput,
  CreateTaskInput,
  UpdateTaskInput,
  CreateDepartmentInput,
  UpdateDepartmentInput,
} from '../lib/api'

// Teammate type - not yet implemented in backend
export interface Teammate {
  id: string
  teamId: string
  departmentId: string
  templateId: string
  name: string
  status: 'idle' | 'running' | 'stopped' | 'error'
}
