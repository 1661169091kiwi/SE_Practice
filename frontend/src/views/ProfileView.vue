<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { get, post } from '@/utils/http'
import { useAuthStore } from '@/stores/auth'
import TeamManagementModal from '../components/TeamManagementModal.vue'
import TeamChat from '../components/TeamChat.vue'

const router = useRouter()
const authStore = useAuthStore()

const fileInput = ref(null)

// 触发文件选择
const triggerFileInput = () => {
  fileInput.value.click()
}

// 处理文件选择
const handleFileChange = async (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  // 验证文件类型
  if (!file.type.startsWith('image/')) {
    showToast('请选择图片文件', 'error')
    return
  }
  
  // 验证文件大小 (例如限制为 10MB)
  if (file.size > 10 * 1024 * 1024) {
    showToast('图片大小不能超过 10MB', 'error')
    return
  }

  // 预览
  const reader = new FileReader()
  reader.onload = (e) => {
    userInfo.avatar = e.target.result
  }
  reader.readAsDataURL(file)

  // 上传文件到服务器
  const formData = new FormData()
  formData.append('avatar', file)

  try {
    const res = await post(`/user/avatar/${userInfo.studentId}`, formData)
    if (res.code === 200 && res.data) {
      userInfo.avatar = res.data.avatar_url ? `http://localhost:8080${res.data.avatar_url}` : ''
      showToast('头像更新成功')
    } else {
      showToast(res.message || res.msg || '头像更新失败', 'error')
    }
  } catch (err) {
    console.error('Upload avatar error:', err)
    showToast('头像上传失败，请重试', 'error')
  }
}

// 用户信息
const userInfo = reactive({
  name: '',
  studentId: '', 
  avatar: '', 
  totalCollected: 0,
  completedTasks: 0,
  pendingTasks: 0,
  joinDate: '',
  department: '',
  role: '',
  roles: [],
  applyStatus: 'none' // none, pending, approved
})

const myTeams = ref([])
const myAthleteInfo = ref([]) // 运动员信息列表
const teamCaptainStatus = ref({}) // 队伍ID -> 是否是队长
const myMatches = ref([]) // 运动员参加的比赛列表
const isLoadingMatches = ref(false)
const availableMatches = ref([]) // 可报名的比赛列表
const isLoadingAvailableMatches = ref(false)
const showJoinMatchModal = ref(false)
const selectedMatch = ref(null)
const joinMatchForm = reactive({
  position: '',
  jerseyNumber: '',
  isStarting: false
})
const showUpdateAthleteModal = ref(false)
const currentAthleteTeam = ref(null)
const athleteUpdateForm = reactive({
  jerseyNumber: ''
})

// 密码修改表单
const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// UI状态
const uiState = reactive({
  showPasswordForm: false,
  showAthleteModal: false,
  showCreateTeamModal: false,
  isLoading: false,
  isSubmitting: false,
  error: '',
  success: '',
  showCurrentPassword: false,
  showNewPassword: false,
  showConfirmPassword: false,
  showLogoutModal: false,
  showCollectorApplyModal: false,
  isEditingName: false
})

// 修改姓名表单
const editNameForm = reactive({
  name: ''
})

const handleEditName = () => {
  console.log('handleEditName called')
  editNameForm.name = userInfo.name
  uiState.isEditingName = true
}

const cancelEditName = () => {
  uiState.isEditingName = false
  editNameForm.name = ''
}

const saveName = async () => {
  if (!editNameForm.name.trim()) {
    showToast('姓名不能为空', 'error')
    return
  }
  
  if (editNameForm.name === userInfo.name) {
    uiState.isEditingName = false
    return
  }

  try {
    const res = await post('/user/update-name', {
      student_id: authStore.studentId,
      name: editNameForm.name
    })
    
    if (res.code === 200) {
      showToast('姓名修改成功')
      userInfo.name = editNameForm.name
      uiState.isEditingName = false
    } else {
      showToast(res.msg || '修改失败', 'error')
    }
  } catch (err) {
    console.error('Update name error:', err)
    showToast('请求失败，请稍后重试', 'error')
  }
}

// 运动员申请表单
const athleteForm = reactive({
  sportId: '',
  teamId: '',
  jerseyNumber: '',
  isCaptain: false
})

// 创建队伍表单
const createTeamForm = reactive({
  teamName: '',
  sportId: '',
  college: '',
  teamType: '',
  avatarUrl: ''
})

const teamList = ref([])
// 队伍管理弹窗状态
const showTeamManagementModal = ref(false)
const currentManagementTeam = ref(null)

// 队内聊天弹窗状态
const showTeamChatModal = ref(false)
const currentChatTeam = ref(null)

// 打开队伍管理
const openTeamManagement = async (team) => {
  currentManagementTeam.value = team
  // 检查是否是队长
  await checkIfCaptain(team.team_id)
  showTeamManagementModal.value = true
}

// 打开队内聊天
const openTeamChat = async (team) => {
  currentChatTeam.value = team
  // 检查是否是队长
  await checkIfCaptain(team.team_id || team.id)
  showTeamChatModal.value = true
}

// 检查是否是队长
const checkIfCaptain = async (teamId) => {
  try {
    const res = await get(`/athlete/check-captain?team_id=${teamId}`)
    if (res.code === 200 && res.data) {
      teamCaptainStatus.value[teamId] = res.data.is_captain || false
    }
  } catch (err) {
    console.error('Check captain error:', err)
    teamCaptainStatus.value[teamId] = false
  }
}

// 获取运动员信息
const fetchMyAthleteInfo = async () => {
  if (!userInfo.roles.includes('athlete')) {
    return
  }
  try {
    const res = await get('/athlete/info')
    if (res.code === 200 && res.data) {
      myAthleteInfo.value = res.data || []
      // 为每个队伍检查是否是队长
      for (const athlete of myAthleteInfo.value) {
        if (athlete.team_id) {
          await checkIfCaptain(athlete.team_id)
        }
      }
    }
  } catch (err) {
    console.error('Fetch athlete info error:', err)
  }
}

// 打开更新运动员信息弹窗
const openUpdateAthleteModal = async (team) => {
  currentAthleteTeam.value = team
  // 获取当前运动员在该队伍中的信息
  try {
    const res = await get(`/athlete/info/team?team_id=${team.team_id}`)
    if (res.code === 200 && res.data) {
      athleteUpdateForm.jerseyNumber = res.data.jersey_number || ''
    }
  } catch (err) {
    console.error('Fetch athlete info by team error:', err)
  }
  showUpdateAthleteModal.value = true
}

// 更新运动员信息
const handleUpdateAthleteInfo = async () => {
  if (!currentAthleteTeam.value) return
  
  try {
    const res = await post('/athlete/update', {
      team_id: currentAthleteTeam.value.team_id,
      jersey_number: athleteUpdateForm.jerseyNumber
    })
    if (res.code === 200) {
      showToast('运动员信息更新成功')
      showUpdateAthleteModal.value = false
      await fetchMyAthleteInfo()
    } else {
      showToast(res.msg || '更新失败', 'error')
    }
  } catch (err) {
    console.error('Update athlete info error:', err)
    showToast('更新失败，请稍后重试', 'error')
  }
}

