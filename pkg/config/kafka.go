package config

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}
