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
                <select v-model="teamForm.team_type" required>
                  <option value="" disabled>选择队伍类型</option>
                  <option value="college">院队</option>
                  <option value="school">校队</option>
                </select>
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
                    {{ sports.find(s => s.value === (team.sport_id === 1 ? 'football' : team.sport_id === 2 ? 'basketball' : team.sport_id === 3 ? 'badminton' : team.sport_id === 4 ? 'volleyball' : ''))?.label || team.sport_id }}
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
                <select v-model="eventForm.format_type" required>
                  <option value="points">积分制</option>
                  <option value="knockout">淘汰赛</option>
                </select>
              </div>
            </div>
            <div class="form-row">
              <div class="form-group full-width">
                <label>参赛队伍</label>
                <div class="team-picker">
                  <div class="team-picker-toolbar">
                    <input v-model="eventTeamSearch" class="team-picker-search" type="text" placeholder="搜索队伍名称/ID" />
                    <div class="team-picker-actions">
                      <button type="button" class="mini-btn" @click="selectAllFilteredTeams">全选</button>
                      <button type="button" class="mini-btn danger" @click="clearSelectedTeams">清空</button>
                    </div>
                  </div>
                  <div class="team-picker-summary">
                    已选 {{ (eventForm.team_ids || []).length }} 支
                    <span v-if="eventForm.sport" class="team-picker-summary-sub">（已按项目筛选，可搜索）</span>
                  </div>
                  <div class="team-picker-list">
                    <label v-for="team in filteredTeamsForEvent" :key="team.team_id" class="team-picker-item">
                      <input v-model="eventForm.team_ids" type="checkbox" :value="Number(team.team_id)" />
                      <span class="team-picker-name">{{ team.team_name }}</span>
                      <span class="team-picker-meta">#{{ team.team_id }}</span>
                    </label>
                  </div>
                </div>
                <small v-if="eventForm.format_type === 'knockout'" class="hint">
                  需要选择与队伍数量一致的队伍，并为每个种子位分配队伍
                </small>
                <small v-else class="hint">
                  赛事创建后，创建比赛时只能从这些队伍中选择
                </small>
              </div>
            </div>

            <div v-if="isKnockoutConfigEnabled" class="form-row">
              <div class="form-group">
                <label>队伍数量</label>
                <input type="number" :value="eventForm.team_count" disabled />
                <small class="hint">根据所选队伍自动计算，建议为 2 的幂</small>
              </div>
            </div>

            <div v-if="isKnockoutConfigEnabled && knockoutPairs.length > 0" class="knockout-config">
              <!-- Round 1 -->
              <div class="stage-section">
                <h4>第 1 轮</h4>
                <div class="knockout-grid">
                  <div v-for="p in knockoutPairs" :key="p.key" class="knockout-match-card">
                    <div class="knockout-match-title">第 {{ p.index }} 场</div>
                    <div class="knockout-match-row">
                      <div class="knockout-slot">
                        <div class="knockout-slot-label">种子 {{ p.slotA }}</div>
                        <select v-model.number="eventForm.knockout_slots[p.slotA - 1]">
                          <option :value="0" disabled>选择队伍</option>
                          <option v-for="t in selectedEventTeams" :key="t.team_id" :value="t.team_id" :disabled="isTeamUsed(t.team_id, p.slotA - 1)">
                            {{ t.team_name }}
                          </option>
                        </select>
                      </div>
                      <div class="knockout-vs">VS</div>
                      <div class="knockout-slot">
                        <div class="knockout-slot-label">种子 {{ p.slotB }}</div>
                        <select v-model.number="eventForm.knockout_slots[p.slotB - 1]">
                          <option :value="0" disabled>选择队伍</option>
                          <option v-for="t in selectedEventTeams" :key="t.team_id" :value="t.team_id" :disabled="isTeamUsed(t.team_id, p.slotB - 1)">
                            {{ t.team_name }}
                          </option>
                        </select>
                      </div>
                    </div>
                    <div class="schedule-inputs" style="margin-top: 10px; border-top: 1px dashed #eee; padding-top: 10px;">
                      <div class="form-group" style="margin-bottom: 0;">
                        <label style="font-size: 0.9em; margin-bottom: 4px;">比赛时间</label>
                        <input type="datetime-local" v-model="scheduleInputMap['r1-m' + p.index]" required />
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Subsequent Rounds -->
              <div v-for="stage in knockoutSchedule.slice(1)" :key="stage.stage_id" class="stage-section">
                <h4>{{ stage.stage_name }}</h4>
                <div class="knockout-grid">
                  <div v-for="m in stage.matches" :key="m.knockout_match_id" class="knockout-match-card schedule-card">
                    <div class="knockout-match-title">{{ m.match_name || '未命名比赛' }}</div>
                    <div class="schedule-inputs">
                      <div class="match-teams-preview">
                        <span>{{ m.team_a_name || 'TBD' }}</span>
                        <span class="vs">VS</span>
                        <span>{{ m.team_b_name || 'TBD' }}</span>
                      </div>
                      <div class="form-group">
                        <label>比赛时间</label>
                        <input type="datetime-local" v-model="scheduleInputMap[m.knockout_match_id]" placeholder="选择时间" />
                      </div>
                    </div>
                  </div>
                </div>
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
                     {{ sports.find(s => s.value === (event.sport_id === 1 ? 'football' : event.sport_id === 2 ? 'basketball' : event.sport_id === 3 ? 'badminton' : event.sport_id === 4 ? 'volleyball' : ''))?.label || event.sport_id }}
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
                <label>运动项目</label>
                <select v-model="matchEventSportFilter" @change="onMatchEventSportFilterChange">
                  <option value="">全部</option>
                  <option v-for="sport in sports" :key="sport.value" :value="sport.value">
                    {{ sport.label }}
                  </option>
                </select>
              </div>
              <div class="form-group">
                <label>赛事</label>
                <div v-if="events.length > 0">
                   <select v-model="matchForm.event_id" required>
                    <option value="" disabled>选择赛事</option>
                    <option v-for="e in filteredEventsForMatch" :key="e.event_id" :value="e.event_id">
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
                     <option v-for="t in availableTeamsForMatch" :key="t.team_id" :value="t.team_id">
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
                     <option v-for="t in availableTeamsForMatch" :key="t.team_id" :value="t.team_id">
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

      <!-- Carousel Management -->
      <div v-if="currentTab === 'carousel'" class="form-container slide-in">
        <h2>首页轮播图片管理</h2>
        
        <div class="collapsible-section">
            <button class="submit-btn" style="width: auto; margin-bottom: 20px;" @click="openCarouselModal()">添加新图片</button>
        </div>

        <!-- Carousel List -->
        <div class="list-section">
          <h3>现有轮播图片</h3>
          <div v-if="carouselImages.length === 0" class="no-data">暂无图片</div>
          <div v-else class="data-table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th style="width: 60px; text-align: center;">排序</th>
                  <th style="width: 100px;">预览</th>
                  <th>标题</th>
                  <th>链接</th>
                  <th style="width: 80px;">状态</th>
                  <th style="width: 140px; text-align: center;">操作</th>
                </tr>
              </thead>
              <tbody class="carousel-list-body">
                <tr v-for="(img, index) in carouselImages" :key="img.id"
                    draggable="true"
                    @dragstart="onDragStart(index)"
                    @dragover="onDragOver($event, index)"
                    @drop="onDrop(index)"
                    :class="{ 'dragging': draggedItemIndex === index, 'drag-over': dragOverItemIndex === index }"
                    @dragenter="dragOverItemIndex = index"
                    @dragleave="dragOverItemIndex = null"
                    style="transition: all 0.2s;"
                >
                  <td class="drag-handle" style="cursor: grab; text-align: center; font-size: 1.5rem; color: #999;" title="按住拖动排序">
                    ☰
                  </td>
                  <td>
                    <img :src="img.image_url" alt="preview" style="height: 50px; width: 80px; object-fit: cover; border-radius: 4px;" />
                  </td>
                  <td>{{ img.title || '-' }}</td>
                  <td>
                    <a v-if="img.link_url" :href="img.link_url" target="_blank" style="max-width: 150px; display: inline-block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
                      {{ img.link_url }}
                    </a>
                    <span v-else>-</span>
                  </td>
                  <td>
                    <span :class="['status-badge', img.is_active ? 'active' : 'inactive']" style="background: #e6f7ff; color: #1890ff; padding: 2px 8px; border-radius: 4px; font-size: 12px;">
                      {{ img.is_active ? '启用' : '禁用' }}
                    </span>
                  </td>
                  <td class="actions-cell" style="text-align: center;">
                    <button class="action-btn view-btn" @click="editCarousel(img)">编辑</button>
                    <button class="action-btn delete-btn" @click="deleteCarousel(img)">删除</button>
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
                <select v-model="eventForm.format_type" required>
                  <option value="points">积分制</option>
                  <option value="knockout">淘汰赛</option>
                </select>
              </div>
            </div>
            <div class="form-row">
              <div class="form-group full-width">
                <label>参赛队伍</label>
                <select v-model="eventForm.team_ids" multiple required>
                  <option v-for="team in teams" :key="team.team_id" :value="team.team_id">
                    {{ team.team_name }}
                  </option>
                </select>
              </div>
            </div>

            <div v-if="isKnockoutConfigEnabled" class="form-row">
              <div class="form-group">
                <label>队伍数量</label>
                <input type="number" :value="eventForm.team_count" disabled />
                <small class="hint">根据所选队伍自动计算，建议为 2 的幂</small>
              </div>
            </div>

            <div v-if="isKnockoutConfigEnabled && knockoutPairs.length > 0" class="knockout-config">
              <!-- Round 1 -->
              <div class="stage-section">
                <h4>第 1 轮</h4>
                <div class="knockout-grid">
                  <div v-for="p in knockoutPairs" :key="p.key" class="knockout-match-card">
                    <div class="knockout-match-title">第 {{ p.index }} 场</div>
                    <div class="knockout-match-row">
                      <div class="knockout-slot">
                        <div class="knockout-slot-label">种子 {{ p.slotA }}</div>
                        <select v-model.number="eventForm.knockout_slots[p.slotA - 1]">
                          <option :value="0" disabled>选择队伍</option>
                          <option v-for="t in selectedEventTeams" :key="t.team_id" :value="t.team_id" :disabled="isTeamUsed(t.team_id, p.slotA - 1)">
                            {{ t.team_name }}
                          </option>
                        </select>
                      </div>
                      <div class="knockout-vs">VS</div>
                      <div class="knockout-slot">
                        <div class="knockout-slot-label">种子 {{ p.slotB }}</div>
                        <select v-model.number="eventForm.knockout_slots[p.slotB - 1]">
                          <option :value="0" disabled>选择队伍</option>
                          <option v-for="t in selectedEventTeams" :key="t.team_id" :value="t.team_id" :disabled="isTeamUsed(t.team_id, p.slotB - 1)">
                            {{ t.team_name }}
                          </option>
                        </select>
                      </div>
                    </div>
                    <div class="schedule-inputs" style="margin-top: 10px; border-top: 1px dashed #eee; padding-top: 10px;">
                      <div class="form-group" style="margin-bottom: 0;">
                        <label style="font-size: 0.9em; margin-bottom: 4px;">比赛时间</label>
                        <input type="datetime-local" v-model="scheduleInputMap['r1-m' + p.index]" required />
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Subsequent Rounds -->
              <div v-for="stage in knockoutSchedule.slice(1)" :key="stage.stage_id" class="stage-section">
                <h4>{{ stage.stage_name }}</h4>
                <div class="knockout-grid">
                  <div v-for="m in stage.matches" :key="m.knockout_match_id" class="knockout-match-card schedule-card">
                    <div class="knockout-match-title">{{ m.match_name || '未命名比赛' }}</div>
                    <div class="schedule-inputs">
                      <div class="match-teams-preview">
                        <span>{{ m.team_a_name || 'TBD' }}</span>
                        <span class="vs">VS</span>
                        <span>{{ m.team_b_name || 'TBD' }}</span>
                      </div>
                      <div class="form-group">
                        <label>比赛时间</label>
                        <input type="datetime-local" v-model="scheduleInputMap[m.knockout_match_id]" placeholder="选择时间" />
                      </div>
                    </div>
                  </div>
                </div>
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
                <label>运动项目</label>
                <select v-model="matchEventSportFilter" @change="onMatchEventSportFilterChange">
                  <option value="">全部</option>
                  <option v-for="sport in sports" :key="sport.value" :value="sport.value">
                    {{ sport.label }}
                  </option>
                </select>
              </div>
              <div class="form-group">
                <label>赛事</label>
                <div v-if="events.length > 0">
                   <select v-model="matchForm.event_id" required>
                    <option value="" disabled>选择赛事</option>
                    <option v-for="e in filteredEventsForMatch" :key="e.event_id" :value="e.event_id">
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
                     <option v-for="t in availableTeamsForMatch" :key="t.team_id" :value="t.team_id">
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
                     <option v-for="t in availableTeamsForMatch" :key="t.team_id" :value="t.team_id">
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

    <!-- Carousel Edit Modal -->
    <div v-if="showCarouselModal" class="modal-overlay" @click="closeCarouselModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3 class="modal-title">{{ isEditingCarousel ? '编辑图片' : '添加图片' }}</h3>
          <button @click="closeCarouselModal" class="modal-close">×</button>
        </div>
        <div class="modal-body">
          <form id="carouselForm" @submit.prevent="saveCarousel" class="compact-form">
            <div class="form-group">
              <label>图片</label>
              <div style="display: flex; flex-direction: column; gap: 8px;">
                 <input type="file" accept="image/*" @change="handleFileUpload" />
                 <span v-if="isUploading" style="font-size: 0.9em; color: #666;">正在上传...</span>
                 <input v-model="carouselForm.image_url" placeholder="或输入图片 URL" />
                 <img v-if="carouselForm.image_url" :src="carouselForm.image_url" style="max-height: 100px; object-fit: contain; margin-top: 5px;" />
              </div>
            </div>
            <div class="form-group">
              <label>标题 (可选)</label>
              <input v-model="carouselForm.title" placeholder="例如：2024 春季运动会" />
            </div>
            <div class="form-group">
              <label>跳转链接 (可选)</label>
              <input v-model="carouselForm.link_url" placeholder="https://..." />
            </div>
            <div class="form-row">
              <div class="form-group" style="display: flex; align-items: center; padding-top: 10px;">
                <label style="display: flex; align-items: center; cursor: pointer;">
                  <input type="checkbox" v-model="carouselForm.is_active" style="width: auto; margin-right: 8px;" />
                  启用
                </label>
              </div>
            </div>
          </form>
        </div>
        <div class="modal-footer">
          <button @click="closeCarouselModal" class="btn-secondary">取消</button>
          <button type="submit" form="carouselForm" class="btn-primary" :disabled="loading">
            {{ loading ? '保存中...' : (isEditingCarousel ? '保存修改' : '添加图片') }}
          </button>
        </div>
      </div>
    </div>

    <!-- Reused Team Management Modal -->
    <TeamManagementModal 
      v-model:visible="showMembersModal"
      :team-id="currentTeamId"
      :team-name="currentTeamName"
      :avatar-url="currentTeamAvatar"
      current-user-role="admin"
      @refresh="fetchTeams"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
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
  { id: 'carousel', label: '首页图片管理' },
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
  { value: 'volleyball', label: '排球' }
]

