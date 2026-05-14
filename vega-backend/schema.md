# Schema

Conversation: id, state[idle, processing, completed, errored], type[task, conversation], created_at, updated_at

Context: id, conversation_id, timestamp

Message: serial id ,context_id, role, content, timestamp, worker_id

JobQueue: id, type, content, timestamp, worker_id, state[pending, processing, processed, errored], locked_until

Device: id, context, deviceToken

Workers: id

OpenSockets: id, device_id, worker_id

DeviceActionQueue: id, device_id,
