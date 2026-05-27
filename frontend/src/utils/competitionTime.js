export function getCompetitionPhase(competition, currentTime) {
  const startTime = new Date(competition?.start_time).getTime()
  const endTime = new Date(competition?.end_time).getTime()

  if (!Number.isFinite(startTime) || !Number.isFinite(endTime)) {
    return 'pending'
  }

  if (currentTime < startTime) {
    return 'upcoming'
  }

  if (currentTime <= endTime) {
    return 'running'
  }

  return 'ended'
}

export function formatCountdown(milliseconds) {
  const totalSeconds = Math.max(0, Math.floor(milliseconds / 1000))
  const days = Math.floor(totalSeconds / 86400)
  const hours = Math.floor((totalSeconds % 86400) / 3600)
  const minutes = Math.floor((totalSeconds % 3600) / 60)
  const seconds = totalSeconds % 60
  const pad = (value) => String(value).padStart(2, '0')

  return days
    ? `${days}天 ${pad(hours)}:${pad(minutes)}:${pad(seconds)}`
    : `${pad(hours)}:${pad(minutes)}:${pad(seconds)}`
}

export function formatCompetitionSchedule(competition) {
  const start = formatTime(competition.start_time)
  const end = formatTime(competition.end_time)
  return start && end ? `${start} 至 ${end}` : '比赛时间待定'
}

export function formatCompetitionDuration(competition) {
  const startTime = new Date(competition.start_time).getTime()
  const endTime = new Date(competition.end_time).getTime()
  const milliseconds = endTime - startTime

  if (!Number.isFinite(milliseconds) || milliseconds <= 0) {
    return ''
  }

  const totalMinutes = Math.round(milliseconds / 60000)
  const days = Math.floor(totalMinutes / 1440)
  const hours = Math.floor((totalMinutes % 1440) / 60)
  const minutes = totalMinutes % 60
  const parts = []

  if (days) {
    parts.push(`${days}天`)
  }
  if (hours) {
    parts.push(`${hours}小时`)
  }
  if (minutes) {
    parts.push(`${minutes}分钟`)
  }

  return parts.join('')
}

export function formatParticipantCount(value) {
  return Number(value || 0).toLocaleString('zh-CN')
}

export function getCompetitionAccess(competition) {
  const requiresPassword =
    Boolean(competition.password_hash) ||
    competition.visibility === false ||
    Number(competition.visibility) === 0

  return requiresPassword
    ? { label: '需要密码', icon: 'bi-lock-fill' }
    : { label: '无需密码', icon: 'bi-unlock-fill' }
}

function formatTime(value) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return ''
  }

  const pad = (part) => String(part).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`
}
