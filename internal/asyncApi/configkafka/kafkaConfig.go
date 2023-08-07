package configkafka

const (
	Host   = "bootstrap.servers"
	Group  = "group.id"
	Offset = "auto.offset.reset"
	Topic  = "kafka-topic-name"
)

var Kafka struct {
	Host   string "localhost:9092"
	Group  string "my-consumer-group"
	Offset string "earliest"
}
