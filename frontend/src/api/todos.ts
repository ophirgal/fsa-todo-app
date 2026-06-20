import { api } from './client'

export interface Todo {
  id: number
  title: string
  done: boolean
  created_at: string
}

export const fetchTodos = () => api.get<Todo[]>('/todos')

export const deleteTodo = (id: number) => api.delete<void>('/todos/' + id)

export const updateTodoDone = (id: number, done: boolean) =>
  api.patch<Todo>('/todos/' + id + '/done', { done })

export const createTodo = (title: string) => api.post<Todo>('/todos', { title })
