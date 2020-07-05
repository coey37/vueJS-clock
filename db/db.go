
package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/VolticFroogo/Froogo/db/dbCredentials"
	"github.com/VolticFroogo/Froogo/helpers"
	"github.com/VolticFroogo/Froogo/models"
	_ "github.com/go-sql-driver/mysql" // Necessary for connecting to MySQL.
)

/*
	Structs and variables
*/

var (
	db *sql.DB
)

// InitDB initializes the Database.
func InitDB() (err error) {
	db, err = sql.Open(dbCredentials.Type, dbCredentials.ConnString)
	if err != nil {
		return
	}

	go jtiGarbageCollector()
	return
}

/*
	Helper functions
*/

func rowExists(query string, args ...interface{}) (exists bool, err error) {
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err = db.QueryRow(query, args...).Scan(&exists)
	return
}

/*
	MySQL DataBase related functions
*/

// StoreRefreshToken generates, stores and then returns a JTI.
func StoreRefreshToken() (jti models.JTI, err error) {
	// No need to duplication check as the JTI takes input from time and are unique.
	jti.JTI, err = helpers.GenerateRandomString(32)
	if err != nil {
		return
	}

	jti.Expiry = time.Now().Add(models.RefreshTokenValidTime).Unix()

	_, err = db.Exec("INSERT INTO jti (jti, expiry) VALUES (?, ?)", jti.JTI, jti.Expiry)
	if err != nil {
		return
	}

	rows, err := db.Query("SELECT id FROM jti WHERE jti=? AND expiry=?", jti.JTI, jti.Expiry)
	if err != nil {
		return
	}

	defer rows.Close()

	rows.Next()
	err = rows.Scan(&jti.ID) // Scan data from query.
	return
}

// GetJTI takes a JTI string and returns the JTI struct.
func GetJTI(jti string) (jtiStruct models.JTI, err error) {
	rows, err := db.Query("SELECT id, expiry FROM jti WHERE jti=?", jti)
	if err != nil {
		return
	}

	defer rows.Close()

	jtiStruct.JTI = jti
	rows.Next()
	err = rows.Scan(&jtiStruct.ID, &jtiStruct.Expiry) // Scan data from query.
	return
}

// CheckJTI returns the validity of a JTI.
func CheckJTI(jti models.JTI) (valid bool, err error) {
	if jti.Expiry > time.Now().Unix() { // Check if token has expired.
		return true, nil // Token is valid.
	}

	_, err = db.Exec("DELETE FROM jti WHERE id=?", jti.ID)
	if err != nil {
		return false, err
	}

	return false, nil // Token is invalid.
}

// DeleteJTI deletes a JTI based on a jti key.
func DeleteJTI(jti string) (err error) {
	_, err = db.Exec("DELETE FROM jti WHERE jti=?", jti)
	return
}

func jtiGarbageCollector() {
	ticker := time.NewTicker(5 * time.Minute) // Tick every five minutes.
	for {
		<-ticker.C
		rows, err := db.Query("SELECT id, jti, expiry FROM jti")
		if err != nil {
			log.Printf("Error querying JTI DB in JTI garbage collector: %v", err)
			return
		}

		defer rows.Close()

		jti := models.JTI{} // Create struct to store a JTI in.
		for rows.Next() {
			err = rows.Scan(&jti.ID, &jti.JTI, &jti.Expiry) // Scan data from query.
			if err != nil {
				log.Printf("Error scanning rows in JTI garbage collector: %v", err)
				return
			}

			_, err := CheckJTI(jti)
			if err != nil {
				log.Printf("Error checking in JTI garbage collector: %v", err)
				return
			}
		}
	}
}

// GetUserFromID retrieves a user from the MySQL database.
func GetUserFromID(uuid int) (user models.User, err error) {
	rows, err := db.Query("SELECT email, username, password, fname, lname, priv, points, steamid, creation FROM users WHERE uuid=?", uuid)
	if err != nil {
		return
	}

	defer rows.Close()

	user.UUID = uuid
	for rows.Next() {
		err = rows.Scan(&user.Email, &user.Username, &user.Password, &user.Fname, &user.Lname, &user.Priv, &user.Points, &user.SteamID, &user.Creation) // Scan data from query.
		if err != nil {
			return
		}
	}

	return
}

// GetUserFromEmail retrieves a user from the MySQL database.
func GetUserFromEmail(email string) (user models.User, err error) {
	rows, err := db.Query("SELECT uuid, username, password, fname, lname, priv, points, steamid, creation FROM users WHERE email=?", email)
	if err != nil {
		return
	}

	defer rows.Close()

	user.Email = email
	for rows.Next() {
		err = rows.Scan(&user.UUID, &user.Username, &user.Password, &user.Fname, &user.Lname, &user.Priv, &user.Points, &user.SteamID, &user.Creation) // Scan data from query.
		if err != nil {
			return
		}
	}

	return
}

