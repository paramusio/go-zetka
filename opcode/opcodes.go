package opcode

type Opcode int

var (
	// Dispatch [Receive] An event was dispatched.
	Dispatch Opcode = 0

	// Heartbeat [Send/Receive] Fired periodically by the client to keep the connection alive.
	Heartbeat Opcode = 1
	// Identify [Send] Starts a new session during the initial handshake.
	Identify Opcode = 2

	// PresenceUpdate [Send] Update the client's presence.
	PresenceUpdate Opcode = 3

	// VoiceStateUpdate [Send] Used to join/leave or move between voice channels.
	VoiceStateUpdate Opcode = 4

	// ResumeSend [Send] Resume a previous session that was disconnected.
	ResumeSend Opcode = 6

	// Reconnect [Receive] You should attempt to reconnect and resume immediately.
	Reconnect Opcode = 7

	// RequestGuildMembers [Send] Request information about offline guild members in a large guild.
	RequestGuildMembers Opcode = 8

	// InvalidSession [Receive] The session has been invalidated. You should reconnect and identify/resume accordingly.
	InvalidSession Opcode = 9

	// Hello [Receive] Sent immediately after connecting, contains the heartbeat_interval to use.
	Hello Opcode = 10

	// HeartbeatAck [Receive] Sent in response to receiving a heartbeat to acknowledge that it has been received.
	HeartbeatACK Opcode = 11
)
