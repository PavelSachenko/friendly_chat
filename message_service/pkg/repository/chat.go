package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/pavel/message_service/pkg/db"
	"github.com/pavel/message_service/pkg/logger"
	"github.com/pavel/message_service/pkg/model"
	"log"
	"time"
)

type Chat interface {
	CreatePrivate(userTitleChats []model.UserTitleChats) (int, error)
	CreateGroup(description, title string, UsersOwners []model.UsersOwners) (int, error)
	All(filterChat model.FilterChat) ([]*model.AllChat, error)
	GetChatUserIds(chatId int) ([]int, error)
}

type ChatPostgreSQL struct {
	db  db.DB
	log logger.Logger
}

func InitChatPostgreSQL(db db.DB, log logger.Logger) *ChatPostgreSQL {
	return &ChatPostgreSQL{
		db:  db,
		log: log,
	}
}

func (c *ChatPostgreSQL) GetChatUserIds(chatId int) ([]int, error) {
	rawSql := fmt.Sprintf("SELECT user_id FROM %s WHERE chat_id = $1", model.USERS_CHATS_TABLE)
	rows, err := c.db.Queryx(rawSql, chatId)
	if err != nil {
		c.log.Errorf("Bad query sql request: %v", err)
		return nil, err
	}
	var userIds []int
	for rows.Next() {
		var userId int
		err := rows.Scan(&userId)
		c.log.Errorf("Scan error: %s", err)
		userIds = append(userIds, userId)
	}
	return userIds, nil
}

func (c *ChatPostgreSQL) All(filterChat model.FilterChat) ([]*model.AllChat, error) {
	getStartPositition := fmt.Sprintf("(select count(cmc.id)-(select count(id) from %s cuum WHERE cuum.chat_id = c.id AND cuum.user_id = $1) FROM %s cmc WHERE cmc.chat_id = c.id) as start_position", model.USER_UNREAD_MESSAGE, model.MESSAGE_CHAT_TABLE)
	countUnreadMessageSubSql := fmt.Sprintf("(SELECT count(uum.id) as count FROM %s uum WHERE c.id = uum.chat_id AND uum.user_id = $1) as count_unread_message", model.USER_UNREAD_MESSAGE)
	lastMessageSubSql := fmt.Sprintf("(select json_build_object('body', mc.body, 'created_at', mc.created_at, 'is_sent', "+
		"(select exists(select id from %s WHERE mc.id = messages_chats.id AND messages_chats.sender_user_id = $1)))"+
		" FROM %s as mc WHERE mc.chat_id = c.id GROUP BY mc.id, mc.body, mc.created_at ORDER BY mc.created_at DESC LIMIT 1)", model.MESSAGE_CHAT_TABLE, model.MESSAGE_CHAT_TABLE)
	chatSql := fmt.Sprintf("SELECT c.id, lc.title, c.description, c.type, c.created_at, c.updated_at, %s, %s, %s as last_message FROM %s as c LEFT JOIN users_chats uc on c.id = uc.chat_id"+
		" LEFT JOIN labels_chats lc on c.id = lc.chat_id WHERE uc.user_id = $1 AND lc.user_id = $1 ORDER BY ", getStartPositition, countUnreadMessageSubSql, lastMessageSubSql, model.CHATS_TABLE)
	chatSql = fmt.Sprintf("%s %s->>'created_at' DESC LIMIT %d OFFSET %d", chatSql, lastMessageSubSql, filterChat.Limit, filterChat.Offset)
	rows, err := c.db.Query(chatSql, filterChat.UserId)
	if err != nil {
		c.log.Warnf("Get all chats. Err: %v", err)
		return nil, err
	}
	var allChats []*model.AllChat
	for rows.Next() {
		var allChat model.AllChat
		var lastMessageJSON *string
		err := rows.Scan(&allChat.ID, &allChat.Title, &allChat.Description, &allChat.Type, &allChat.CreatedAt, &allChat.UpdateAt, &allChat.StartPosition, &allChat.CountUnreadMessage, &lastMessageJSON)
		if err != nil {
			c.log.Warnf("Get all chats scan error. Err: %v", err)
			return nil, err
		}
		if lastMessageJSON != nil {
			err = json.Unmarshal([]byte(*lastMessageJSON), &allChat.LastMessage)
			if err != nil {
				c.log.Warnf("Can't unmarshal last_message. Err: %v", err)
				return nil, err
			}
		}
		allChats = append(allChats, &allChat)
	}

	return allChats, nil
}

