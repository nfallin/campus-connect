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
        <a className="navbar-brand" href="#events" aria-label="Campus Connect home">
          <img src="/favicon.png" alt="" className="navbar-logo" />
          <span className="navbar-title">Campus Connect</span>
        </a>

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

        {currentUser && (
          <div className="navbar-account">
            <button type="button" className="navbar-logout" onClick={onLogout}>
              Log out
            </button>
          </div>
        )}
      </div>
    </header>
  )
}

export default AppHeader