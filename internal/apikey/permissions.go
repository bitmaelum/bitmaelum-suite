package apikey

const (
	// PermFlush Permission to restart/reload the system including flushing/forcing the queues
	PermFlush string = "flush"
	// PermGenerateInvites Permission to generate invites remotely
	PermGenerateInvites string = "invite"
	// PermAPIKeys Permission to create api keys
	PermAPIKeys string = "apikey"
	// PermMail Permission to send email
	PermMail string = "mail"
)

// AllPermissons is a list of all permissions available for API keys
var AllPermissons = []string{
	PermAPIKeys,
	PermFlush,
	PermGenerateInvites,
	PermMail,
}
