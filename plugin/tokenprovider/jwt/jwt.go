package jwt

import (
	"flag"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/haohmaru3000/go_sdk/common"
	tokenprovider "github.com/haohmaru3000/go_sdk/plugin/tokenprovider"
)

type jwtProvider struct {
	name   string
	secret string
}

func NewJWTProvider(name string) *jwtProvider {
	return &jwtProvider{name: name}
}

func (p *jwtProvider) GetPrefix() string {
	return p.Name()
}

func (p *jwtProvider) Get() interface{} {
	return p
}

func (p *jwtProvider) Name() string {
	return p.name
}

func (p *jwtProvider) InitFlags() {
	flag.StringVar(&p.secret, "jwt-secret", "haohmaru3000", "Secret key for generating JWT")
}

func (p *jwtProvider) Configure() error {
	return nil
}

func (p *jwtProvider) Run() error {
	return nil
}

func (p *jwtProvider) Stop() <-chan bool {
	c := make(chan bool)
	go func() {
		c <- true
	}()
	return c
}

func (j *jwtProvider) SecretKey() string {
	return j.secret
}

func (j *jwtProvider) Generate(data tokenprovider.TokenPayload, expiry int) (tokenprovider.Token, error) {
	convertedData, _ := data.(*common.TokenPayload)

	// generate the JWT
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims{
		*convertedData,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Second * time.Duration(expiry))),
			IssuedAt:  jwt.NewNumericDate(time.Now().Local()),
		},
	})

	myToken, err := t.SignedString([]byte(j.secret)) // Convert key to []byte and send in SignedString()
	if err != nil {
		return nil, tokenprovider.ErrEncodingToken
	}

	// return the token
	return &token{
		Token:     myToken,
		Expiry:    expiry,
		CreatedAt: time.Now(),
	}, nil
}

func (j *jwtProvider) Validate(myToken string) (tokenprovider.TokenPayload, error) {
	res, err := jwt.ParseWithClaims(myToken, &myClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, tokenprovider.ErrInvalidToken
	}

	// validate the token
	if !res.Valid {
		return nil, tokenprovider.ErrInvalidToken
	}

	claims, ok := res.Claims.(*myClaims)
	if !ok {
		return nil, tokenprovider.ErrInvalidToken
	}

	// return the token
	return &claims.Payload, nil
}

func (j *jwtProvider) String() string {
	return "JWT implement Provider"
}

type myClaims struct {
	Payload common.TokenPayload `json:"payload"`
	jwt.RegisteredClaims
}

type token struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	Expiry    int       `json:"expiry"`
}

func (t *token) GetToken() string {
	return t.Token
}
