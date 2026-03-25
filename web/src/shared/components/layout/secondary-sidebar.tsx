import { useUIStore } from '../../stores'
import { cn } from '../../lib/utils'

const tabs = [
  { id: 'tasks' as const, label: 'Tasks' },
  { id: 'departments' as const, label: 'Departments' },
  { id: 'lead' as const, label: 'Lead' },
]

export function SecondarySidebar() {
  const { activeTab, setActiveTab } = useUIStore()

  return (
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
      <div className="flex-1 overflow-auto p-4">
        <p className="text-sm text-muted-foreground">
          {activeTab === 'tasks' && 'Task list will appear here'}
          {activeTab === 'departments' && 'Departments will appear here'}
          {activeTab === 'lead' && 'Lead info will appear here'}
        </p>
      </div>
    </div>
  )
}
