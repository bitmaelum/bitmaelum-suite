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
	// PermGetHeaders allows you to fetch header and catalog from messages
	PermGetHeaders string = "get-headers"
)

// ManagementPermissons is a list of all permissions available for remote management
var ManagementPermissons = []string{
	PermAPIKeys,
	PermFlush,
	PermGenerateInvites,
}

// AccountPermissions is a set of permissions for specific accounts
var AccountPermissions = []string{
	PermGetHeaders,
	PermMail,
}
