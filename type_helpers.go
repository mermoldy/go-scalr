package scalr

// Bool returns a pointer to the given bool
func Bool(v bool) *bool {
	return &v
}

// Category returns a pointer to the given category type.
func Category(v CategoryType) *CategoryType {
	return &v
}

// EnforcementMode returns a pointer to the given enforcement level.
func EnforcementMode(v EnforcementLevel) *EnforcementLevel {
	return &v
}

// Int returns a pointer to the given int.
func Int(v int) *int {
	return &v
}

// Int64 returns a pointer to the given int64.
func Int64(v int64) *int64 {
	return &v
}

// NotificationDestination returns a pointer to the given notification configuration destination type
func NotificationDestination(v NotificationDestinationType) *NotificationDestinationType {
	return &v
}

// PlanExportType returns a pointer to the given plan export data type.
func PlanExportType(v PlanExportDataType) *PlanExportDataType {
	return &v
}

// String returns a pointer to the given string.
func String(v string) *string {
	return &v
}
