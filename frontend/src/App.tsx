import TodoList from './components/TodoList'

export default function App() {
  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center">
      <div className="bg-white rounded-xl shadow p-8 w-full max-w-md">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">Todos</h1>
        <TodoList />
      </div>
    </div>
  )
}
