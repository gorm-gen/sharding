package list

type Option func(*List)

func WithOffset(offset uint64) Option {
	return func(l *List) {
		l.offset = int64(offset)
	}
}

func WithPage(page uint64) Option {
	return func(l *List) {
		l.page = int64(page)
	}
}

func WithPageSize(pageSize uint64) Option {
	return func(l *List) {
		l.pageSize = int64(pageSize)
	}
}

func WithDesc() Option {
	return func(l *List) {
		l.desc = true
		l.asc = false
	}
}

func WithAsc() Option {
	return func(l *List) {
		l.asc = true
		l.desc = false
	}
}
