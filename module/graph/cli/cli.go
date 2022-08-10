package cli

type Service interface {
	CliLink(applicationId int) (int, error)
}
