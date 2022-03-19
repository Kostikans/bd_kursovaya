package main

import (
	"context"
	"log"
	"math/rand"
	"sync"

	desc "github.com/kostikan/bd_kursovaya/internal/pb/api/forum"
	"google.golang.org/grpc"
)

const (
	workerCount      = 100
	accountCount     = 10
	postCount        = 15
	tagCount         = 20
	commentCount     = 30
	likeCountPost    = 30
	likeCountComment = 30
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:8079", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	client := desc.NewForumClient(conn)
	wg := &sync.WaitGroup{}

	accountIDs := make([]uint64, 0, accountCount)
	for i := 0; i < accountCount; i++ {
		wg.Add(1)
		go func(ctx context.Context, c desc.ForumClient) {
			defer wg.Done()
			id, err := createAccount(ctx, c)
			if err != nil {
				log.Fatalln(err)
			}
			accountIDs = append(accountIDs, id)
		}(ctx, client)
	}
	wg.Wait()

	tagsIDs := make([]uint64, 0, tagCount)
	for i := 0; i < tagCount; i++ {
		wg.Add(1)
		go func(ctx context.Context, c desc.ForumClient) {
			defer wg.Done()
			id, err := createTag(ctx, c, accountIDs[rand.Intn(len(accountIDs))])
			if err != nil {
				log.Fatalln(err)
			}
			tagsIDs = append(tagsIDs, id)
		}(ctx, client)
	}
	wg.Wait()

	postIDs := make([]uint64, 0, postCount)
	for i := 0; i < postCount; i++ {
		wg.Add(1)
		go func(ctx context.Context, c desc.ForumClient) {
			defer wg.Done()
			id, err := createPost(ctx, c, accountIDs[rand.Intn(len(accountIDs))])
			if err != nil {
				log.Fatalln(err)
			}
			postIDs = append(postIDs, id)
		}(ctx, client)
	}
	wg.Wait()

	for i := 0; i < postCount; i++ {
		wg.Add(1)
		go func(ctx context.Context, c desc.ForumClient) {
			defer wg.Done()
			postID := postIDs[rand.Intn(len(postIDs))]
			tagID := tagsIDs[rand.Intn(len(tagsIDs))]
			err := assignTags(ctx, c, postID, []uint64{tagID})
			if err != nil {
				log.Fatalln(err)
			}
		}(ctx, client)
	}
	wg.Wait()

	postCommMap := map[uint64]uint64{}
	commentsIDs := make([]uint64, 0, commentCount)
	mtx := &sync.Mutex{}
	for i := 0; i < commentCount/4; i++ {
		wg.Add(1)
		go func(ctx context.Context, c desc.ForumClient) {
			defer wg.Done()
			authorID := accountIDs[rand.Intn(len(accountIDs))]
			postID := postIDs[rand.Intn(len(postIDs))]
			var parent uint64
			id, err := createComment(ctx, c, authorID, postID, parent)
			if err != nil {
				log.Fatalln(err)
			}
			commentsIDs = append(commentsIDs, id)
			mtx.Lock()
			postCommMap[id] = postID
			mtx.Unlock()
		}(ctx, client)
	}
	wg.Wait()

	for i := 0; i < 3*commentCount/4; i++ {
		wg.Add(1)
		go func(ctx context.Context, c desc.ForumClient) {
			defer wg.Done()
			authorID := accountIDs[rand.Intn(len(accountIDs))]
			postID := postIDs[rand.Intn(len(postIDs))]
			var parent uint64
			if len(commentsIDs) > 0 {
				parent = commentsIDs[rand.Intn(len(commentsIDs))]
			}
			id, err := createComment(ctx, c, authorID, postID, parent)
			if err != nil {
				log.Fatalln(err)
			}
			commentsIDs = append(commentsIDs, id)
			mtx.Lock()
			postCommMap[id] = postID
			mtx.Unlock()
		}(ctx, client)
	}
	wg.Wait()

	postLikeIDs := make([]uint64, 0, likeCountPost)
	for i := 0; i < likeCountPost; i++ {
		wg.Add(1)
		go func(ctx context.Context, c desc.ForumClient) {
			defer wg.Done()
			authorID := accountIDs[rand.Intn(len(accountIDs))]
			postID := postIDs[rand.Intn(len(postIDs))]
			id, err := createPostVote(ctx, c, authorID, postID)
			if err != nil {
				log.Fatalln(err)
			}
			postLikeIDs = append(postLikeIDs, id)
		}(ctx, client)
	}
	wg.Wait()

	commentLikeIDs := make([]uint64, 0, likeCountComment)
	for i := 0; i < likeCountComment; i++ {
		wg.Add(1)
		go func(ctx context.Context, c desc.ForumClient) {
			defer wg.Done()
			authorID := accountIDs[rand.Intn(len(accountIDs))]
			commentID := commentsIDs[rand.Intn(len(commentsIDs))]
			postID, ok := postCommMap[commentID]
			if !ok {
				log.Fatalln("cannot find post by comment")
			}
			id, err := createCommentVote(ctx, c, authorID, postID, commentID)
			if err != nil {
				log.Fatalln(err)
			}
			commentLikeIDs = append(commentLikeIDs, id)
		}(ctx, client)
	}
	wg.Wait()
}

