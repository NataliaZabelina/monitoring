package monitoring

import "context"

type Monitoring struct {
	data   string
}

func New(data string) *Monitoring {
	return &Monitoring{data: data}
}


func (s *Monitoring) Run(ctx context.Context) error {
	return nil
}

