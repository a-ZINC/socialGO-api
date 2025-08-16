package db

import (
	"context"
	"fmt"
	"math/rand"
	"social-api/internal/store"
	"social-api/model"
	"strconv"
)

var names = []string{"Alice", "Bob", "Charlie", "Diana", "Eve", "Frank", "Grace", "Hank", "Ivy", "Jack"}
var sampleTags = [][]string{
	{"tech", "go"},
	{"life", "travel"},
	{"news", "updates"},
	{"programming", "backend"},
	{"design", "ux"},
}
var commentPhrases = []string{
	"Great point!",
	"Thanks for sharing!",
	"I totally agree with this.",
	"Could you elaborate more?",
	"This is really helpful.",
	"Interesting take 🔍",
	"I never thought of it that way 🤯",
	"Nice write-up! 👏",
	"Any sources for this?",
	"Made my day 😂",
	"This deserves more attention 🚀",
	"Solid perspective.",
	"🔥🔥🔥",
}

var titles = []string{
	"Exploring the Future of AI",
	"Lessons from My Startup Journey",
	"The Hidden Power of Simplicity",
	"Why Rust Might Replace C++",
	"Thoughts on Remote Work Culture",
	"The Art of Problem Solving",
	"Debugging Nightmares and How to Escape",
	"Functional vs OOP: Final Verdict",
	"Building Scalable Microservices",
	"Top 5 Productivity Hacks for Devs",
}

var contents = []string{
	"In this post, I share insights on how technology is reshaping our lives...",
	"This week, I struggled with an elusive bug that turned into a breakthrough...",
	"Some thoughts on how focusing on fundamentals can solve even the hardest problems...",
	"Here’s what I learned after switching from monoliths to microservices...",
	"Productivity is not about working more, it’s about working smarter...",
	"Rust's safety and performance made me rethink my C++ habits...",
	"Sometimes, the best solutions come from taking a walk away from the keyboard...",
	"Here’s a simple technique I use to manage technical debt in large systems...",
	"Why building a personal brand online matters for developers...",
	"An honest take on impostor syndrome and how I deal with it...",
}

func Seed(store *store.Store) error {
	// Seed users
	users := generateUser(100)
	for i := range users {
		if err := store.Users.Create(context.Background(), &users[i]); err != nil {
			return err
		}
	}
	fmt.Println("Users seeded successfully.")
	fmt.Printf("Total Users: %d\n", len(users))
	fmt.Printf("Sample User: %d, Email: %s\n", users[0].ID, users[0].Email)

	posts := generatePost(200, users)
	for i := range posts {
		if err := store.Posts.Create(context.Background(), &posts[i]); err != nil {
			return err
		}
	}

	comments := generateComments(500, posts, users)
	for i := range comments {
		if err := store.Comments.Create(context.Background(), &comments[i]); err != nil {
			return err
		}
	}

	fmt.Println("Database seeded successfully with users, posts, and comments.")
	fmt.Printf("Total Users: %d, Posts: %d, Comments: %d\n", len(users), len(posts), len(comments))
	fmt.Printf("Sample User: %d, Sample Post: %s, Sample Comment: %s\n",
		users[0].ID, posts[0].Title, comments[0].Content)

	return nil
}

func generateUser(count int) []model.User {
	users := make([]model.User, count)
	for i := 0; i < count; i++ {
		uniqueSuffix := strconv.Itoa(i) + "-" + strconv.Itoa(rand.Intn(100000))
		name := names[rand.Intn(len(names))] + uniqueSuffix
		users[i] = model.User{
			Name:     name,
			Email:    name + "@example.com",
			Password: model.Password{Text: &name, Hash: []byte("hashed_password")}, // Replace with real hash in production
		}
	}
	return users
}
func generatePost(count int, users []model.User) []model.Post {
	posts := make([]model.Post, count)
	for i := 0; i < count; i++ {
		user := users[rand.Intn(len(users))]
		title := titles[rand.Intn(len(titles))] + " #" + strconv.Itoa(i)
		content := contents[rand.Intn(len(contents))] + " #" + strconv.Itoa(i)

		posts[i] = model.Post{
			Title:   title,
			Content: content,
			UserID:  user.ID,
			Tags:    sampleTags[rand.Intn(len(sampleTags))],
		}
	}
	return posts
}
func generateComments(count int, posts []model.Post, users []model.User) []model.Comment {
	comments := make([]model.Comment, count)
	for i := 0; i < count; i++ {
		post := posts[rand.Intn(len(posts))]
		user := users[rand.Intn(len(users))]
		comments[i] = model.Comment{
			Content: commentPhrases[rand.Intn(len(commentPhrases))] + " (on post #" + strconv.Itoa(int(post.ID)) + ")",
			PostID:  post.ID,
			UserID:  user.ID,
		}
	}
	return comments
}