func createAccount(ctx context.Context, c desc.ForumClient) (id uint64, err error) {
	resp, err := c.CreateAccount(ctx, &desc.CreateAccountRequest{
		Account: &desc.Account{
			Nickname:    RandStringRunes(10),
			Avatar:      RandStringRunes(10),
			Description: RandStringRunes(10),
		},
	})
	if err != nil {
		return
	}

	return resp.GetId(), nil
}

func createTag(ctx context.Context, c desc.ForumClient, authorID uint64) (id uint64, err error) {
	resp, err := c.CreateTag(ctx, &desc.CreateTagRequest{
		AuthorId: authorID,
		Name:     RandStringRunes(10),
	})
	if err != nil {
		return
	}

	return resp.GetId(), nil
}

func createPost(ctx context.Context, c desc.ForumClient, authorID uint64) (id uint64, err error) {
	resp, err := c.CreatePost(ctx, &desc.CreatePostRequest{
		AuthorId: authorID,
		Title:    RandStringRunes(10),
		Text:     RandStringRunes(10),
	})
	if err != nil {
		return
	}

	return resp.GetId(), nil
}

func createComment(ctx context.Context, c desc.ForumClient, authorID uint64, postID uint64, parentID uint64) (id uint64, err error) {
	resp, err := c.CreateComment(ctx, &desc.CreateCommentRequest{
		AuthorId: authorID,
		PostId:   postID,
		Text:     RandStringRunes(10),
		ParentId: parentID,
	})
	if err != nil {
		return
	}

	return resp.GetId(), nil
}

func createPostVote(ctx context.Context, c desc.ForumClient, authorID uint64, postID uint64) (id uint64, err error) {
	resp, err := c.CreatePostVote(ctx, &desc.CreatePostVoteRequest{
		AuthorId: authorID,
		PostId:   postID,
		Vote:     rand.Intn(2) != 0,
	})
	if err != nil {
		return
	}

	return resp.GetId(), nil
}

func createCommentVote(ctx context.Context, c desc.ForumClient, authorID uint64, postID uint64, commentID uint64) (id uint64, err error) {
	resp, err := c.CreateCommentVote(ctx, &desc.CreateCommentVoteRequest{
		AuthorId:  authorID,
		PostId:    postID,
		Vote:      rand.Intn(2) != 0,
		CommentId: commentID,
	})
	if err != nil {
		return
	}

	return resp.GetId(), nil
}

func assignTags(ctx context.Context, c desc.ForumClient, postID uint64, tags []uint64) (err error) {
	_, err = c.AssignTagsToPost(ctx, &desc.AssignTagsToPostRequest{
		PostId: postID,
		TagsId: tags,
	})
	return
}