// 退出队伍
const handleLeaveTeam = async (team) => {
  if (!confirm(`确定要退出队伍 "${team.team_name}" 吗？`)) {
    return
  }
  
  try {
    const res = await post('/athlete/leave-team', {
      team_id: team.team_id
    })
    if (res.code === 200) {
      showToast('已成功退出队伍')
      await fetchMyTeams()
      await fetchMyAthleteInfo()
    } else {
      showToast(res.msg || '退出失败', 'error')
    }
  } catch (err) {
    console.error('Leave team error:', err)
    const errorMsg = err?.response?.data?.message || err?.message || '退出失败，请稍后重试'
    showToast(errorMsg, 'error')
  }
}

// 判断是否可以管理队伍（创建者或队长）
const canManageTeam = (team) => {
  const teamId = team.team_id || team.id
  return team.created_by === userInfo.studentId || teamCaptainStatus.value[teamId] === true
}

// 获取队伍管理角色
const getTeamManagementRole = (team) => {
  if (userInfo.role === 'admin') {
    return 'admin'
  }
  const teamId = team.team_id || team.id
  if (teamCaptainStatus.value[teamId] === true) {
    return 'captain'
  }
  if (team.created_by === userInfo.studentId) {
    return 'captain' // 创建者也有管理权限
  }
  return 'member'
}

// 队伍更新回调
const handleTeamUpdated = async () => {
  await fetchMyTeams()
  // Update currentManagementTeam to point to the new object so avatarUrl updates
  if (currentManagementTeam.value) {
     const newTeam = myTeams.value.find(t => (t.team_id || t.id) === (currentManagementTeam.value.team_id || currentManagementTeam.value.id))
     if (newTeam) {
        currentManagementTeam.value = newTeam
     }
  }
}

// 加载状态
const isLoading = ref(true)

// 数据统计
const stats = computed(() => [
  { label: '已采集场次', value: userInfo.totalCollected, icon: '📊' },
  { label: '已完成任务', value: userInfo.completedTasks, icon: '✅' },
  { label: '待处理任务', value: userInfo.pendingTasks, icon: '⏳' }
])

// 菜单项
const menuItems = computed(() => {
  const items = [
    {
      id: 'account',
      title: '账号信息',
      icon: '👤',
      items: [
        { label: '学号', value: userInfo.studentId },
        { label: '所属专业', value: userInfo.department },
        { label: '注册日期', value: userInfo.joinDate },
        { label: '当前角色', value: getRoleNames(userInfo.roles) }
      ]
    },
    {
      id: 'settings',
      title: '设置',
      icon: '⚙️',
      actions: [
        { label: '修改密码', action: () => uiState.showPasswordForm = true },
        // { label: '通知设置', action: () => showToast('通知设置功能开发中') },
        // { label: '隐私设置', action: () => showToast('隐私设置功能开发中') }
      ]
    }
    /*
    {
      id: 'help',
      title: '帮助与支持',
      icon: '❓',
      actions: [
        { label: '使用指南', action: () => showToast('使用指南功能开发中') },
        { label: '常见问题', action: () => showToast('常见问题功能开发中') },
        { label: '联系客服', action: () => showToast('联系客服功能开发中') }
      ]
    }
    */
  ]
  
  // 如果是普通学生（或运动员，反正不是采集员/管理员），添加申请采集员选项
  if (!userInfo.roles.includes('collector') && !userInfo.roles.includes('admin')) {
    if (userInfo.applyStatus === 'none') {
      items[1].actions.push({
        label: '申请成为采集员',
        action: handleApplyCollector
      })
    } else if (userInfo.applyStatus === 'pending') {
      items[1].actions.push({
        label: '采集员资格审核中...',
        action: () => showToast('您的申请正在审核中，请耐心等待'),
        disabled: true
      })
    }
  }

  // 申请成为运动员 (学生和采集员都可以申请)
  if (!userInfo.roles.includes('admin')) {
    items[1].actions.push({
      label: '申请成为运动员',
      action: openAthleteModal
    })
  }

  // 如果是采集员，添加进入采集员界面选项
  if (userInfo.roles.includes('collector')) {
    items[1].actions.push({
      label: '进入采集员工作台',
      action: () => router.push('/events')
    })
  }

  // 如果是管理员，添加管理后台入口
  if (userInfo.roles.includes('admin') || userInfo.roles.includes('super_admin')) {
      items[1].actions.push({
        label: '进入管理后台',
        action: () => router.push('/admin/dashboard')
      })
  }

  // 队伍创建申请 (任何学生都可以申请创建队伍，或者限制)
  // 这里假设所有学生都可以申请创建
  items[1].actions.push({
    label: '申请创建新队伍',
    action: openCreateTeamModal
  })
  
  return items
})

// 显示提示信息
const showToast = (message, type = 'info') => {
  if (type === 'error') {
    uiState.error = message
    setTimeout(() => { uiState.error = '' }, 3000)
  } else {
    uiState.success = message
    setTimeout(() => { uiState.success = '' }, 3000)
  }
}

// 加载用户信息
const loadUserInfo = async () => {
  isLoading.value = true
  
  try {
    const studentId = authStore.studentId
    if (!studentId) {
      showToast('未登录', 'error')
      router.push('/login')
      return
    }

    const res = await get(`/user/profile/${studentId}`)
    if (res.code === 200 && res.data) {
      const data = res.data
      userInfo.name = data.name
      userInfo.studentId = data.student_id
      userInfo.department = data.college || '未设置'
      userInfo.avatar = data.avatar_url ? `http://localhost:8080${data.avatar_url}` : ''
      userInfo.role = data.role
      userInfo.roles = data.roles || (data.role ? [data.role] : ['student'])
      userInfo.applyStatus = data.apply_status || 'none'
      
      if (data.created_at) {
         userInfo.joinDate = new Date(data.created_at).toLocaleDateString()
      }
      
      // 如果不是采集员，隐藏相关统计或设为0 (这里暂时保留默认值或设为0)
      if (data.role !== 'collector' && data.role !== 'admin') {
         userInfo.totalCollected = 0
         userInfo.completedTasks = 0
         userInfo.pendingTasks = 0
      }

      // 获取我的队伍
      // 只要是学生就可以尝试获取队伍（包括自己创建的待审核队伍）
      fetchMyTeams()
      
      // 如果是运动员，获取参加的比赛和可报名的比赛，以及运动员信息
      if (userInfo.roles.includes('athlete')) {
        fetchMyMatches()
        fetchAvailableMatches()
        fetchMyAthleteInfo()
      }
      
    } else {
      showToast(res.msg || '获取用户信息失败', 'error')
    }
  } catch (error) {
    console.error('Load profile error:', error)
    if (error.status === 404) {
      showToast('用户信息不存在，请重新登录', 'error')
      authStore.clearAuth()
      router.push('/login')
      return
    }
    showToast('加载失败，请重试', 'error')
  } finally {
    isLoading.value = false
  }
}

// 验证密码格式
const validatePassword = (password) => {
  if (password.length < 6) {
    return '密码长度不能少于6位'
  }
  return null
}

// 提交密码修改
const handlePasswordSubmit = async () => {
  // 验证表单
  if (!passwordForm.currentPassword) {
    showToast('请输入当前密码', 'error')
    return
  }
  
  const newPasswordError = validatePassword(passwordForm.newPassword)
  if (newPasswordError) {
    showToast(newPasswordError, 'error')
    return
  }
  
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    showToast('两次输入的新密码不一致', 'error')
    return
  }
  
  uiState.isSubmitting = true
  
  try {
    const res = await post('/user/change-password', {
      student_id: authStore.studentId,
      old_password: passwordForm.currentPassword,
      new_password: passwordForm.newPassword
    })

    if (res.code === 200) {
      showToast('密码修改成功')
      handlePasswordCancel()
    } else {
      showToast(res.msg || '密码修改失败', 'error')
    }
  } catch (error) {
    console.error('Change password error:', error)
    showToast('请求失败，请稍后重试', 'error')
  } finally {
    uiState.isSubmitting = false
  }
}

