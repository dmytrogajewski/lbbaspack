package components

// BackendAssignment component assigns an entity to a backend
// and tracks backend-specific data
type BackendAssignment struct {
	BackendID int
	Counter   int
}

func NewBackendAssignment(id int) *BackendAssignment {
	return &BackendAssignment{BackendID: id, Counter: 0}
}

// GetType implements Component interface
func (ba *BackendAssignment) GetType() string {
	return "BackendAssignment"
}

// GetBackendID implements BackendAssignmentComponent interface
func (ba *BackendAssignment) GetBackendID() int {
	return ba.BackendID
}

// GetAssignedPackets implements BackendAssignmentComponent interface
func (ba *BackendAssignment) GetAssignedPackets() int {
	return ba.Counter
}

// SetBackendID implements BackendAssignmentComponent interface
func (ba *BackendAssignment) SetBackendID(id int) {
	ba.BackendID = id
}

// IncrementAssignedPackets implements BackendAssignmentComponent interface
func (ba *BackendAssignment) IncrementAssignedPackets() {
	ba.Counter++
}
