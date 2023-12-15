package models

// Конфиг для сборки строки подключения.
type ConfigDB struct {
	Host     string
	Port     string
	UserName string
	Password string
	DbName   string
	SslMode  string
}

// Конфиг для подключение к консьюмеру натс.
type ConfigNatsPub struct {
	NatsURL        string
	NatsCluster    string
	NatsClient     string
	NatsSubjectPub string
}

// Конфиг для подключения к подписчику натс.
type ConfigNatsSub struct {
	NatsURL        string
	NatsCluster    string
	NatsClient     string
	NatsSubjectSub string
}
