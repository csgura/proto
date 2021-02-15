package actor

import (
	"github.com/AsynkronIT/protoactor-go/actor"
)

type RootContext struct {
	*actor.RootContext
}

func messageConverter(next actor.ReceiverFunc) actor.ReceiverFunc {
	return func(c actor.ReceiverContext, envelope *actor.MessageEnvelope) {
		switch msg := envelope.Message.(type) {
		case *actor.Terminated:
			envelope.Message = &Terminated{
				Who: &PID{
					context: c.(actor.SenderContext),
					PID:     msg.Who,
				},
				AddressTerminated: msg.AddressTerminated,
			}
			next(c, envelope)
			return
		}
		next(c, envelope)
	}
}

func (rc *RootContext) SpawnNamed(props *Props, name string) (*PID, error) {
	pid, err := rc.RootContext.SpawnNamed(props.Props, name)
	if pid == nil {
		return nil, err
	}

	return &PID{
		rc.RootContext,
		pid,
	}, err
}

type ActorSystem struct {
	*actor.ActorSystem
	Root *RootContext
}

func NewActorSystem() *ActorSystem {
	ret := actor.NewActorSystem()

	return wrapActorSystem(ret)
}

func wrapActorSystem(asys *actor.ActorSystem) *ActorSystem {

	return &ActorSystem{ActorSystem: asys, Root: &RootContext{asys.Root}}
}
