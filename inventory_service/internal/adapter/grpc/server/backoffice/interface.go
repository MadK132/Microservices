package backoffice

type Server interface {
	Run(errCh chan<- error)
	Stop() error
}
