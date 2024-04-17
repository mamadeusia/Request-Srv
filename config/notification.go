package config

type Notification struct {
	Name string
}

func GetNotificationServiceName() string {
	return cfg.Notification.Name
}