// GetUserFromUsername retrieves a user from the MySQL database.
func GetUserFromUsername(username string) (user models.User, err error) {
	rows, err := db.Query("SELECT uuid, email, password, fname, lname, priv, points, steamid, creation FROM users WHERE username=?", username)
	if err != nil {
		return
	}

	defer rows.Close()

	user.Username = username
	for rows.Next() {
		err = rows.Scan(&user.UUID, &user.Email, &user.Password, &user.Fname, &user.Lname, &user.Priv, &user.Points, &user.SteamID, &user.Creation) // Scan data from query.
		if err != nil {
			return
		}
	}

	return
}

// EditUser updates a user.
func EditUser(ID int, Email, Password, Fname, Lname string, Privileges int) (err error) {
	_, err = db.Exec("UPDATE users SET email=?, password=?, fname=?, lname=?, priv=? WHERE uuid=?", Email, Password, Fname, Lname, Privileges, ID)
	return
}

// EditUserNoPassword updates a user without changing the password.
func EditUserNoPassword(ID int, Email, Fname, Lname string, Privileges int) (err error) {
	_, err = db.Exec("UPDATE users SET email=?, fname=?, lname=?, priv=? WHERE uuid=?", Email, Fname, Lname, Privileges, ID)
	return
}

// EditSelf updates a user from settings.
func EditSelf(ID int, Password, Fname, Lname string) (err error) {
	_, err = db.Exec("UPDATE users SET password=?, fname=?, lname=? WHERE uuid=?", Password, Fname, Lname, ID)
	return
}

// EditSelfNoPassword updates a user from settings without changing the password.
func EditSelfNoPassword(ID int, Fname, Lname string) (err error) {
	_, err = db.Exec("UPDATE users SET fname=?, lname=? WHERE uuid=?", Fname, Lname, ID)
	return
}

// NewUser creates a new user.
func NewUser(Email, Password, Fname, Lname string, Privileges int) (id int, err error) {
	creation := time.Now().Unix()

	_, err = db.Exec("INSERT INTO users (email, password, fname, lname, priv, creation) VALUES (?, ?, ?, ?, ?, ?)", Email, Password, Fname, Lname, Privileges, creation)
	if err != nil {
		return
	}

	rows, err := db.Query("SELECT uuid FROM users WHERE email=? AND creation=? ORDER BY uuid DESC", Email, creation)
	if err != nil {
		return
	}

	defer rows.Close()

	rows.Next()
	err = rows.Scan(&id)
	return
}

// DeleteUser deletes a user.
func DeleteUser(ID int) (err error) {
	_, err = db.Exec("DELETE FROM users WHERE uuid=?", ID)
	return
}

// AddEmailVerification adds an email verification code to the DB.
func AddEmailVerification(id string, userUUID int, email string) (err error) {
	exists, err := rowExists("SELECT id FROM email WHERE useruuid=?", userUUID)
	if err != nil {
		return
	}
	if exists {
		_, err = db.Exec("DELETE FROM email WHERE useruuid=?", userUUID)
		if err != nil {
			return
		}
	}

	_, err = db.Exec("INSERT INTO email (uuid, useruuid, email) VALUES (?, ?, ?)", id, userUUID, email)
	return
}

// GetEmailVerification retrieves an email verification information.
func GetEmailVerification(id string) (userUUID int, email string, err error) {
	rows, err := db.Query("SELECT useruuid, email FROM email WHERE uuid=?", id)
	if err != nil {
		return
	}

	defer rows.Close()

	rows.Next()
	err = rows.Scan(&userUUID, &email)
	if err != nil {
		return
	}

	if userUUID != 0 && email != "" {
		_, err = db.Exec("DELETE FROM email WHERE uuid=?", id)
	}

	return
}

// EditSelfEmail updates a user's email after verification.
func EditSelfEmail(uuid int, email string) (err error) {
	_, err = db.Exec("UPDATE users SET email=? WHERE uuid=?", email, uuid)
	return
}

// AddRecovery adds a password recovery code to the DB.
func AddRecovery(id string, userUUID int, email string) (err error) {
	exists, err := rowExists("SELECT id FROM recovery WHERE useruuid=?", userUUID)
	if err != nil {
		return
	}
	if exists {
		_, err = db.Exec("DELETE FROM recovery WHERE useruuid=?", userUUID)
		if err != nil {
			return
		}
	}

	_, err = db.Exec("INSERT INTO recovery (uuid, useruuid, email) VALUES (?, ?, ?)", id, userUUID, email)
	return
}

