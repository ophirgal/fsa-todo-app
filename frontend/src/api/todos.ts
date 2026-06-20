import { api } from './client'

export interface Todo {
  id: number
  title: string
  done: boolean
  created_at: string
}

export const fetchTodos = () => api.get<Todo[]>('/todos')

export const deleteTodo = (id: number) => api.delete<void>('/todos/' + id)
