package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/pavel/message_service/pkg/db"
	"github.com/pavel/message_service/pkg/logger"
	"github.com/pavel/message_service/pkg/model"
	"time"
)

type Message interface {
	CreateChatMessage(senderUserId, chatId int, body string) (int, error)
	All(filterMessage model.FilterMessage) ([]*model.MessageChat, error)
}

type MessagePostgreSQL struct {
	db  *db.DB
	log *logger.Logger
}

func InitMessagePostgreSQL(db *db.DB, log *logger.Logger) *MessagePostgreSQL {
	return &MessagePostgreSQL{
		db:  db,
		log: log,
	}
}

func (m MessagePostgreSQL) All(filterMessage model.FilterMessage) ([]*model.MessageChat, error) {

	isReadMessageSQL := fmt.Sprintf("(SELECT CASE WHEN "+
		"(SELECT count(uum.id) from %s uum LEFT JOIN %s uc on uum.chat_id = uc.chat_id WHERE uum.chat_id = $1 AND uum.message_id = mc.id AND uc.user_id <> $2)"+
		" = "+
		"(SELECT count(id) as count_user FROM users_chats uc WHERE uc.chat_id = $1 AND uc.user_id <> $2)"+
		" THEN false ELSE true END) as is_read", model.USER_UNREAD_MESSAGE, model.USERS_CHATS_TABLE)
	messagesRows, err := m.db.Query(fmt.Sprintf("SELECT mc.id, mc.chat_id, mc.sender_user_id, mc.body, mc.created_at, mc.updated_at, %s "+
		"FROM %s mc WHERE chat_id = $1 ORDER BY created_at DESC LIMIT %d OFFSET %d",
		isReadMessageSQL, model.MESSAGE_CHAT_TABLE, filterMessage.Limit, filterMessage.Offset), filterMessage.ChatId, filterMessage.UserId)
	var messages []*model.MessageChat
	for messagesRows.Next() {
		var message model.MessageChat
		err = messagesRows.Scan(&message.ID, &message.ChatId, &message.SenderUserId, &message.Body, &message.CreatedAt, &message.UpdatedAt, &message.IsRead)
		if err != nil {
			m.log.Warnf("Can't scan message. Err %v", err.Error())
			return nil, err
		}
		messages = append(messages, &message)
	}
	return messages, nil
}

func (m MessagePostgreSQL) CreateChatMessage(senderUserId, chatId int, body string) (int, error) {
	if !m.checkIfChatExist(chatId) {
		return 0, errors.New("not found chat")
	}
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		m.log.Warnf("Transaction error: %v", err)
		return 0, err
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
		m.log.Warnf("%s insert error: %v", model.MESSAGE_CHAT_TABLE, err)
		tx.Rollback()
		return 0, err
	}
	err = m.addUnreadMessage(ctx, tx, chatId, senderUserId, messageId)
	if err != nil {
		tx.Rollback()
		m.log.Warnf("Can't create unread message")
	}
	tx.Commit()
	m.log.Printf("Created message. Message id: %d", messageId)
	return messageId, nil
}

func (m MessagePostgreSQL) checkIfChatExist(chatId int) bool {
	var chat bool
	err := m.db.QueryRow(fmt.Sprintf("SELECT exists(SELECT id FROM %s WHERE id = $1)", model.CHATS_TABLE), chatId).Scan(&chat)
	if err != nil || chat == false {
		return false
	}
	return true
}

func (m MessagePostgreSQL) addUnreadMessage(ctx context.Context, tx *sql.Tx, chatId, senderUserId, messageId int) error {
	//get user from chats
	var receiversUsersIds []int
	chatUsers, err := tx.QueryContext(ctx, fmt.Sprintf("SELECT user_id FROM %s WHERE chat_id = $1 AND user_id <> $2", model.USERS_CHATS_TABLE), chatId, senderUserId)
	if chatUsers.Next() {
		var userId int
		err := chatUsers.Scan(&userId)
		if err != nil {
			m.log.Warnf("Cannot scan userId in table %s, err: %v", model.USERS_CHATS_TABLE, err)
			return err
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
		m.log.Warnf("%s insert error %v", model.USER_UNREAD_MESSAGE, err)
		return err
	}
	return nil
}
