package forum

import desc "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"

type Implementation struct {
	desc.UnimplementedForumServer
}


// Opts - configuration service dependencies
type Opts struct {
}

// NewForum return new instance of Implementation.
func NewForum(opts Opts) desc.ForumServer {
	return &Implementation{}
}
