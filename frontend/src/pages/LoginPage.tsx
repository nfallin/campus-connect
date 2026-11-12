
import { useState } from 'react'
import '../components/Form.css'

type LoginPageProps = {
  onLogin: (identifier: string, password: string) => Promise<void>
  onRegister: (username: string, email: string, password: string) => Promise<void>
}

function LoginPage({ onLogin, onRegister }: LoginPageProps) {
  const [mode, setMode] = useState<'login' | 'register'>('login')
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  async function handleSubmit(e: React.SyntheticEvent<HTMLFormElement>) {
    e.preventDefault()

    if (!password || password.length < 8 || (mode === 'login' ? !username.trim() : !username.trim() || !email.trim())) {
      setError(mode === 'login'
        ? 'Enter a username or email and a password of at least 8 characters.'
        : 'Enter a username, email, and a password of at least 8 characters.')
      return
    }

    setSubmitting(true)
    setError('')

    try {
      if (mode === 'login') {
        await onLogin(username.trim(), password)
      } else {
        await onRegister(username.trim(), email.trim(), password)
      }
    } catch (err) {
      console.error('Authentication failed', err)
      setError(err instanceof Error ? err.message : 'Authentication failed. Please try again.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <main className="card">
      <form onSubmit={handleSubmit}>
        <div className="login-intro">
          <h2>{mode === 'login' ? 'Log in' : 'Create an account'}</h2>
        </div>

        <div className="field">
          <label htmlFor="username">{mode === 'login' ? 'Username or email' : 'Username'}</label>
          <input
            id="username"
            type="text"
            placeholder={mode === 'login' ? 'Username or email' : 'Your username'}
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </div>

        {mode === 'register' && <div className="field">
          <label htmlFor="email">Email</label>
          <input
            id="email"
            type="email"
            placeholder="you@example.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </div>}

        <div className="field">
          <label htmlFor="password">Password</label>
          <input
            id="password"
            type="password"
            minLength={8}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>

        <button
          type="submit"
          className="create-btn"
          disabled={submitting}
        >
          {submitting ? 'Please wait...' : mode === 'login' ? 'Log in' : 'Register'}
        </button>

        {error && (
          <p className="form-message" role="alert">
            {error}
          </p>
        )}

        <button type="button" className="login-toggle" onClick={() => setMode(mode === 'login' ? 'register' : 'login')}>
          {mode === 'login' ? 'Create a new account' : 'I already have an account'}
        </button>
      </form>
    </main>
  )
}

export default LoginPage
