package copydb

// Stats contains basic metrics of the database.
type Stats struct {
	ItemsApplied           int
	ItemsFailed            int
	ItemsEvicted           int
	ItemsReplicated        int
	VersionConfictDetected int
	DBScanned              int
}
