import { useState } from 'react'
import { Plus, ChevronDown, ChevronRight, Users, Loader2 } from 'lucide-react'
import { Button } from '../ui/button'
import { CreateDepartmentModal } from '../ui/create-department-modal'
import { useUIStore, useTeamStore } from '../../stores'
import { cn } from '../../lib/utils'
import type { Department } from '../../lib/api'

const tabs = [
  { id: 'tasks' as const, label: 'Tasks' },
  { id: 'departments' as const, label: 'Departments' },
  { id: 'lead' as const, label: 'Lead' },
]

export function SecondarySidebar() {
  const { activeTab, setActiveTab } = useUIStore()
  const { departments, selectedTeamId, isLoading, addDepartment, error, clearError } = useTeamStore()
  const [expandedDepts, setExpandedDepts] = useState<Set<string>>(new Set())
  const [isDeptModalOpen, setIsDeptModalOpen] = useState(false)

  const toggleDept = (deptId: string) => {
    const newExpanded = new Set(expandedDepts)
    if (newExpanded.has(deptId)) {
      newExpanded.delete(deptId)
    } else {
      newExpanded.add(deptId)
    }
    setExpandedDepts(newExpanded)
  }

  const handleCreateDepartment = async (data: { name: string; description: string }) => {
    await addDepartment(data)
    // Refresh departments list
    if (selectedTeamId) {
      useTeamStore.getState().fetchDepartments(selectedTeamId)
    }
  }

  return (
    <>
      <div className="flex w-64 flex-col border-r">
        <div className="flex border-b">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={cn(
                'flex-1 px-3 py-2 text-sm font-medium transition-colors',
                activeTab === tab.id
                  ? 'border-b-2 border-primary text-primary'
                  : 'text-muted-foreground hover:text-foreground'
              )}
            >
              {tab.label}
            </button>
          ))}
        </div>
        <div className="flex-1 overflow-auto">
          {activeTab === 'tasks' && (
            <div className="p-4">
              <p className="text-sm text-muted-foreground">Task list will appear here</p>
            </div>
          )}

          {activeTab === 'departments' && (
            <div className="p-2">
              <div className="flex items-center justify-between mb-2">
                <span className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
                  Departments
                </span>
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-6 w-6"
                  onClick={() => setIsDeptModalOpen(true)}
                  disabled={!selectedTeamId}
                >
                  <Plus className="h-3 w-3" />
                </Button>
              </div>

              {!selectedTeamId ? (
                <p className="px-2 py-4 text-center text-sm text-muted-foreground">
                  Select a team first
                </p>
              ) : isLoading ? (
                <div className="flex items-center justify-center py-4">
                  <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
                </div>
              ) : departments.length === 0 ? (
                <p className="px-2 py-4 text-center text-sm text-muted-foreground">
                  No departments yet
                </p>
              ) : (
                <div className="space-y-1">
                  {departments.map((dept) => (
                    <DepartmentItem
                      key={dept.id}
                      department={dept}
                      isExpanded={expandedDepts.has(dept.id)}
                      onToggle={() => toggleDept(dept.id)}
                    />
                  ))}
                </div>
              )}
            </div>
          )}

          {activeTab === 'lead' && (
            <div className="p-4">
              <p className="text-sm text-muted-foreground">Lead info will appear here</p>
            </div>
          )}
        </div>
      </div>

      <CreateDepartmentModal
        isOpen={isDeptModalOpen}
        onClose={() => {
          setIsDeptModalOpen(false)
          clearError()
        }}
        onSubmit={handleCreateDepartment}
        error={error}
      />
    </>
  )
}

interface DepartmentItemProps {
  department: Department
  isExpanded: boolean
  onToggle: () => void
}

function DepartmentItem({ department, isExpanded, onToggle }: DepartmentItemProps) {
  // Mock teammates for now - will be replaced with real data
  const teammates: Array<{ id: string; name: string; status: string }> = []

  return (
    <div className="rounded-md border">
      <button
        onClick={onToggle}
        className="flex w-full items-center justify-between px-2 py-1.5 text-sm hover:bg-accent rounded-md"
      >
        <div className="flex items-center gap-2">
          {isExpanded ? (
            <ChevronDown className="h-3 w-3 text-muted-foreground" />
          ) : (
            <ChevronRight className="h-3 w-3 text-muted-foreground" />
          )}
          <span>{department.name}</span>
        </div>
        <Button variant="ghost" size="icon" className="h-5 w-5" onClick={(e) => {
          e.stopPropagation()
          // TODO: Open add teammate modal
        }}>
          <Plus className="h-3 w-3" />
        </Button>
      </button>

      {isExpanded && (
        <div className="border-t px-2 py-1">
          {teammates.length === 0 ? (
            <p className="py-2 text-xs text-muted-foreground text-center">
              No teammates yet
            </p>
          ) : (
            <div className="space-y-1">
              {teammates.map((tm) => (
                <button
                  key={tm.id}
                  className="flex w-full items-center gap-2 rounded px-2 py-1 text-xs hover:bg-accent"
                >
                  <Users className="h-3 w-3 text-muted-foreground" />
                  <span>{tm.name}</span>
                  <span className={cn(
                    'ml-auto rounded-full px-1.5 py-0.5 text-[10px]',
                    tm.status === 'running' ? 'bg-green-500/20 text-green-500' :
                    tm.status === 'idle' ? 'bg-yellow-500/20 text-yellow-500' :
                    'bg-gray-500/20 text-gray-500'
                  )}>
                    {tm.status}
                  </span>
                </button>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  )
}
