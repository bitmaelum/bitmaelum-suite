// Copyright (c) 2021 BitMaelum Authors
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

package account

// Folder constants
const (
	BoxInbox  = iota + 1 // BoxInbox is the mandatory inbox where all incoming messages are stored
	BoxOutbox            // BoxOutbox is the mandatory outbox where send messages are stored
	BoxTrash             // BoxTrash is the mandatory trashcan where deleted messages are stored (before actual deletion)
)

// MandatoryBoxes is a list of all boxes that are mandatory. Makes it easier to range on them
var MandatoryBoxes = []int{
	BoxInbox,
	BoxOutbox,
	BoxTrash,
}

// MaxMandatoryBoxID is the largest box that must be present. Everything below this box (including this box) is
// mandatory and cannot be removed.
const MaxMandatoryBoxID = 99
