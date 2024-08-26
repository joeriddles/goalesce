package builder

type StringBuilder interface {
	Append(value string, repeatCount int) StringBuilder
	AppendLine() StringBuilder
	Remove(start, length int)
	String() string
	GetLength() int
}
