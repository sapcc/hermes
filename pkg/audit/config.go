package audit

//CADFConfiguration contains configuration parameters for audit trail.
type CADFConfiguration struct {
	Enabled  bool `yaml:"enabled"`
	RabbitMQ struct {
		URL       string `yaml:"url"`
		QueueName string `yaml:"queue_name"`
	} `yaml:"rabbitmq"`
}