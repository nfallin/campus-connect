import UserCard from './UserCard'
import type { User, Page } from '../types'
import './AppHeader.css'

type AppHeaderProps = {
  page: Page
  currentUser: User | null
  welcomeMessage: string
  onLogout: () => void
}

function AppHeader({
  page,
  currentUser,
  welcomeMessage,
  onLogout,
}: AppHeaderProps) {
  const subtitle =
    page === 'login' ? 'Sign in or create an account'
    : page === 'create' ? 'Create an event'
    : page === 'rsvps' ? 'See your events'
    : 'Browse events'

  return (
    <header className="page-header">
      <div className="header-top">
        <div className="header-title">
          <h1>Campus Connect</h1>
          <p className="subtitle">{subtitle}</p>
        </div>

        {currentUser && (
          <UserCard user={currentUser} onLogout={onLogout} />
        )}
      </div>

      <div className="welcome-message-slot">
        <div className="welcome-message">
          {currentUser ? welcomeMessage : ''}
        </div>
      </div>
    </header>
  )
}

export default AppHeader