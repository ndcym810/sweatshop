import { useEffect, useRef, useCallback } from 'react'
import { useClaudeTeamStore } from '../stores/claude-team-store'

interface WSMessage {
  event: 'team:discovered' | 'team:updated' | 'message:new' | 'message:read'
  timestamp: string
  data: unknown
}

export function useWebSocket(url: string = 'ws://localhost:8000/ws') {
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const { setWsConnected, fetchTeams, selectedTeam, refreshCurrentInbox, selectedMember } = useClaudeTeamStore()

  const connect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) return

    const ws = new WebSocket(url)

    ws.onopen = () => {
      setWsConnected(true)
      console.log('WebSocket connected')
    }

    ws.onclose = () => {
      setWsConnected(false)
      console.log('WebSocket disconnected, reconnecting...')
      // Reconnect after 3 seconds
      reconnectTimeoutRef.current = setTimeout(connect, 3000)
    }

    ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }

    ws.onmessage = (event) => {
      try {
        const msg: WSMessage = JSON.parse(event.data)
        handleWSMessage(msg)
      } catch (e) {
        console.error('Failed to parse WebSocket message:', e)
      }
    }

    wsRef.current = ws
  }, [url, setWsConnected])

  const handleWSMessage = useCallback(
    (msg: WSMessage) => {
      switch (msg.event) {
        case 'team:discovered':
          // Refresh team list
          fetchTeams()
          break
        case 'team:updated':
          // Refresh if it's the current team
          if (selectedTeam && (msg.data as { name: string }).name === selectedTeam.name) {
            fetchTeams()
          }
          break
        case 'message:new': {
          const data = msg.data as { team: string; agent: string }
          if (
            selectedTeam?.name === data.team &&
            selectedMember?.name === data.agent
          ) {
            refreshCurrentInbox()
          }
          break
        }
        case 'message:read':
          // Already handled optimistically in store
          break
      }
    },
    [fetchTeams, selectedTeam, selectedMember, refreshCurrentInbox]
  )

  useEffect(() => {
    connect()

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
      }
      wsRef.current?.close()
    }
  }, [connect])

  return {
    isConnected: wsRef.current?.readyState === WebSocket.OPEN,
  }
}
