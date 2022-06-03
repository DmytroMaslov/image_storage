package config

const (
	storageFolder         = "storage"
	QueueName             = "images"
	QueueConnectionString = "amqp://guest:guest@localhost:5672/"
	ServerPort            = ":8080"
)

type Config struct {
	StorageFolderName     string
	QueueName             string
	QueueConnectionString string
	ServerPort            string
}

func GetConfig() *Config {
	return &Config{
		StorageFolderName:     storageFolder,
		QueueName:             QueueName,
		QueueConnectionString: QueueConnectionString,
		ServerPort:            ServerPort,
	}
}
