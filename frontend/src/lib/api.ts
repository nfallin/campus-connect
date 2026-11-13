const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? '/api'
const TOKEN_KEY = 'campus-connect-token'

function getStoredToken() {
  return typeof localStorage === 'undefined' ? null : localStorage.getItem(TOKEN_KEY)
}

export function setAuthToken(token: string) {
  if (typeof localStorage !== 'undefined') localStorage.setItem(TOKEN_KEY, token)
}

export function clearAuthToken() {
  if (typeof localStorage !== 'undefined') localStorage.removeItem(TOKEN_KEY)
}

export class ApiError extends Error {
  status: number

  constructor(status: number, message: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(getStoredToken()
        ? { Authorization: `Bearer ${getStoredToken()}` }
        : {}),
      ...options.headers,
    },
    ...options,
  })

  if (!res.ok) {
    const body = await res.text().catch(() => '')
    let message = body || res.statusText
    try {
      const parsed = JSON.parse(body) as { message?: string }
      if (parsed.message) message = parsed.message
    } catch {
      // Keep the response body when it is not a JSON error response.
    }
    throw new ApiError(res.status, message)
  }

  const text = await res.text()

  if (!text) return undefined as T

  return JSON.parse(text) as T
}


export const api = {
  get: <T>(path: string) => request<T>(path),
  post: <T>(path: string, data?: unknown) =>
    request<T>(path, { method: 'POST', body: data ? JSON.stringify(data) : undefined }),
  put: <T>(path: string, data?: unknown) =>
    request<T>(path, { method: 'PUT', body: data ? JSON.stringify(data) : undefined }),
  patch: <T>(path: string, data?: unknown) =>
    request<T>(path, { method: 'PATCH', body: data ? JSON.stringify(data) : undefined }),
  delete: <T>(path: string) => request<T>(path, { method: 'DELETE' }),
}
