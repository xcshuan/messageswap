package messageswap_pb

// Loggable turns a Message into machine-readable log output
func (m *Message) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"message": map[string]string{
			"type": Message_MessageType_name[int32(m.GetType())],
		},
	}
}
