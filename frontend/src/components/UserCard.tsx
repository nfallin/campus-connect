
import type { User } from '../types'
import './UserCard.css'

type UserCardProps = {
  user: User
  onLogout: () => void
}

function UserCard({ user, onLogout }: UserCardProps) {
  return (
    <div className="user-card">
      <div className="user-info">
        <strong>{user.username}</strong>
        <span>{user.email}</span>
      </div>

      <button className="button button-secondary" onClick={onLogout}>
        Log out
      </button>
    </div>
  )
}

export default UserCard