// 取消密码修改
const handlePasswordCancel = () => {
  passwordForm.currentPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
  uiState.showPasswordForm = false
}

// 退出登录
const handleLogout = () => {
  uiState.showLogoutModal = true
}

// 确认退出
const confirmLogout = () => {
  authStore.clearAuth()
  router.replace('/login')
}

// 申请成为采集员
const handleApplyCollector = () => {
  uiState.showCollectorApplyModal = true
}

// 确认申请成为采集员
const confirmApplyCollector = async () => {
  uiState.showCollectorApplyModal = false
  
  try {
    const res = await post('/user/apply-collector', {
      student_id: authStore.studentId
    })
    if (res.code === 200) {
      showToast('申请已提交，请等待审核')
      userInfo.applyStatus = 'pending'
    } else {
      showToast(res.msg || '申请失败', 'error')
    }
  } catch (err) {
    console.error('Apply collector error:', err)
    showToast('请求失败，请稍后重试', 'error')
  }
}

// 打开运动员申请弹窗
const openAthleteModal = async () => {
  uiState.showAthleteModal = true
  // 获取队伍列表
  fetchTeams()
}

// 获取队伍列表
const fetchTeams = async () => {
  try {
    const res = await get('/teams/list')
    if (res.code === 200) {
      teamList.value = res.data
    }
  } catch (err) {
    console.error('Fetch teams error:', err)
  }
}

// 获取我的队伍
const fetchMyTeams = async () => {
  try {
    const res = await get(`/teams/my?student_id=${authStore.studentId}`)
    if (res.code === 200) {
      myTeams.value = (res.data || []).map(t => ({
        ...t,
        avatar_url: t.avatar_url ? `http://localhost:8080${t.avatar_url}` : ''
      }))
    }
  } catch (err) {
    console.error('Fetch my teams error:', err)
  }
}

// 获取我参加的比赛
const fetchMyMatches = async () => {
  // 只有运动员才需要获取比赛列表
  if (!userInfo.roles.includes('athlete')) {
    return
  }
  
  isLoadingMatches.value = true
  try {
    const res = await get(`/athlete/matches?student_id=${authStore.studentId}`)
    if (res.code === 200 && res.data) {
      myMatches.value = res.data || []
    }
  } catch (err) {
    console.error('Fetch my matches error:', err)
    myMatches.value = []
  } finally {
    isLoadingMatches.value = false
  }
}

// 格式化比赛时间
const formatMatchTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取比赛状态文本
const getMatchStatusText = (status) => {
  const statusMap = {
    'not_started': '未开始',
    'ongoing': '进行中',
    'finished': '已结束',
    'cancelled': '已取消'
  }
  return statusMap[status] || status
}

// 跳转到比赛详情
const goToMatchDetail = (matchId) => {
  router.push(`/match/${matchId}`)
}

// 获取可报名的比赛列表
const fetchAvailableMatches = async () => {
  // 只有运动员才需要获取可报名比赛列表
  if (!userInfo.roles.includes('athlete')) {
    return
  }
  
  isLoadingAvailableMatches.value = true
  try {
    const res = await get(`/athlete/available-matches?student_id=${authStore.studentId}`)
    if (res.code === 200 && res.data) {
      availableMatches.value = res.data || []
    }
  } catch (err) {
    console.error('Fetch available matches error:', err)
    availableMatches.value = []
  } finally {
    isLoadingAvailableMatches.value = false
  }
}

// 打开报名弹窗
const openJoinMatchModal = (match) => {
  selectedMatch.value = match
  // 确定运动员属于哪个队伍
  const myTeamIds = myTeams.value.map(t => t.team_id)
  const teamId = myTeamIds.includes(match.team_a_id) ? match.team_a_id : match.team_b_id
  selectedMatch.value.myTeamId = teamId
  selectedMatch.value.myTeamName = teamId === match.team_a_id ? match.team_a_name : match.team_b_name
  joinMatchForm.position = ''
  joinMatchForm.jerseyNumber = ''
  joinMatchForm.isStarting = false
  showJoinMatchModal.value = true
}

// 提交报名
const handleJoinMatch = async () => {
  if (!selectedMatch.value) {
    return
  }
  
  uiState.isSubmitting = true
  try {
    const res = await post('/athlete/join-match', {
      student_id: authStore.studentId,
      match_id: selectedMatch.value.match_id,
      team_id: selectedMatch.value.myTeamId,
      position: joinMatchForm.position,
      jersey_number: joinMatchForm.jerseyNumber,
      is_starting: joinMatchForm.isStarting
    })
    
    if (res.code === 200) {
      showToast('报名成功')
      showJoinMatchModal.value = false
      // 刷新比赛列表
      fetchMyMatches()
      fetchAvailableMatches()
    } else {
      showToast(res.message || res.msg || '报名失败', 'error')
    }
  } catch (err) {
    console.error('Join match error:', err)
    showToast('报名失败，请稍后重试', 'error')
  } finally {
    uiState.isSubmitting = false
  }
}

// 提交运动员申请
const handleApplyAthleteSubmit = async () => {
  if (!athleteForm.sportId) {
    showToast('请选择运动项目', 'error')
    return
  }
  if (!athleteForm.teamId) {
    showToast('请选择队伍', 'error')
    return
  }
  
  uiState.isSubmitting = true
  const sportIdMap = {
    '1': 'football',
    '2': 'basketball',
    '3': 'badminton',
    '4': 'volleyball'
  }

  try {
    const res = await post('/user/apply-athlete', {
      student_id: authStore.studentId,
      sport_type: sportIdMap[athleteForm.sportId] || 'football',
      team_id: parseInt(athleteForm.teamId),
      jersey_number: athleteForm.jerseyNumber,
      is_captain: athleteForm.isCaptain
    })
    
    if (res.code === 200) {
      showToast('申请提交成功')
      uiState.showAthleteModal = false
      // 刷新用户信息以更新运动员状态
      loadUserInfo()
    } else {
      // 后端返回的是 message 字段
      const errorMsg = res.message || res.msg || '申请失败'
      showToast(errorMsg, 'error')
      console.error('Apply athlete failed:', res)
    }
  } catch (err) {
    console.error('Apply athlete error:', err)
    // 尝试从错误对象中提取消息
    const errorMsg = err?.message || err?.response?.data?.message || '请求失败，请稍后重试'
    showToast(errorMsg, 'error')
  } finally {
    uiState.isSubmitting = false
  }
}

// 打开创建队伍弹窗
const openCreateTeamModal = () => {
  uiState.showCreateTeamModal = true
  // 重置表单
  createTeamForm.teamName = ''
  createTeamForm.sportId = ''
  createTeamForm.teamType = ''
  // 预填学院
  createTeamForm.college = userInfo.department
}

