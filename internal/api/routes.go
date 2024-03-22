package api

// import (
// 	"github.com/GhostDrew11/vigor-api/internal/middlewares"
// 	"github.com/GhostDrew11/vigor-api/internal/utils"
// 	"github.com/gin-gonic/gin"
// )

// func SetupRoutes(router *gin.Engine, ts utils.TokenService) {
// 	// API root
// 	apiRoot := router.Group("/api/v1")

// 	// Auth routes
// 	authRoutes := apiRoot.Group("/auth")
// 	authRoutes.POST("/refresh", middlewares.RefreshHandler(ts))

// 	// Admin routes
// 	// Workout Plan routes
// 	adminRoutes := apiRoot.Group("/admin")
// 	adminRoutes.POST("/wokout-plans", createWorkoutPlan)
// 	adminRoutes.PUT("/workout-plans/:id", updateWorkoutPlan)
// 	adminRoutes.DELETE("/workout-plans/:id", deleteWorkoutPlan)
// 	adminRoutes.GET("/workout-plans", getWorkoutPlans)
// 	// Meal Plan routes
// 	adminRoutes.POST("/meal-plans", createMealPlan)
// 	adminRoutes.PUT("/meal-plans/:id", updateMealPlan)
// 	adminRoutes.DELETE("/meal-plans/:id", deleteMealPlan)
// 	adminRoutes.GET("/meal-plans", getMealPlans)
// 	// Exercise routes
// 	adminRoutes.POST("/exercises", createExercise)
// 	adminRoutes.PUT("/exercises/:id", updateExercise)
// 	adminRoutes.DELETE("/exercises/:id", deleteExercise)
// 	adminRoutes.GET("/exercises", getExercises)
// 	// Meal routes
// 	adminRoutes.POST("/meals", createMeal)
// 	adminRoutes.PUT("/meals/:id", updateMeal)
// 	adminRoutes.DELETE("/meals/:id", deleteMeal)
// 	adminRoutes.GET("/meals", getMeals)
// 	// Add more admin routes as needed (e.g., get user workout plans, get user meal plans, etc.)

// 	// User routes
// 	userRoutes := apiRoot.Group("/users")
// 	userRoutes.POST("/register", registerUser)
// 	userRoutes.POST("/login", loginUser)

// 	// Authenticated user routes
// 	authenticatedUserRoutes := userRoutes.Group("/")
// 	authenticatedUserRoutes.Use(middlewares.Authenticate(ts))
// 	authenticatedUserRoutes.GET("/profile", getUserProfile)
// 	authenticatedUserRoutes.PUT("/profile", updateUserProfile)
// 	authenticatedUserRoutes.GET("/preferences", getUserPreferences)
// 	authenticatedUserRoutes.PUT("/preferences", updateUserPreferences)
// 	authenticatedUserRoutes.GET("/subscription", getUserSubscription)
// 	authenticatedUserRoutes.PUT("/subscription", updateUserSubscription)
// 	authenticatedUserRoutes.DELETE("/subscription", cancelUserSubscription)
// 	authenticatedUserRoutes.DELETE("/account", deleteUserAccount)

// 	// Search for a user via it's username to start a new conversation
// 	authenticatedUserRoutes.GET("/search", searchUser)

// 	// Progress tracking routes
// 	// Workout Plan
// 	authenticatedUserRoutes.POST("/workout-plans/:workoutPlanId/progress", createWorkoutPlanProgress)
// 	authenticatedUserRoutes.GET("/workout-plans/:workoutPlanId/progress", getWorkoutPlanProgress)
// 	authenticatedUserRoutes.PUT("/workout-plans/:workoutPlanId/progress", updateWorkoutPlanProgress)
// 	// Meal Plan
// 	authenticatedUserRoutes.POST("/meal-plans/:mealPlanId/progress", createMealPlanProgress)
// 	authenticatedUserRoutes.GET("/meal-plans/:mealPlanId/progress", getMealPlanProgress)
// 	authenticatedUserRoutes.PUT("/meal-plans/:mealPlanId/progress", updateMealPlanProgress)

// 	// Daily nutritional logs
// 	authenticatedUserRoutes.POST("/nutritional-logs", createNutritionalLog)
// 	// Maybe use them to build some analytics and graphs you can show to the user
// 	authenticatedUserRoutes.GET("/nutritional-logs", getNutritionalLogs)

// 	// Authenticated user direct message conversations
// 	authenticatedUserRoutes.POST("/conversations", createConversation)
// 	// Delete a conversation
// 	authenticatedUserRoutes.DELETE("/conversations/:conversationId", deleteConversation)
// 	// Get all conversations for a user
// 	authenticatedUserRoutes.GET("/conversations", getConversations)
// 	// Send Message in a conversation
// 	authenticatedUserRoutes.POST("/conversations/:conversationId/messages", sendMessage)
// 	// Read a message in a conversation
// 	authenticatedUserRoutes.GET("/conversations/:conversationId/messages/:messageId", readMessage)
// 	// Update message content in a conversation
// 	authenticatedUserRoutes.PUT("/conversations/:conversationId/messages/:messageId", updateMessage)
// 	// Delete a message in a conversation
// 	authenticatedUserRoutes.DELETE("/conversations/:conversationId/messages/:messageId", deleteMessage)
// 	// Get all messages in a conversation
// 	authenticatedUserRoutes.GET("/conversations/:conversationId/messages", getMessages)

// 	// Two group routes because group might have other functionalities that just chatting within the group
// 	// Authenticated user group conversations
// 	authenticatedUserRoutes.POST("/groups", createGroup)
// 	// Delete a group the user created
// 	authenticatedUserRoutes.DELETE("/groups/:groupId", deleteGroup)
// 	// Leave a group
// 	authenticatedUserRoutes.DELETE("/groups/:groupId/leave", leaveGroup)
// 	// Link Group to Conversation
// 	authenticatedUserRoutes.POST("/groups/:groupId/conversations", createGroupChat)
// 	// Send a message in a group conversation
// 	authenticatedUserRoutes.POST("/groups/:groupId/conversations/:conversationId/messages", sendGroupMessage)
// 	// update group message content
// 	authenticatedUserRoutes.PUT("/groups/:groupId/conversations/:conversationId/messages/:messageId", updateGroupMessage)
// 	// Remove a user from a group
// 	authenticatedUserRoutes.DELETE("/groups/:groupId/members/:userId", removeGroupMember)
// 	// Other group functionalities as needed (e.g, add member, join a group, having a group live workout party etc.)
// }
