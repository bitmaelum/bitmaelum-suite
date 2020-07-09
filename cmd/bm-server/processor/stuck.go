package processor

// ProcessStuckIncomingMessages will process stuck message found in the incoming queue.
func ProcessStuckIncomingMessages() {
	// How do we know if something in incoming is stuck?
	//  Time? not really. It can take a while to upload everything
	//  proof-of-work expired. Same issue
	//  .completed file in the directory?
	// @TODO implement
}

// ProcessStuckProcessingMessages will process stuck message found in the processing queue.
func ProcessStuckProcessingMessages() {
	// If the message is NOT found in the processingList, then it's a message that is currently not being processed and we can process the message.
	// @TODO implement
}
