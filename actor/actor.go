package actor

import "github.com/AsynkronIT/protoactor-go/actor"

type Actor interface {
	Receive(c Context)
}

type actorWrapper struct {
	actor Actor
}

func (r *actorWrapper) Receive(c actor.Context) {
	r.actor.Receive(ContextWrapper{c})
}

// The ReceiveFunc type is an adapter to allow the use of ordinary functions as actors to process messages
type ReceiveFunc func(c Context)

// Receive calls f(c)
func (f ReceiveFunc) Receive(c Context) {
	f(c)
}

type ReceiverFunc func(c ReceiverContext, envelope *MessageEnvelope)

type ReceiverMiddleware func(next ReceiverFunc) ReceiverFunc

type Props struct {
	*actor.Props
}

func (props *Props) WithReceiverMiddleware(middleware ...ReceiverMiddleware) *Props {

	m := make([]actor.ReceiverMiddleware, len(middleware))

	for i := range middleware {
		m[i] = func(next actor.ReceiverFunc) actor.ReceiverFunc {

			cn := func(c ReceiverContext, envelope *MessageEnvelope) {
				next(c.protoReceiverContext(), envelope)
			}

			ret := middleware[i](cn)

			return func(c actor.ReceiverContext, envelope *actor.MessageEnvelope) {

				if fc, ok := c.(actor.Context); ok {
					ret(ContextWrapper{fc}, envelope)
				} else {
					ret(ReceiverContextWrapper{c}, envelope)
				}
			}
		}
	}

	return &Props{props.Props.WithReceiverMiddleware(m...)}
}

type ActorFunc = ReceiveFunc
