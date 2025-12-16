<template>
  <div v-if="visible" class="modal-overlay" @click.self="close">
    <div class="modal-content large-modal">
      <div class="modal-header">
        <h3 v-if="!isEditingName">队伍管理: {{ teamName }}</h3>
        <div v-else class="edit-name-container">
          <input 
            v-model="editNameValue" 
            class="edit-name-input"
            placeholder="输入新队伍名称"
          />
          <button class="action-btn save-btn small" @click="saveName">保存</button>
          <button class="action-btn cancel-btn small" @click="cancelEditName">取消</button>
        </div>
        
        <div class="header-actions">
           <button v-if="!isEditingName && canManage" class="action-btn edit-btn small" @click="startEditName">修改队名</button>
           <button @click="close" class="modal-close">×</button>
        </div>
      </div>

      <div class="modal-body">
        <div v-if="loading" class="loading">加载中...</div>
        
        <!-- 待审核成员列表 -->
        <div v-if="canManage && pendingMembers.length > 0" class="pending-section">
          <h4 class="section-title">待审核申请</h4>
          <div class="data-table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>学号</th>
                  <th>姓名</th>
                  <th>学院</th>
                  <th>项目</th>
                  <th>号码</th>
                  <th>申请时间</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="member in pendingMembers" :key="member.team_member_id">
                  <td>{{ member.student_id }}</td>
                  <td>{{ member.name }}</td>
                  <td>{{ member.college }}</td>
                  <td>{{ member.sport_type }}</td>
                  <td>{{ member.jersey_number || '-' }}</td>
                  <td>{{ new Date(member.join_date).toLocaleDateString() }}</td>
                  <td class="actions-cell">
                    <button class="action-btn approve-btn small" @click="approveMember(member)">同意</button>
                    <button class="action-btn reject-btn small" @click="rejectMember(member)">拒绝</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="divider"></div>
        </div>

        <div v-if="members.length === 0 && (!canManage || pendingMembers.length === 0)" class="no-data">暂无成员</div>
        <div v-else-if="members.length > 0" class="data-table-container">
          <h4 v-if="canManage && pendingMembers.length > 0" class="section-title">正式成员</h4>
          <table class="data-table">
            <thead>
              <tr>
                <th>学号</th>
                <th>姓名</th>
                <th>学院</th>
                <th>项目</th>
                <th>号码</th>
                <th>加入时间</th>
                <th v-if="canManage">操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="member in members" :key="member.team_member_id">
                <td>{{ member.student_id }}</td>
                <td>{{ member.name }}</td>
                <td>{{ member.college }}</td>
                <td>{{ member.sport_type }}</td>
                <td>{{ member.jersey_number || '-' }}</td>
                <td>{{ new Date(member.join_date).toLocaleDateString() }}</td>
                <td v-if="canManage">
                  <button class="action-btn delete-btn small" @click="removeMember(member)">移除</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      
      <div class="modal-footer">
        <button class="modal-btn cancel" @click="close">关闭</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed, onMounted } from 'vue'
import http from '../utils/http'

const props = defineProps({
  visible: Boolean,
  teamId: Number,
  teamName: String,
  currentUserRole: String // 'admin', 'captain' or 'member'
})

const emit = defineEmits(['update:visible', 'refresh'])

const members = ref([])
const pendingMembers = ref([])
const loading = ref(false)
const isEditingName = ref(false)
const editNameValue = ref('')

const canManage = computed(() => {
  return props.currentUserRole === 'admin' || props.currentUserRole === 'captain'
})

const fetchPendingMembers = async () => {
  if (!props.teamId || !canManage.value) return
  try {
    const res = await http.get(`/teams/members/pending?team_id=${props.teamId}`)
    if (Array.isArray(res)) {
      pendingMembers.value = res
    } else if (res && Array.isArray(res.data)) {
      pendingMembers.value = res.data
    } else {
      pendingMembers.value = []
    }
  } catch (e) {
    console.error('Fetch pending members error:', e)
    pendingMembers.value = []
  }
}

const fetchMembers = async () => {
  if (!props.teamId) return
  loading.value = true
  try {
    const res = await http.get(`/teams/members?team_id=${props.teamId}`)
    if (Array.isArray(res)) {
      members.value = res
    } else if (res && Array.isArray(res.data)) {
      members.value = res.data
    } else {
      members.value = []
    }
    
    if (canManage.value) {
      await fetchPendingMembers()
    }
  } catch (e) {
    console.error('Fetch members error:', e)
    members.value = []
  } finally {
    loading.value = false
  }
}

const approveMember = async (member) => {
  try {
    await http.post('/teams/members/approve', {
      team_member_id: member.team_member_id
    })
    fetchMembers()
  } catch (e) {
    alert('批准失败: ' + (e.response?.data?.message || e.message))
  }
}

