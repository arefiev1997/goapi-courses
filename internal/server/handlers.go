package server

import "github.com/gin-gonic/gin"

func (s *Server) GetClasses(c *gin.Context) {
	classes, err := s.db.GetClasses(c.Request.Context())
	if err != nil {
		c.Error(err)
	}
	c.JSON(200, gin.H{
		"result": classes,
	})
}
