package audit

//Config contains configuration parameters for audit trail.
type Config struct {
	Enabled  bool `yaml:"enabled"`
	RabbitMQ struct {
		URL       string `yaml:"url"`
		QueueName string `yaml:"queue_name"`
	} `yaml:"rabbitmq"`

	// A single, consistent UUID for an OpenStack Service.
	ObserverUUID string `yaml:"observer_uuid"`
}