// GetRecovery retrieves a password recovery code from the DB.
func GetRecovery(id string) (userUUID int, email string, err error) {
	rows, err := db.Query("SELECT useruuid, email FROM recovery WHERE uuid=?", id)
	if err != nil {
		return
	}

	defer rows.Close()

	rows.Next()
	err = rows.Scan(&userUUID, &email)
	if err != nil {
		return
	}

	if userUUID != 0 && email != "" {
		_, err = db.Exec("DELETE FROM recovery WHERE uuid=?", id)
	}

	return
}

// EditPassword updates a user's password after password recovery.
func EditPassword(uuid int, password string) (err error) {
	_, err = db.Exec("UPDATE users SET password=? WHERE uuid=?", password, uuid)
	return
}

func LinkSteam(uuid int, steamID int64) (err error) {
	_, err = db.Exec("UPDATE users SET steamid=? WHERE uuid=?", steamID, uuid)
	return
}

func AddOffer(offer models.Offer) (err error) {
	item, err := json.Marshal(offer.Item)
	if err != nil {
		return
	}

	floatAPI, err := json.Marshal(offer.FloatAPI)
	if err != nil {
		return
	}

	_, err = db.Exec("INSERT INTO offer (senderuuid, receiveruuid, points, item, floatapi, timestamp) VALUES (?, ?, ?, ?, ?, ?)",
		offer.UserUUID, offer.ReceiverUUID, offer.Points, string(item[:]), string(floatAPI[:]), offer.Timestamp)
	return
}

func GetOffers(userUUID int) (offers []models.Offer, err error) {
	rows, err := db.Query("SELECT id, status, senderuuid, receiveruuid, points, item, floatapi, timestamp FROM offer WHERE senderuuid=? OR receiveruuid=? ORDER BY id DESC", userUUID, userUUID)
	if err != nil {
		return
	}

	defer rows.Close()

	var usernames map[int]string
	usernames = make(map[int]string)

	i := 0
	for rows.Next() {
		offers = append(offers, models.Offer{})
		err = rows.Scan(&offers[i].ID, &offers[i].Status, &offers[i].UserUUID, &offers[i].ReceiverUUID, &offers[i].Points, &offers[i].ItemJSON, &offers[i].FloatAPIJSON, &offers[i].Timestamp) // Scan data from query.
		if err != nil {
			return
		}

		err = json.Unmarshal([]byte(offers[i].ItemJSON), &offers[i].Item)
		if err != nil {
			return
		}

		err = json.Unmarshal([]byte(offers[i].FloatAPIJSON), &offers[i].FloatAPI)
		if err != nil {
			return
		}

		// Free RAM from useless JSON.
		offers[i].ItemJSON = ""
		offers[i].FloatAPIJSON = ""

		var otherUserUUID int
		if userUUID == offers[i].UserUUID {
			otherUserUUID = offers[i].ReceiverUUID
		} else {
			otherUserUUID = offers[i].UserUUID
		}

		if username, ok := usernames[otherUserUUID]; ok {
			offers[i].OtherUsername = username
		} else {
			user, err := GetUserFromID(otherUserUUID)
			if err != nil {
				return offers, err
			}

			offers[i].OtherUsername = user.Username
			usernames[offers[i].UserUUID] = user.Username
		}

		i++
	}

	return
}

func GetOffer(offerID, userUUID int) (offer models.Offer, err error) {
	rows, err := db.Query("SELECT status, senderuuid, receiveruuid, points, item, floatapi, timestamp FROM offer WHERE id=?", offerID)
	if err != nil {
		return
	}

	defer rows.Close()

	rows.Next()

	offer.ID = offerID
	err = rows.Scan(&offer.Status, &offer.UserUUID, &offer.ReceiverUUID, &offer.Points, &offer.ItemJSON, &offer.FloatAPIJSON, &offer.Timestamp) // Scan data from query.
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(offer.ItemJSON), &offer.Item)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(offer.FloatAPIJSON), &offer.FloatAPI)
	if err != nil {
		return
	}

	// Free RAM from useless JSON.
	offer.ItemJSON = ""
	offer.FloatAPIJSON = ""

	var otherUserUUID int
	if userUUID == offer.UserUUID {
		otherUserUUID = offer.ReceiverUUID
	} else {
		otherUserUUID = offer.UserUUID
	}

	user, err := GetUserFromID(otherUserUUID)
	if err != nil {
		return
	}

	offer.OtherUsername = user.Username
	return
}

// DeleteOffer deletes an offer.
func DeleteOffer(id int) (err error) {
	_, err = db.Exec("DELETE FROM offer WHERE id=?", id)
	return
}

func SetOfferStatus(status, id int) (err error) {
	_, err = db.Exec("UPDATE offer SET status=? WHERE id=?", status, id)
	return
}