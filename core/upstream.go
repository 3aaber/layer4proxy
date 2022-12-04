package core

/**
 * Upstream host and port
 */
type Upstream struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

/**
 * Compare to other target
 */
func (t *Upstream) EqualTo(other Upstream) bool {
	return t.Host == other.Host &&
		t.Port == other.Port
}

/**
 * Get target full address
 * host:port
 */
func (t *Upstream) Address() string {
	return t.Host + ":" + t.Port
}

/**
 * To String conversion
 */
func (t *Upstream) String() string {
	return t.Address()
}
