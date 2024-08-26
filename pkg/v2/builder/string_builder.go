package builder

type StringBuilder interface {
	Append(value string) StringBuilder
	AppendN(value string, repeatCount int) StringBuilder
	AppendLine() StringBuilder
	Remove(start, length int) StringBuilder
	String() string
	GetLength() int
}

type stringBuilder struct {
	value string
}

func NewStringBuilder() StringBuilder {
	return &stringBuilder{}
}

func (s *stringBuilder) Append(value string) StringBuilder {
	s.value += value
	return s
}

func (s *stringBuilder) AppendLine() StringBuilder {
	s.value += "\n"
	return s
}

func (s *stringBuilder) AppendN(value string, repeatCount int) StringBuilder {
	for i := 0; i < repeatCount; i++ {
		s.Append(value)
	}
	return s
}

func (s *stringBuilder) GetLength() int {
	return len(s.value)
}

func (s *stringBuilder) Remove(start int, length int) StringBuilder {
	s.value = s.value[0:start] + s.value[start+length:]
	return s
}

func (s *stringBuilder) String() string {
	return s.value
}
