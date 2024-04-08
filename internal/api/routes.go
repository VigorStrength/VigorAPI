package api

import (
	"github.com/GhostDrew11/vigor-api/internal/controllers"
	"github.com/GhostDrew11/vigor-api/internal/middlewares"
	"github.com/GhostDrew11/vigor-api/internal/services"
	"github.com/GhostDrew11/vigor-api/internal/utils"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, ts utils.TokenService, userService services.UserService, adminService services.AdminService) {
	// API root
	apiRoot := router.Group("/api/v1")
	adminController := controllers.NewAdminController(adminService, ts)
	userController := controllers.NewUserController(userService, ts)

	// Auth routes
	authRoutes := apiRoot.Group("/auth")
	authRoutes.POST("/admin/register", adminController.Register) 
    authRoutes.POST("/admin/login", adminController.Login)       
    authRoutes.POST("/user/register", userController.Register)   
    authRoutes.POST("/user/login", userController.Login)        
	authRoutes.POST("/refresh", middlewares.RefreshHandler(ts))

	// Admin routes
	adminRoutes := apiRoot.Group("/admin")
	adminRoutes.Use(middlewares.RequireRole(ts, "admin"))
	// // CRUD Exercises
	adminRoutes.POST("/exercises", adminController.CreateExercise)
	adminRoutes.POST("/exercises/bulk", adminController.CreateMultipleExercises)
	// TODO: Implement sort, filter, and pagination for exercises
	adminRoutes.GET("/exercises", adminController.GetExercises)
	adminRoutes.GET("/exercises/:id", adminController.GetExerciseByID)
	adminRoutes.GET("/exercises/search", adminController.SearchExercisesByName)
	adminRoutes.PUT("/exercises/:id", adminController.UpdateExercise)
	adminRoutes.DELETE("/exercises/:id", adminController.DeleteExercise)
	// // CRUD Workout Plans
	// adminRoutes.POST("/workout-plans", createWorkoutPlan)
	// adminRoutes.PUT("/workout-plans/:id", updateWorkoutPlan)
	// adminRoutes.DELETE("/workout-plans/:id", deleteWorkoutPlan)
	// adminRoutes.GET("/workout-plans", getWorkoutPlans)
	
	
	

	// User routes
	userRoutes := apiRoot.Group("/users")
	userRoutes.Use(middlewares.RequireRole(ts, "user"))
	// CRUD User data
	// userRoutes.GET("/profile", getUserProfile)
	// userRoutes.PUT("/profile", updateUserProfile)
	// userRoutes.GET("/preferences", getUserPreferences)
	// userRoutes.PUT("/preferences", updateUserPreferences)
	// userRoutes.GET("/subscription", getUserSubscription)
	// userRoutes.PUT("/subscription", updateUserSubscription)
	// userRoutes.DELETE("/subscription", cancelUserSubscription)
	// userRoutes.DELETE("/account", deleteUserAccount)
	// // other user routes as needed(eg list user workout plans, list user meal plans, list user progress, other analytics etc.)

	// // Progress tracking routes
	// // Workout Plan
	// userRoutes.POST("/workout-plans/:workoutPlanId/progress", createWorkoutPlanProgress)
	// userRoutes.GET("/workout-plans/:workoutPlanId/progress", getWorkoutPlanProgress)
	// userRoutes.PUT("/workout-plans/:workoutPlanId/progress", updateWorkoutPlanProgress)
	// // Meal Plan
	// userRoutes.POST("/meal-plans/:mealPlanId/progress", createMealPlanProgress)
	// userRoutes.GET("/meal-plans/:mealPlanId/progress", getMealPlanProgress)
	// userRoutes.PUT("/meal-plans/:mealPlanId/progress", updateMealPlanProgress)
	// // Daily nutritional logs
	// userRoutes.POST("/nutritional-logs", createNutritionalLog)
	// // Maybe use them to build some analytics and graphs you can show to the user
	// userRoutes.GET("/nutritional-logs", getNutritionalLogs)
	// userRoutes.PUT("/nutritional-logs/:id", updateNutritionalLog)

	// // Interactions with other users routes(search, chat, etc.)
	// userRoutes.GET("/search", searchUser)
	// userRoutes.POST("/conversations", createConversation)
	// userRoutes.GET("/conversations", getConversations)
	// // Not sure if we need to update a conversation
	// userRoutes.PUT("/conversations/:conversationId", updateConversation)
	// userRoutes.DELETE("/conversations/:conversationId", deleteConversation)
	// // Send Message in a conversation
	// userRoutes.POST("/conversations/:conversationId/messages", sendMessage)
	// // Read a message in a conversation
	// userRoutes.GET("/conversations/:conversationId/messages/:messageId", readMessage)
	// // Update message content in a conversation
	// userRoutes.PUT("/conversations/:conversationId/messages/:messageId", updateMessage)
	// // Delete a message in a conversation just for the sender
	// userRoutes.DELETE("/conversations/:conversationId/messages/:messageId", deleteMessage)
	// // Get all messages in a conversation
	// userRoutes.GET("/conversations/:conversationId/messages", getMessages)
	// // Interactions within a group
	// userRoutes.POST("/groups", createGroup)
	// // Update a group
	// userRoutes.PUT("/groups/:groupId", updateGroup)
	// // Delete a group the user created
	// userRoutes.DELETE("/groups/:groupId", deleteGroup)
	// // Leave a group
	// userRoutes.DELETE("/groups/:groupId/leave", leaveGroup)
	// // Create a group conversation
	// userRoutes.POST("/groups/:groupId/conversations", createGroupChat)
	// // Send a message in a group conversation
	// userRoutes.POST("/groups/:groupId/conversations/:conversationId/messages", sendGroupMessage)
	// // Read a message in a group conversation
	// userRoutes.GET("/groups/:groupId/conversations/:conversationId/messages/:messageId", readGroupMessage)
	// // Update message content in a group conversation
	// userRoutes.PUT("/groups/:groupId/conversations/:conversationId/messages/:messageId", updateGroupMessage)
	// // Delete a message in a group conversation just for the sender
	// userRoutes.DELETE("/groups/:groupId/conversations/:conversationId/messages/:messageId", deleteGroupMessage)
	// // Get all messages in a group conversation
	// userRoutes.GET("/groups/:groupId/conversations/:conversationId/messages", getGroupMessages)
	// // Add a member to a group
	// userRoutes.POST("/groups/:groupId/members", addGroupMember)
	// // Remove a user from a group
	// userRoutes.DELETE("/groups/:groupId/members/:userId", removeGroupMember)
	// Other group functionalities as needed (e.g, add member, join a group, having a group live workout party etc.)

	// // CRUD Meals
	// adminRoutes.POST("/meals", createMeal)
	// adminRoutes.PUT("/meals/:id", updateMeal)
	// adminRoutes.DELETE("/meals/:id", deleteMeal)
	// adminRoutes.GET("/meals", getMeals)
	// // CRUD Meal Plans
	// adminRoutes.POST("/meal-plans", createMealPlan)
	// adminRoutes.PUT("/meal-plans/:id", updateMealPlan)
	// adminRoutes.DELETE("/meal-plans/:id", deleteMealPlan)
	// adminRoutes.GET("/meal-plans", getMealPlans)
	// // CRUD Admins
	// adminRoutes.POST("/admins", createAdmin)
	// adminRoutes.PUT("/admins/:id", updateAdmin)
	// adminRoutes.DELETE("/admins/:id", deleteAdmin)
	// adminRoutes.GET("/admins", getAdmins)
	// // CRUD Users
	// adminRoutes.GET("/users", getUsers)
	// other admin routes as needed(eg list users with active subscriptions, list users with pending subscriptions, list of sales, other analytics etc.)
}
