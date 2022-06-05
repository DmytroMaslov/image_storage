package config

const (
	storageFolder         = "storage"
	QueueName             = "images"
	QueueConnectionString = "amqp://guest:guest@localhost:5672/"
	ServerPort            = ":8080"
	WorkerPoolSize        = 1
	Quality_100           = 100
	Quality_75            = 75
	Quality_50            = 50
)

type Config struct {
	StorageFolderName     string
	QueueName             string
	QueueConnectionString string
	ServerPort            string
	WorkerPoolSize        int
	InitialQuality        int
	QualityArray          []int
}

func GetConfig() *Config {
	return &Config{
		StorageFolderName:     storageFolder,
		QueueName:             QueueName,
		QueueConnectionString: QueueConnectionString,
		ServerPort:            ServerPort,
		WorkerPoolSize:        WorkerPoolSize,
		InitialQuality:        Quality_100,
		QualityArray:          []int{Quality_75, Quality_50},
	}
}
