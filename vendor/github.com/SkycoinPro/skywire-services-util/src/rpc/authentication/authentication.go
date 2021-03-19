package authentication

// Request for authentication between the services
type Request struct {
	Username string
	Password string
	OtpToken string
}

// Response for authentication request between the services
type Response struct {
	Success bool
}
