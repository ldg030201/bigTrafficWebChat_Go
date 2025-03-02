package repository

import (
	"chat_server/config"
	"chat_server/repository/kafka"
	"chat_server/types/schema"
	"database/sql"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	cfg *config.Config

	db *sql.DB

	Kafka *kafka.Kafka
}

const (
	room       = " bigTrafficWebChat.room "
	chat       = " bigTrafficWebChat.chat "
	serverInfo = " bigTrafficWebChat.serverInfo "
)

func (s *Repository) ServerSet(ip string, available bool) error {
	_, err := s.db.Exec("INSERT INTO serverInfo(`ip`, `available`) VALUES(?, ?) ON DUPLICATE KEY UPDATE `available` = VALUES(`available`)", ip, available)
	return err
}

func NewRepository(cfg *config.Config) (*Repository, error) {
	r := &Repository{cfg: cfg}
	var err error

	if r.db, err = sql.Open(cfg.DB.Database, cfg.DB.URL); err != nil {
		return nil, err
	} else if r.Kafka, err = kafka.NewKafka(cfg); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (s *Repository) InsertChatting(user, message, roomName string) error {
	log.Println("InsertChatting 실행")

	_, err := s.db.Exec("INSERT INTO bigTrafficWebChat.chat(room, name, message) VALUES(?, ?, ?)", roomName, user, message)

	return err
}

func (s *Repository) GetChatList(roomName string) ([]*schema.Chat, error) {
	qs := query([]string{"SELECT * FROM", chat, "WHERE room = ? ORDER BY `when` DESC LIMIT 10"})

	if cursor, err := s.db.Query(qs, roomName); err != nil {
		return nil, err
	} else {
		defer cursor.Close()

		var result []*schema.Chat

		for cursor.Next() {
			d := new(schema.Chat)

			if err = cursor.Scan(
				&d.ID,
				&d.Room,
				&d.Name,
				&d.Message,
				&d.When,
			); err != nil {
				return nil, err
			} else {
				result = append(result, d)
			}
		}

		if len(result) == 0 {
			return []*schema.Chat{}, nil
		} else {
			return result, nil
		}
	}
}

func (s *Repository) MakeRoom(name string) error {
	_, err := s.db.Exec("INSERT INTO bigTrafficWebChat.room(name) VALUES(?)", name)
	return err
}

func (s *Repository) RoomList() ([]*schema.Room, error) {
	qs := query([]string{"SELECT * FROM", room})

	if cursor, err := s.db.Query(qs); err != nil {
		return nil, err
	} else {
		defer cursor.Close()

		var result []*schema.Room

		for cursor.Next() {
			d := new(schema.Room)

			if err = cursor.Scan(
				&d.ID,
				&d.Name,
				&d.CreateAt,
				&d.UpdateAt,
			); err != nil {
				return nil, err
			} else {
				result = append(result, d)
			}
		}

		if len(result) == 0 {
			return []*schema.Room{}, nil
		} else {
			return result, nil
		}
	}
}

func (s *Repository) Room(name string) (*schema.Room, error) {
	d := new(schema.Room)
	qs := query([]string{"SELECT * FROM", room, "WHERE name = ?"})

	err := s.db.QueryRow(qs, name).Scan(
		&d.ID,
		&d.Name,
		&d.CreateAt,
		&d.UpdateAt,
	)

	if err = noResult(err); err != nil {
		return nil, err
	} else {
		return nil, nil
	}
}

func query(qs []string) string {
	return strings.Join(qs, "") + ";"
}

func noResult(err error) error {
	if strings.Contains(err.Error(), "sql: no rows in result set") {
		return nil
	} else {
		return err
	}
}
