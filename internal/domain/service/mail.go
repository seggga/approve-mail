package service

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"text/template"

	"github.com/seggga/approve-mail/internal/domain/models"
)

var (
	sender      = "approver_service@corp.com"
	approveText = `Dear approver, 

there is a task you are to approve.
Task: {{.TaskName}} ({{.TaskID}})
Description: {{.TaskDescr}}

To approve the task please follow link {{.AcceptLink}}
To decline the task choose {{.DeclineLink}}

Best wishes, 
approvers service.`

	declText = `Dear approver, 

there is a task that has been declined:
Task: {{.TaskName}} ({{.TaskID}})
Description: {{.TaskDescr}}

Best wishes, 
approvers service.`

	approveTmpl, declTmpl *template.Template
)

// Service implements main analytics logic
type Service struct {
	chanIn  chan models.MsgTask
	chanOut chan models.MsgSMTP
	workers int
}

// New creates a new auth service
func New(chIn chan models.MsgTask, chOut chan models.MsgSMTP, wrk int) (*Service, error) {
	var err error
	approveTmpl, err = template.New("approveReq").Parse(approveText)
	if err != nil {
		return nil, fmt.Errorf("error parsing smtp-message body template: %w", err)
	}
	declTmpl, err = template.New("declineReq").Parse(declText)
	if err != nil {
		return nil, fmt.Errorf("error parsing smtp-message body template: %w", err)
	}

	return &Service{
		chanIn:  chIn,
		chanOut: chOut,
		workers: wrk,
	}, nil
}

// Start reads MsgTask from chanIn, composes MsgSMTP and sends them
// to chanOut. The number of simultaneously working go-routines
// is defined by Service.workers
func (s *Service) Start(ctx context.Context) error {
	wg := new(sync.WaitGroup)
	for i := 0; i < s.workers; i += 1 {
		wg.Add(1)
		go s.compose(ctx, wg)
	}
	wg.Wait()
	return nil
}

func (s *Service) compose(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return

		case taskMsg := <-s.chanIn:
			var tmpl *template.Template
			switch taskMsg.EventType {

			case models.Approval:
				tmpl = approveTmpl

			case models.Decline:
				tmpl = declTmpl
			}

			if tmpl == nil {
				// incorrect MsgTask.EventType recieved: expected approval / decine
				continue
			}

			smtp := &models.MsgSMTP{
				Sender:   sender,
				Receiver: taskMsg.Approver,
			}
			var body bytes.Buffer
			err := tmpl.Execute(&body, taskMsg)
			if err != nil {
				// s.Logger
				continue
			}
			smtp.Body = body.Bytes()
			s.chanOut <- *smtp
		}
	}
}
