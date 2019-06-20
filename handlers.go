package main

import (
  "fmt"
  "log"

  "encoding/json"
  "encoding/base64"
  "net/http"
  "crypto/rand"

  "golang.org/x/crypto/argon2"
  "github.com/shekodn/auth-app/model"
)

type params struct {
    memory      uint32
    iterations  uint32
    parallelism uint8
    saltLength  uint32
    keyLength   uint32
}

func Signup(w http.ResponseWriter, r *http.Request) {

  // Establish the parameters to use for Argon2.
  p := &params{
    memory:      64 * 1024,
    iterations:  3,
    parallelism: 2,
    saltLength:  16,
    keyLength:   32,
  }

  auxUser := model.User{}
  err := json.NewDecoder(r.Body).Decode(&auxUser)

  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  if len(auxUser.Password) < 6 {
    w.WriteHeader(http.StatusPreconditionFailed)
    return
  }

  // Pass the plaintext password and parameters to our generateFromPassword
  // helper function.
  hashedPassword, err := generateFromPassword(auxUser.Password, p)

  if err != nil {
    log.Fatal(err)
  }

  // Insert the username, along with the hashed password into the database
  if err := model.GetDB().Create(&model.User{Username: auxUser.Username, Password: string(hashedPassword)}).Error; err != nil {
    w.WriteHeader(http.StatusInternalServerError)
  }
}

func generateRandomBytes(n uint32) ([]byte, error) {
    b := make([]byte, n)
    _, err := rand.Read(b)
    if err != nil {
        return nil, err
    }

    return b, nil
}

func generateFromPassword(password string, p *params) (encodedHash string, err error) {
    salt, err := generateRandomBytes(p.saltLength)
    if err != nil {
        return "", err
    }

    hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

    // Base64 encode the salt and hashed password.
    b64Salt := base64.RawStdEncoding.EncodeToString(salt)
    b64Hash := base64.RawStdEncoding.EncodeToString(hash)

    // Return a string using the standard encoded hash representation.
    encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

    return encodedHash, nil
}
