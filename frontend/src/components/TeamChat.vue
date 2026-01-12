<template>
  <div class="team-chat-container">
    <div class="chat-header">
      <h3>{{ teamName }}</h3>
      <div class="header-tabs">
        <button 
          :class="['tab-btn', { active: currentTab === 'chat' }]"
          @click="currentTab = 'chat'"
        >
          聊天
        </button>
        <button 
          :class="['tab-btn', { active: currentTab === 'votes' }]"
          @click="currentTab = 'votes'"
        >
          投票
        </button>
        <button 
          :class="['tab-btn', { active: currentTab === 'notifications' }]"
          @click="currentTab = 'notifications'"
        >
          通知
        </button>
        <button 
          :class="['tab-btn', { active: currentTab === 'leave' }]"
          @click="currentTab = 'leave'"
        >
          请假
        </button>
      </div>
    </div>

    <div class="chat-content">
      <!-- 聊天标签页 -->
      <div v-if="currentTab === 'chat'" class="chat-tab">
        <div class="messages-container" ref="messagesContainer">
          <div v-if="loadingMessages" class="loading">加载中...</div>
          <div v-else-if="messages.length === 0" class="no-data">暂无消息</div>
          <div v-else class="messages-list">
            <div 
              v-for="msg in messages" 
              :key="msg.message_id"
              :class="['message-item', { 'is-read': msg.is_read }, { 'clickable': msg.message_type === 'vote' || msg.message_type === 'notification' }]"
              @click.stop="handleMessageClick(msg)"
            >
              <div class="message-header" @click.stop>
                <span class="sender-name">{{ msg.sender_name }}</span>
                <span class="message-time">{{ formatTime(msg.created_at) }}</span>
                <span v-if="msg.message_type === 'vote'" class="message-type-badge vote-badge">投票</span>
                <span v-if="msg.message_type === 'notification'" class="message-type-badge notification-badge">通知</span>
                <span v-if="msg.message_type === 'leave_request'" class="message-type-badge leave-badge">请假</span>
                <span v-if="!msg.is_read" class="unread-badge">未读</span>
                <span class="read-count">{{ msg.read_count }}/{{ msg.total_count }} 已读</span>
                <button 
                  v-if="msg.sender_id === authStore.studentId && (msg.message_type === 'text' || msg.message_type === 'vote' || msg.message_type === 'notification' || msg.message_type === 'leave_request')"
                  @click.stop="deleteMessage(msg)"
                  class="delete-message-btn"
                  title="删除消息"
                >
                  ×
                </button>
              </div>
              <div class="message-content">
                <span v-if="msg.message_type === 'vote'" class="vote-link">📊 {{ msg.content }} (点击查看详情)</span>
                <span v-else-if="msg.message_type === 'leave_request'" class="leave-link">{{ msg.content }}</span>
                <span v-else class="text-content">{{ msg.content }}</span>
              </div>
              <!-- 请假消息的审核按钮（仅队长可见） -->
              <div v-if="msg.message_type === 'leave_request' && isCaptain" class="message-actions" @click.stop>
                <button 
                  @click="handleLeaveReviewFromMessage(msg)"
                  class="review-btn"
                >
                  审核请假
                </button>
              </div>
            </div>
          </div>
        </div>
        <div class="message-input-container">
          <input 
            v-model="newMessage"
            @keyup.enter="sendMessage"
            placeholder="输入消息..."
            class="message-input"
          />
          <button @click="sendMessage" class="send-btn">发送</button>
        </div>
      </div>

      <!-- 投票标签页 -->
      <div v-if="currentTab === 'votes'" class="votes-tab">
        <div class="tab-header">
          <button @click="showCreateVoteModal = true" class="create-btn">创建投票</button>
        </div>
        <div v-if="loadingVotes" class="loading">加载中...</div>
        <div v-else-if="votes.length === 0" class="no-data">暂无投票</div>
        <div v-else class="votes-list">
          <div v-for="vote in votes" :key="vote.vote_id" class="vote-item">
            <div class="vote-header">
              <h4>{{ vote.title }}</h4>
              <div class="vote-header-right">
                <span :class="['vote-status', vote.status]">{{ vote.status === 'active' ? '进行中' : '已结束' }}</span>
                <button 
                  v-if="vote.creator_id === authStore.studentId"
                  @click="deleteVote(vote)"
                  class="delete-vote-btn"
                  title="删除投票"
                >
                  删除
                </button>
              </div>
            </div>
            <div v-if="vote.description" class="vote-description">{{ vote.description }}</div>
            <div class="vote-options">
              <label 
                v-for="(option, idx) in vote.options" 
                :key="idx"
                :class="['option-item', { 'selected': vote.user_vote && vote.user_vote.includes(idx) }]"
              >
                <input 
                  type="checkbox" 
                  v-if="vote.is_multiple"
                  :checked="vote.user_vote && vote.user_vote.includes(idx)"
                  @change="handleVoteOptionChange(vote, idx, $event)"
                  :disabled="vote.status === 'closed' || (vote.deadline && new Date(vote.deadline) < new Date())"
                />
                <input 
                  type="radio" 
                  v-else
                  :name="'vote_' + vote.vote_id"
                  :checked="vote.user_vote && vote.user_vote.includes(idx)"
                  @change="handleVoteOptionChange(vote, idx, $event)"
                  :disabled="vote.status === 'closed' || (vote.deadline && new Date(vote.deadline) < new Date())"
                />
                <span class="option-text">{{ option }}</span>
                <span v-if="vote.results && vote.results.length > 0" class="option-result">
                  {{ vote.results[idx].vote_count }}票 ({{ vote.results[idx].percentage.toFixed(1) }}%)
                </span>
              </label>
            </div>
            <div class="vote-footer">
              <span>共 {{ vote.vote_count }} 人投票</span>
              <span v-if="vote.deadline">截止时间: {{ formatTime(vote.deadline) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 通知标签页 -->
      <div v-if="currentTab === 'notifications'" class="notifications-tab">
        <div class="tab-header" v-if="isCaptain">
          <button @click="showCreateNotificationModal = true" class="create-btn">发布通知</button>
        </div>
        <div v-if="loadingNotifications" class="loading">加载中...</div>
        <div v-else-if="notifications.length === 0" class="no-data">暂无通知</div>
        <div v-else class="notifications-list">
          <div 
            v-for="notif in notifications" 
            :key="notif.notification_id"
            :class="['notification-item', notif.notification_type, { 'is-read': notif.is_read }]"
          >
            <div class="notification-header">
              <h4>{{ notif.title }}</h4>
              <div class="notification-header-right">
                <span :class="['notification-type', notif.notification_type]">
                  {{ getNotificationTypeText(notif.notification_type) }}
                </span>
                <button 
                  v-if="notif.sender_id === authStore.studentId"
                  @click="deleteNotification(notif)"
                  class="delete-notification-btn"
                  title="删除通知"
                >
                  删除
                </button>
              </div>
            </div>
            <div class="notification-content">{{ notif.content }}</div>
            <div class="notification-footer">
              <span>{{ notif.sender_name }}</span>
              <span>{{ formatTime(notif.created_at) }}</span>
              <span>{{ notif.read_count }}/{{ notif.total_count }} 已读</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 请假标签页 -->
      <div v-if="currentTab === 'leave'" class="leave-tab">
        <div class="tab-header">
          <button @click="showCreateLeaveModal = true" class="create-btn">申请请假</button>
        </div>
        <div v-if="loadingLeaveRequests" class="loading">加载中...</div>
        <div v-else-if="leaveRequests.length === 0" class="no-data">暂无请假申请</div>
        <div v-else class="leave-requests-list">
          <div 
            v-for="req in leaveRequests" 
            :key="req.leave_id"
            :class="['leave-request-item', req.status]"
          >
            <div class="leave-header">
              <div>
                <span class="applicant-name">{{ req.applicant_name }}</span>
                <span :class="['leave-status', req.status]">{{ getLeaveStatusText(req.status) }}</span>
              </div>
              <span class="leave-type">{{ getLeaveTypeText(req.leave_type) }}</span>
            </div>
            <div class="leave-dates">
              {{ formatDate(req.start_date) }} 至 {{ formatDate(req.end_date) }}
            </div>
            <div class="leave-reason">{{ req.reason }}</div>
            <div v-if="req.reviewer_name" class="leave-review">
              <span>审核人: {{ req.reviewer_name }}</span>
              <span v-if="req.review_comment"> - {{ req.review_comment }}</span>
            </div>
            <div v-if="isCaptain && req.status === 'pending'" class="leave-actions">
              <button @click="reviewLeave(req, 'approved')" class="approve-btn">批准</button>
              <button @click="reviewLeave(req, 'rejected')" class="reject-btn">拒绝</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建投票弹窗 -->
    <div v-if="showCreateVoteModal" class="modal-overlay" @click.self="showCreateVoteModal = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3>创建投票</h3>
          <button @click="showCreateVoteModal = false" class="modal-close">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>标题 *</label>
            <input v-model="voteForm.title" class="form-input" />
          </div>
          <div class="form-group">
            <label>描述</label>
            <textarea v-model="voteForm.description" class="form-textarea"></textarea>
          </div>
          <div class="form-group">
            <label>选项 *</label>
            <div v-for="(option, idx) in voteForm.options" :key="idx" class="option-input-row">
              <input v-model="voteForm.options[idx]" class="form-input" :placeholder="`选项 ${idx + 1}`" />
              <button v-if="voteForm.options.length > 2" @click="removeOption(idx)" class="remove-btn">删除</button>
            </div>
            <button @click="addOption" class="add-option-btn">添加选项</button>
          </div>
          <div class="form-group">
            <label>
              <input type="checkbox" v-model="voteForm.is_multiple" />
              允许多选
            </label>
          </div>
          <div class="form-group">
            <label>截止时间（可选）</label>
            <input type="datetime-local" v-model="voteForm.deadline" class="form-input" />
          </div>
        </div>
        <div class="modal-footer">
          <button @click="showCreateVoteModal = false" class="modal-btn cancel">取消</button>
          <button @click="createVote" class="modal-btn primary" :disabled="submitting">创建</button>
        </div>
      </div>
    </div>

    <!-- 创建通知弹窗 -->
    <div v-if="showCreateNotificationModal" class="modal-overlay" @click.self="showCreateNotificationModal = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3>发布通知</h3>
          <button @click="showCreateNotificationModal = false" class="modal-close">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>标题 *</label>
            <input v-model="notificationForm.title" class="form-input" />
          </div>
          <div class="form-group">
            <label>内容 *</label>
            <textarea v-model="notificationForm.content" class="form-textarea"></textarea>
          </div>
          <div class="form-group">
            <label>类型</label>
            <select v-model="notificationForm.notification_type" class="form-input">
              <option value="info">普通</option>
              <option value="warning">警告</option>
              <option value="urgent">紧急</option>
            </select>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="showCreateNotificationModal = false" class="modal-btn cancel">取消</button>
          <button @click="createNotification" class="modal-btn primary" :disabled="submitting">发布</button>
        </div>
      </div>
    </div>

    <!-- 投票详情弹窗 -->
    <div v-if="showVoteDetailModal && currentVoteDetail" class="modal-overlay" @click.self="showVoteDetailModal = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3>投票详情</h3>
          <button @click="showVoteDetailModal = false" class="modal-close">×</button>
        </div>
        <div class="modal-body">
          <div class="vote-detail">
            <h4>{{ currentVoteDetail.title }}</h4>
            <div v-if="currentVoteDetail.description" class="vote-description">{{ currentVoteDetail.description }}</div>
            <div class="vote-options">
              <label 
                v-for="(option, idx) in currentVoteDetail.options" 
                :key="idx"
                :class="['option-item', { 'selected': currentVoteDetail.user_vote && currentVoteDetail.user_vote.includes(idx) }]"
              >
                <input 
                  type="checkbox" 
                  v-if="currentVoteDetail.is_multiple"
                  :checked="currentVoteDetail.user_vote && currentVoteDetail.user_vote.includes(idx)"
                  @change="handleVoteOptionChange(currentVoteDetail, idx, $event)"
                  :disabled="currentVoteDetail.status === 'closed' || (currentVoteDetail.deadline && new Date(currentVoteDetail.deadline) < new Date())"
                />
                <input 
                  type="radio" 
                  v-else
                  :name="'vote_detail_' + currentVoteDetail.vote_id"
                  :checked="currentVoteDetail.user_vote && currentVoteDetail.user_vote.includes(idx)"
                  @change="handleVoteOptionChange(currentVoteDetail, idx, $event)"
                  :disabled="currentVoteDetail.status === 'closed' || (currentVoteDetail.deadline && new Date(currentVoteDetail.deadline) < new Date())"
                />
                <span class="option-text">{{ option }}</span>
                <span v-if="currentVoteDetail.results && currentVoteDetail.results.length > 0" class="option-result">
                  {{ currentVoteDetail.results[idx].vote_count }}票 ({{ currentVoteDetail.results[idx].percentage.toFixed(1) }}%)
                </span>
              </label>
            </div>
            <div class="vote-footer">
              <span>共 {{ currentVoteDetail.vote_count }} 人投票</span>
              <span v-if="currentVoteDetail.deadline">截止时间: {{ formatTime(currentVoteDetail.deadline) }}</span>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="showVoteDetailModal = false" class="modal-btn primary">关闭</button>
        </div>
      </div>
    </div>

    <!-- 创建请假申请弹窗 -->
    <div v-if="showCreateLeaveModal" class="modal-overlay" @click.self="showCreateLeaveModal = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3>申请请假</h3>
          <button @click="showCreateLeaveModal = false" class="modal-close">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>请假类型 *</label>
            <select v-model="leaveForm.leave_type" class="form-input">
              <option value="sick">病假</option>
              <option value="personal">事假</option>
              <option value="other">其他</option>
            </select>
          </div>
          <div class="form-group">
            <label>开始日期 *</label>
            <input type="date" v-model="leaveForm.start_date" class="form-input" />
          </div>
          <div class="form-group">
            <label>结束日期 *</label>
            <input type="date" v-model="leaveForm.end_date" class="form-input" />
          </div>
          <div class="form-group">
            <label>请假原因 *</label>
            <textarea v-model="leaveForm.reason" class="form-textarea"></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="showCreateLeaveModal = false" class="modal-btn cancel">取消</button>
          <button @click="createLeaveRequest" class="modal-btn primary" :disabled="submitting">提交</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch, nextTick } from 'vue'
import { get, post, del } from '@/utils/http'
import { useAuthStore } from '@/stores/auth'

const props = defineProps({
  teamId: {
    type: Number,
    required: true
  },
  teamName: {
    type: String,
    required: true
  },
  isCaptain: {
    type: Boolean,
    default: false
  }
})

const authStore = useAuthStore()

const currentTab = ref('chat')
const loadingMessages = ref(false)
const loadingVotes = ref(false)
const loadingNotifications = ref(false)
const loadingLeaveRequests = ref(false)
const submitting = ref(false)

const messages = ref([])
const votes = ref([])
const notifications = ref([])
const leaveRequests = ref([])

const newMessage = ref('')
const messagesContainer = ref(null)

const showCreateVoteModal = ref(false)
const showCreateNotificationModal = ref(false)
const showCreateLeaveModal = ref(false)
const showVoteDetailModal = ref(false)
const currentVoteDetail = ref(null)
const currentLeaveRequest = ref(null)

const voteForm = reactive({
  title: '',
  description: '',
  options: ['', ''],
  is_multiple: false,
  deadline: ''
})

const notificationForm = reactive({
  title: '',
  content: '',
  notification_type: 'info'
})

const leaveForm = reactive({
  leave_type: 'personal',
  start_date: '',
  end_date: '',
  reason: ''
})

// 获取消息列表
const fetchMessages = async () => {
  loadingMessages.value = true
  try {
    const res = await get(`/team/chat/messages/list?team_id=${props.teamId}&limit=50`)
    if (res.code === 200) {
      messages.value = (res.data || []).reverse() // 反转以显示最新消息在底部
      console.log('Messages loaded:', messages.value)
      // 检查投票消息
      const voteMessages = messages.value.filter(m => m.message_type === 'vote')
      console.log('Vote messages:', voteMessages)
      nextTick(() => {
        scrollToBottom()
      })
    }
  } catch (err) {
    console.error('Fetch messages error:', err)
  } finally {
    loadingMessages.value = false
  }
}

// 发送消息
const sendMessage = async () => {
  if (!newMessage.value.trim()) return
  
  submitting.value = true
  try {
    const res = await post('/team/chat/messages', {
      team_id: props.teamId,
      message_type: 'text',
      content: newMessage.value
    })
    if (res.code === 200) {
      newMessage.value = ''
      await fetchMessages()
    }
  } catch (err) {
    console.error('Send message error:', err)
    alert('发送失败: ' + (err.response?.data?.message || err.message))
  } finally {
    submitting.value = false
  }
}

// 处理消息点击
const handleMessageClick = async (msg) => {
  console.log('Message clicked:', msg)
  console.log('Message type:', msg.message_type)
  
  // 如果是请假消息，获取请假申请详情
  if (msg.message_type === 'leave_request') {
    try {
      const res = await get(`/team/chat/leave-requests/by-message?message_id=${msg.message_id}`)
      if (res.code === 200 && res.data) {
        currentLeaveRequest.value = res.data
      }
    } catch (err) {
      console.error('Get leave request error:', err)
    }
  }

  // 如果是投票消息，显示投票详情
  if (msg.message_type === 'vote') {
    console.log('Opening vote detail for message_id:', msg.message_id)
    try {
      const res = await get(`/team/chat/votes/by-message?message_id=${msg.message_id}`)
      console.log('Vote detail response:', res)
      if (res.code === 200 && res.data) {
        currentVoteDetail.value = res.data
        showVoteDetailModal.value = true
        console.log('Vote detail modal opened, vote:', res.data)
      } else {
        console.error('Invalid response:', res)
        // 如果通过message_id找不到，尝试从投票列表中查找
        await fetchVotes()
        const voteFromList = votes.value.find(v => v.message_id === msg.message_id)
        if (voteFromList) {
          currentVoteDetail.value = voteFromList
          showVoteDetailModal.value = true
          console.log('Found vote from list:', voteFromList)
        } else {
          alert('获取投票详情失败: ' + (res.message || '投票不存在'))
        }
      }
    } catch (err) {
      console.error('Get vote detail error:', err)
      console.error('Error details:', err.response)
      // 如果API调用失败，尝试从投票列表中查找
      try {
        await fetchVotes()
        const voteFromList = votes.value.find(v => v.message_id === msg.message_id)
        if (voteFromList) {
          currentVoteDetail.value = voteFromList
          showVoteDetailModal.value = true
          console.log('Found vote from list after error:', voteFromList)
        } else {
          alert('获取投票详情失败: ' + (err.response?.data?.message || err.message || '网络错误'))
        }
      } catch {
        alert('获取投票详情失败: ' + (err.response?.data?.message || err.message || '网络错误'))
      }
    }
  }

  // 标记为已读（无论是否投票消息都标记）
  if (!msg.is_read) {
    try {
      await post('/team/chat/messages/read', {
        message_id: msg.message_id
      })
      // 更新本地状态
      msg.is_read = true
      msg.read_count = (msg.read_count || 0) + 1
    } catch (err) {
      console.error('Mark message as read error:', err)
    }
  }
}

// 获取投票列表
const fetchVotes = async () => {
  loadingVotes.value = true
  try {
    const res = await get(`/team/chat/votes/list?team_id=${props.teamId}`)
    if (res.code === 200) {
      votes.value = res.data || []
    }
  } catch (err) {
    console.error('Fetch votes error:', err)
  } finally {
    loadingVotes.value = false
  }
}

// 处理投票选项变化
const handleVoteOptionChange = async (vote, optionIdx, event) => {
  if (vote.status === 'closed') return
  if (vote.deadline && new Date(vote.deadline) < new Date()) return

  let selectedOptions = vote.user_vote ? [...vote.user_vote] : []
  
  if (vote.is_multiple) {
    if (event.target.checked) {
      if (!selectedOptions.includes(optionIdx)) {
        selectedOptions.push(optionIdx)
      }
    } else {
      selectedOptions = selectedOptions.filter(idx => idx !== optionIdx)
    }
  } else {
    selectedOptions = [optionIdx]
  }

  try {
    const res = await post('/team/chat/votes/vote', {
      vote_id: vote.vote_id,
      selected_options: selectedOptions
    })
    if (res.code === 200) {
      await fetchVotes()
      // 如果是在详情弹窗中，更新详情数据
      if (showVoteDetailModal.value && currentVoteDetail.value && currentVoteDetail.value.vote_id === vote.vote_id) {
        const detailRes = await get(`/team/chat/votes/by-message?message_id=${vote.message_id}`)
        if (detailRes.code === 200) {
          currentVoteDetail.value = detailRes.data
        }
      }
    }
  } catch (err) {
    console.error('Vote error:', err)
    alert('投票失败: ' + (err.response?.data?.message || err.message))
  }
}

// 创建投票
const createVote = async () => {
  if (!voteForm.title.trim()) {
    alert('请输入标题')
    return
  }
  if (voteForm.options.filter(opt => opt.trim()).length < 2) {
    alert('至少需要2个选项')
    return
  }

  submitting.value = true
  try {
    // 处理截止时间格式
    let deadlineValue = null
    if (voteForm.deadline) {
      // datetime-local 格式转换为 ISO 8601 格式
      const date = new Date(voteForm.deadline)
      if (!isNaN(date.getTime())) {
        deadlineValue = date.toISOString()
      }
    }
    
    const res = await post('/team/chat/votes', {
      team_id: props.teamId,
      title: voteForm.title,
      description: voteForm.description,
      options: voteForm.options.filter(opt => opt.trim()),
      is_multiple: voteForm.is_multiple,
      deadline: deadlineValue
    })
    if (res.code === 200) {
      alert('投票发布成功！')
      showCreateVoteModal.value = false
      Object.assign(voteForm, {
        title: '',
        description: '',
        options: ['', ''],
        is_multiple: false,
        deadline: ''
      })
      await fetchVotes()
      // 刷新消息列表以显示新创建的投票消息
      await fetchMessages()
    } else {
      alert('创建失败: ' + (res.message || '未知错误'))
    }
  } catch (err) {
    console.error('Create vote error:', err)
    alert('创建失败: ' + (err.response?.data?.message || err.message))
  } finally {
    submitting.value = false
  }
}

const addOption = () => {
  voteForm.options.push('')
}

const removeOption = (idx) => {
  voteForm.options.splice(idx, 1)
}

// 获取通知列表
const fetchNotifications = async () => {
  loadingNotifications.value = true
  try {
    const res = await get(`/team/chat/notifications/list?team_id=${props.teamId}&limit=20`)
    if (res.code === 200) {
      notifications.value = res.data || []
    }
  } catch (err) {
    console.error('Fetch notifications error:', err)
  } finally {
    loadingNotifications.value = false
  }
}

// 创建通知
const createNotification = async () => {
  if (!notificationForm.title.trim() || !notificationForm.content.trim()) {
    alert('请填写标题和内容')
    return
  }

  submitting.value = true
  try {
    const res = await post('/team/chat/notifications', {
      team_id: props.teamId,
      title: notificationForm.title,
      content: notificationForm.content,
      notification_type: notificationForm.notification_type
    })
    if (res.code === 200) {
      showCreateNotificationModal.value = false
      Object.assign(notificationForm, {
        title: '',
        content: '',
        notification_type: 'info'
      })
      await fetchNotifications()
    }
  } catch (err) {
    console.error('Create notification error:', err)
    alert('发布失败: ' + (err.response?.data?.message || err.message))
  } finally {
    submitting.value = false
  }
}

// 获取请假申请列表
const fetchLeaveRequests = async () => {
  loadingLeaveRequests.value = true
  try {
    const res = await get(`/team/chat/leave-requests/list?team_id=${props.teamId}`)
    if (res.code === 200) {
      leaveRequests.value = res.data || []
    }
  } catch (err) {
    console.error('Fetch leave requests error:', err)
  } finally {
    loadingLeaveRequests.value = false
  }
}

// 创建请假申请
const createLeaveRequest = async () => {
  if (!leaveForm.start_date || !leaveForm.end_date || !leaveForm.reason.trim()) {
    alert('请填写完整信息')
    return
  }

  submitting.value = true
  try {
    const res = await post('/team/chat/leave-requests', {
      team_id: props.teamId,
      leave_type: leaveForm.leave_type,
      start_date: leaveForm.start_date,
      end_date: leaveForm.end_date,
      reason: leaveForm.reason
    })
    if (res.code === 200) {
      alert('请假申请提交成功！')
      showCreateLeaveModal.value = false
      Object.assign(leaveForm, {
        leave_type: 'personal',
        start_date: '',
        end_date: '',
        reason: ''
      })
      await fetchLeaveRequests()
      // 刷新消息列表以显示新创建的请假申请消息
      await fetchMessages()
    } else {
      alert('申请失败: ' + (res.message || '未知错误'))
    }
  } catch (err) {
    console.error('Create leave request error:', err)
    alert('申请失败: ' + (err.response?.data?.message || err.message))
  } finally {
    submitting.value = false
  }
}

// 审核请假申请
const reviewLeave = async (req, status, comment = null) => {
  if (comment === null) {
    comment = prompt(status === 'approved' ? '请输入批准意见（可选）' : '请输入拒绝原因（可选）')
    if (comment === null) return // 用户取消
  }

  try {
    const res = await post('/team/chat/leave-requests/review', {
      leave_id: req.leave_id,
      status: status,
      review_comment: comment || ''
    })
    if (res.code === 200) {
      alert(status === 'approved' ? '已批准请假申请' : '已拒绝请假申请')
      await fetchLeaveRequests()
      // 刷新消息列表以显示审核结果
      await fetchMessages()
      currentLeaveRequest.value = null
    }
  } catch (err) {
    console.error('Review leave request error:', err)
    alert('审核失败: ' + (err.response?.data?.message || err.message))
  }
}

// 删除消息
const deleteMessage = async (msg) => {
  if (!confirm('确定要删除这条消息吗？删除后无法恢复。')) {
    return
  }

  try {
    const res = await del(`/team/chat/messages/delete?message_id=${msg.message_id}`)
    if (res.code === 200) {
      // 从列表中移除消息
      messages.value = messages.value.filter(m => m.message_id !== msg.message_id)
      // 如果是投票、通知或请假，也从对应列表中移除
      if (msg.message_type === 'vote') {
        votes.value = votes.value.filter(v => v.message_id !== msg.message_id)
      } else if (msg.message_type === 'notification') {
        notifications.value = notifications.value.filter(n => n.message_id !== msg.message_id)
      } else if (msg.message_type === 'leave_request') {
        leaveRequests.value = leaveRequests.value.filter(l => l.message_id !== msg.message_id)
      }
      alert('消息已删除')
    } else {
      alert('删除失败: ' + (res.message || '未知错误'))
    }
  } catch (err) {
    console.error('Delete message error:', err)
    const errorMsg = err.message || '未知错误'
    // 尝试从错误消息中提取有用的信息
    if (errorMsg.includes('只能删除')) {
      alert('删除失败: ' + errorMsg)
    } else if (errorMsg.includes('JSON parse error')) {
      alert('删除失败: 服务器响应格式错误，请检查网络连接')
    } else {
      alert('删除失败: ' + errorMsg)
    }
  }
}

// 删除投票
const deleteVote = async (vote) => {
  if (!confirm('确定要删除这个投票吗？删除后无法恢复。')) {
    return
  }

  try {
    const res = await del(`/team/chat/messages/delete?message_id=${vote.message_id}`)
    if (res.code === 200) {
      votes.value = votes.value.filter(v => v.vote_id !== vote.vote_id)
      messages.value = messages.value.filter(m => m.message_id !== vote.message_id)
      alert('投票已删除')
    } else {
      alert('删除失败: ' + (res.message || '未知错误'))
    }
  } catch (err) {
    console.error('Delete vote error:', err)
    alert('删除失败: ' + (err.response?.data?.message || err.message))
  }
}

// 删除通知
const deleteNotification = async (notif) => {
  if (!confirm('确定要删除这个通知吗？删除后无法恢复。')) {
    return
  }

  try {
    const res = await del(`/team/chat/messages/delete?message_id=${notif.message_id}`)
    if (res.code === 200) {
      notifications.value = notifications.value.filter(n => n.notification_id !== notif.notification_id)
      messages.value = messages.value.filter(m => m.message_id !== notif.message_id)
      alert('通知已删除')
    } else {
      alert('删除失败: ' + (res.message || '未知错误'))
    }
  } catch (err) {
    console.error('Delete notification error:', err)
    alert('删除失败: ' + (err.response?.data?.message || err.message))
  }
}

// 从消息中处理请假审核
const handleLeaveReviewFromMessage = async (msg) => {
  try {
    // 确保已加载请假申请列表
    if (!leaveRequests.value.length) {
      await fetchLeaveRequests()
    }

    // 1) 优先根据 message_id 精确匹配（本地列表）
    let leaveReq = leaveRequests.value.find(
      l => l.message_id !== undefined && l.message_id !== null && Number(l.message_id) === Number(msg.message_id)
    )

    // 2) 如果没有找到，尝试根据申请人姓名 + 待审核状态进行模糊匹配（兼容旧数据未写入 message_id 的情况）
    if (!leaveReq) {
      const pendingForSender = leaveRequests.value.filter(
        l => l.status === 'pending' && l.applicant_name === msg.sender_name
      )
      if (pendingForSender.length === 1) {
        leaveReq = pendingForSender[0]
      }
    }

    // 3) 仍然没有找到时，调用后端接口通过 message_id 精确查询一遍
    if (!leaveReq) {
      try {
        const res = await get(`/team/chat/leave-requests/by-message?message_id=${msg.message_id}`)
        if (res.code === 200 && res.data) {
          leaveReq = res.data
        }
      } catch (e) {
        console.error('Get leave request by message_id error:', e)
      }
    }

    if (!leaveReq) {
      alert('未找到对应的请假申请记录，请到“请假”页中查看/审核')
      currentTab.value = 'leave'
      return
    }

    if (leaveReq.status !== 'pending') {
      alert('该请假申请已处理')
      return
    }

    // 显示审核弹窗
    const action = confirm('批准请假申请？\n点击"确定"批准，点击"取消"拒绝')
    const comment = prompt(action ? '请输入批准意见（可选）' : '请输入拒绝原因（可选）')
    if (comment === null && !action) return // 用户取消拒绝操作

    const status = action ? 'approved' : 'rejected'
    await reviewLeave(leaveReq, status, comment || '')
  } catch (err) {
    console.error('Get leave request error:', err)
    const errorMsg = err.message || '未知错误'
    alert('获取请假申请失败: ' + errorMsg)
  }
}

// 工具函数
const formatTime = (timeStr) => {
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

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}

const getNotificationTypeText = (type) => {
  const map = {
    info: '普通',
    warning: '警告',
    urgent: '紧急'
  }
  return map[type] || type
}

const getLeaveStatusText = (status) => {
  const map = {
    pending: '待审核',
    approved: '已批准',
    rejected: '已拒绝'
  }
  return map[status] || status
}

const getLeaveTypeText = (type) => {
  const map = {
    sick: '病假',
    personal: '事假',
    other: '其他'
  }
  return map[type] || type
}

const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

// 监听标签页切换
watch(currentTab, (newTab) => {
  if (newTab === 'chat') {
    fetchMessages()
  } else if (newTab === 'votes') {
    fetchVotes()
  } else if (newTab === 'notifications') {
    fetchNotifications()
  } else if (newTab === 'leave') {
    fetchLeaveRequests()
  }
})

onMounted(() => {
  fetchMessages()
})
</script>

<style scoped>
.team-chat-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: white;
  border-radius: 12px;
  overflow: hidden;
}

