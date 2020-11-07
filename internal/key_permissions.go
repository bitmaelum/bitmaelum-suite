// Copyright (c) 2020 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package internal

const (
	// PermFlush Permission to restart/reload the system including flushing/forcing the queues
	PermFlush string = "flush"
	// PermGenerateInvites Permission to generate invites remotely
	PermGenerateInvites string = "invite"
	// PermAPIKeys Permission to create api keys
	PermAPIKeys string = "apikey"
	// PermGetHeaders allows you to fetch header and catalog from messages
	PermGetHeaders string = "get-headers"
)

// ManagementPermissions is a list of all permissions available for remote management
var ManagementPermissions = []string{
	PermAPIKeys,
	PermFlush,
	PermGenerateInvites,
}

// AccountPermissions is a set of permissions for specific accounts
var AccountPermissions = []string{
	PermGetHeaders,
}
