package actor

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type SupervisorStrategy interface {
	HandleFailure(actorSystem *ActorSystem, supervisor Supervisor, child *PID, rs *RestartStatistics, reason interface{}, message interface{})
}

type SupervisorStrategyFunc func(actorSystem *ActorSystem, supervisor Supervisor, child *PID, rs *RestartStatistics, reason interface{}, message interface{})

func (r SupervisorStrategyFunc) HandleFailure(actorSystem *ActorSystem, supervisor Supervisor, child *PID, rs *RestartStatistics, reason interface{}, message interface{}) {
	r(actorSystem, supervisor, child, rs, reason, message)
}

type Supervisor = actor.Supervisor

type PID struct {
	context actor.SenderContext
	*actor.PID
}

func (r *PID) Tell(message interface{}) {
	r.context.Send(r.PID, message)
}

type Future = actor.Future

type protoStopperPart interface {
	// Stop will stop actor immediately regardless of existing user messages in mailbox.
	Stop(pid *actor.PID)

	// StopFuture will stop actor immediately regardless of existing user messages in mailbox, and return its future.
	StopFuture(pid *actor.PID) *Future

	// Poison will tell actor to stop after processing current user messages in mailbox.
	Poison(pid *actor.PID)

	// PoisonFuture will tell actor to stop after processing current user messages in mailbox, and return its future.
	PoisonFuture(pid *actor.PID) *Future
}

func (r *PID) StopFuture() *Future {
	return r.context.(protoStopperPart).StopFuture(r.PID)
}

func (r *PID) Stop() {
	r.context.(protoStopperPart).Stop(r.PID)
}

func (r *PID) GracefulStop() {
	r.context.(protoStopperPart).StopFuture(r.PID).Wait()
}

func (r *PID) RequestFuture(message interface{}, timeout time.Duration) *Future {
	return r.context.RequestFuture(r.PID, message, timeout)
}

type RestartStatistics = actor.RestartStatistics

func DefaultSupervisorStrategy() SupervisorStrategy {
	f := actor.DefaultSupervisorStrategy()
	return SupervisorStrategyFunc(func(actorSystem *ActorSystem, supervisor Supervisor, child *PID, rs *RestartStatistics, reason, message interface{}) {
		f.HandleFailure(actorSystem.ActorSystem, supervisor, child.PID, rs, reason, message)
	})

}

var UnwrapEnvelopeMessage = actor.UnwrapEnvelopeMessage

type MessageEnvelope = actor.MessageEnvelope

type Terminated struct {
	Who               *PID
	AddressTerminated bool
}

func (*Terminated) SystemMessage() {}

type Started = actor.Started

type Stopping = actor.Stopping

type Restarting = actor.Restarting

type Stopped = actor.Stopped

type DeadLetterEvent = actor.DeadLetterEvent

type Producer func() Actor

var ErrTimeout = actor.ErrTimeout

func FromProducer(producer Producer) *Props {
	props := actor.PropsFromProducer(func() actor.Actor {
		a := producer()
		return &actorWrapper{a}
	})
	return &Props{props.WithReceiverMiddleware(messageConverter)}
}

func FromFunc(f ReceiveFunc) *Props {
	props := FromProducer(func() Actor {
		return f
	})
	return props

}

type ReceiveTimeout = actor.ReceiveTimeout

type SystemMessage = actor.SystemMessage
type AutoReceiveMessage = actor.AutoReceiveMessage
