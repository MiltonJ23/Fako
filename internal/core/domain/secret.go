package domain

// Secret is a secure wrapper for managing our sensible data
type Secret struct {
	value []byte
}

func NewSecret(s string) *Secret {
	return &Secret{value: []byte(s)}
}

// Reveal will expose the value for an immediate use
func (s *Secret) Reveal() []byte {
	return s.value
}

// Wipe will erase the value by crushing it with zeros
func (s *Secret) Wipe() {
	// If the value is already nil, there no wipe to do
	if s.value == nil {
		return
	} else {
		for i := range s.value {
			s.value[i] = 0
		}
	}
}

// now let's write the string method and ensure the value is not logged by accident

func (s *Secret) String() string {
	return "******NOTHING TO SEE HERE******"
}