.chat-header {
  padding: 16px 20px;
  border-bottom: 1px solid #eee;
  background: #fafafa;
}

.chat-header h3 {
  margin: 0 0 12px 0;
  font-size: 18px;
  color: #333;
}

.header-tabs {
  display: flex;
  gap: 8px;
}

.tab-btn {
  padding: 6px 16px;
  border: none;
  background: #f0f0f0;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.tab-btn:hover {
  background: #e0e0e0;
}

.tab-btn.active {
  background: #1890ff;
  color: white;
}

.chat-content {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* 聊天标签页 */
.chat-tab {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.messages-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.message-item {
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
  border-left: 3px solid #1890ff;
}

.message-item.is-read {
  opacity: 0.7;
}

.message-item.clickable {
  cursor: pointer;
  transition: background-color 0.2s;
  user-select: none;
}

.message-item.clickable:hover {
  background: #f0f0f0;
}

.message-item.clickable:active {
  background: #e8e8e8;
}

.message-type-badge {
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 11px;
  font-weight: 500;
}

.vote-badge {
  background: #1890ff;
  color: white;
}

.notification-badge {
  background: #faad14;
  color: white;
}

.leave-badge {
  background: #722ed1;
  color: white;
}

.vote-link {
  color: #1890ff;
  text-decoration: underline;
  cursor: pointer;
}

.leave-link {
  color: #722ed1;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
  font-size: 12px;
  color: #666;
}

.sender-name {
  font-weight: 600;
  color: #333;
}

.message-time {
  color: #999;
}

.unread-badge {
  background: #ff4d4f;
  color: white;
  padding: 2px 6px;
  border-radius: 10px;
  font-size: 10px;
}

.read-count {
  margin-left: auto;
  color: #999;
}

.delete-message-btn {
  margin-left: 8px;
  padding: 2px 8px;
  background: #ff4d4f;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 18px;
  line-height: 1;
  opacity: 0.7;
  transition: opacity 0.2s;
}

.delete-message-btn:hover {
  opacity: 1;
}

.delete-vote-btn, .delete-notification-btn {
  padding: 4px 12px;
  background: #ff4d4f;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  transition: background 0.2s;
}

.delete-vote-btn:hover, .delete-notification-btn:hover {
  background: #ff7875;
}

.message-actions {
  margin-top: 8px;
  display: flex;
  gap: 8px;
}

.review-btn {
  padding: 6px 12px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  transition: background 0.2s;
}

.review-btn:hover:not(:disabled) {
  background: #40a9ff;
}

.review-btn:disabled {
  background: #d9d9d9;
  cursor: not-allowed;
}

.message-content {
  color: #333;
  line-height: 1.5;
}

.message-content .text-content {
  white-space: pre-wrap;
  word-break: break-word;
}

.message-input-container {
  display: flex;
  gap: 8px;
  padding: 16px;
  border-top: 1px solid #eee;
}

.message-input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
}

.send-btn {
  padding: 10px 20px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

/* 投票标签页 */
.votes-tab, .notifications-tab, .leave-tab {
  padding: 16px;
  overflow-y: auto;
}

.tab-header {
  margin-bottom: 16px;
}

.create-btn {
  padding: 8px 16px;
  background: #52c41a;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

.votes-list, .notifications-list, .leave-requests-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.vote-item, .notification-item, .leave-request-item {
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.vote-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.vote-header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.vote-header h4 {
  margin: 0;
  font-size: 16px;
}

.vote-status {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.vote-status.active {
  background: #52c41a;
  color: white;
}

.vote-status.closed {
  background: #999;
  color: white;
}

.vote-description {
  color: #666;
  margin-bottom: 12px;
}

.vote-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
}

.option-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  background: white;
  border-radius: 4px;
  cursor: pointer;
}

.option-item.selected {
  background: #e6f7ff;
  border: 1px solid #1890ff;
}

.option-text {
  flex: 1;
}

.option-result {
  color: #666;
  font-size: 12px;
}

.vote-footer {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #999;
}

/* 通知标签页 */
.notification-item.info {
  border-left: 3px solid #1890ff;
}

.notification-item.warning {
  border-left: 3px solid #faad14;
}

.notification-item.urgent {
  border-left: 3px solid #ff4d4f;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.notification-header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.notification-header h4 {
  margin: 0;
  font-size: 16px;
}

.notification-type {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.notification-type.info {
  background: #e6f7ff;
  color: #1890ff;
}

.notification-type.warning {
  background: #fff7e6;
  color: #faad14;
}

.notification-type.urgent {
  background: #fff1f0;
  color: #ff4d4f;
}

.notification-content {
  color: #333;
  margin-bottom: 8px;
  line-height: 1.5;
}

.notification-footer {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #999;
}

/* 请假标签页 */
.leave-request-item {
  border-left: 3px solid #1890ff;
}

.leave-request-item.approved {
  border-left-color: #52c41a;
}

.leave-request-item.rejected {
  border-left-color: #ff4d4f;
}

.leave-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.applicant-name {
  font-weight: 600;
  margin-right: 8px;
}

.leave-status {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.leave-status.pending {
  background: #fff7e6;
  color: #faad14;
}

.leave-status.approved {
  background: #f6ffed;
  color: #52c41a;
}

.leave-status.rejected {
  background: #fff1f0;
  color: #ff4d4f;
}

.leave-type {
  font-size: 12px;
  color: #666;
}

.leave-dates {
  color: #666;
  margin-bottom: 8px;
}

.leave-reason {
  color: #333;
  margin-bottom: 8px;
  line-height: 1.5;
}

.leave-review {
  font-size: 12px;
  color: #999;
  margin-bottom: 8px;
}

.leave-actions {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}

.approve-btn, .reject-btn {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.approve-btn {
  background: #52c41a;
  color: white;
}

.reject-btn {
  background: #ff4d4f;
  color: white;
}

/* 弹窗样式 */
.modal-overlay {
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

.modal-content {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.modal-header {
  padding: 16px 20px;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
}

.modal-close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #999;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 14px;
  color: #333;
}

.form-input, .form-textarea {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.form-textarea {
  min-height: 80px;
  resize: vertical;
}

.option-input-row {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.option-input-row .form-input {
  flex: 1;
}

.remove-btn {
  padding: 8px 12px;
  background: #ff4d4f;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.add-option-btn {
  padding: 6px 12px;
  background: #f0f0f0;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.modal-footer {
  padding: 16px 20px;
  border-top: 1px solid #eee;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.modal-btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.modal-btn.cancel {
  background: #f0f0f0;
  color: #333;
}

.modal-btn.primary {
  background: #1890ff;
  color: white;
}

.modal-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.loading, .no-data {
  text-align: center;
  padding: 40px;
  color: #999;
}
</style>
