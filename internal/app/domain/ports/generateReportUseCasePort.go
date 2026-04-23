package ports

type GenerateReportUseCasePort interface {
	Execute(packageId uint) error
}
