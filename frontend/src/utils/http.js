import router from '../router'
import { useAuthStore } from '../stores/auth'

const BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

function buildUrl(path, params) {
  const url = /^https?:\/\//.test(path) ? path : `${BASE_URL}${path.startsWith('/') ? '' : '/'}${path}`
  if (params && typeof params === 'object') {
    const usp = new URLSearchParams()
    Object.entries(params).forEach(([k, v]) => {
      if (v === undefined || v === null) return
      if (Array.isArray(v)) v.forEach((item) => usp.append(k, item))
      else usp.append(k, String(v))
    })
    const sep = url.includes('?') ? '&' : '?'
    return `${url}${usp.toString() ? sep + usp.toString() : ''}`
  }
  return url
}

function getAuthHeader() {
  const authStore = useAuthStore()
  const token = authStore.token
  return token ? { Authorization: `Bearer ${token}` } : {}
}

async function handleResponse(res) {
  if (res.status === 401) {
    const authStore = useAuthStore()
    authStore.clearAuth()
    router.push('/login')
    throw new Error('Unauthorized')
  }
  return res
}

async function request(path, { method = 'GET', headers = {}, body, params, auth = true } = {}) {
  const url = buildUrl(path, params)
  const isFormData = body instanceof FormData
  const finalHeaders = {
    Accept: 'application/json',
    ...(method !== 'GET' && !isFormData ? { 'Content-Type': 'application/json' } : {}),
    ...(auth ? getAuthHeader() : {}),
    ...headers,
  }
  const res = await fetch(url, { method, headers: finalHeaders, body: body !== undefined ? (isFormData ? body : JSON.stringify(body)) : undefined })
  await handleResponse(res)
  if (!res.ok) {
    const contentType = res.headers.get('content-type') || ''
    let message = `请求失败 (${res.status})`
    if (contentType.includes('application/json')) {
      try {
        const payload = await res.json()
        message = payload?.message || payload?.msg || message
      } catch {
        message = `请求失败 (${res.status})`
      }
    } else {
      try {
        const text = await res.text()
        if (text) message = text.substring(0, 200)
      } catch {
        message = `请求失败 (${res.status})`
      }
    }
    const err = new Error(message)
    err.status = res.status
    throw err
  }
  return res
}

export async function get(path, options = {}) {
  const res = await request(path, { ...options, method: 'GET' })
  const contentType = res.headers.get('content-type') || ''
  if (!contentType.includes('application/json')) {
    const text = await res.text()
    throw new Error(`Invalid response format: ${text.substring(0, 200)}`)
  }
  try {
    return await res.json()
  } catch (err) {
    const text = await res.text()
    throw new Error(`JSON parse error: ${err.message}, response: ${text.substring(0, 200)}`)
  }
}

export async function post(path, body, options = {}) {
  const res = await request(path, { ...options, method: 'POST', body })
  return res.json()
}

export async function put(path, body, options = {}) {
  const res = await request(path, { ...options, method: 'PUT', body })
  return res.json()
}

export async function del(path, body, options = {}) {
  const res = await request(path, { ...options, method: 'DELETE', body })
  const contentType = res.headers.get('content-type') || ''
  if (!contentType.includes('application/json')) {
    const text = await res.text()
    throw new Error(`Invalid response format: ${text.substring(0, 200)}`)
  }
  try {
    return await res.json()
  } catch (err) {
    const text = await res.text()
    throw new Error(`JSON parse error: ${err.message}, response: ${text.substring(0, 200)}`)
  }
}

const httpDelete = del
export default { get, post, put, delete: httpDelete }
