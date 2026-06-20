import { useMutation, useQueryClient } from '@tanstack/react-query'
import TodoList from './components/TodoList'
import AddTodo from './components/AddTodo'
import { createTodo, Todo } from './api/todos'

export default function App() {
  const queryClient = useQueryClient()

  const { mutate: handleAdd } = useMutation({
    mutationFn: createTodo,
    onSuccess: (newTodo: Todo) => {
      queryClient.setQueryData<Todo[]>(['todos'], (prev: Todo[] | undefined) =>
        [...(prev ?? []), newTodo],
      )
    },
  })

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center">
      <div className="bg-white rounded-xl shadow p-8 w-full max-w-md">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">Todos</h1>
        <TodoList />
        <div className="mt-4">
          <AddTodo onAdd={handleAdd} />
        </div>
      </div>
    </div>
  )
}
