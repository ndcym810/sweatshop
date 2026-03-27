import { useEffect, useState, useRef } from 'react'
import { Send, Loader2, Bot, MessageSquare, CheckCircle } from 'lucide-react'
import { useClaudeTeamStore } from '../../stores/claude-team-store'
import { Button } from '../ui/button'
import { cn } from '../../lib/utils'

export function ClaudeTeamDashboard() {
  const {
    selectedTeam,
    selectedMember,
    inbox,
    isLoading,
    error,
    wsConnected,
    selectMember,
    sendMessage,
    markRead,
    fetchInbox,
  } = useClaudeTeamStore()

  const [messageText, setMessageText] = useState('')
  const [recipient, setRecipient] = useState('')
  const messagesEndRef = useRef<HTMLDivElement>(null)

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [inbox])

  // Refresh inbox periodically when viewing a member
  useEffect(() => {
    if (!selectedTeam || !selectedMember) return
    const interval = setInterval(() => {
      fetchInbox(selectedTeam.name, selectedMember.name)
    }, 5000)
    return () => clearInterval(interval)
  }, [selectedTeam, selectedMember, fetchInbox])

  if (!selectedTeam) {
    return (
      <div className="flex flex-1 items-center justify-center text-muted-foreground">
        <div className="text-center">
          <Bot className="mx-auto mb-2 h-12 w-12 opacity-50" />
          <p>Select a Claude team to view details</p>
        </div>
      </div>
    )
  }

  const handleSendMessage = async () => {
    if (!messageText.trim() || !recipient.trim()) return

    try {
      await sendMessage({ to: recipient, message: messageText })
      setMessageText('')
    } catch (e) {
      // Error handled by store
    }
  }

  return (
    <div className="flex flex-1">
      {/* Members list */}
      <div className="w-64 border-r p-4">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="font-semibold">{selectedTeam.name}</h2>
          <div
            className={cn(
              'h-2 w-2 rounded-full',
              wsConnected ? 'bg-green-500' : 'bg-gray-400'
            )}
            title={wsConnected ? 'Connected' : 'Disconnected'}
          />
        </div>

        {selectedTeam.description && (
          <p className="mb-4 text-sm text-muted-foreground">
            {selectedTeam.description}
          </p>
        )}

        <h3 className="mb-2 text-xs font-medium uppercase text-muted-foreground">
          Members
        </h3>
        <div className="space-y-1">
          {selectedTeam.members.map((member) => (
            <button
              key={member.agentId}
              onClick={() => selectMember(member)}
              className={cn(
                'flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors',
                selectedMember?.agentId === member.agentId
                  ? 'bg-accent text-accent-foreground'
                  : 'hover:bg-accent'
              )}
            >
              <Bot className="h-4 w-4" />
              <span className="truncate">{member.name}</span>
              <span className="ml-auto text-xs text-muted-foreground">
                {member.status}
              </span>
            </button>
          ))}
        </div>
      </div>

      {/* Inbox and messaging */}
      <div className="flex flex-1 flex-col">
        {selectedMember ? (
          <>
            {/* Inbox header */}
            <div className="border-b p-4">
              <h3 className="font-medium">{selectedMember.name}'s Inbox</h3>
              <p className="text-sm text-muted-foreground">
                {selectedMember.agentType} • {selectedMember.model}
              </p>
            </div>

            {/* Messages */}
            <div className="flex-1 overflow-auto p-4">
              {inbox.length === 0 ? (
                <div className="flex h-full items-center justify-center text-muted-foreground">
                  <div className="text-center">
                    <MessageSquare className="mx-auto mb-2 h-8 w-8 opacity-50" />
                    <p>No messages yet</p>
                  </div>
                </div>
              ) : (
                <div className="space-y-3">
                  {inbox.map((msg, i) => (
                    <div
                      key={i}
                      className={cn(
                        'rounded-lg border p-3',
                        msg.read ? 'bg-background' : 'bg-accent/50'
                      )}
                    >
                      <div className="mb-1 flex items-center justify-between">
                        <span className="text-sm font-medium">{msg.from}</span>
                        <div className="flex items-center gap-2">
                          {!msg.read && (
                            <button
                              onClick={() => markRead(msg.timestamp)}
                              className="text-xs text-primary hover:underline"
                            >
                              <CheckCircle className="h-3 w-3" />
                            </button>
                          )}
                          <span className="text-xs text-muted-foreground">
                            {new Date(msg.timestamp).toLocaleTimeString()}
                          </span>
                        </div>
                      </div>
                      <p className="text-sm">{msg.summary || msg.text}</p>
                    </div>
                  ))}
                  <div ref={messagesEndRef} />
                </div>
              )}
            </div>

            {/* Send message form */}
            <div className="border-t p-4">
              {error && (
                <div className="mb-2 text-sm text-destructive">{error}</div>
              )}
              <div className="flex gap-2">
                <input
                  type="text"
                  placeholder="Recipient (agent name)"
                  value={recipient}
                  onChange={(e) => setRecipient(e.target.value)}
                  className="flex-1 rounded-md border bg-background px-3 py-2 text-sm"
                />
                <input
                  type="text"
                  placeholder="Message..."
                  value={messageText}
                  onChange={(e) => setMessageText(e.target.value)}
                  onKeyDown={(e) => e.key === 'Enter' && handleSendMessage()}
                  className="flex-1 rounded-md border bg-background px-3 py-2 text-sm"
                />
                <Button onClick={handleSendMessage} disabled={isLoading || !messageText.trim() || !recipient.trim()}>
                  {isLoading ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <Send className="h-4 w-4" />
                  )}
                </Button>
              </div>
            </div>
          </>
        ) : (
          <div className="flex flex-1 items-center justify-center text-muted-foreground">
            <p>Select a member to view their inbox</p>
          </div>
        )}
      </div>
    </div>
  )
}
