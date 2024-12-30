package db

import (
	"context"
	"log"
	"math/rand"
	"strconv"

	"github.com/karthik446/social/internal/store"
)

var usernames = []string{
	"cosmic_voyager",
	"digital_nomad",
	"pixel_pioneer",
	"cyber_sage",
	"quantum_dreamer",
	"techie_wizard",
	"data_drifter",
	"neo_navigator",
	"byte_wanderer",
	"cloud_chaser",
	"echo_explorer",
	"future_finder",
	"matrix_maven",
	"swift_seeker",
	"code_catalyst",
	"wave_walker",
	"binary_builder",
	"neural_knight",
	"pulse_pathfinder",
	"vector_voyager",
	"logic_leaper",
	"omega_oracle",
	"delta_drifter",
	"azure_artist",
	"crypto_creator",
	"meta_maker",
	"quantum_quester",
	"sonic_scholar",
	"cyber_sentinel",
	"pixel_prophet",
	"data_dynamo",
	"echo_enigma",
	"flux_phoenix",
	"grid_guardian",
	"nova_nomad",
	"pulse_pilot",
	"spark_spectre",
	"tech_templar",
	"void_vector",
	"wave_weaver",
	"zero_zenith",
	"byte_baron",
	"code_corsair",
	"data_duchess",
	"echo_empress",
	"flux_fighter",
	"grid_ghost",
	"nova_ninja",
	"pulse_pirate",
	"tech_titan",
}

var titles = []string{
	"The Future of Artificial Intelligence in Healthcare",
	"Exploring the Hidden Gems of Southeast Asia",
	"A Beginner's Guide to Urban Gardening",
	"Understanding Quantum Computing Basics",
	"The Art of Mindful Living in a Digital Age",
	"Sustainable Architecture: Building for Tomorrow",
	"Essential Tips for Remote Work Success",
	"The Evolution of Electric Vehicles",
	"Mastering the Basics of Digital Photography",
	"The Impact of Social Media on Modern Society",
	"Secrets of Effective Time Management",
	"Understanding Blockchain Technology",
	"The Science Behind Climate Change",
	"Modern Web Development Practices",
	"The Psychology of Habit Formation",
	"Space Exploration: Past, Present, and Future",
	"Essential Skills for Data Science",
	"The Art of Public Speaking",
	"Understanding Cryptocurrency Markets",
	"Healthy Eating in the Modern World",
}

var contents = []string{
	"AI is transforming healthcare through improved diagnosis and treatment planning. The future looks promising.",
	"Hidden beaches, ancient temples, and vibrant street markets make Southeast Asia a traveler's paradise.",
	"Start your urban garden today with these simple tips for growing herbs and vegetables on your balcony.",
	"Quantum computing leverages quantum mechanics to solve complex problems faster than traditional computers.",
	"Practice mindfulness in the digital age by setting boundaries and creating tech-free zones in your life.",
	"Green buildings and sustainable materials are shaping the future of architecture worldwide.",
	"Success in remote work requires discipline, communication, and a well-organized workspace.",
	"Electric vehicles are revolutionizing transportation with improved range and charging infrastructure.",
	"Learn composition, lighting, and camera settings to take your photography skills to the next level.",
	"Social media has changed how we communicate, share information, and build relationships.",
	"Effective time management starts with clear priorities and smart scheduling techniques.",
	"Blockchain technology offers secure, transparent transactions without traditional intermediaries.",
	"Understanding climate change requires knowledge of global weather patterns and human impact.",
	"Modern web development embraces responsive design and progressive enhancement.",
	"Form better habits by understanding the psychology behind behavior change.",
	"Space exploration continues to push boundaries with new missions to Mars and beyond.",
	"Data science combines statistics, programming, and domain expertise for insights.",
	"Master public speaking through practice, preparation, and audience engagement.",
	"Cryptocurrency markets operate 24/7 with unique opportunities and risks.",
	"Balance nutrition and convenience in your daily diet for optimal health.",
}

