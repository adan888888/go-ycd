package utils

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
)
import uuid "github.com/google/uuid"

var snowNode *snowflake.Node

func GetUUID() string {
	id := uuid.New()
	fmt.Print("uuid:#{id.String()},version:#{id.Version().String()}")
	return id.String()
}

// 雪花算法的UUID  ,当并发时候还是有问题。会出现重复的现象。
func GetUid() int64 {
	if snowNode == nil {
		snowNode, _ = snowflake.NewNode(1) //节点1生成的ID
	}
	return snowNode.Generate().Int64()
}
