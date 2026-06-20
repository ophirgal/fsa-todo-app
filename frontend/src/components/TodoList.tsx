import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { fetchTodos, deleteTodo, Todo } from '../api/todos'
import TodoItem from './TodoItem'

export default function TodoList() {
  const queryClient = useQueryClient()

  const { data: todos = [], error } = useQuery({
    queryKey: ['todos'],
    queryFn: fetchTodos,
  })

  const { mutate: handleDelete } = useMutation({
    mutationFn: deleteTodo,
    onSuccess: (_: void, id: number) => {
      queryClient.setQueryData<Todo[]>(['todos'], (prev: Todo[] | undefined) =>
        (prev ?? []).filter((t: Todo) => t.id !== id),
      )
    },
  })

  if (error) {
    return <p className="text-red-500">{(error as Error).message}</p>
  }

  return (
    <ul className="space-y-3">
      {todos && todos.map((todo: Todo) => (
        <TodoItem key={todo.id} todo={todo} onDelete={handleDelete} />
      ))}
    </ul>
  )
}
