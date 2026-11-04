// Unit tests written by Kaavya Jethwa

import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { api, ApiError } from './api'

describe('api client (Kaavya Jethwa)', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn())
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('resolves with the parsed JSON body on a successful GET', async () => {
    vi.mocked(fetch).mockResolvedValue(
      new Response(JSON.stringify({ id: 1 }), { status: 200 }),
    )

    const result = await api.get<{ id: number }>('/events/1')

    expect(result).toEqual({ id: 1 })
  })

  it('sends a JSON body and the correct method on POST', async () => {
    vi.mocked(fetch).mockResolvedValue(new Response(JSON.stringify({ ok: true }), { status: 200 }))

    await api.post('/events', { name: 'Hackathon' })

    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining('/events'),
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ name: 'Hackathon' }),
      }),
    )
  })

  it('sends a JSON body and the correct method on PATCH', async () => {
    vi.mocked(fetch).mockResolvedValue(new Response(JSON.stringify({ ok: true }), { status: 200 }))

    await api.patch('/events/1', { name: 'Updated' })

    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining('/events/1'),
      expect.objectContaining({
        method: 'PATCH',
        body: JSON.stringify({ name: 'Updated' }),
      }),
    )
  })

  it('resolves to undefined for an empty response body', async () => {
    vi.mocked(fetch).mockResolvedValue(new Response(null, { status: 204 }))

    const result = await api.delete('/events/1')

    expect(result).toBeUndefined()
  })

  it('throws an ApiError using the JSON message field when the response body is JSON', () => {
    vi.mocked(fetch).mockResolvedValue(
      new Response(JSON.stringify({ message: 'event not found' }), { status: 404 }),
    )

    return expect(api.get('/events/999')).rejects.toMatchObject(
      new ApiError(404, 'event not found'),
    )
  })

  it('throws an ApiError using the raw body text when the response is not JSON', () => {
    vi.mocked(fetch).mockResolvedValue(new Response('internal server error', { status: 500 }))

    return expect(api.get('/events/1')).rejects.toMatchObject(
      new ApiError(500, 'internal server error'),
    )
  })
})
