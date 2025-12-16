<template>
  <div class="admin-dashboard fade-in">
    <div class="header">
      <div class="header-content">
        <h1>管理控制台</h1>
        <p>管理队伍、赛事和比赛。</p>
      </div>
      <button type="button" class="logout-btn" @click.stop.prevent="logout">退出登录</button>
    </div>

    <div class="tabs">
      <button 
        v-for="tab in tabs" 
        :key="tab.id"
        :class="['tab-btn', { active: currentTab === tab.id }]"
        @click="currentTab = tab.id"
      >
        {{ tab.label }}
      </button>
    </div>

    <div class="content-area">
      <!-- Create Team -->
      <div v-if="currentTab === 'team'" class="form-container slide-in">
        <h2>队伍管理</h2>
        
        <!-- Create Form -->
        <div class="collapsible-section">
          <h3>创建新队伍</h3>
          <form @submit.prevent="createTeam" class="compact-form">
            <div class="form-row">
              <div class="form-group">
                <label>队伍名称</label>
                <input v-model="teamForm.name" required placeholder="例如：计算机学院 A 队" />
              </div>
              <div class="form-group">
                <label>运动项目</label>
                <select v-model="teamForm.sport" required>
                  <option value="" disabled>选择项目</option>
                  <option v-for="sport in sports" :key="sport.value" :value="sport.value">
                    {{ sport.label }}
                  </option>
                </select>
              </div>
              <div class="form-group">
                <label>所属学院</label>
                <input v-model="teamForm.college" required placeholder="例如：计算机学院" />
              </div>
              <div class="form-group">
                <label>队伍类型</label>
                <input v-model="teamForm.team_type" placeholder="例如：男子组" />
              </div>
            </div>
            <button type="submit" class="submit-btn" :disabled="loading">
              {{ loading ? '创建中...' : '创建队伍' }}
            </button>
          </form>
        </div>

        <!-- Team List -->
        <div class="list-section">
          <h3>现有队伍列表</h3>
          <div v-if="teams.length === 0" class="no-data">暂无队伍</div>
          <div v-else class="data-table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>队伍名称</th>
                  <th>学院</th>
                  <th>项目</th>
                  <th>类型</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="team in teams" :key="team.team_id">
                  <td>{{ team.team_id }}</td>
                  <td>{{ team.team_name }}</td>
                  <td>{{ team.college }}</td>
                  <td>
                    {{ sports.find(s => s.value === (team.sport_id === 1 ? 'football' : team.sport_id === 2 ? 'basketball' : team.sport_id === 3 ? 'badminton' : team.sport_id === 4 ? 'volleyball' : 'water-sports'))?.label || team.sport_id }}
                  </td>
                  <td>{{ team.team_type }}</td>
                  <td class="actions-cell">
                    <button class="action-btn view-btn" @click="openMembersModal(team)">管理</button>
                    <button class="action-btn delete-btn" @click="deleteTeam(team)">删除</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Create Event -->
      <div v-if="currentTab === 'event'" class="form-container slide-in">
        <h2>赛事管理</h2>
        
        <!-- Create Form -->
        <div class="collapsible-section">
          <h3>{{ isEditingEvent ? '编辑赛事' : '创建新赛事' }}</h3>
          <form @submit.prevent="createEvent" class="compact-form">
            <div class="form-row">
              <div class="form-group">
                <label>赛事名称</label>
                <input v-model="eventForm.name" required placeholder="例如：2024 春季足球杯" />
              </div>
              <div class="form-group">
                <label>运动项目</label>
                <select v-model="eventForm.sport" required>
                  <option value="" disabled>选择项目</option>
                  <option v-for="sport in sports" :key="sport.value" :value="sport.value">
                    {{ sport.label }}
                  </option>
                </select>
              </div>
              <div class="form-group">
                <label>开始日期</label>
                <input type="date" v-model="eventForm.start_date" required />
              </div>
              <div class="form-group">
                <label>赛制类型</label>
                <input v-model="eventForm.format_type" placeholder="例如：淘汰赛" />
              </div>
            </div>
            <div class="form-actions">
              <button type="submit" class="submit-btn" :disabled="loading">
                {{ loading ? (isEditingEvent ? '保存中...' : '创建中...') : (isEditingEvent ? '保存修改' : '创建赛事') }}
              </button>
              <button v-if="isEditingEvent" type="button" class="cancel-btn" @click="cancelEditEvent">取消</button>
            </div>
          </form>
        </div>

        <!-- Event List -->
        <div class="list-section">
          <h3>现有赛事列表</h3>
          <div v-if="events.length === 0" class="no-data">暂无赛事</div>
          <div v-else class="data-table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>赛事名称</th>
                  <th>项目</th>
                  <th>状态</th>
                  <th>开始时间</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="event in events" :key="event.event_id">
                  <td>{{ event.event_id }}</td>
                  <td>{{ event.event_name }}</td>
                  <td>
                     {{ sports.find(s => s.value === (event.sport_id === 1 ? 'football' : event.sport_id === 2 ? 'basketball' : event.sport_id === 3 ? 'badminton' : event.sport_id === 4 ? 'volleyball' : 'water-sports'))?.label || event.sport_id }}
                  </td>
                  <td>
                    <span :class="['status-badge', event.status]">{{ event.status }}</span>
                  </td>
                  <td>{{ new Date(event.start_date).toLocaleDateString() }}</td>
                  <td class="actions-cell">
                    <button class="action-btn view-btn" @click="editEvent(event)">编辑</button>
                    <button class="action-btn delete-btn" @click="deleteEvent(event)">删除</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Create Match -->
      <div v-if="currentTab === 'match'" class="form-container slide-in">
        <h2>比赛管理</h2>
        
        <!-- Create Form -->
        <div class="collapsible-section">
          <h3>{{ isEditingMatch ? '编辑比赛' : '创建新比赛' }}</h3>
          <form @submit.prevent="createMatch" class="compact-form">
            <div class="form-row">
              <div class="form-group">
                <label>赛事</label>
                <div v-if="events.length > 0">
                   <select v-model="matchForm.event_id" required>
                    <option value="" disabled>选择赛事</option>
                    <option v-for="e in events" :key="e.event_id" :value="e.event_id">
                      {{ e.event_name }}
                    </option>
                  </select>
                </div>
                <div v-else>
                   <input v-model="matchForm.event_id" required placeholder="赛事 ID" />
                </div>
              </div>
              
              <div class="form-group">
                <label>队伍 A</label>
                <div v-if="teams.length > 0">
                  <select v-model="matchForm.team_a_id" required>
                     <option value="" disabled>选择队伍 A</option>
                     <option v-for="t in teams" :key="t.team_id" :value="t.team_id">
                       {{ t.team_name }}
                     </option>
                  </select>
                </div>
                 <div v-else>
                   <input v-model="matchForm.team_a_id" required placeholder="队伍 A ID" />
                </div>
              </div>

              <div class="form-group">
                <label>队伍 B</label>
                 <div v-if="teams.length > 0">
                  <select v-model="matchForm.team_b_id" required>
                     <option value="" disabled>选择队伍 B</option>
                     <option v-for="t in teams" :key="t.team_id" :value="t.team_id">
                       {{ t.team_name }}
                     </option>
                  </select>
                </div>
                <div v-else>
                   <input v-model="matchForm.team_b_id" required placeholder="队伍 B ID" />
                </div>
              </div>

              <div class="form-group">
                <label>比赛名称</label>
                <input v-model="matchForm.name" required placeholder="例如：决赛" />
              </div>

              <div class="form-group">
                <label>时间</label>
                <input type="datetime-local" v-model="matchForm.time" required />
              </div>
            </div>

            <div class="form-actions">
              <button type="submit" class="submit-btn" :disabled="loading">
                {{ loading ? (isEditingMatch ? '保存中...' : '创建中...') : (isEditingMatch ? '保存修改' : '创建比赛') }}
              </button>
              <button v-if="isEditingMatch" type="button" class="cancel-btn" @click="cancelEditMatch">取消</button>
            </div>
          </form>
        </div>

        <!-- Match List -->
        <div class="list-section">
          <h3>现有比赛列表</h3>
          <div v-if="matches.length === 0" class="no-data">暂无比赛</div>
          <div v-else class="data-table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>比赛名称</th>
                  <th>赛事ID</th>
                  <th>队伍 A</th>
                  <th>队伍 B</th>
                  <th>时间</th>
                  <th>状态</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="match in matches" :key="match.match_id">
                  <td>{{ match.match_id }}</td>
                  <td>{{ match.match_name }}</td>
                  <td>{{ match.event_id }}</td>
                  <td>{{ match.team_a_name || match.team_a_id }}</td>
                  <td>{{ match.team_b_name || match.team_b_id }}</td>
                  <td>{{ new Date(match.match_time).toLocaleString() }}</td>
                  <td>
                    <span :class="['status-badge', match.status]">{{ match.status }}</span>
                  </td>
                  <td class="actions-cell">
                    <button class="action-btn view-btn" @click="editMatch(match)">编辑</button>
                    <button class="action-btn delete-btn" @click="deleteMatch(match)">删除</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Collector Requests -->
      <div v-if="currentTab === 'collectors'" class="form-container slide-in">
        <h2>待审核采集员申请</h2>
        <div v-if="collectorRequests.length === 0" class="no-data">
          暂无待审核申请。
        </div>
        <div v-else class="requests-list">
          <div v-for="req in collectorRequests" :key="req.student_id" class="request-item">
            <div class="request-info">
              <h3>{{ req.name }} ({{ req.student_id }})</h3>
              <p>{{ req.college }}</p>
              <small>申请时间: {{ new Date(req.apply_time).toLocaleString() }}</small>
            </div>
            <button type="button" class="approve-btn" @click.stop.prevent="approveCollector(req.student_id)">
              批准
            </button>
          </div>
        </div>

        <h2 style="margin-top: 30px;">现有采集员列表</h2>
        <div v-if="collectors.length === 0" class="no-data">暂无采集员</div>
        <div v-else class="data-table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>学号</th>
                <th>姓名</th>
                <th>学院</th>
                <th>加入时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="c in collectors" :key="c.student_id">
                <td>{{ c.student_id }}</td>
                <td>{{ c.name }}</td>
                <td>{{ c.college }}</td>
                <td>{{ new Date(c.apply_time).toLocaleString() }}</td>
                <td>
                  <button class="action-btn delete-btn" @click="deleteCollector(c.student_id)">删除</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Athlete Requests -->
      <div v-if="currentTab === 'athletes'" class="form-container slide-in">
        <h2>待审核运动员申请</h2>
        <div v-if="athleteRequests.length === 0" class="no-data">
          暂无待审核申请。
        </div>
        <div v-else class="requests-list">
          <div v-for="req in athleteRequests" :key="req.team_member_id" class="request-item">
            <div class="request-info">
              <h3>{{ req.name }} ({{ req.student_id }})</h3>
              <p>申请加入: {{ req.team_name }} ({{ req.sport_type }})</p>
              <small>申请时间: {{ new Date(req.apply_time).toLocaleString() }}</small>
            </div>
            <div class="actions-cell">
              <button type="button" class="approve-btn" @click.stop.prevent="approveAthlete(req.team_member_id)">
                批准
              </button>
              <button type="button" class="reject-btn" @click.stop.prevent="rejectAthlete(req.team_member_id)">
                拒绝
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Team Requests -->
      <div v-if="currentTab === 'teams'" class="form-container slide-in">
        <h2>待审核队伍申请</h2>
        <div v-if="teamRequests.length === 0" class="no-data">
          暂无待审核申请。
        </div>
        <div v-else class="requests-list">
          <div v-for="req in teamRequests" :key="req.team_id" class="request-item">
            <div class="request-info">
              <h3>{{ req.team_name }}</h3>
              <p>创建人: {{ req.created_by }} | 学院: {{ req.college }}</p>
              <p>项目: {{ req.sport_id }} | 类型: {{ req.team_type }}</p>
              <small>申请时间: {{ new Date(req.created_at).toLocaleString() }}</small>
            </div>
            <div class="actions-cell">
              <button type="button" class="approve-btn" @click.stop.prevent="approveTeam(req.team_id)">
                批准
              </button>
              <button type="button" class="reject-btn" @click.stop.prevent="rejectTeam(req.team_id)">
                拒绝
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Confirmation Modal -->
    <div v-if="showModal" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3 class="modal-title">{{ modalTitle }}</h3>
          <button @click="closeModal" class="modal-close">×</button>
        </div>
        <div class="modal-body">
          <p>{{ modalMessage }}</p>
        </div>
        <div class="modal-footer">
          <button @click="closeModal" class="btn-secondary">取消</button>
          <button @click="confirmModal" class="btn-primary">确认</button>
        </div>
      </div>
    </div>

    <!-- Event Edit Modal -->
    <div v-if="showEventEditModal" class="modal-overlay" @click.self="cancelEditEvent">
      <div class="modal-content large-modal" @click.stop>
        <div class="modal-header">
          <h3 class="modal-title">编辑赛事</h3>
          <button @click="cancelEditEvent" class="modal-close">×</button>
        </div>
        <div class="modal-body">
          <form id="eventEditForm" @submit.prevent="createEvent" class="compact-form">
            <div class="form-row">
              <div class="form-group">
                <label>赛事名称</label>
                <input v-model="eventForm.name" required placeholder="例如：2024 春季足球杯" />
              </div>
              <div class="form-group">
                <label>运动项目</label>
                <select v-model="eventForm.sport" required>
                  <option value="" disabled>选择项目</option>
                  <option v-for="sport in sports" :key="sport.value" :value="sport.value">
                    {{ sport.label }}
                  </option>
                </select>
              </div>
              <div class="form-group">
                <label>开始日期</label>
                <input type="date" v-model="eventForm.start_date" required />
              </div>
              <div class="form-group">
                <label>赛制类型</label>
                <input v-model="eventForm.format_type" placeholder="例如：淘汰赛" />
              </div>
            </div>
          </form>
        </div>
        <div class="modal-footer">
          <button @click="cancelEditEvent" class="btn-secondary">取消</button>
          <button type="submit" form="eventEditForm" class="btn-primary" :disabled="loading">
            {{ loading ? '保存中...' : '保存修改' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Match Edit Modal -->
    <div v-if="showMatchEditModal" class="modal-overlay" @click.self="cancelEditMatch">
      <div class="modal-content large-modal" @click.stop>
        <div class="modal-header">
          <h3 class="modal-title">编辑比赛</h3>
          <button @click="cancelEditMatch" class="modal-close">×</button>
        </div>
        <div class="modal-body">
          <form id="matchEditForm" @submit.prevent="createMatch" class="compact-form">
            <div class="form-row">
              <div class="form-group">
                <label>赛事</label>
                <div v-if="events.length > 0">
                   <select v-model="matchForm.event_id" required>
                    <option value="" disabled>选择赛事</option>
                    <option v-for="e in events" :key="e.event_id" :value="e.event_id">
                      {{ e.event_name }}
                    </option>
                  </select>
                </div>
                <div v-else>
                   <input v-model="matchForm.event_id" required placeholder="赛事 ID" />
                </div>
              </div>
              
              <div class="form-group">
                <label>队伍 A</label>
                <div v-if="teams.length > 0">
                  <select v-model="matchForm.team_a_id" required>
                     <option value="" disabled>选择队伍 A</option>
                     <option v-for="t in teams" :key="t.team_id" :value="t.team_id">
                       {{ t.team_name }}
                     </option>
                  </select>
                </div>
                 <div v-else>
                   <input v-model="matchForm.team_a_id" required placeholder="队伍 A ID" />
                </div>
              </div>

              <div class="form-group">
                <label>队伍 B</label>
                 <div v-if="teams.length > 0">
                  <select v-model="matchForm.team_b_id" required>
                     <option value="" disabled>选择队伍 B</option>
                     <option v-for="t in teams" :key="t.team_id" :value="t.team_id">
                       {{ t.team_name }}
                     </option>
                  </select>
                </div>
                <div v-else>
                   <input v-model="matchForm.team_b_id" required placeholder="队伍 B ID" />
                </div>
              </div>

              <div class="form-group">
                <label>比赛名称</label>
                <input v-model="matchForm.name" required placeholder="例如：决赛" />
              </div>

              <div class="form-group">
                <label>时间</label>
                <input type="datetime-local" v-model="matchForm.time" required />
              </div>
            </div>
          </form>
        </div>
        <div class="modal-footer">
          <button @click="cancelEditMatch" class="btn-secondary">取消</button>
          <button type="submit" form="matchEditForm" class="btn-primary" :disabled="loading">
            {{ loading ? '保存中...' : '保存修改' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Reused Team Management Modal -->
    <TeamManagementModal 
      v-model:visible="showMembersModal"
      :team-id="currentTeamId"
      :team-name="currentTeamName"
      current-user-role="admin"
      @refresh="fetchTeams"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import http from '../utils/http'
import { useAuthStore } from '@/stores/auth'
import TeamManagementModal from '../components/TeamManagementModal.vue'

const router = useRouter()
const authStore = useAuthStore()

const tabs = [
  { id: 'team', label: '队伍管理' },
  { id: 'event', label: '赛事管理' },
  { id: 'match', label: '比赛管理' },
  { id: 'collectors', label: '采集员审核' },
  { id: 'athletes', label: '运动员审核' },
  { id: 'teams', label: '队伍审核' }
]

const currentTab = ref('team')
const loading = ref(false)
const sports = [
  { value: 'football', label: '足球' },
  { value: 'basketball', label: '篮球' },
  { value: 'badminton', label: '羽毛球' },
  { value: 'volleyball', label: '排球' },
  { value: 'water-sports', label: '水上运动' }
]

// Forms
const teamForm = ref({
  name: '',
  sport: '',
  college: '',
  team_type: ''
})

const eventForm = ref({
  name: '',
  sport: '',
  start_date: '',
  format_type: ''
})

const matchForm = ref({
  event_id: '',
  team_a_id: '',
  team_b_id: '',
  name: '',
  time: ''
})

// Data Sources
const teams = ref([])
const events = ref([])
const matches = ref([]) // Added for match management
const collectorRequests = ref([])
const collectors = ref([])
const athleteRequests = ref([])
const teamRequests = ref([])

// Team Members Modal State
const showMembersModal = ref(false)
const currentTeamName = ref('')
const currentTeamId = ref(null)

// Edit State
const isEditingEvent = ref(false)
const currentEventId = ref(null)
const isEditingMatch = ref(false)
const currentMatchId = ref(null)

// Modal State
const showModal = ref(false)
const showEventEditModal = ref(false)
const showMatchEditModal = ref(false)
const modalTitle = ref('')
const modalMessage = ref('')
const modalConfirmAction = ref(null)

const openModal = (title, message, action) => {
  modalTitle.value = title
  modalMessage.value = message
  modalConfirmAction.value = action
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  modalConfirmAction.value = null
}

const confirmModal = async () => {
  if (modalConfirmAction.value) {
    await modalConfirmAction.value()
  }
  closeModal()
}

// Members Modal Actions
const openMembersModal = (team) => {
  currentTeamName.value = team.team_name
  currentTeamId.value = team.team_id
  showMembersModal.value = true
}

const deleteTeam = (team) => {
  openModal(
    '确认删除队伍',
    `警告：确定要删除队伍 "${team.team_name}" 吗？此操作将移除所有队员且不可恢复！`,
    async () => {
      try {
        await http.delete(`/teams/delete?id=${team.team_id}`)
        alert('队伍已删除')
        fetchTeams()
      } catch (e) {
        alert('删除失败: ' + e.message)
      }
    }
  )
}

const editEvent = (event) => {
  isEditingEvent.value = true
  currentEventId.value = event.event_id
  eventForm.value = {
    name: event.event_name,
    sport: sports.find(s => s.value === (event.sport_id === 1 ? 'football' : event.sport_id === 2 ? 'basketball' : event.sport_id === 3 ? 'badminton' : event.sport_id === 4 ? 'volleyball' : 'water-sports'))?.value || '',
    start_date: event.start_date.split('T')[0],
    format_type: event.format_type || ''
  }
  showEventEditModal.value = true
}

const cancelEditEvent = () => {
  isEditingEvent.value = false
  showEventEditModal.value = false
  currentEventId.value = null
  eventForm.value = { name: '', sport: '', start_date: '', format_type: '' }
}

const deleteEvent = (event) => {
  openModal(
    '确认删除赛事',
    `警告：确定要删除赛事 "${event.event_name}" 吗？此操作将删除相关的所有比赛和数据！`,
    async () => {
      try {
        await http.delete(`/events/delete?id=${event.event_id}`)
        alert('赛事已删除')
        fetchEvents()
      } catch (e) {
        alert('删除失败: ' + e.message)
      }
    }
  )
}

const editMatch = (match) => {
  isEditingMatch.value = true
  currentMatchId.value = match.match_id
  
  // Format datetime-local: YYYY-MM-DDTHH:mm
  // Assuming match.match_time is ISO string from backend
  const date = new Date(match.match_time)
  // Adjust to local ISO string for input
  const localIso = new Date(date.getTime() - (date.getTimezoneOffset() * 60000)).toISOString().slice(0, 16)
  
  matchForm.value = {
    event_id: match.event_id,
    team_a_id: match.team_a_id,
    team_b_id: match.team_b_id,
    name: match.match_name,
    time: localIso
  }
  
  showMatchEditModal.value = true
}

const cancelEditMatch = () => {
  isEditingMatch.value = false
  showMatchEditModal.value = false
  currentMatchId.value = null
  matchForm.value = { event_id: '', team_a_id: '', team_b_id: '', name: '', time: '' }
}

const deleteMatch = (match) => {
  openModal(
    '确认删除比赛',
    `警告：确定要删除比赛 "${match.match_name}" 吗？`,
    async () => {
      try {
        await http.delete(`/matches/delete?id=${match.match_id}`)
        alert('比赛已删除')
        fetchMatches() // Need to implement fetchMatches or refresh list
      } catch (e) {
        alert('删除失败: ' + e.message)
      }
    }
  )
}

const fetchCollectorRequests = async () => {
  try {
    const res = await http.get('/admin/collectors/pending')
    if (Array.isArray(res)) {
      collectorRequests.value = res
    } else if (res && res.code === 200) {
      collectorRequests.value = Array.isArray(res.data) ? res.data : []
    } else if (res && Array.isArray(res.data)) {
      collectorRequests.value = res.data
    }
  } catch (e) {
    console.warn('获取采集员申请失败', e)
    collectorRequests.value = [] 
  }
}

const fetchCollectors = async () => {
  try {
    const res = await http.get('/admin/collectors/list')
    if (Array.isArray(res)) {
      collectors.value = res
    } else if (res && res.code === 200) {
      collectors.value = Array.isArray(res.data) ? res.data : []
    } else if (res && Array.isArray(res.data)) {
      collectors.value = res.data
    }
  } catch (e) {
    console.warn('获取采集员列表失败', e)
    collectors.value = [] 
  }
}

const deleteCollector = (studentId) => {
  openModal(
    '确认删除',
    `确定要删除采集员 ${studentId} 吗？`,
    async () => {
      try {
        await http.delete(`/admin/collectors/delete?student_id=${studentId}`)
        alert('采集员已删除')
        fetchCollectors()
      } catch (e) {
        alert('删除失败: ' + e.message)
      }
    }
  )
}

const fetchAthleteRequests = async () => {
  try {
    const res = await http.get('/admin/athletes/pending')
    if (Array.isArray(res)) {
      athleteRequests.value = res
    } else if (res && res.code === 200) {
      athleteRequests.value = Array.isArray(res.data) ? res.data : []
    } else if (res && Array.isArray(res.data)) {
      athleteRequests.value = res.data
    }
  } catch (e) {
    console.warn('获取运动员申请失败', e)
    athleteRequests.value = [] 
  }
}

const fetchTeamRequests = async () => {
  try {
    const res = await http.get('/admin/teams/pending')
    if (Array.isArray(res)) {
      teamRequests.value = res
    } else if (res && res.code === 200) {
      teamRequests.value = Array.isArray(res.data) ? res.data : []
    } else if (res && Array.isArray(res.data)) {
      teamRequests.value = res.data
    }
  } catch (e) {
    console.warn('获取队伍申请失败', e)
    teamRequests.value = [] 
  }
}

const approveCollector = (studentId) => {
  openModal(
    '确认批准',
    `确认批准 ${studentId} 的采集员申请？`,
    async () => {
      try {
        await http.post('/admin/collectors/approve', { student_id: studentId })
        alert('申请已批准')
        fetchCollectorRequests()
        fetchCollectors()
      } catch (e) {
        alert('批准失败: ' + e.message)
      }
    }
  )
}

const approveAthlete = (teamMemberId) => {
  openModal(
    '确认批准',
    '确认批准该运动员申请？',
    async () => {
      try {
        await http.post('/admin/athletes/approve', { team_member_id: teamMemberId })
        alert('申请已批准')
        fetchAthleteRequests()
      } catch (e) {
        alert('批准失败: ' + e.message)
      }
    }
  )
}

const approveTeam = (teamId) => {
  openModal(
    '确认批准',
    '确认批准该队伍创建申请？',
    async () => {
      try {
        await http.post('/admin/teams/approve', { team_id: teamId })
        alert('申请已批准')
        fetchTeamRequests()
        fetchTeams() // refresh teams list
      } catch (e) {
        alert('批准失败: ' + e.message)
      }
    }
  )
}

const rejectAthlete = (teamMemberId) => {
  openModal(
    '确认拒绝',
    '确认拒绝并删除该运动员申请？',
    async () => {
      try {
        await http.post('/teams/remove-member', { team_member_id: teamMemberId })
        alert('已拒绝并删除申请')
        fetchAthleteRequests()
      } catch (e) {
        alert('操作失败: ' + e.message)
      }
    }
  )
}

const rejectTeam = (teamId) => {
  openModal(
    '确认拒绝',
    '确认拒绝并删除该队伍创建申请？',
    async () => {
      try {
        await http.delete(`/teams/delete?id=${teamId}`)
        alert('已拒绝并删除申请')
        fetchTeamRequests()
      } catch (e) {
        alert('操作失败: ' + e.message)
      }
    }
  )
}

const fetchTeams = async () => {
  try {
    const res = await http.get('/teams/list')
    if (res && Array.isArray(res.data)) {
        teams.value = res.data
    } else if (Array.isArray(res)) {
        teams.value = res
    }
  } catch (e) {
    console.warn('获取队伍失败，将使用手动输入', e)
  }
}

const fetchEvents = async () => {
  try {
    const res = await http.get('/events/list')
     if (res && Array.isArray(res.data)) {
        events.value = res.data
    } else if (Array.isArray(res)) {
        events.value = res
    }
  } catch (e) {
    console.warn('获取赛事失败，将使用手动输入', e)
  }
}

const fetchMatches = async () => {
  try {
    const res = await http.get('/matches')
    if (res && Array.isArray(res.data)) {
        matches.value = res.data
    } else if (Array.isArray(res)) {
        matches.value = res
    }
  } catch (e) {
    console.warn('获取比赛失败', e)
    matches.value = []
  }
}

onMounted(() => {
  fetchTeams()
  fetchEvents()
  fetchMatches()
  fetchCollectorRequests()
  fetchCollectors()
  fetchAthleteRequests()
  fetchTeamRequests()
})

// Logout
const logout = () => {
  openModal(
    '确认退出',
    '确定要退出登录吗？',
    () => {
      authStore.clearAuth()
      router.push('/login')
    }
  )
}

// Actions
const createTeam = async () => {
  if (!teamForm.value.name || !teamForm.value.sport || !teamForm.value.college) {
    alert('请填写所有必填字段')
    return
  }
  loading.value = true
  try {
    const sportId = sports.find(s => s.value === teamForm.value.sport)?.value === 'football' ? 1 : 
                    sports.find(s => s.value === teamForm.value.sport)?.value === 'basketball' ? 2 : 
                    sports.find(s => s.value === teamForm.value.sport)?.value === 'badminton' ? 3 : 
                    sports.find(s => s.value === teamForm.value.sport)?.value === 'volleyball' ? 4 : 5
    
    await http.post('/teams/create', {
      team_name: teamForm.value.name,
      sport_id: sportId,
      college: teamForm.value.college,
      team_type: teamForm.value.team_type,
      created_by: authStore.studentId
    })
    alert('队伍创建成功')
    teamForm.value = { name: '', sport: '', college: '', team_type: '' }
    fetchTeams()
  } catch (e) {
    alert('创建失败: ' + (e.response?.data?.message || e.message))
  } finally {
    loading.value = false
  }
}

const createEvent = async () => {
  if (!eventForm.value.name || !eventForm.value.sport || !eventForm.value.start_date) {
    alert('请填写必填字段')
    return
  }
  loading.value = true
  try {
    const sportId = sports.find(s => s.value === eventForm.value.sport)?.value === 'football' ? 1 : 
                    sports.find(s => s.value === eventForm.value.sport)?.value === 'basketball' ? 2 : 
                    sports.find(s => s.value === eventForm.value.sport)?.value === 'badminton' ? 3 : 
                    sports.find(s => s.value === eventForm.value.sport)?.value === 'volleyball' ? 4 : 5
    
    const payload = {
      event_name: eventForm.value.name,
      sport_id: sportId,
      start_date: new Date(eventForm.value.start_date).toISOString(),
      end_date: new Date(eventForm.value.start_date).toISOString(), // Default end date to start date
      format_type: eventForm.value.format_type || 'group_knockout',
      season: '2024',
      round: 'Regular'
    }

    if (isEditingEvent.value) {
      await http.put(`/events/${currentEventId.value}`, payload)
      alert('赛事更新成功')
      isEditingEvent.value = false
      showEventEditModal.value = false
      currentEventId.value = null
    } else {
      await http.post('/events/create', payload)
      alert('赛事创建成功')
    }

    eventForm.value = { name: '', sport: '', start_date: '', format_type: '' }
    fetchEvents()
  } catch (e) {
    alert((isEditingEvent.value ? '更新' : '创建') + '失败: ' + (e.response?.data?.message || e.message))
  } finally {
    loading.value = false
  }
}

const createMatch = async () => {
  if (!matchForm.value.event_id || !matchForm.value.team_a_id || !matchForm.value.team_b_id || !matchForm.value.name || !matchForm.value.time) {
    alert('请填写所有字段')
    return
  }
  loading.value = true
  try {
    const payload = {
      event_id: Number(matchForm.value.event_id),
      team_a_id: Number(matchForm.value.team_a_id),
      team_b_id: Number(matchForm.value.team_b_id),
      match_name: matchForm.value.name,
      match_time: new Date(matchForm.value.time).toISOString(),
      round: 'Regular'
    }

    if (isEditingMatch.value) {
      await http.put(`/matches/${currentMatchId.value}`, payload)
      alert('比赛更新成功')
      isEditingMatch.value = false
      showMatchEditModal.value = false
      currentMatchId.value = null
    } else {
      await http.post('/matches/create', payload)
      alert('比赛创建成功')
    }

    matchForm.value = { event_id: '', team_a_id: '', team_b_id: '', name: '', time: '' }
    fetchMatches()
  } catch (e) {
    alert((isEditingMatch.value ? '更新' : '创建') + '失败: ' + (e.response?.data?.message || e.message))
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.admin-dashboard {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-content {
  text-align: left;
}

.logout-btn {
  background-color: transparent;
  border: 1px solid var(--danger-color, #ff4d4f);
  color: var(--danger-color, #ff4d4f);
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.logout-btn:hover {
  background-color: var(--danger-color, #ff4d4f);
  color: white;
}

.tabs {
  display: flex;
  margin-bottom: 20px;
  border-bottom: 2px solid var(--border-color);
}

.tab-btn {
  flex: 1;
  padding: 12px;
  font-weight: 600;
  color: var(--text-secondary);
  border-bottom: 3px solid transparent;
  transition: all 0.3s;
}

.tab-btn.active {
  color: var(--primary-color);
  border-bottom-color: var(--primary-color);
  background-color: var(--primary-light);
}

.form-container {
  background: white;
  padding: 25px;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0,0,0,0.05);
  animation: fadeIn 0.5s ease;
}

.collapsible-section {
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eee;
}

.list-section {
  margin-top: 20px;
}

.compact-form .form-row {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

.compact-form .form-group {
  flex: 1;
  min-width: 200px;
}

.data-table-container {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 15px;
}

.data-table th, .data-table td {
  padding: 12px 15px;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.data-table th {
  background-color: #f8f9fa;
  font-weight: 600;
  color: #2c3e50;
}

.action-btn {
  padding: 6px 12px;
  border-radius: 4px;
  font-size: 0.9em;
  cursor: pointer;
  border: none;
  transition: background 0.3s;
  margin-right: 5px;
}

.action-btn.view-btn {
  background-color: #3498db;
  color: white;
}

.action-btn.delete-btn {
  background-color: #e74c3c;
  color: white;
}

.action-btn.delete-btn:hover {
  background-color: #c0392b;
}

.action-btn.small {
  padding: 4px 8px;
  font-size: 0.8em;
}

.large-modal {
  max-width: 900px;
  width: 90%;
}

.form-group {
  margin-bottom: 20px;
}

label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #2c3e50;
}

input, select {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 1rem;
}

input:focus, select:focus {
  border-color: #3498db;
  outline: none;
}

.submit-btn {
  width: 100%;
  padding: 12px;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.3s;
}

.submit-btn:disabled {
  background-color: #bdc3c7;
}

.hint {
  color: var(--text-secondary);
  font-size: 12px;
  margin-top: 4px;
  display: block;
}

.requests-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.request-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  background: white;
}

.request-info h3 {
  margin: 0 0 4px 0;
  font-size: 16px;
}

.request-info p {
  margin: 0 0 4px 0;
  color: var(--text-secondary);
}

.request-info small {
  color: var(--text-disabled);
}

.approve-btn {
  background-color: var(--success-color, #52c41a);
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.approve-btn:hover {
  background-color: #73d13d;
}

.reject-btn {
  background-color: var(--danger-color, #ff4d4f);
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.reject-btn:hover {
  background-color: #ff7875;
}

.actions-cell {
  display: flex;
  gap: 8px;
  align-items: center;
}

.no-data {
  text-align: center;
  color: var(--text-secondary);
  padding: 20px;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
  animation: fadeIn 0.3s ease;
}

.modal-content {
  background-color: white;
  border-radius: 12px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.2);
  animation: slideUp 0.3s ease;
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
  font-weight: 600;
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
  transition: all 0.2s;
}

.modal-close:hover {
  background-color: #f5f5f5;
  color: #333;
}

.modal-body {
  padding: 24px 20px;
  font-size: 16px;
  color: #666;
  line-height: 1.5;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px;
  border-top: 1px solid #f0f0f0;
  background-color: #fafafa;
  border-radius: 0 0 12px 12px;
}

.btn-secondary,
.btn-primary {
  padding: 8px 20px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary {
  border: 1px solid #d9d9d9;
  background-color: white;
  color: #666;
}

.btn-secondary:hover {
  border-color: #40a9ff;
  color: #40a9ff;
}

.btn-primary {
  border: 1px solid var(--primary-color);
  background-color: var(--primary-color);
  color: white;
}

.btn-primary:hover {
  opacity: 0.9;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from { transform: translateY(20px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}
</style>