var tags = []string{
	"technology",
	"programming",
	"health",
	"science",
	"travel",
	"food",
	"lifestyle",
	"fitness",
	"education",
	"business",
	"finance",
	"art",
	"music",
	"sports",
	"books",
	"movies",
	"photography",
	"design",
	"gaming",
	"nature",
	"environment",
	"politics",
	"history",
	"philosophy",
	"psychology",
	"space",
	"ai",
	"blockchain",
	"crypto",
	"web3",
	"coding",
	"data",
	"cloud",
	"security",
	"mobile",
	"social",
	"startup",
	"innovation",
	"research",
	"career",
	"productivity",
	"mindfulness",
	"culture",
	"development",
	"marketing",
	"writing",
	"learning",
	"leadership",
	"community",
	"future",
}

var comments = []string{
	"Great insights! This really helped me understand the topic better.",
	"I never thought about it from this perspective before. Thanks for sharing!",
	"Would love to see a follow-up post on this topic.",
	"This is exactly what I've been looking for. Well explained!",
	"Very informative article. Keep up the great work!",
	"Thanks for breaking this down in such a clear way.",
	"Interesting points, though I slightly disagree with the third one.",
	"This has given me some great ideas to try out.",
	"Can you elaborate more on the technical aspects?",
	"Brilliant analysis! Sharing this with my team.",
	"The examples really helped clarify the concepts.",
	"Looking forward to more content like this.",
	"Your writing style makes complex topics easy to understand.",
	"This is a game-changer for my workflow.",
	"Really appreciate the practical tips!",
	"Would be great to see some code examples next time.",
	"This matches my experience exactly.",
	"Great article! Just what I needed today.",
	"The visuals really help explain the concept.",
	"I've implemented these ideas with great results.",
	"Simple yet effective explanation.",
	"This should be required reading for beginners.",
	"Love how you approached this topic.",
	"Adding this to my bookmarks for future reference.",
	"The step-by-step breakdown is super helpful.",
	"This explains why my previous attempts failed!",
	"Excellent resource for anyone starting out.",
	"I learned something new today, thanks!",
	"Clear, concise, and practical advice.",
	"This deserves more attention.",
	"Finally, an explanation that makes sense!",
	"Been waiting for someone to address this.",
	"The research behind this is solid.",
	"Perfect timing - I was just studying this topic.",
	"This has changed my perspective completely.",
	"Saving this for later reference.",
	"Your expertise really shows in this post.",
	"Can't wait to try these techniques.",
	"This is going to help a lot of people.",
	"Such a comprehensive overview!",
	"The real-world applications are spot-on.",
	"This is becoming my go-to resource.",
	"Love the practical approach here.",
	"Well-researched and thoughtfully presented.",
	"These insights are invaluable.",
	"This clarified several misconceptions I had.",
	"Exactly what I needed to read today.",
	"The quality of content here is outstanding.",
	"This sparked some interesting ideas.",
	"Really appreciate the depth of this analysis.",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, u := range users {
		if err := store.Users.Create(ctx, u); err != nil {
			log.Println("Error creating user:", err)
			return
		}
	}

	posts := generatePosts(100, users)

	for _, p := range posts {
		if err := store.Posts.Create(ctx, p); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(100, users, posts)

	for _, c := range comments {
		if err := store.Comments.Create(ctx, c); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}
	log.Println("Successfully seeded the database.")
}

func generateUsers(n int) []*store.User {
	users := make([]*store.User, n)
	for i := 0; i < n; i++ {
		u := usernames[i%len(usernames)]
		users[i] = &store.User{
			Username: u + strconv.Itoa(i),
			Email:    u + strconv.Itoa(i) + "@example.com",
			Password: "password",
		}
	}
	return users
}

func generatePosts(n int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, n)
	for i := 0; i < n; i++ {
		u := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			Title:   titles[rand.Intn(len(titles))] + strconv.Itoa(i),
			Content: contents[rand.Intn(len(titles))] + strconv.Itoa(i),
			Tags:    []string{tags[rand.Intn(len(tags))], tags[rand.Intn(len(tags))]},
			UserID:  u.ID,
		}
	}
	return posts
}

func generateComments(n int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, n)
	for i := 0; i < n; i++ {
		u := users[rand.Intn(len(users))]
		p := posts[rand.Intn(len(posts))]
		cms[i] = &store.Comment{
			Content: comments[rand.Intn(len(comments))],
			UserID:  u.ID,
			PostID:  p.ID,
		}
	}
	return cms
}
