package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/Abdillah-Epi/fizz-buzz/internal/config"
	"github.com/ClickHouse/clickhouse-go/v2"
)

type Client struct {
	conn clickhouse.Conn
}

func New(cfg config.ClickHouseConfig) (*Client, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{
			fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.Username,
			Password: cfg.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		return nil, fmt.Errorf("open clickhouse connection: %w", err)
	}

	client := &Client{
		conn: conn,
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
		conn_err := conn.Close()
		if conn_err != nil {
			return nil, fmt.Errorf("ping clickhouse: %w close conn: %w", err, conn_err)
		}

		return nil, fmt.Errorf("ping clickhouse: %w", err)
	}

	return client, nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.conn.Ping(ctx)
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Conn() clickhouse.Conn {
	return c.conn
}
