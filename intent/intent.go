package intent

type Intent int

var (
	Guilds Intent = 1 << 0
	// GuildMembers is a privileged intent.
	GuildMembers      Intent = 1 << 1
	GuildBans         Intent = 1 << 2
	GuildEmojis       Intent = 1 << 3
	GuildIntegrations Intent = 1 << 4
	GuildWebhooks     Intent = 1 << 5
	GuildInvites      Intent = 1 << 6
	GuildVoiceStates  Intent = 1 << 7
	// GuildPresences is a privileged intent
	GuildPresences         Intent = 1 << 8
	GuildMessages          Intent = 1 << 9
	GuildMessageReactions  Intent = 1 << 10
	GuildMessageTyping     Intent = 1 << 11
	DirectMessages         Intent = 1 << 12
	DirectMessageReactions Intent = 1 << 13
	DirectMessageTyping    Intent = 1 << 14

	// All
	All Intent = 32767
)
