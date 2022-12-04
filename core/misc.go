package core

/**
 * Next r/w operation data counters
 */
type ReadWriteCount struct {

	/* Read bytes count */
	CountRead uint

	/* Write bytes count */
	CountWrite uint

	Target Upstream
}

func (rw ReadWriteCount) IsZero() bool {
	return rw.CountRead == 0 && rw.CountWrite == 0
}
