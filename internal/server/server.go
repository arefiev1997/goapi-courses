package server

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/arefiev1997/goapi/internal/config"
	"github.com/arefiev1997/goapi/internal/database"
	"github.com/gin-gonic/gin"
)

// Server struct
type Server struct {
	router *gin.Engine
	port   int
	db     database.Database
}

// New create new Server instance
func New(cfg config.ServerConfig, db database.Database) *Server {
	return &Server{
		router: gin.Default(),
		port:   cfg.Port,
		db:     db,
	}
}

// Run server
func (s *Server) Run() {
	s.initHandlers()
	s.router.Run(fmt.Sprintf(":%d", s.port))
}

func (s *Server) initHandlers() {
	group := s.router.Group("/api")
	group.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"hello": "world",
		})
	})

	classGroup := group.Group("/class")
	classGroup.GET("/list", s.GetClasses)
	classGroup.POST("/create", func(c *gin.Context) {
		var input database.Class
		c.ShouldBindJSON(&input)
		if err := s.db.CreateClass(c.Request.Context(), input); err != nil {
			c.Error(err)
		}
	})

	classGroup.DELETE("/:id", func(c *gin.Context) {
		classIDint, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Error(errors.New("text string"))
		}
		if err := s.db.DeleteClass(c.Request.Context(), classIDint); err != nil {
			c.Error(err)
		}
		c.Status(200)
	})

	studentsGroup := group.Group("students")
	studentsGroup.GET("/list", func(c *gin.Context) {
		classIDint, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			c.Error(errors.New("text string"))
		}
		result, err := s.db.GetStudentsByClass(c.Request.Context(), classIDint)
		if err != nil {
			c.Error(err)
		}
		c.JSON(200, gin.H{
			"result": result,
		})
	})

	studentsGroup.POST("/create", func(c *gin.Context) {
		var input database.Student
		if err := c.ShouldBindJSON(&input); err != nil {
			c.Status(400)
		}
		if err := s.db.CreateStudent(c.Request.Context(), input); err != nil {
			c.Status(400)
		}
		c.Status(200)
	})
}
