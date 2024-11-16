package redis

type UserCaching struct {
	redisClient *redis.Client
	logger      *slog.Logger
}
