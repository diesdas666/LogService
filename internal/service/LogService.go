package service

//go:generate asyncApiCodeGen generate -c true -i ../../apis/asyncApi/identity-service.yaml -o ../../gensrc/asyncApi/events/events.gen.go -p events
type LogService interface {
	UpdateStatusOfLoadingStation(ev Event)
}