// 提交创建队伍
const handleCreateTeamSubmit = async () => {
  if (!createTeamForm.teamName) {
    showToast('请输入队伍名称', 'error')
    return
  }
  if (!createTeamForm.sportId) {
    showToast('请选择运动项目', 'error')
    return
  }
  if (!createTeamForm.college) {
    showToast('请输入所属学院', 'error')
    return
  }
  if (!createTeamForm.teamType) {
    showToast('请选择队伍类型', 'error')
    return
  }
  
  uiState.isSubmitting = true
  try {
    const res = await post('/teams/apply-create', {
      team_name: createTeamForm.teamName,
      sport_id: parseInt(createTeamForm.sportId),
      college: createTeamForm.college,
      team_type: createTeamForm.teamType,
      avatar_url: createTeamForm.avatarUrl
    })
    
    if (res.code === 200) {
      showToast('队伍创建申请已提交，请等待审核')
      uiState.showCreateTeamModal = false
      // 刷新队伍列表
      fetchTeams()
      // 刷新我的队伍列表（显示新创建的队伍）
      fetchMyTeams()
      // 可选：自动选中新队伍（如果后端返回了ID且已批准，但这里未批准）
    } else {
      showToast(res.msg || '创建失败', 'error')
    }
  } catch (err) {
    console.error('Create team error:', err)
    showToast('请求失败', 'error')
  } finally {
    uiState.isSubmitting = false
  }
}

const getRoleNames = (roles) => {
  if (!roles || roles.length === 0) return '普通学生'
  const names = roles.map(role => {
    switch(role) {
      case 'admin': return '管理员'
      case 'collector': return '赛事采集员'
      case 'athlete': return '校队运动员'
      case 'student': return '普通学生'
      default: return role
    }
  })
  
  // 如果有其他角色，过滤掉普通学生
  const filtered = names.filter(n => n !== '普通学生')
  return filtered.length > 0 ? filtered.join(', ') : '普通学生'
}

// 生命周期钩子
onMounted(() => {
  loadUserInfo()
})
</script>

