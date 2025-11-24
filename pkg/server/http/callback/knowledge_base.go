package callback

import (
	"github.com/gin-gonic/gin"

	"sdk-demo-go/pkg/invoker"
	"sdk-demo-go/pkg/models/db"
)

type KnowledgeBase struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetKnowledgeBases(c *gin.Context) {
	kbs, err := db.FindAllKnowledgeBases(invoker.DB)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to get knowledge bases", "error": err.Error()})
		return
	}
	var res []KnowledgeBase
	for _, kb := range kbs {
		res = append(res, KnowledgeBase{
			ID:   kb.Guid,
			Name: kb.Name,
		})
	}
	c.JSON(200, res)
}
