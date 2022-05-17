package auth

import (
	"fmt"
	"testing"
)

type MockCache struct {
	token string
}

func (m MockCache) GetAccessToken() string {
	return m.token
}

func (m MockCache) SetAccessToken(token string, expiresIn int) {
	panic("not implemented") // TODO: Implement
}

func TestAuthPhoneNumber(t *testing.T) {
	auth := &Auth{
		cache: MockCache{
			token: "56_Hv1lKyrBZIvfsfSfNhSiHaRbAYT2oI-YqHaHxHc0W4lE5ZzlJDvV5boX7bykR3_fY5IU7pKu9hfvClvi2UtxlGz7__pDee6XdaJARP3usbstle-Z93SJzE3nK81y5k5qI386zXJGQjXhzPsGJSLeAFACLE",
		},
	}

	r, err := auth.GetPhoneNumber("b10dd98ad08333cd95682b14271d27243cda142b28bf198b443ab2eeb4a435c9")
	if err != nil {
		t.Fatalf("%+v\n", err)
	}
	fmt.Printf("%+v\n", r)
}