<template>
  <div class="profile-view">
    
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <button class="back-btn" @click="router.back()">←</button>
        <h1 class="page-title">个人中心</h1>
      </div>
    </div>
    
    <!-- 加载状态 -->
    <div v-if="isLoading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>
    
    <!-- 主内容区域 -->
    <div v-else class="profile-content">
      <!-- 用户信息卡片 -->
      <div class="user-profile-card">
        <div class="avatar-container">
          <div v-if="!userInfo.avatar" class="avatar-placeholder">
            {{ userInfo.name.substring(0, 1) }}
          </div>
          <img v-else :src="userInfo.avatar" alt="用户头像" class="avatar">
          <button class="edit-avatar-btn" @click="triggerFileInput" title="修改头像">
            📷
          </button>
          <input 
            type="file" 
            ref="fileInput" 
            class="hidden-file-input" 
            accept="image/*"
            @change="handleFileChange"
          >
        </div>
        
        <div class="user-info">
          <div v-if="uiState.isEditingName" class="edit-name-container">
            <input 
              v-model="editNameForm.name" 
              class="edit-name-input"
              @keyup.enter="saveName"
              @keyup.esc="cancelEditName"
              ref="nameInput"
              placeholder="请输入姓名"
            >
            <div class="edit-name-actions">
              <button @click="saveName" class="btn-icon save-btn" title="保存">✓</button>
              <button @click="cancelEditName" class="btn-icon cancel-btn" title="取消">✕</button>
            </div>
          </div>
          <div v-else>
            <h2 class="user-name">{{ userInfo.name }}</h2>
          </div>
          <p class="user-username">{{ userInfo.studentId }}</p>
        </div>
        
        <!-- 绝对定位的修改按钮，移出 user-info 以避免布局干扰 -->
         <button 
           v-if="!uiState.isEditingName"
           class="edit-name-btn-fixed" 
           @click.stop="handleEditName" 
           title="修改姓名" 
           type="button"
         >
           <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 4px;"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
           修改名字
         </button>
      </div>
      
      <!-- 数据统计 (仅采集员可见) -->
      <!-- <div v-if="userInfo.role === 'collector' || userInfo.role === 'admin'" class="stats-container">
        <div 
          v-for="stat in stats" 
          :key="stat.label"
          class="stat-card"
        >
          <div class="stat-icon">{{ stat.icon }}</div>
          <div class="stat-content">
            <div class="stat-value">{{ stat.value }}</div>
            <div class="stat-label">{{ stat.label }}</div>
          </div>
        </div>
      </div> -->

      <!-- 我的队伍 (仅运动员可见) -->
      <div v-if="myTeams.length > 0" class="my-teams-container">
        <h3 class="section-title">我的队伍</h3>
        <div class="teams-grid">
          <div v-for="team in myTeams" :key="team.team_id || team.id" class="team-card">
            <div class="team-avatar-container">
              <div v-if="!team.avatar_url" class="team-avatar-placeholder">
                {{ team.team_name.substring(0, 1) }}
              </div>
              <img v-else :src="team.avatar_url" :alt="team.team_name" class="team-avatar">
              <span v-if="teamCaptainStatus[team.team_id || team.id]" class="captain-badge">队长</span>
            </div>
            <div class="team-info">
              <div class="team-name">{{ team.team_name }}</div>
              <div class="team-college">{{ team.college }}</div>
              <div class="team-actions">
                <!-- <button 
                  class="chat-team-btn" 
                  @click.stop="openTeamChat(team)"
                >
                  队内聊天
                </button> -->
                <button 
                  v-if="canManageTeam(team)" 
                  class="manage-team-btn" 
                  @click.stop="openTeamManagement(team)"
                >
                  管理
                </button>
                <button 
                  class="update-athlete-btn" 
                  @click.stop="openUpdateAthleteModal(team)"
                >
                  更新信息
                </button>
                <button 
                  v-if="!teamCaptainStatus[team.team_id || team.id] || team.created_by !== userInfo.studentId"
                  class="leave-team-btn" 
                  @click.stop="handleLeaveTeam(team)"
                >
                  退出
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 我参加的比赛 (仅运动员可见) -->
      <!-- <div v-if="userInfo.roles.includes('athlete')" class="my-matches-container">
        <h3 class="section-title">我参加的比赛</h3>
        <div v-if="isLoadingMatches" class="loading-matches">
          <div class="loading-spinner"></div>
          <p>加载中...</p>
        </div>
        <div v-else-if="myMatches.length === 0" class="no-matches">
          <p>暂无参加的比赛</p>
        </div>
        <div v-else class="matches-list">
          <div 
            v-for="match in myMatches" 
            :key="match.match_id" 
            class="match-card"
            @click="goToMatchDetail(match.match_id)"
          >
            <div class="match-header">
              <div class="match-name">{{ match.match_name || '比赛' }}</div>
              <div class="match-status" :class="'status-' + match.status">
                {{ getMatchStatusText(match.status) }}
              </div>
            </div>
            <div class="match-teams">
              <div class="team-info">
                <div class="team-name">{{ match.team_a_name || '主队' }}</div>
                <div class="team-score">{{ match.score_team_a }}</div>
              </div>
              <div class="vs-divider">VS</div>
              <div class="team-info">
                <div class="team-name">{{ match.team_b_name || '客队' }}</div>
                <div class="team-score">{{ match.score_team_b }}</div>
              </div>
            </div>
            <div class="match-time">
              <span>🕐</span>
              {{ formatMatchTime(match.match_time) }}
            </div>
            <div v-if="match.round" class="match-round">
              {{ match.round }}
            </div>
          </div>
        </div>
      </div> -->

      <!-- 可报名的比赛 (仅运动员可见) -->
      <!-- <div v-if="userInfo.roles.includes('athlete')" class="available-matches-container">
        <h3 class="section-title">可报名的比赛</h3>
        <div v-if="isLoadingAvailableMatches" class="loading-matches">
          <div class="loading-spinner"></div>
          <p>加载中...</p>
        </div>
        <div v-else-if="availableMatches.length === 0" class="no-matches">
          <p>暂无可报名的比赛</p>
        </div>
        <div v-else class="matches-list">
          <div 
            v-for="match in availableMatches" 
            :key="match.match_id" 
            class="match-card available-match-card"
          >
            <div class="match-header">
              <div class="match-name">{{ match.match_name || '比赛' }}</div>
              <div class="match-status" :class="'status-' + match.status">
                {{ getMatchStatusText(match.status) }}
              </div>
            </div>
            <div class="match-teams">
              <div class="team-info">
                <div class="team-name">{{ match.team_a_name || '主队' }}</div>
                <div class="team-score">{{ match.score_team_a }}</div>
              </div>
              <div class="vs-divider">VS</div>
              <div class="team-info">
                <div class="team-name">{{ match.team_b_name || '客队' }}</div>
                <div class="team-score">{{ match.score_team_b }}</div>
              </div>
            </div>
            <div class="match-time">
              <span>🕐</span>
              {{ formatMatchTime(match.match_time) }}
            </div>
            <div v-if="match.round" class="match-round">
              {{ match.round }}
            </div>
            <div class="match-actions">
              <button class="btn-join-match" @click.stop="openJoinMatchModal(match)">
                报名参加
              </button>
              <button class="btn-view-detail" @click.stop="goToMatchDetail(match.match_id)">
                查看详情
              </button>
            </div>
          </div>
        </div>
      </div> -->
      
      <!-- 菜单列表 -->
      <div class="menu-container">
        <div 
          v-for="menu in menuItems" 
          :key="menu.id"
          class="menu-section"
        >
          <div class="menu-title">
            <span class="menu-icon">{{ menu.icon }}</span>
            <span class="menu-text">{{ menu.title }}</span>
          </div>
          
          <div class="menu-content">
            <!-- 信息项 -->
            <div v-if="menu.items" class="info-items">
              <div 
                v-for="item in menu.items" 
                :key="item.label"
                class="info-item"
              >
                <span class="info-label">{{ item.label }}</span>
                <span class="info-value">{{ item.value }}</span>
              </div>
            </div>
            
            <!-- 操作项 -->
            <div v-if="menu.actions" class="action-items">
              <div 
                v-for="action in menu.actions" 
                :key="action.label"
                class="action-item"
                :class="{ 'disabled': action.disabled }"
                @click="!action.disabled && action.action()"
              >
                <span class="action-label" :class="{ 'text-disabled': action.disabled }">{{ action.label }}</span>
                <span class="action-arrow" v-if="!action.disabled">›</span>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 退出登录按钮 -->
        <button 
          @click="handleLogout"
          class="logout-button"
        >
          退出登录
        </button>
      </div>
    </div>
    
    <!-- 退出确认弹窗 -->
    <div v-if="uiState.showLogoutModal" class="modal-overlay" @click="uiState.showLogoutModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3 class="modal-title">确认退出</h3>
          <button @click="uiState.showLogoutModal = false" class="modal-close">×</button>
        </div>
        <div class="modal-body">
          <p>确定要退出登录吗？</p>
        </div>
        <div class="modal-footer">
          <button @click="uiState.showLogoutModal = false" class="btn-secondary">取消</button>
          <button @click="confirmLogout" class="btn-primary danger-btn">确认退出</button>
        </div>
      </div>
    </div>

    <!-- 申请采集员确认弹窗 -->
    <div v-if="uiState.showCollectorApplyModal" class="modal-overlay" @click="uiState.showCollectorApplyModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3 class="modal-title">申请成为采集员</h3>
          <button @click="uiState.showCollectorApplyModal = false" class="modal-close">×</button>
        </div>
        <div class="modal-body">
          <p>确定要申请成为采集员吗？这将提交给管理员审核。</p>
        </div>
        <div class="modal-footer">
          <button @click="uiState.showCollectorApplyModal = false" class="btn-secondary">取消</button>
          <button @click="confirmApplyCollector" class="btn-primary">确认申请</button>
        </div>
      </div>
    </div>

    <!-- 密码修改弹窗 -->
    <div v-if="uiState.showPasswordForm" class="modal-overlay" @click="handlePasswordCancel">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3 class="modal-title">修改密码</h3>
          <button @click="handlePasswordCancel" class="modal-close">×</button>
        </div>
        
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">当前密码</label>
            <div class="password-input-container">
              <input 
                v-model="passwordForm.currentPassword"
                :type="uiState.showCurrentPassword ? 'text' : 'password'"
                placeholder="请输入当前密码"
                class="form-input"
              />
              <button 
                type="button"
                @click="uiState.showCurrentPassword = !uiState.showCurrentPassword"
                class="toggle-password"
              >
                {{ uiState.showCurrentPassword ? '👁️' : '👁️‍🗨️' }}
              </button>
            </div>
          </div>
          
          <div class="form-group">
            <label class="form-label">新密码</label>
            <div class="password-input-container">
              <input 
                v-model="passwordForm.newPassword"
                :type="uiState.showNewPassword ? 'text' : 'password'"
                placeholder="请输入新密码"
                class="form-input"
              />
              <button 
                type="button"
                @click="uiState.showNewPassword = !uiState.showNewPassword"
                class="toggle-password"
              >
                {{ uiState.showNewPassword ? '👁️' : '👁️‍🗨️' }}
              </button>
            </div>
            <p class="form-hint">密码长度不少于6位</p>
          </div>
          
          <div class="form-group">
            <label class="form-label">确认新密码</label>
            <div class="password-input-container">
              <input 
                v-model="passwordForm.confirmPassword"
                :type="uiState.showConfirmPassword ? 'text' : 'password'"
                placeholder="请再次输入新密码"
                class="form-input"
              />
              <button 
                type="button"
                @click="uiState.showConfirmPassword = !uiState.showConfirmPassword"
                class="toggle-password"
              >
                {{ uiState.showConfirmPassword ? '👁️' : '👁️‍🗨️' }}
              </button>
            </div>
          </div>
        </div>
        
        <div class="modal-footer">
          <button 
            @click="handlePasswordCancel"
            class="btn-secondary"
            :disabled="uiState.isSubmitting"
          >
            取消
          </button>
          <button 
            @click="handlePasswordSubmit"
            class="btn-primary"
            :disabled="uiState.isSubmitting"
          >
            <span v-if="!uiState.isSubmitting">确定修改</span>
            <span v-else>修改中...</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 运动员申请弹窗 -->
    <div v-if="uiState.showAthleteModal" class="modal-overlay" @click="uiState.showAthleteModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3 class="modal-title">申请成为运动员</h3>
          <button @click="uiState.showAthleteModal = false" class="modal-close">×</button>
        </div>
        
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">运动项目</label>
            <select v-model="athleteForm.sportId" class="form-input">
              <option value="" disabled>请选择项目</option>
              <option value="1">足球</option>
              <option value="2">篮球</option>
              <option value="3">羽毛球</option>
              <option value="4">排球</option>
            </select>
          </div>
          
          <div class="form-group">
            <label class="form-label">选择队伍</label>
            <div class="select-with-action">
              <select v-model="athleteForm.teamId" class="form-input">
                <option value="" disabled>请选择队伍</option>
                <option v-for="team in teamList" :key="team.team_id" :value="team.team_id">
                  {{ team.team_name }} {{ team.is_approved ? '' : '(待审核)' }}
                </option>
              </select>
              <button @click="openCreateTeamModal" class="btn-text">
                找不到队伍？创建新队伍
              </button>
            </div>
          </div>
          
          <div class="form-group">
            <label class="form-label">球衣号码 (可选)</label>
            <input v-model="athleteForm.jerseyNumber" type="text" class="form-input" placeholder="例如: 10">
          </div>
          
          <div class="form-group">
            <label class="checkbox-label">
              <input v-model="athleteForm.isCaptain" type="checkbox">
              申请成为队长
            </label>
          </div>
        </div>
        
        <div class="modal-footer">
          <button @click="uiState.showAthleteModal = false" class="btn-secondary">取消</button>
          <button @click="handleApplyAthleteSubmit" class="btn-primary" :disabled="uiState.isSubmitting">
            提交申请
          </button>
        </div>
      </div>
    </div>

    <!-- 报名比赛弹窗 -->
    <div v-if="showJoinMatchModal" class="modal-overlay" style="z-index: 1003;" @click="showJoinMatchModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3 class="modal-title">报名参加比赛</h3>
          <button @click="showJoinMatchModal = false" class="modal-close">×</button>
        </div>
        
        <div class="modal-body" v-if="selectedMatch">
          <div class="match-info-summary">
            <h4>{{ selectedMatch.match_name || '比赛' }}</h4>
            <p class="match-teams-info">
              {{ selectedMatch.team_a_name || '主队' }} VS {{ selectedMatch.team_b_name || '客队' }}
            </p>
            <p class="match-time-info">
              <span>🕐</span> {{ formatMatchTime(selectedMatch.match_time) }}
            </p>
            <p class="match-team-info">
              您的队伍：<strong>{{ selectedMatch.myTeamName }}</strong>
            </p>
          </div>
          
          <div class="form-group">
            <label class="form-label">位置 (可选)</label>
            <input 
              v-model="joinMatchForm.position" 
              type="text" 
              class="form-input" 
              placeholder="例如: 前锋、后卫"
            >
          </div>
          
          <div class="form-group">
            <label class="form-label">球衣号码 (可选)</label>
            <input 
              v-model="joinMatchForm.jerseyNumber" 
              type="text" 
              class="form-input" 
              placeholder="例如: 10"
            >
          </div>
          
          <div class="form-group">
            <label class="checkbox-label">
              <input 
                v-model="joinMatchForm.isStarting" 
                type="checkbox"
              >
              申请首发
            </label>
          </div>
        </div>
        
        <div class="modal-footer">
          <button @click="showJoinMatchModal = false" class="btn-secondary">取消</button>
          <button 
            @click="handleJoinMatch" 
            class="btn-primary" 
            :disabled="uiState.isSubmitting"
          >
            <span v-if="!uiState.isSubmitting">确认报名</span>
            <span v-else>报名中...</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 创建队伍弹窗 -->
    <div v-if="uiState.showCreateTeamModal" class="modal-overlay" style="z-index: 1002;" @click="uiState.showCreateTeamModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3 class="modal-title">创建新队伍</h3>
          <button @click="uiState.showCreateTeamModal = false" class="modal-close">×</button>
        </div>
        
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">队伍名称</label>
            <input v-model="createTeamForm.teamName" type="text" class="form-input" placeholder="例如: 计算机学院足球队">
          </div>
          
          <div class="form-group">
            <label class="form-label">运动项目</label>
            <select v-model="createTeamForm.sportId" class="form-input">
              <option value="" disabled>请选择项目</option>
              <option value="1">足球</option>
              <option value="2">篮球</option>
              <option value="3">羽毛球</option>
              <option value="4">排球</option>
            </select>
          </div>

          <div class="form-group">
            <label class="form-label">所属学院/单位</label>
            <input v-model="createTeamForm.college" type="text" class="form-input" placeholder="例如: 计算机学院">
          </div>
          
          <div class="form-group">
            <label class="form-label">队伍类型</label>
            <select v-model="createTeamForm.teamType" class="form-input">
              <option value="" disabled>请选择队伍类型</option>
              <option value="college">院队</option>
              <option value="school">校队</option>
            </select>
          </div>
        </div>
        
        <div class="modal-footer">
          <button @click="uiState.showCreateTeamModal = false" class="btn-secondary">取消</button>
          <button @click="handleCreateTeamSubmit" class="btn-primary" :disabled="uiState.isSubmitting">
            提交创建申请
          </button>
        </div>
      </div>
    </div>
    
    <!-- 错误提示 -->
    <div v-if="uiState.error" class="toast error-toast">
      <span class="toast-icon">❌</span>
      <span class="toast-message">{{ uiState.error }}</span>
    </div>
    
    <!-- 成功提示 -->
    <div v-if="uiState.success" class="toast success-toast">
      <span class="toast-icon">✅</span>
      <span class="toast-message">{{ uiState.success }}</span>
    </div>

    <!-- 队伍管理弹窗 -->
    <TeamManagementModal 
      v-if="currentManagementTeam"
      v-model:visible="showTeamManagementModal"
      :team-id="currentManagementTeam.team_id || currentManagementTeam.id"
      :team-name="currentManagementTeam.team_name"
      :avatar-url="currentManagementTeam.avatar_url"
      :current-user-role="getTeamManagementRole(currentManagementTeam)"
      @refresh="handleTeamUpdated"
    />

    <!-- 队内聊天弹窗 -->
    <div v-if="showTeamChatModal && currentChatTeam" class="team-chat-modal-overlay" @click.self="showTeamChatModal = false">
      <div class="team-chat-modal-content">
        <div class="team-chat-modal-header">
          <h3>队内聊天 - {{ currentChatTeam.team_name }}</h3>
          <button @click="showTeamChatModal = false" class="modal-close">×</button>
        </div>
        <TeamChat 
          :team-id="currentChatTeam.team_id || currentChatTeam.id"
          :team-name="currentChatTeam.team_name"
          :is-captain="teamCaptainStatus[currentChatTeam.team_id || currentChatTeam.id] === true"
        />
      </div>
    </div>

    <!-- 更新运动员信息弹窗 -->
    <div v-if="showUpdateAthleteModal" class="modal-overlay" @click.self="showUpdateAthleteModal = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">更新运动员信息</h3>
          <button @click="showUpdateAthleteModal = false" class="modal-close">×</button>
        </div>
        <div class="modal-body">
          <div v-if="currentAthleteTeam" class="form-group">
            <label>队伍名称</label>
            <input type="text" :value="currentAthleteTeam.team_name" disabled class="form-input" />
          </div>
          <div class="form-group">
            <label>球衣号码</label>
            <input 
              v-model="athleteUpdateForm.jerseyNumber" 
              type="text" 
              placeholder="请输入球衣号码"
              class="form-input"
            />
          </div>
        </div>
        <div class="modal-footer">
          <button class="modal-btn cancel" @click="showUpdateAthleteModal = false">取消</button>
          <button class="modal-btn primary" @click="handleUpdateAthleteInfo" :disabled="uiState.isSubmitting">
            {{ uiState.isSubmitting ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.danger-btn {
  background-color: #ff4d4f;
  border-color: #ff4d4f;
  color: white;
}
.danger-btn:hover {
  background-color: #ff7875;
  border-color: #ff7875;
}

.profile-view {
  min-height: 100vh;
  background-color: #f5f5f5;
  display: flex;
  flex-direction: column;
  padding-bottom: 40px;
}

.page-header {
  padding: 20px;
  background-color: white;
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.header-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.back-btn {
  background: none;
  border: none;
  font-size: 24px;
  color: #666;
  cursor: pointer;
  padding: 0;
  line-height: 1;
  transition: color 0.3s;
}

.back-btn:hover {
  color: #333;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

/* 加载状态 */
.loading-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 16px;
  color: #666;
  padding: 40px;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #1890ff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 主内容区域 */
.profile-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 0 20px;
}

/* 用户信息卡片 */
.user-profile-card {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 24px;
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: relative;
}

.avatar-container {
  position: relative;
}

.avatar-placeholder {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 32px;
  font-weight: 600;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.edit-avatar-btn {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background-color: white;
  border: 1px solid #ddd;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  font-size: 14px;
  transition: all 0.3s;
}

.edit-avatar-btn:hover {
  background-color: #f0f0f0;
  transform: scale(1.1);
}

.hidden-file-input {
  display: none;
}

.user-info {
  flex: 1;
}

.user-name {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin: 0 0 8px 0;
}

.edit-name-btn-fixed {
  position: absolute;
  top: 24px;
  right: 24px;
  display: flex;
  align-items: center;
  padding: 6px 12px;
  background-color: transparent;
  border: 1px solid #d9d9d9;
  border-radius: 20px;
  color: #666;
  font-size: 13px;
  font-weight: 400;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
  white-space: nowrap;
  z-index: 999;
}

.edit-name-btn-fixed:hover {
  color: #1890ff;
  border-color: #1890ff;
  background-color: #e6f7ff;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.edit-name-container {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.edit-name-input {
  font-size: 18px;
  font-weight: 500;
  padding: 6px 10px;
  border: 1px solid #1890ff;
  border-radius: 4px;
  width: 200px;
  outline: none;
}

.edit-name-actions {
  display: flex;
  gap: 6px;
}

.btn-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  border: 1px solid #ddd;
  background: white;
  cursor: pointer;
  font-size: 16px;
  transition: all 0.2s;
}

.save-btn {
  color: #52c41a;
  border-color: #b7eb8f;
  background: #f6ffed;
}

.save-btn:hover {
  background: #d9f7be;
}

.cancel-btn {
  color: #ff4d4f;
  border-color: #ffa39e;
  background: #fff1f0;
}

.cancel-btn:hover {
  background: #ffccc7;
}

.user-username {
  font-size: 14px;
  color: #999;
  margin: 0;
}

/* 数据统计 */
.stats-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s, box-shadow 0.3s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.stat-icon {
  font-size: 32px;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #999;
}

.my-teams-container {
  margin-top: 24px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
  padding-left: 4px;
}

.teams-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 12px;
}

.team-card {
  background: #fff;
  border-radius: 12px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  transition: transform 0.2s;
  position: relative;
}

.team-card:active {
  transform: scale(0.98);
}

.team-avatar-container {
  position: relative;
  margin-bottom: 12px;
}

.captain-badge {
  position: absolute;
  top: -8px;
  right: -8px;
  background: #ff4d4f;
  color: white;
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 10px;
  font-weight: 600;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.team-info {
  width: 100%;
  text-align: center;
}

.team-name {
  font-size: 14px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.team-college {
  font-size: 12px;
  color: #999;
  margin-bottom: 8px;
}

.team-actions {
  display: flex;
  flex-direction: column;
  gap: 6px;
  width: 100%;
  margin-top: 8px;
}

.manage-team-btn,
.update-athlete-btn,
.leave-team-btn {
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
  width: 100%;
}

.manage-team-btn {
  background: #1890ff;
  color: white;
}

.manage-team-btn:hover {
  background: #40a9ff;
}

.update-athlete-btn {
  background: #52c41a;
  color: white;
}

.update-athlete-btn:hover {
  background: #73d13d;
}

.leave-team-btn {
  background: #ff4d4f;
  color: white;
}

.leave-team-btn:hover {
  background: #ff7875;
}

.chat-team-btn {
  background: #722ed1;
  color: white;
}

.chat-team-btn:hover {
  background: #9254de;
}

/* 我参加的比赛样式 */
.my-matches-container {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.loading-matches {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  color: #999;
}

.no-matches {
  text-align: center;
  padding: 40px;
  color: #999;
}

.matches-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.match-card {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid #e9ecef;
}

.match-card:hover {
  background: #e9ecef;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.match-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.match-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.match-status {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-not_started {
  background: #e6f7ff;
  color: #1890ff;
}

.status-ongoing {
  background: #fff7e6;
  color: #fa8c16;
}

.status-finished {
  background: #f6ffed;
  color: #52c41a;
}

.status-cancelled {
  background: #fff1f0;
  color: #ff4d4f;
}

.match-teams {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  gap: 16px;
}

.team-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.team-name {
  font-size: 14px;
  color: #666;
  margin-bottom: 4px;
}

.team-score {
  font-size: 24px;
  font-weight: 700;
  color: #333;
}

.vs-divider {
  font-size: 14px;
  color: #999;
  font-weight: 500;
}

.match-time {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #999;
  margin-bottom: 4px;
}

.match-round {
  font-size: 12px;
  color: #666;
  font-style: italic;
}

/* 可报名比赛样式 */
.available-matches-container {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-top: 24px;
}

.available-match-card {
  border-left: 4px solid #1890ff;
}

.match-actions {
  display: flex;
  gap: 8px;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #e9ecef;
}

.btn-join-match {
  flex: 1;
  padding: 8px 16px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-join-match:hover {
  background: #40a9ff;
  transform: translateY(-1px);
}

.btn-join-match:active {
  transform: translateY(0);
}

.btn-view-detail {
  flex: 1;
  padding: 8px 16px;
  background: #f0f0f0;
  color: #666;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-view-detail:hover {
  background: #e6e6e6;
  border-color: #bfbfbf;
}

/* 报名弹窗样式 */
.match-info-summary {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 20px;
}

.match-info-summary h4 {
  margin: 0 0 12px 0;
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.match-teams-info {
  margin: 8px 0;
  font-size: 14px;
  color: #666;
}

.match-time-info {
  margin: 8px 0;
  font-size: 13px;
  color: #999;
  display: flex;
  align-items: center;
  gap: 6px;
}

.match-team-info {
  margin: 8px 0 0 0;
  font-size: 14px;
  color: #333;
}

.match-team-info strong {
  color: #1890ff;
  font-weight: 600;
}

.team-avatar-container {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  /* overflow: hidden; 移除此属性以允许徽章显示 */
  margin-bottom: 8px;
  background-color: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.team-avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%; /* 确保图片是圆形的 */
}

.team-avatar-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #4caf50;
  color: white;
  font-size: 20px;
  font-weight: bold;
  border-radius: 50%; /* 确保占位符是圆形的 */
}

.team-info {
  text-align: center;
}

.team-name {
  font-size: 14px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 120px;
}

.team-college {
  font-size: 12px;
  color: #666;
}

/* 菜单区域 */
.menu-container {
  margin-top: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.menu-section {
  background-color: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.menu-title {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  background-color: #fafafa;
  border-bottom: 1px solid #f0f0f0;
}

.menu-icon {
  font-size: 20px;
}

.menu-text {
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.menu-content {
  padding: 8px 0;
}

.info-items,
.action-items {
  display: flex;
  flex-direction: column;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  border-bottom: 1px solid #f5f5f5;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  font-size: 14px;
  color: #666;
}

.info-value {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.action-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f5f5f5;
  cursor: pointer;
  transition: background-color 0.3s;
}

.action-item:last-child {
  border-bottom: none;
}

.action-item:hover {
  background-color: #fafafa;
}

.action-item.disabled {
  cursor: not-allowed;
  opacity: 0.7;
}

.action-item.disabled:hover {
  background-color: transparent;
}

.text-disabled {
  color: #999;
}

.action-label {
  font-size: 14px;
  color: #333;
}

.action-arrow {
  font-size: 18px;
  color: #999;
}

/* 退出登录按钮 */
.logout-button {
  padding: 16px;
  background-color: white;
  border: 1px solid #ff4d4f;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 500;
  color: #ff4d4f;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.logout-button:hover {
  background-color: #ff4d4f;
  color: white;
  box-shadow: 0 4px 16px rgba(255, 77, 79, 0.3);
}

/* 队内聊天弹窗 */
.team-chat-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.team-chat-modal-content {
  width: 90%;
  max-width: 900px;
  height: 80vh;
  background: white;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.team-chat-modal-header {
  padding: 16px 20px;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #fafafa;
}

.team-chat-modal-header h3 {
  margin: 0;
  font-size: 18px;
}

/* 密码修改弹窗 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  background-color: white;
  border-radius: 12px;
  width: 100%;
  max-width: 480px;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.2);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #f0f0f0;
}

.modal-title {
  font-size: 18px;
  font-weight: 500;
  color: #333;
  margin: 0;
}

.modal-close {
  width: 32px;
  height: 32px;
  border: none;
  background: none;
  font-size: 24px;
  color: #999;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  transition: all 0.3s;
}

.modal-close:hover {
  background-color: #f5f5f5;
  color: #333;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.password-input-container {
  position: relative;
}

.form-input {
  width: 100%;
  padding: 12px 40px 12px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.3s;
}

.form-input:focus {
  outline: none;
  border-color: #1890ff;
}

.toggle-password {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  border: none;
  background: none;
  font-size: 18px;
  cursor: pointer;
}

.form-hint {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
  margin-bottom: 0;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px;
  border-top: 1px solid #f0f0f0;
}

.btn-secondary,
.btn-primary {
  padding: 10px 20px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-secondary {
  border: 1px solid #ddd;
  background-color: white;
  color: #333;
}

.btn-secondary:hover:not(:disabled) {
  border-color: #1890ff;
  color: #1890ff;
}

.btn-primary {
  border: 1px solid #1890ff;
  background-color: #1890ff;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #40a9ff;
  border-color: #40a9ff;
}

.btn-primary:disabled {
  background-color: #bae7ff;
  border-color: #bae7ff;
  cursor: not-allowed;
}

/* Name Edit Styles */
.edit-name-container {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.edit-name-input {
  font-size: 20px;
  font-weight: 600;
  padding: 4px 8px;
  border: 1px solid #1890ff;
  border-radius: 4px;
  width: 180px;
  color: #333;
}

.edit-name-input:focus {
  outline: none;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.edit-name-actions {
  display: flex;
  gap: 4px;
}

.btn-icon {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 16px;
  padding: 4px;
  border-radius: 4px;
  transition: background-color 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
}

.save-btn {
  color: #52c41a;
  border: 1px solid #b7eb8f;
  background-color: #f6ffed;
}

.save-btn:hover {
  background-color: #d9f7be;
}

.cancel-btn {
  color: #ff4d4f;
  border: 1px solid #ffa39e;
  background-color: #fff1f0;
}

.cancel-btn:hover {
  background-color: #ffccc7;
}

.edit-name-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 16px;
  color: #999;
  opacity: 0;
  transition: all 0.2s;
  margin-left: 8px;
  padding: 2px 6px;
  border-radius: 4px;
}

.edit-name-btn:hover {
  background-color: #f0f0f0;
  color: #1890ff;
}

.user-name:hover .edit-name-btn {
  opacity: 1;
}

.btn-primary:hover:not(:disabled) {
  background-color: #40a9ff;
  border-color: #40a9ff;
}

.btn-secondary:disabled,
.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 提示信息 */
.toast {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 1001;
  animation: slideDown 0.3s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translate(-50%, -20px);
  }
  to {
    opacity: 1;
    transform: translate(-50%, 0);
  }
}

.error-toast {
  background-color: #fff2f0;
  border: 1px solid #ffccc7;
  color: #f5222d;
}

.success-toast {
  background-color: #f6ffed;
  border: 1px solid #b7eb8f;
  color: #52c41a;
}

.toast-icon {
  font-size: 18px;
}

.toast-message {
  font-size: 14px;
  font-weight: 500;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .page-header {
    padding: 16px;
  }
  
  .page-title {
    font-size: 20px;
  }
  
  .profile-content {
    padding: 0 16px;
    gap: 16px;
  }
  
  .user-profile-card {
    flex-direction: column;
    text-align: center;
    padding: 20px;
  }
  
  .avatar-placeholder,
  .avatar {
    width: 64px;
    height: 64px;
  }
  
  .avatar-placeholder {
    font-size: 24px;
  }
  
  .user-name {
    font-size: 20px;
  }
  
  .stats-container {
    grid-template-columns: repeat(3, 1fr);
  }
  
  .stat-card {
    flex-direction: column;
    text-align: center;
    padding: 16px;
    gap: 12px;
  }
  
  .stat-icon {
    font-size: 24px;
  }
  
  .stat-value {
    font-size: 20px;
  }
  
  .menu-section {
    border-radius: 8px;
  }
  
  .menu-title {
    padding: 14px 16px;
  }
  
  .info-item {
    padding: 12px 16px;
  }
  
  .action-item {
    padding: 14px 16px;
  }
  
  .logout-button {
    border-radius: 8px;
  }
  
  .modal-overlay {
    padding: 16px;
  }
  
  .modal-content {
    border-radius: 8px;
  }
  
  .modal-header,
  .modal-body,
  .modal-footer {
    padding: 16px;
  }
  
  .modal-title {
    font-size: 16px;
  }
  
  .modal-footer {
    flex-direction: column;
  }
  
  .btn-secondary,
  .btn-primary {
    width: 100%;
  }
  
  .toast {
    left: 16px;
    right: 16px;
    transform: none;
  }
  
  @keyframes slideDown {
    from {
      opacity: 0;
      transform: translateY(-20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
}

.select-with-action {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.btn-text {
  background: none;
  border: none;
  color: #1890ff;
  cursor: pointer;
  padding: 0;
  font-size: 14px;
  text-align: left;
  align-self: flex-start;
}

.btn-text:hover {
  text-decoration: underline;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  user-select: none;
}

.manage-team-btn {
  margin-top: 8px;
  padding: 6px 12px;
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.3s;
  width: 100%;
}

.manage-team-btn:hover {
  background-color: #40a9ff;
  transform: translateY(-1px);
}
</style>
