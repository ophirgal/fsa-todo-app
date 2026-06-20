import { useState } from 'react'

export default function AddTodo({ onAdd }: { onAdd: (title: string) => void }) {
  const [title, setTitle] = useState('')

  const submit = () => {
    const trimmed = title.trim()
    if (!trimmed) return
    onAdd(trimmed)
    setTitle('')
  }

  return (
    <input
      type="text"
      value={title}
      onChange={e => setTitle(e.target.value)}
      onKeyDown={e => e.key === 'Enter' && submit()}
      placeholder="Type and hit Enter to add a todo..."
      className="w-full border border-gray-300 rounded-lg px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
    />
  )
}
