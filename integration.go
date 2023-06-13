package scalr

type IntegrationStatus string

const (
	IntegrationStatusActive   IntegrationStatus = "active"
	IntegrationStatusDisabled IntegrationStatus = "disabled"
	IntegrationStatusFailed   IntegrationStatus = "failed"
)
