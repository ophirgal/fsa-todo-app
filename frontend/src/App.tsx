import { useEffect, useState } from 'react'
import { api } from './api/client'

interface Todo {
  id: number
  title: string
  done: boolean
  created_at: string
}

function App() {
  const [todos, setTodos] = useState<Todo[]>([])
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    api.get<Todo[]>('/todos')
      .then(setTodos)
      .catch((e: Error) => setError(e.message))
  }, [])

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center">
      <div className="bg-white rounded-xl shadow p-8 w-full max-w-md">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">Todos</h1>
        {error && <p className="text-red-500 mb-4">{error}</p>}
        <ul className="space-y-3">
          {todos.map((todo) => (
            <li key={todo.id} className="flex items-center gap-3">
              <span
                className={todo.done ? 'line-through text-gray-400' : 'text-gray-800'}
              >
                {todo.title}
              </span>
            </li>
          ))}
        </ul>
      </div>
    </div>
  )
}

export default App
