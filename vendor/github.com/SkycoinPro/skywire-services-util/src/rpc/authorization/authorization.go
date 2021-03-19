package authorization

// GetRequest contains Username for User we want to receive access Rights
type GetRequest struct {
	Username string
}

// GetResponse contains returned Rights after GetRequest has been made
type GetResponse struct {
	Rights []Right
}

// SetRequest that will update Rights for the User with given Username
type SetRequest struct {
	Username string
	Rights   []Right
}

// SetResponse received after SetRequest has been made
type SetResponse struct {
}

// Right definition holding the value, name
type Right struct {
	Name  string
	Label string
	Value bool
}
