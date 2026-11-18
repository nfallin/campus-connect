import UserCard from './UserCard'
import type { User, Page } from '../types'
import './AppHeader.css'

type AppHeaderProps = {
  page: Page
  currentUser: User | null
  onLogout: () => void
}

function AppHeader({
  page,
  currentUser,
  onLogout,
}: AppHeaderProps) {
  return (
    <header className="navbar">
      <div className="navbar-row">
        <div className="navbar-brand">
          <img src="/favicon.png" alt="" className="navbar-logo" />
          <span className="navbar-title">Campus Connect</span>
        </div>

        <nav className="navbar-links" aria-label="Main navigation">
          {!currentUser && (
            <a className="button button-secondary" href="#login" aria-current={page === 'login' ? 'page' : undefined}>
              Login
            </a>
          )}
          {currentUser && (
            <a className="button button-secondary" href="#create" aria-current={page === 'create' ? 'page' : undefined}>
              Create event
            </a>
          )}
          <a className="button button-secondary" href="#events" aria-current={page === 'events' ? 'page' : undefined}>
            View events
          </a>
          {currentUser && (
            <a className="button button-secondary" href="#rsvps" aria-current={page === 'rsvps' ? 'page' : undefined}>
              My events
            </a>
          )}
        </nav>

        <div className="navbar-user">
          {currentUser && (
            <UserCard user={currentUser} onLogout={onLogout} />
          )}
        </div>
      </div>
    </header>
  )
}

export default AppHeader