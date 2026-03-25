import { Folder, File } from 'lucide-react'

const mockFiles = [
  { name: 'src', type: 'folder' as const },
  { name: 'tests', type: 'folder' as const },
  { name: 'docs', type: 'folder' as const },
  { name: 'main.go', type: 'file' as const },
  { name: 'go.mod', type: 'file' as const },
]

export function FileBrowser() {
  return (
    <div className="flex w-56 flex-col border-l">
      <div className="border-b p-2">
        <span className="text-sm font-medium">Files</span>
      </div>
      <div className="flex-1 overflow-auto p-2">
        {mockFiles.map((file) => (
          <div
            key={file.name}
            className="flex items-center gap-2 rounded px-2 py-1 text-sm hover:bg-accent"
          >
            {file.type === 'folder' ? (
              <Folder className="h-4 w-4 text-muted-foreground" />
            ) : (
              <File className="h-4 w-4 text-muted-foreground" />
            )}
            <span className="truncate">{file.name}</span>
          </div>
        ))}
      </div>
    </div>
  )
}
