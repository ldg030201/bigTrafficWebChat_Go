package service

import (
	"chat_server/repository"
	"chat_server/types/schema"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	s := &Service{repository: repository}

	return s
}

func (s *Service) ServerSet(ip string, available bool) error {
	if err := s.repository.ServerSet(ip, available); err != nil {
		log.Println("ServerSet err:", err)
		return err
	} else {
		return nil
	}
}

func (s *Service) EnterRoom(roomName string) ([]*schema.Chat, error) {
	if res, err := s.repository.GetChatList(roomName); err != nil {
		log.Println("GetChatList 에러", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}

func (s *Service) MakeRoom(name string) error {
	if err := s.repository.MakeRoom(name); err != nil {
		log.Println("MakeRoom 에러", "err", err.Error())
		return err
	} else {
		return nil
	}
}

func (s *Service) RoomList() ([]*schema.Room, error) {
	if res, err := s.repository.RoomList(); err != nil {
		log.Println("RoomList 에러", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}

func (s *Service) Room(name string) (*schema.Room, error) {
	if res, err := s.repository.Room(name); err != nil {
		log.Println("Room 에러", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}

func (s *Service) InsertChatting(user, message, roomName string) {
	if err := s.repository.InsertChatting(user, message, roomName); err != nil {
		log.Println("InsertChatting 에러", "err", err.Error())
	}

}

func (s *Service) PublishEvent(topic string, value []byte, ch chan kafka.Event) (kafka.Event, error) {
	return s.repository.Kafka.PublishEvent(topic, value, ch)
}

func (s *Service) PublishServerStatusEvent(ip string, status bool) {
	type ServerInfoEvent struct {
		IP     string
		Status bool
	}

	e := &ServerInfoEvent{IP: ip, Status: status}
	ch := make(chan kafka.Event)

	if v, err := json.Marshal(e); err != nil {
		log.Printf("json marshal error")
	} else if result, err := s.PublishEvent("chat", v, ch); err != nil {
		log.Println("end event kafka error")
	} else {
		log.Println("success event", result)
	}
}
