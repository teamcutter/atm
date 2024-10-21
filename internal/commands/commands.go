package commands

type Command int

const (
	SET Command = iota
	GET
	LPUSH
	RPUSH
	LPOP
)

func (c Command) String() string {
	return [...]string{"SET", "GET", "LPUSH", "RPUSH", "LPOP"}[c]
}
