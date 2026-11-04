import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import App from './App'
import { api } from './lib/api'

vi.mock('./lib/api', () => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
    delete: vi.fn(),
  },
}))

describe('RSVP functionality', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    window.location.hash = '#login'
  })

  it('allows a logged-in user to RSVP to an event', async () => {
    const user = userEvent.setup()

    // Mock the events returned by the backend
    vi.mocked(api.get).mockImplementation(async (url: string) => {
      if (url === '/events') {
        return {
          events: [
            {
              id: 1,
              name: 'Coding Workshop',
              dateTime: '2027-01-15T18:00:00Z',
              location: 'Room 101',
              clubName: 'Coding Club',
              description: 'Learn coding!',
              contactEmail: 'coding@example.com',
              websiteURL: '',
              hostUserId: 1,
            },
          ],
        } as never
      }

      if (url === '/users/1/rsvps') {
        return { events: [] } as never
      }

      throw new Error(`Unexpected GET request: ${url}`)
    })

    // Mock login and RSVP requests
    vi.mocked(api.post).mockImplementation(async (url: string) => {
      if (url === '/users') {
        return {
          user: {
            id: 1,
            name: 'Test User',
            email: 'test@example.com',
          },
          isNewUser: true,
        } as never
      }

      if (url === '/users/1/rsvps/1') {
        return {} as never
      }

      throw new Error(`Unexpected POST request: ${url}`)
    })

    render(<App />)

    // Log in
    await user.type(
      screen.getByRole('textbox', { name: /name/i }),
      'Test User'
    )

    await user.type(
      screen.getByRole('textbox', { name: /email/i }),
      'test@example.com'
    )

    await user.click(
      screen.getByRole('button', { name: /login/i })
    )

    // Wait for events page
    expect(
      await screen.findByText('Coding Workshop')
    ).toBeInTheDocument()

    // Click RSVP
    const rsvpButton = screen.getByRole('button', {
      name: /rsvp/i,
    })

    await user.click(rsvpButton)

    // Check that the correct API request was made
    await waitFor(() => {
      expect(api.post).toHaveBeenCalledWith(
        '/users/1/rsvps/1'
      )
    })

    // Check that the UI reflects the RSVP
    expect(
      await screen.findByRole('button', {
        name: /remove rsvp|cancel rsvp|rsvped/i,
      })
    ).toBeInTheDocument()
  })
})