const sportIdToValue = (sportID) => {
  if (sportID === 1) return 'football'
  if (sportID === 2) return 'basketball'
  if (sportID === 3) return 'badminton'
  if (sportID === 4) return 'volleyball'
  return 'football'
}

// Forms
const teamForm = ref({
  name: '',
  sport: '',
  college: '',
  team_type: ''
})

const matchEventSportFilter = ref('')

const filteredEventsForMatch = computed(() => {
  if (!matchEventSportFilter.value) return events.value
  return events.value.filter(e => sportIdToValue(e.sport_id) === matchEventSportFilter.value)
})

const onMatchEventSportFilterChange = () => {
  const selectedID = Number(matchForm.value.event_id)
  if (!selectedID) return
  const stillVisible = filteredEventsForMatch.value.some(e => e.event_id === selectedID)
  if (!stillVisible) {
    matchForm.value.event_id = ''
  }
}

const normalizeFormatType = (formatType) => {
  const v = String(formatType || '').trim()
  if (v === 'knockout' || v === '淘汰赛') return 'knockout'
  if (v === 'points' || v === '积分制') return 'points'
  return 'points'
}

const sportValueToId = (sportValue) => {
  const v = String(sportValue || '').trim()
  if (v === 'football') return 1
  if (v === 'basketball') return 2
  if (v === 'badminton') return 3
  if (v === 'volleyball') return 4
  if (v) return 5
  return 0
}