func (c *ChatPostgreSQL) CreateGroup(description, title string, usersOwners []model.UsersOwners) (int, error) {
	ctx := context.Background()
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		c.log.Warnf("Transaction error: %v", err)
		return 0, err
	}
	chatId, err := c.createChat(ctx, tx, description, model.GROUP_CHAT)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	userQueryBuilder := c.db.QueryBuilder.NewQueryBuilder()
	labelChatSqlBuilder := c.db.QueryBuilder.Insert(fmt.Sprintf("%s (chat_id, user_id, title)", model.LABELS_CHATS_TABLE)).Values("")
	userChatSqlBuilder := userQueryBuilder.Insert(fmt.Sprintf("%s (chat_id, user_id, is_owner)", model.USERS_CHATS_TABLE)).Values("")
	for count, user := range usersOwners {
		lastSymbol := ","
		if count == len(usersOwners)-1 {
			lastSymbol = ";"
		}
		labelChatSqlBuilder.AddValue("(@, @, @)"+lastSymbol, chatId, user.UserId, title)
		userChatSqlBuilder.AddValue("(@, @, @)"+lastSymbol, chatId, user.UserId, user.IsOwner)
	}

	err = c.insertRelationToChat(labelChatSqlBuilder, tx, model.LABELS_CHATS_TABLE)
	if err != nil {
		return 0, err
	}
	err = c.insertRelationToChat(userQueryBuilder, tx, model.USERS_CHATS_TABLE)
	if err != nil {
		return 0, err
	}

	tx.Commit()
	log.Printf("Created chat, insert users to chat and labels for user. Chat id: %d", chatId)
	return chatId, nil
}

func (c *ChatPostgreSQL) CreatePrivate(userTitleChats []model.UserTitleChats) (int, error) {
	ctx := context.Background()
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		c.log.Warnf("Transaction error: %v", err)
		return 0, err
	}

	chatIdIfExist, err := tx.QueryContext(
		ctx,
		fmt.Sprintf("SELECT chat_id FROM %s uc LEFT JOIN %s c on c.id = uc.chat_id WHERE user_id in ($1, $2) AND c.type = %d GROUP BY chat_id HAVING (count(chat_id)) >= 2 LIMIT 1",
			model.USERS_CHATS_TABLE,
			model.CHATS_TABLE,
			model.PRIVATE_CHAT,
		),
		userTitleChats[0].UserId, userTitleChats[1].UserId,
	)
	var chatId int
	if chatIdIfExist.Next() {
		err = chatIdIfExist.Scan(&chatId)
		if err != nil {
			c.log.Warnf("Can't scan chatId. err: %v", err)
			tx.Rollback()
			return 0, err
		}
	}
	if chatId != 0 {
		tx.Commit()
		return chatId, nil
	}

	chatId, err = c.createChat(ctx, tx, "", model.PRIVATE_CHAT)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	userQueryBuilder := c.db.QueryBuilder.NewQueryBuilder()
	labelChatSqlBuilder := c.db.QueryBuilder.Insert(fmt.Sprintf("%s (chat_id, user_id, title)", model.LABELS_CHATS_TABLE)).Values("")
	userChatSqlBuilder := userQueryBuilder.Insert(fmt.Sprintf("%s (chat_id, user_id, is_owner)", model.USERS_CHATS_TABLE)).Values("")
	for count, user := range userTitleChats {
		lastSymbol := ","
		if count == len(userTitleChats)-1 {
			lastSymbol = ";"
		}
		labelChatSqlBuilder.AddValue("(@, @, @)"+lastSymbol, chatId, user.UserId, user.Title)
		userChatSqlBuilder.AddValue("(@, @, @)"+lastSymbol, chatId, user.UserId, true)
	}

	err = c.insertRelationToChat(labelChatSqlBuilder, tx, model.LABELS_CHATS_TABLE)
	if err != nil {
		return 0, err
	}
	err = c.insertRelationToChat(userQueryBuilder, tx, model.USERS_CHATS_TABLE)
	if err != nil {
		return 0, err
	}

	tx.Commit()
	c.log.Printf("Created chat, insert users to chat and labels for user. Chat id: %d", chatId)
	return chatId, nil
}

func (c *ChatPostgreSQL) insertRelationToChat(builder db.QueryBuilder, tx *sql.Tx, tableName string) error {
	userChatSql, userChatArgs := builder.ToSql()
	_, err := tx.Exec(userChatSql, userChatArgs...)
	if err != nil {
		tx.Rollback()
		c.log.Warnf("%s insert error: %v", tableName, err)
		return err
	}

	return nil
}

func (c *ChatPostgreSQL) createChat(ctx context.Context, tx *sql.Tx, description string, chatType int) (chatId int, err error) {
	dateNow := time.Now()
	chatRes, err := tx.QueryContext(ctx, fmt.Sprintf("INSERT INTO %s (description, type, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id", model.CHATS_TABLE), description, chatType, dateNow, dateNow)
	if err != nil {
		c.log.Warnf("Transaction error: %v", err)
		tx.Rollback()
		return 0, err
	}
	if chatRes.Next() {
		err = chatRes.Scan(&chatId)
		if err != nil {
			tx.Rollback()
			c.log.Warnf("Scan chats error: %v", err)
			return 0, err
		}
	}
	chatRes.Close()

	return chatId, nil
}