const rejectMember = async (member) => {
  if (!confirm(`确定要拒绝 ${member.name} 的申请吗？`)) return
  try {
    await http.post('/teams/remove-member', {
      team_member_id: member.team_member_id
    })
    fetchMembers()
  } catch (e) {
    alert('拒绝失败: ' + (e.response?.data?.message || e.message))
  }
}

const removeMember = async (member) => {
  if (!confirm(`确定要移除成员 ${member.name} 吗？`)) return
  try {
    await http.post('/teams/remove-member', {
      team_member_id: member.team_member_id
    })
    // alert('移除成功')
    fetchMembers()
  } catch (e) {
    alert('移除失败: ' + (e.response?.data?.message || e.message))
  }
}

const startEditName = () => {
  editNameValue.value = props.teamName
  isEditingName.value = true
}

const cancelEditName = () => {
  isEditingName.value = false
  editNameValue.value = ''
}

const saveName = async () => {
  if (!editNameValue.value.trim()) {
    alert('队伍名称不能为空')
    return
  }
  try {
    await http.post('/teams/update-name', {
      team_id: props.teamId,
      team_name: editNameValue.value
    })
    // alert('修改成功')
    isEditingName.value = false
    emit('refresh') // Notify parent to refresh team list/name
  } catch (e) {
    alert('修改失败: ' + (e.response?.data?.message || e.message))
  }
}

const close = () => {
  emit('update:visible', false)
  isEditingName.value = false
}

watch(() => props.visible, (newVal) => {
  if (newVal) {
    fetchMembers()
  }
})

// 组件挂载时如果可见，立即获取数据
onMounted(() => {
  if (props.visible) {
    fetchMembers()
  }
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease;
}

.modal-content.large-modal {
  width: 90%;
  max-width: 800px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  max-height: 85vh;
  overflow: hidden;
  animation: slideUp 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.actions-cell {
  display: flex;
  gap: 8px;
  align-items: center;
}

.modal-header {
  padding: 20px 24px;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #fafafa;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.25rem;
  color: #333;
  font-weight: 600;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.edit-name-container {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.edit-name-input {
  padding: 6px 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 1rem;
  outline: none;
  width: 200px;
}

.edit-name-input:focus {
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.modal-close {
  background: none;
  border: none;
  font-size: 24px;
  color: #999;
  cursor: pointer;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.2s;
}

.modal-close:hover {
  background: #eee;
  color: #333;
}

.modal-body {
  padding: 24px;
  overflow-y: auto;
  flex: 1;
}

.modal-footer {
  padding: 16px 24px;
  border-top: 1px solid #eee;
  display: flex;
  justify-content: flex-end;
  background: #fafafa;
}

/* Table Styles */
.data-table-container {
  border: 1px solid #eee;
  border-radius: 8px;
  overflow: hidden;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.data-table th {
  background: #f8f9fa;
  font-weight: 600;
  color: #555;
  font-size: 0.9rem;
}

.data-table tbody tr:hover {
  background: #f0f7ff;
}

.data-table tbody tr:last-child td {
  border-bottom: none;
}

/* Buttons */
.section-title {
  margin: 0 0 12px 0;
  font-size: 1rem;
  color: #666;
  font-weight: 600;
}

.pending-section {
  margin-bottom: 24px;
}

.divider {
  height: 1px;
  background: #eee;
  margin: 24px 0;
}

.actions-cell {
  display: flex;
  gap: 8px;
}

.action-btn {
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.85rem;
  transition: all 0.2s;
}

.action-btn.small {
  padding: 4px 8px;
  font-size: 0.8rem;
}

.edit-btn {
  background-color: #f0f0f0;
  color: #333;
}

.delete-btn, .reject-btn {
  background-color: #fff1f0;
  color: #ff4d4f;
}

.delete-btn:hover, .reject-btn:hover {
  background-color: #ffccc7;
}

.approve-btn {
  background-color: #f6ffed;
  color: #52c41a;
}

.approve-btn:hover {
  background-color: #d9f7be;
}

.save-btn {
  background: #f6ffed;
  color: #52c41a;
  border: 1px solid #b7eb8f;
}

.save-btn:hover {
  background: #d9f7be;
}

.cancel-btn {
  background: #f5f5f5;
  color: #666;
  border: 1px solid #d9d9d9;
}

.cancel-btn:hover {
  background: #e8e8e8;
}

.modal-btn {
  padding: 8px 20px;
  border-radius: 8px;
  border: none;
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.2s;
  font-weight: 500;
}

.modal-btn.cancel {
  background: #fff;
  border: 1px solid #d9d9d9;
  color: #666;
}

.modal-btn.cancel:hover {
  border-color: #40a9ff;
  color: #40a9ff;
}

.loading, .no-data {
  text-align: center;
  padding: 40px;
  color: #999;
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
