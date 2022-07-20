package repository

import (
	"context"
	"fmt"
	"github.com/pavel/message_service/pkg/db"
	"github.com/pavel/message_service/pkg/logger"
	"github.com/pavel/message_service/pkg/model"
	"log"
	"time"
)

type Message interface {
	CreateMessage(senderUserId, chatId int, body string) error
	GetAll(filterMessage model.FilterMessage) ([]*model.MessageChat, error)
}

type MessagePostgreSQL struct {
	db  *db.DB
	log *logger.Logger
}

func InitMessagePostgeSQL(db *db.DB, log *logger.Logger) *MessagePostgreSQL {
	return &MessagePostgreSQL{
		db:  db,
		log: log,
	}
}

func (m MessagePostgreSQL) GetAll(filterMessage model.FilterMessage) ([]*model.MessageChat, error) {

	messagesRows, err := m.db.Query(fmt.Sprintf("SELECT id, chat_id, sender_user_id, body, created_at, updated_at FROM %s WHERE chat_id = $1 ORDER BY created_at DESC LIMIT %d OFFSET %d",
		model.MESSAGE_CHAT_TABLE, filterMessage.Limit, filterMessage.Offset), filterMessage.ChatId)
	var messages []*model.MessageChat
	if messagesRows.Next() {
		var message model.MessageChat
		err = messagesRows.Scan(&message.ID, &message.ChatId, &message.SenderUserId, &message.Body, &message.CreatedAt, &message.UpdatedAt)
		if err != nil {
			log.Fatalf("Can't scan message. Err %v", err.Error())
			return nil, err
		}
		messages = append(messages, &message)
	}
	return messages, nil
}

func (m MessagePostgreSQL) CreateMessage(senderUserId, chatId int, body string) error {
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		m.log.Fatalf("Transaction error: %v", err)
		return err
	}
	dateNow := time.Now()

	//create message
	messageSQL, messageArguments := m.db.QueryBuilder.NewQueryBuilder().
		Insert(fmt.Sprintf("%s (chat_id, sender_user_id, body, created_at, updated_at)", model.MESSAGE_CHAT_TABLE)).Values("").
		AddValue("(@, @, @, @, @)", chatId, senderUserId, body, dateNow, dateNow).
		ToSql()
	var messageId int
	err = tx.QueryRowContext(ctx, messageSQL+" RETURNING id", messageArguments...).Scan(&messageId)
	if err != nil {
		m.log.Fatalf("%s insert error: %v", model.MESSAGE_CHAT_TABLE, err)
		tx.Rollback()
		return err
	}

	//get user from chats
	var receiversUsersIds []int
	chatUsers, err := tx.QueryContext(ctx, fmt.Sprintf("SELECT user_id FROM %s WHERE chat_id = $1 AND user_id <> $2", model.USERS_CHATS_TABLE), chatId, senderUserId)
	if chatUsers.Next() {
		var userId int
		err := chatUsers.Scan(&userId)
		if err != nil {
			log.Fatalf("Cannot scan userId in table %s, err: %v", model.USERS_CHATS_TABLE, err)
			tx.Rollback()
			return nil
		}
		receiversUsersIds = append(receiversUsersIds, userId)
	}
	chatUsers.Close()

	//add unread message for users
	userUnreadMessageBuilder := m.db.QueryBuilder.NewQueryBuilder().
		Insert(fmt.Sprintf("%s (chat_id, user_id, message_id)", model.USER_UNREAD_MESSAGE)).Values("")

	for count, receiverUserId := range receiversUsersIds {
		lastSymbol := ","
		if count == len(receiversUsersIds)-1 {
			lastSymbol = ";"
		}
		userUnreadMessageBuilder.AddValue("(@, @, @)"+lastSymbol, chatId, receiverUserId, messageId)
	}
	userUnreadMessageSQL, userUnreadMessageArguments := userUnreadMessageBuilder.ToSql()
	_, err = tx.Exec(userUnreadMessageSQL, userUnreadMessageArguments...)
	if err != nil {
		log.Fatalf("%s insert error %v", model.USER_UNREAD_MESSAGE, err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	m.log.Printf("Created message. Message id: %d", messageId)
	return nil
}
