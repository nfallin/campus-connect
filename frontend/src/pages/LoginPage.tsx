
import { useState } from 'react'
import '../components/Form.css'

type LoginPageProps = {
  onLogin: (name: string, email: string) => Promise<void>
}

function LoginPage({ onLogin }: LoginPageProps) {
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  async function handleSubmit(e: React.SyntheticEvent<HTMLFormElement>) {
    e.preventDefault()

    if (!name.trim() || !email.trim()) {
      setError('Please enter your name and email.')
      return
    }

    setSubmitting(true)
    setError('')

    try {
      await onLogin(name.trim(), email.trim())
    } catch (err) {
      console.error('Login failed', err)
      setError('Unable to log in. Please check your information and try again.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <main className="card">
      <form onSubmit={handleSubmit}>
        <div className="field">
          <label htmlFor="name">Name</label>
          <input
            id="name"
            type="text"
            placeholder="Your name"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </div>

        <div className="field">
          <label htmlFor="email">Email</label>
          <input
            id="email"
            type="email"
            placeholder="you@example.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </div>

        <button
          type="submit"
          className="create-btn"
          disabled={submitting}
        >
          {submitting ? 'Logging in...' : 'Continue'}
        </button>

        {error && (
          <p className="form-message" role="alert">
            {error}
          </p>
        )}
      </form>
    </main>
  )
}

export default LoginPage
