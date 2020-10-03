package message

type User struct {
	Username      string    `json:"username"`
	MFAEnabled    bool      `json:"mfa_enabled"`
	Verified      bool      `json:"verified"`
	ID            Snowflake `json:"id"`
	Flags         int64     `json:"flags"`
	Email         string    `json:"email"`
	Discriminator string    `json:"discriminator"`
	Bot           bool      `json:"bot"`
	Avatar        string    `json:"avatar"`
}
