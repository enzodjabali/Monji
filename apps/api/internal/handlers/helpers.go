package handlers

// joinUpdates is a helper to join SQL update fields with a given separator.
func joinUpdates(updates []string, sep string) string {
	out := ""
	for i, u := range updates {
		if i > 0 {
			out += sep
		}
		out += u
	}
	return out
}
