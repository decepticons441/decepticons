package sessions

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

//RedisStore represents a session.Store backed by redis.
type RedisStore struct {
	//Redis client used to talk to redis server.
	Client *redis.Client
	//Used for key expiry time on redis.
	SessionDuration time.Duration
}

//NewRedisStore constructs a new RedisStore
func NewRedisStore(client *redis.Client, sessionDuration time.Duration) *RedisStore {
	//initialize and return a new RedisStore struct
	return &RedisStore{
		Client:          client,
		SessionDuration: sessionDuration,
	}
}

//Store implementation

//Save saves the provided `sessionState` and associated SessionID to the store.
//The `sessionState` parameter is typically a pointer to a struct containing
//all the data you want to associated with the given SessionID.
func (rs *RedisStore) Save(sid SessionID, sessionState interface{}) error {
	//TODO: marshal the `sessionState` to JSON and save it in the redis database,
	//using `sid.getRedisKey()` for the key.
	//return any errors that occur along the way.

	j, err := json.Marshal(sessionState)
	if err != nil {
		return err
	}
	session := rs.Client.Set(sid.getRedisKey(), j, rs.SessionDuration)
	if session.Err() != nil {
		return ErrStateNotFound
	}
	// err = statusErr.Err()
	// if err != nil {
	// 	return ErrStateNotFound
	// }
	return nil
}

//Get populates `sessionState` with the data previously saved
//for the given SessionID
func (rs *RedisStore) Get(sid SessionID, sessionState interface{}) error {
	//TODO: get the previously-saved session state data from redis,
	//unmarshal it back into the `sessionState` parameter
	//and reset the expiry time, so that it doesn't get deleted until
	//the SessionDuration has elapsed.

	p := rs.Client.Pipeline()
	j := p.Get(sid.getRedisKey())
	boolcmd := p.Expire(sid.getRedisKey(), rs.SessionDuration)
	p.Exec()
	// j := rs.Client.Get(sid.getRedisKey())

	found := j.Err()
	if found != nil {
		return ErrStateNotFound
	}

	result, err := j.Result()

	if err != nil {
		return fmt.Errorf("error when getting JSON object result: %v", err)
	}
	err = json.Unmarshal([]byte(result), sessionState)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %v", err)
	}
	// boolcmd := rs.Client.Expire(sid.getRedisKey(), rs.SessionDuration)

	err = boolcmd.Err()
	if err != nil {
		return fmt.Errorf("error when resetting time to live: %s", err)
	}

	//reset TTL
	// boolcmd := rs.Client.Expire(sid.String(), rs.SessionDuration)
	// err = boolcmd.Err()
	// if err != nil {
	// 	return fmt.Errorf("error when resetting time to live: %s", err)
	// }

	return nil
	//for extra-credit using the Pipeline feature of the redis
	//package to do both the get and the reset of the expiry time
	//in just one network round trip!
}

//Delete deletes all state data associated with the SessionID from the store.
func (rs *RedisStore) Delete(sid SessionID) error {
	//TODO: delete the data stored in redis for the provided SessionID
	intcmd := rs.Client.Del(sid.String())
	err := intcmd.Err()
	if err != nil {
		return fmt.Errorf("error when deleting state data associated with session id: %s", err)
	}
	return nil
}

//getRedisKey() returns the redis key to use for the SessionID
func (sid SessionID) getRedisKey() string {
	//convert the SessionID to a string and add the prefix "sid:" to keep
	//SessionID keys separate from other keys that might end up in this
	//redis instance
	return "sid:" + sid.String()
}
