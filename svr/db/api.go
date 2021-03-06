package db

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adayswait/mojo/global"
	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

type UserInfo struct {
	User  string `json:"user"`
	Group int    `json:"group"`
}

func createAdminAccount(passwd string) {
	tempMD5 := md5.New()
	tempMD5.Write([]byte(passwd))
	tempMD5.Write([]byte(global.MD5_SALT))
	passwdMD5 := hex.EncodeToString(tempMD5.Sum(nil))
	user := "admin"
	rwLock.Lock()
	/*createAdminErr := */ localDB.Update(func(tx *bolt.Tx) error {
		bup := tx.Bucket([]byte(global.BUCKET_USR_PASSWD))
		if bup == nil {
			return fmt.Errorf("bucket:%s is nil", global.BUCKET_USR_PASSWD)
		}
		eup := bup.Put([]byte(user), []byte(passwdMD5))
		if eup != nil {
			return eup
		}

		but := tx.Bucket([]byte(global.BUCKET_USER_TOKEN))
		if but == nil {
			return fmt.Errorf("bucket:%s is nil", global.BUCKET_USER_TOKEN)
		}
		var token []byte
		oldToken := but.Get([]byte(user))
		if len(oldToken) != 0 {
			token = oldToken
		} else {
			token = []byte(uuid.New().String())
			eut := but.Put([]byte(user), token)
			if eut != nil {
				return eut
			}
		}

		info := UserInfo{User: user, Group: int(global.GROUP_WHOSYOURDADDY)}
		infoData, ejn := json.Marshal(info)
		if ejn != nil {
			return ejn
		}
		bti := tx.Bucket([]byte(global.BUCKET_TOKEN_INFO))
		if bti == nil {
			return fmt.Errorf("bucket:%s is nil", global.BUCKET_TOKEN_INFO)
		}

		eti := bti.Put(token, []byte(string(infoData)))
		return eti
	})
	defer rwLock.Unlock()
}

func Register(user, passwd string) error {
	rwLock.Lock()
	registerErr := localDB.Update(func(tx *bolt.Tx) error {
		bup := tx.Bucket([]byte(global.BUCKET_USR_PASSWD))
		if bup == nil {
			return fmt.Errorf("bucket:%s is nil", global.BUCKET_USR_PASSWD)
		}
		existUser := bup.Get([]byte(user))
		if len(existUser) != 0 {
			return errors.New("user already exist")
		}
		eup := bup.Put([]byte(user), []byte(passwd))
		if eup != nil {
			return eup
		}
		info := UserInfo{User: user, Group: int(global.GROUP_UNDEF)}

		infoData, ejn := json.Marshal(info)
		if ejn != nil {
			return ejn
		}
		bti := tx.Bucket([]byte(global.BUCKET_TOKEN_INFO))
		if bti == nil {
			return fmt.Errorf("bucket:%s is nil", global.BUCKET_TOKEN_INFO)
		}
		token := uuid.New().String()
		eti := bti.Put([]byte(token), []byte(string(infoData)))
		if eti != nil {
			return eti
		}

		but := tx.Bucket([]byte(global.BUCKET_USER_TOKEN))
		if but == nil {
			return fmt.Errorf("bucket:%s is nil", global.BUCKET_USER_TOKEN)
		}
		eut := but.Put([]byte(user), []byte(token))
		return eut
	})
	defer rwLock.Unlock()
	return registerErr
}

func Auth(user, passwd string) error {
	rwLock.Lock()
	authErr := localDB.Update(func(tx *bolt.Tx) error {
		bup := tx.Bucket([]byte(global.BUCKET_USR_PASSWD))
		if bup == nil {
			return fmt.Errorf("bucket:%s is nil", global.BUCKET_USR_PASSWD)
		}
		passwdInDb := string(bup.Get([]byte(user)))
		if passwd != passwdInDb {
			return errors.New("user or password wrong")
		}
		return nil
		// token := uuid.New().String()
		// but := tx.Bucket([]byte(global.BUCKET_USER_TOKEN))
		// eut := but.Put([]byte(user), []byte(token))
		// return eut
	})
	defer rwLock.Unlock()
	return authErr
}

func Login(accessToken string) (UserInfo, error) {
	infoInDb, _ := Get(global.BUCKET_TOKEN_INFO, accessToken)
	if len(infoInDb) == 0 {
		return UserInfo{}, errors.New("access_token error")
	}
	info := UserInfo{}
	json.Unmarshal(infoInDb, &info)
	return info, nil
}

func SetAccessToken(user, access_token string) error {
	return Set(global.BUCKET_USER_TOKEN, user, access_token)
}
func GetAccessToken(user string) (string, error) {
	ret, err := Get(global.BUCKET_USER_TOKEN, user)
	return string(ret), err
}
func SetTokenInfo(token, info string) error {
	return Set(global.BUCKET_TOKEN_INFO, token, info)
}

func Buckets() ([]string, error) {
	ret := []string{}
	rwLock.RLock()
	err := localDB.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			b := []string{string(name)}
			ret = append(ret, b...)
			return nil
		})
	})
	defer rwLock.RUnlock()
	return ret, err
}

func Keys(bucketName string) ([]string, error) {
	ret := []string{}
	rwLock.RLock()
	err := localDB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket:%s is nil", bucketName)
		}
		b.ForEach(func(k, v []byte) error {
			ret = append(ret, string(k))
			ret = append(ret, string(v))
			return nil
		})
		return nil
	})
	defer rwLock.RUnlock()
	return ret, err
}
