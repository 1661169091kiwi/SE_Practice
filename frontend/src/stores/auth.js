import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const role = ref(localStorage.getItem('role') || '')
  const roles = ref(JSON.parse(localStorage.getItem('roles') || '[]'))
  const studentId = ref(localStorage.getItem('studentId') || '')
  const avatar = ref(localStorage.getItem('avatar') || '')

  function setToken(newToken) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function setRole(newRole) {
    role.value = newRole
    localStorage.setItem('role', newRole)
  }

  function setRoles(newRoles) {
    roles.value = newRoles
    localStorage.setItem('roles', JSON.stringify(newRoles))
  }

  function hasRole(checkRole) {
    return roles.value.includes(checkRole)
  }

  function setStudentId(id) {
    studentId.value = id
    localStorage.setItem('studentId', id)
  }

  function setAvatar(url) {
    avatar.value = url
    localStorage.setItem('avatar', url)
  }

  function clearAuth() {
    token.value = ''
    role.value = ''
    roles.value = []
    studentId.value = ''
    avatar.value = ''
    localStorage.removeItem('token')
    localStorage.removeItem('role')
    localStorage.removeItem('roles')
    localStorage.removeItem('studentId')
    localStorage.removeItem('avatar')
    // 清除Kiwi助手会话记录
    sessionStorage.removeItem('kiwi_chat_messages')
    sessionStorage.removeItem('kiwi_selected_model')
  }

  return {
    token,
    role,
    roles,
    studentId,
    avatar,
    setToken,
    setRole,
    setRoles,
    hasRole,
    setStudentId,
    setAvatar,
    clearAuth
  }
})