const eventForm = ref({
  name: '',
  sport: '',
  start_date: '',
  format_type: 'points',
  team_ids: [],
  team_count: 0,
  knockout_slots: []
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

// Carousel State
const carouselImages = ref([])
const showCarouselModal = ref(false)
const isEditingCarousel = ref(false)
const currentCarouselId = ref(null)
const carouselForm = ref({
  image_url: '',
  title: '',
  link_url: '',
  sort_order: 0,
  is_active: true
})

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
const isUploading = ref(false)

const handleFileUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  const formData = new FormData()
  formData.append('file', file)

  isUploading.value = true
  try {
    const res = await http.post('/admin/carousel/upload', formData)
    let url = res.data && res.data.url ? res.data.url : res.url
    if (url && url.startsWith('/')) {
       url = `${API_BASE_URL}${url}`
    }
    carouselForm.value.image_url = url
  } catch (e) {
    alert('上传失败: ' + e.message)
  } finally {
    isUploading.value = false
    event.target.value = ''
  }
}

// Drag and Drop
const draggedItemIndex = ref(null)
const dragOverItemIndex = ref(null)

const onDragStart = (index) => {
  draggedItemIndex.value = index
}

const onDragOver = (event, index) => {
  event.preventDefault()
  if (draggedItemIndex.value === index) return
  dragOverItemIndex.value = index
}

const onDrop = async (index) => {
  const fromIndex = draggedItemIndex.value
  const toIndex = index
  
  draggedItemIndex.value = null
  dragOverItemIndex.value = null
  
  if (fromIndex === null || fromIndex === toIndex) return

  // Move item in local array
  const item = carouselImages.value.splice(fromIndex, 1)[0]
  carouselImages.value.splice(toIndex, 0, item)
  
  // Save order to backend
  await saveOrder()
}

const saveOrder = async () => {
  const ids = carouselImages.value.map(img => img.id)
  try {
    await http.post('/admin/carousel/reorder', ids)
  } catch (e) {
    console.error('Failed to save order', e)
    alert('排序保存失败，请刷新重试')
    fetchCarouselImages()
  }
}

const fetchCarouselImages = async () => {
  try {
    const res = await http.get('/admin/carousel')
    let data = []
    if (Array.isArray(res)) {
      data = res
    } else if (res && Array.isArray(res.data)) {
      data = res.data
    }
    
    carouselImages.value = data.map(img => {
       let url = img.image_url
       if (url && url.startsWith('/')) {
         url = `${API_BASE_URL}${url}`
       }
       return { ...img, image_url: url }
    })
  } catch (e) {
    console.warn('获取轮播图片失败', e)
  }
}

const openCarouselModal = () => {
  isEditingCarousel.value = false
  currentCarouselId.value = null
  carouselForm.value = { image_url: '', title: '', link_url: '', sort_order: 0, is_active: true }
  showCarouselModal.value = true
}

const closeCarouselModal = () => {
  showCarouselModal.value = false
}

const editCarousel = (img) => {
  isEditingCarousel.value = true
  currentCarouselId.value = img.id
  carouselForm.value = {
    image_url: img.image_url,
    title: img.title,
    link_url: img.link_url,
    sort_order: img.sort_order,
    is_active: img.is_active
  }
  showCarouselModal.value = true
}

const saveCarousel = async () => {
  loading.value = true
  try {
    const payload = { ...carouselForm.value }
    if (isEditingCarousel.value) {
      payload.id = currentCarouselId.value
      await http.post('/admin/carousel/update', payload)
      alert('更新成功')
    } else {
      await http.post('/admin/carousel/create', payload)
      alert('创建成功')
    }
    closeCarouselModal()
    fetchCarouselImages()
  } catch (e) {
    alert('操作失败: ' + (e.response?.data?.message || e.message))
  } finally {
    loading.value = false
  }
}

const deleteCarousel = (img) => {
  openModal(
    '确认删除',
    '确定要删除这张轮播图片吗？',
    async () => {
      try {
        await http.post(`/admin/carousel/delete?id=${img.id}`)
        alert('删除成功')
        fetchCarouselImages()
      } catch (e) {
        alert('删除失败: ' + e.message)
      }
    }
  )
}

const eventTeamSearch = ref('')

// Team Members Modal State
const showMembersModal = ref(false)
const currentTeamId = ref(null)

const currentTeamInfo = computed(() => {
  return teams.value.find(t => t.team_id === currentTeamId.value) || {}
})

const currentTeamName = computed(() => currentTeamInfo.value.team_name || '')
const currentTeamAvatar = computed(() => currentTeamInfo.value.avatar_url || '')

// Edit State
const isEditingEvent = ref(false)
const currentEventId = ref(null)
const isEditingMatch = ref(false)
const currentMatchId = ref(null)

const knockoutSchedule = ref([])
const scheduleInputMap = ref({})

const fetchKnockoutSchedule = async (eventID) => {
  try {
    const res = await http.get(`/events/${eventID}/standings/overview`)
    if (res && res.knockout_stages) {
      scheduleInputMap.value = {}
      knockoutSchedule.value = res.knockout_stages.map(stage => ({
        ...stage,
        matches: stage.matches.map(m => {
          let timeVal = ''
          const raw = m.match_time || m.scheduled_time
          if (raw) {
            // Convert to local time for input
            const d = new Date(raw)
            if (!isNaN(d.getTime())) {
               timeVal = new Date(d.getTime() - (d.getTimezoneOffset() * 60000)).toISOString().slice(0, 16)
            } else {
               timeVal = raw.replace(' ', 'T').slice(0, 16)
            }
          }
          scheduleInputMap.value[m.knockout_match_id] = timeVal
          return {
            ...m,
            scheduled_time_input: timeVal
          }
        })
      }))
    } else {
      knockoutSchedule.value = []
      scheduleInputMap.value = {}
    }
  } catch {
    knockoutSchedule.value = []
    scheduleInputMap.value = {}
  }
}

const generateKnockoutSchedulePreview = () => {
  // if (isEditingEvent.value) return 
  
  const n = Number(eventForm.value.team_count || 0)
  if (!n || n < 2 || !isPowerOfTwo(n)) {
    knockoutSchedule.value = []
    scheduleInputMap.value = {}
    return
  }

  const rounds = Math.log2(n)
  const schedule = []
  const newScheduleInputMap = { ...scheduleInputMap.value }
  
  for (let r = 1; r <= rounds; r++) {
    const matchCount = n / Math.pow(2, r)
    const matches = []
    for (let m = 1; m <= matchCount; m++) {
      const matchKey = `r${r}-m${m}`
      
      let teamA = 'TBD'
      let teamB = 'TBD'
      
      if (r === 1) {
        // Match m corresponds to pair m (1-based)
        // Pair logic: slotA = m, slotB = n + 1 - m
        // Array index = slot - 1
        const idxA = m - 1
        const idxB = n - m
        
        const tidA = (eventForm.value.knockout_slots || [])[idxA]
        const tidB = (eventForm.value.knockout_slots || [])[idxB]
        
        if (tidA) {
            const t = teams.value.find(x => x.team_id === tidA)
            if (t) teamA = t.team_name
        }
        if (tidB) {
            const t = teams.value.find(x => x.team_id === tidB)
            if (t) teamB = t.team_name
        }
      }
      
      matches.push({
        knockout_match_id: matchKey,
        match_name: `第 ${r} 轮 第 ${m} 场`,
        team_a_name: teamA,
        team_b_name: teamB
      })
      
      // Initialize input map key if not exists
      if (newScheduleInputMap[matchKey] === undefined) {
        newScheduleInputMap[matchKey] = ''
      }
    }
    schedule.push({
      stage_id: `preview-stage-${r}`,
      stage_name: `第 ${r} 轮`,
      matches: matches
    })
  }
  knockoutSchedule.value = schedule
  scheduleInputMap.value = newScheduleInputMap
}

const matchEventTeams = ref([])
const availableTeamsForMatch = computed(() => {
  if (matchEventTeams.value.length > 0) return matchEventTeams.value
  return teams.value
})

const selectedEventTeams = computed(() => {
  const ids = new Set((eventForm.value.team_ids || []).map(v => Number(v)))
  return teams.value.filter(t => ids.has(Number(t.team_id)))
})

const filteredTeamsForEvent = computed(() => {
  const q = eventTeamSearch.value.trim().toLowerCase()
  const sportID = sportValueToId(eventForm.value.sport)
  const selected = new Set((eventForm.value.team_ids || []).map(v => Number(v)))
  return teams.value.filter(t => {
    const tid = Number(t.team_id)
    if (sportID && Number(t.sport_id) !== sportID && !selected.has(tid)) return false
    if (!q) return true
    const name = String(t.team_name || '').toLowerCase()
    return name.includes(q) || String(tid).includes(q)
  })
})

const selectAllFilteredTeams = () => {
  const next = new Set((eventForm.value.team_ids || []).map(v => Number(v)))
  for (const t of filteredTeamsForEvent.value) {
    next.add(Number(t.team_id))
  }
  let list = Array.from(next)
  if (eventForm.value.format_type === 'knockout') {
    const n = Number(eventForm.value.team_count || 0)
    if (Number.isInteger(n) && n > 0) {
      list = list.slice(0, n)
    }
  }
  eventForm.value.team_ids = list
}

const clearSelectedTeams = () => {
  eventForm.value.team_ids = []
}

const isKnockoutConfigEnabled = computed(() => {
  const ft = eventForm.value.format_type
  return ft === 'knockout'
})

const knockoutPairs = computed(() => {
  if (!isKnockoutConfigEnabled.value) return []
  const n = Number(eventForm.value.team_count || 0)
  if (!n || n < 2) return []
  const out = []
  for (let i = 1; i <= n / 2; i++) {
    out.push({ key: `${n}-${i}`, index: i, slotA: i, slotB: n + 1 - i })
  }
  return out
})

const isPowerOfTwo = (n) => {
  const x = Number(n)
  return Number.isInteger(x) && x >= 2 && (x & (x - 1)) === 0
}

const isTeamUsed = (teamId, currentSlotIndex) => {
  if (!eventForm.value.knockout_slots) return false
  return eventForm.value.knockout_slots.some((id, idx) => idx !== currentSlotIndex && id === teamId)
}

const ensureKnockoutSlots = () => {
  const n = Number(eventForm.value.team_count || 0)
  if (!Number.isInteger(n) || n < 2) {
    eventForm.value.knockout_slots = []
    return
  }
  const next = new Array(n).fill(0)
  const current = Array.isArray(eventForm.value.knockout_slots) ? eventForm.value.knockout_slots : []
  for (let i = 0; i < Math.min(current.length, next.length); i++) {
    next[i] = Number(current[i] || 0)
  }
  const allowed = new Set((eventForm.value.team_ids || []).map(v => Number(v)))
  for (let i = 0; i < next.length; i++) {
    if (!allowed.has(Number(next[i] || 0))) {
      next[i] = 0
    }
  }
  eventForm.value.knockout_slots = next
}

watch(
  () => eventForm.value.format_type,
  (ft) => {
    if (ft === 'knockout') {
      // Auto-set team count if valid
      const currentCount = eventForm.value.team_ids.length
      if (currentCount >= 2) {
        eventForm.value.team_count = currentCount
      }
      
      ensureKnockoutSlots()
      generateKnockoutSchedulePreview()
      return
    }
    eventForm.value.knockout_slots = []
    // Clear schedule if not knockout
    knockoutSchedule.value = []
    scheduleInputMap.value = {}
  }
)

watch(
  () => eventForm.value.team_count,
  () => {
    if (!isKnockoutConfigEnabled.value) return
    ensureKnockoutSlots()
    generateKnockoutSchedulePreview()
  }
)

watch(
  () => eventForm.value.team_ids,
  () => {
    if (!isKnockoutConfigEnabled.value) return
    // Always sync team_count with team_ids length when in knockout mode
    eventForm.value.team_count = eventForm.value.team_ids.length
    ensureKnockoutSlots()
    generateKnockoutSchedulePreview()
  },
  { deep: true }
)

watch(
  () => eventForm.value.knockout_slots,
  () => {
     if (isKnockoutConfigEnabled.value) {
        generateKnockoutSchedulePreview()
     }
  },
  { deep: true }
)

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

const fetchEventTeams = async (eventID) => {
  try {
    const res = await http.get(`/events/${eventID}/teams`)
    if (res && Array.isArray(res.data)) return res.data
    if (Array.isArray(res)) return res
    return []
  } catch {
    return []
  }
}

const loadEventTeamsIntoForm = async (eventID) => {
  const list = await fetchEventTeams(eventID)
  const teamIDs = list.map(it => Number(it.team_id)).filter(Boolean)
  eventForm.value.team_ids = teamIDs

  const maxSlot = list.reduce((m, it) => Math.max(m, Number(it.slot || 0)), 0)
  if (eventForm.value.format_type !== 'knockout') {
    eventForm.value.knockout_slots = []
    return
  }

  const inferredCount = maxSlot > 0 ? maxSlot : (isPowerOfTwo(teamIDs.length) ? teamIDs.length : 0)
  const n = inferredCount || Number(eventForm.value.team_count || 0) || teamIDs.length || 8
  eventForm.value.team_count = n

  const slots = new Array(n).fill(0)
  if (maxSlot > 0) {
    for (const it of list) {
      const slot = Number(it.slot || 0)
      const tid = Number(it.team_id || 0)
      if (slot >= 1 && slot <= slots.length && tid) {
        slots[slot - 1] = tid
      }
    }
  } else {
    for (let i = 0; i < Math.min(teamIDs.length, slots.length); i++) {
      slots[i] = Number(teamIDs[i] || 0)
    }
  }

  eventForm.value.knockout_slots = slots
}

const editEvent = async (event) => {
  isEditingEvent.value = true
  currentEventId.value = event.event_id
  eventForm.value = {
    name: event.event_name,
    sport: sports.find(s => s.value === (event.sport_id === 1 ? 'football' : event.sport_id === 2 ? 'basketball' : event.sport_id === 3 ? 'badminton' : event.sport_id === 4 ? 'volleyball' : ''))?.value || '',
    start_date: String(event.start_date || '').slice(0, 10),
    format_type: normalizeFormatType(event.format_type),
    team_ids: [],
    team_count: 8,
    knockout_slots: []
  }
  await loadEventTeamsIntoForm(event.event_id)
  
  if (normalizeFormatType(event.format_type) === 'knockout') {
    await fetchKnockoutSchedule(event.event_id)
  }

  showEventEditModal.value = true
}

const cancelEditEvent = () => {
  isEditingEvent.value = false
  showEventEditModal.value = false
  currentEventId.value = null
  eventForm.value = { name: '', sport: '', start_date: '', format_type: 'points', team_ids: [], team_count: 8, knockout_slots: [] }
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
  const event = events.value.find(e => e.event_id === match.event_id)
  matchEventSportFilter.value = event ? sportIdToValue(event.sport_id) : ''
  
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

const loadMatchEventTeams = async () => {
  const id = Number(matchForm.value.event_id || 0)
  if (!id) {
    matchEventTeams.value = []
    return
  }
  const list = await fetchEventTeams(id)
  matchEventTeams.value = Array.isArray(list) ? list : []
  const allowed = new Set(matchEventTeams.value.map(t => Number(t.team_id)))
  if (matchForm.value.team_a_id && !allowed.has(Number(matchForm.value.team_a_id))) {
    matchForm.value.team_a_id = ''
  }
  if (matchForm.value.team_b_id && !allowed.has(Number(matchForm.value.team_b_id))) {
    matchForm.value.team_b_id = ''
  }
}

watch(
  () => matchForm.value.event_id,
  () => {
    loadMatchEventTeams()
  }
)

onMounted(() => {
  fetchTeams()
  fetchEvents()
  fetchMatches()
  fetchCollectorRequests()
  fetchCollectors()
  fetchAthleteRequests()
  fetchTeamRequests()
  fetchCarouselImages()
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
  if (!teamForm.value.name || !teamForm.value.sport || !teamForm.value.college || !teamForm.value.team_type) {
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

const updateKnockoutSchedule = async (eventID) => {
  const updates = {}
  for (const [k, v] of Object.entries(scheduleInputMap.value)) {
    if (v) {
      updates[k] = new Date(v).toISOString()
    }
  }

  if (Object.keys(updates).length === 0) return

  try {
    await http.put(`/events/${eventID}/knockout-schedule`, updates)
  } catch (e) {
    console.error('Failed to update knockout schedule', e)
    throw new Error('更新赛程时间失败: ' + (e.response?.data?.message || e.message))
  }
}

const createEvent = async () => {
  if (!eventForm.value.name || !eventForm.value.sport || !eventForm.value.start_date) {
    alert('请填写必填字段')
    return
  }
  const formatType = normalizeFormatType(eventForm.value.format_type)
  const selectedIDs = (eventForm.value.team_ids || []).map(v => Number(v)).filter(Boolean)
  if (selectedIDs.length === 0) {
    alert('请选择参赛队伍')
    return
  }

  let eventTeams = []
  if (formatType === 'points') {
    eventTeams = selectedIDs.map(id => ({ team_id: id }))
  } else if (formatType === 'knockout') {
    const n = Number(eventForm.value.team_count || 0)
    if (!isPowerOfTwo(n)) {
      alert('队伍数量必须为 2 的幂（2/4/8/16/32...）')
      return
    }
    if (selectedIDs.length !== n) {
      alert('选中队伍数量需要与队伍数量一致')
      return
    }
    const slots = Array.isArray(eventForm.value.knockout_slots) ? eventForm.value.knockout_slots : []
    if (slots.length !== n) {
      alert('请先完成对阵配置')
      return
    }
    const allowed = new Set(selectedIDs)
    const seen = new Set()
    for (let i = 0; i < n; i++) {
      const tid = Number(slots[i] || 0)
      if (!tid) {
        alert(`请为种子位 ${i + 1} 选择队伍`)
        return
      }
      if (!allowed.has(tid)) {
        alert('对阵配置中的队伍必须来自已选择的参赛队伍')
        return
      }
      if (seen.has(tid)) {
        alert('同一支队伍不能出现在多个种子位')
        return
      }
      seen.add(tid)
      eventTeams.push({ team_id: tid, slot: i + 1 })
    }
  } else {
    alert('未知赛制类型')
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
      format_type: formatType,
      season: '2024',
      round: 'Regular',
      teams: eventTeams
    }

    if (formatType === 'knockout') {
      const schedule = {}
      for (const [k, v] of Object.entries(scheduleInputMap.value)) {
        if (v && String(k).startsWith('r')) {
           schedule[k] = new Date(v).toISOString()
        }
      }
      payload.knockout_schedule = schedule
    }

    if (isEditingEvent.value) {
      await http.put(`/events/${currentEventId.value}`, payload)
      
      if (formatType === 'knockout') {
        await updateKnockoutSchedule(currentEventId.value)
      }

      alert('赛事更新成功')
      isEditingEvent.value = false
      showEventEditModal.value = false
      currentEventId.value = null
    } else {
      await http.post('/events/create', payload)
      alert('赛事创建成功')
    }

    eventForm.value = { name: '', sport: '', start_date: '', format_type: 'points', team_ids: [], team_count: 0, knockout_slots: [] }
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

/* Drag and Drop Styles */
.dragging {
  opacity: 0.5;
  background-color: #e6f7ff;
}

.drag-over {
  border-top: 2px solid var(--primary-color);
  transform: translateY(2px);
}

.drag-handle:hover {
  opacity: 0.8;
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

input:not([type='checkbox']):not([type='radio']), select {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 1rem;
}

input:not([type='checkbox']):not([type='radio']):focus, select:focus {
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

.form-group.full-width {
  flex-basis: 100%;
}

.team-picker {
  border: 1px solid #eee;
  border-radius: 10px;
  padding: 12px;
  background: white;
}

.team-picker-toolbar {
  display: flex;
  gap: 10px;
  align-items: center;
}

.team-picker-search {
  flex: 1;
}

.team-picker-actions {
  display: flex;
  gap: 8px;
}

.mini-btn {
  padding: 8px 10px;
  border: 1px solid #ddd;
  background: #f8fafc;
  border-radius: 10px;
  cursor: pointer;
  font-size: 13px;
  line-height: 1;
  white-space: nowrap;
}

.mini-btn:hover {
  background: #eef2f7;
}

.mini-btn.danger {
  border-color: rgba(231, 76, 60, 0.35);
  color: #e74c3c;
  background: rgba(231, 76, 60, 0.06);
}

.mini-btn.danger:hover {
  background: rgba(231, 76, 60, 0.1);
}

.team-picker-summary {
  margin-top: 10px;
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-secondary);
  font-size: 13px;
}

.team-picker-summary-sub {
  color: var(--text-tertiary);
  font-size: 12px;
}

.team-picker-list {
  margin-top: 10px;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(210px, 1fr));
  gap: 8px;
  max-height: 240px;
  overflow: auto;
  padding: 8px;
  border: 1px solid #f0f0f0;
  border-radius: 10px;
  background: #fafafa;
}

.team-picker-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  border: 1px solid #eee;
  border-radius: 10px;
  background: white;
  cursor: pointer;
  user-select: none;
}

.team-picker-item:hover {
  border-color: #dbeafe;
  background: #f8fbff;
}

.team-picker-item input[type='checkbox'] {
  width: 16px;
  height: 16px;
  margin: 0;
  padding: 0;
  accent-color: #3498db;
  flex: 0 0 auto;
}

.team-picker-name {
  flex: 1;
  font-weight: 500;
  color: #2c3e50;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.team-picker-meta {
  color: #94a3b8;
  font-size: 12px;
  flex: 0 0 auto;
}

.team-config {
  border: 1px solid #eee;
  border-radius: 10px;
  padding: 12px;
  margin-bottom: 20px;
  background: #fafafa;
}

.team-config-title {
  font-weight: 600;
  margin-bottom: 10px;
  color: #2c3e50;
}

.team-config-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.team-config-row:last-child {
  margin-bottom: 0;
}

.team-config-name {
  flex: 1;
  font-weight: 500;
  color: #2c3e50;
}

.team-config-input {
  width: 160px;
}

.knockout-config {
  border: 1px solid #eee;
  border-radius: 10px;
  padding: 12px;
  margin-bottom: 20px;
  background: #fafafa;
}

.knockout-config-title {
  font-weight: 600;
  margin-bottom: 12px;
  color: #2c3e50;
}

.knockout-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 12px;
}

.knockout-match-card {
  border: 1px solid #ddd;
  border-radius: 10px;
  background: white;
  padding: 12px;
  box-shadow: 0 6px 18px rgba(0, 0, 0, 0.04);
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.knockout-match-card:hover {
  transform: translateY(-1px);
  box-shadow: 0 10px 22px rgba(0, 0, 0, 0.06);
}

.knockout-match-title {
  font-weight: 600;
  margin-bottom: 8px;
  color: #2c3e50;
}

.knockout-match-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.knockout-slot {
  flex: 1;
}

.knockout-slot-label {
  font-size: 12px;
  color: #666;
  margin-bottom: 6px;
}

.knockout-vs {
  width: 36px;
  text-align: center;
  font-weight: 700;
  color: #3498db;
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
  max-height: 70vh;
  overflow-y: auto;
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

.empty-schedule-hint {
  text-align: center;
  padding: 20px;
  background-color: #f8f9fa;
  border-radius: 8px;
  border: 1px dashed #ced4da;
  color: #6c757d;
  margin-top: 15px;
}

.stage-section {
  margin-bottom: 20px;
}
.stage-section h4 {
  margin-bottom: 10px;
  color: #555;
  border-left: 3px solid var(--primary-color);
  padding-left: 8px;
}
.schedule-inputs {
  margin-top: 10px;
  border-top: 1px solid #eee;
  padding-top: 10px;
}
.match-teams-preview {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-weight: 500;
  font-size: 0.9em;
  color: #333;
}
.match-teams-preview .vs {
  color: #999;
  font-size: 0.8em;
  font-weight: bold;
}
</style>
