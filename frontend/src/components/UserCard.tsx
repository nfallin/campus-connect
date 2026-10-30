
import type { User } from '../types'

type UserCardProps = {
  user: User
  onLogout: () => void
}

function UserCard({ user, onLogout }: UserCardProps) {
  return (
    <div className="user-card">
      <div className="user-info">
        <strong>{user.name}</strong>
        <span>{user.email}</span>
      </div>

      <button onClick={onLogout}>
        Log out
      </button>
    </div>
  )
}

export default UserCard
