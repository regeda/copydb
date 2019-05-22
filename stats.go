package copydb

// Stats contains basic metrics of the database.
type Stats struct {
	ItemsApplied           int
	ItemsFailed            int
	ItemsEvicted           int
	VersionConfictDetected int
	DBScanned              int
}
