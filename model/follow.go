package model

type Follow struct {
	FollowerID int64 `db:"follower_id"`
	FollowedID int64 `db:"followed_id"`
}
