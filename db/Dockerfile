FROM redis:alpine

# Copy custom redis configuration file
COPY redis.conf /usr/local/etc/redis/redis.conf

# Expose Redis port
EXPOSE 6379

# Run Redis server with custom configuration
CMD ["redis-server", "/usr/local/etc/redis/redis.conf"]
