import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: HomePage,
})

function HomePage() {
  return (
    <div className="flex h-screen items-center justify-center">
      <h1 className="text-4xl font-bold">Sweatshop</h1>
    </div>
  )
}
