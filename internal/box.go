package internal

const (
	// BoxRoot is the root box when we want to add boxes without a parent ID
	BoxRoot   = iota
	// BoxInbox is the mandatory inbox where all incoming messages are stored
	BoxInbox  // Always box 1
	// BoxOutbox is the mandatory outbox where send messages are stored
	BoxOutbox // Always box 2
	// BoxTrash is the mandatory trashcan where deleted messages are stored (before actual deletion)
	BoxTrash  // Always box 3
)

// MaxMandatoryBoxID is the largest box that must be present. Everything below this box (including this box) is
// mandatory and cannot be removed.
const MaxMandatoryBoxID = 3
