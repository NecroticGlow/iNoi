package model

import "testing"

func TestValidateRawPasswordAcceptsINoiAndOpenListStaticSalts(t *testing.T) {
	const password = "correct horse battery staple"
	for name, staticHash := range map[string]string{
		"inoi":     StaticHash(password),
		"openlist": OpenListStaticHash(password),
	} {
		t.Run(name, func(t *testing.T) {
			user := &User{Salt: "per-user-salt", PwdHash: HashPwd(staticHash, "per-user-salt")}
			if err := user.ValidateRawPassword(password); err != nil {
				t.Fatalf("expected compatible password hash to validate: %v", err)
			}
			if err := user.ValidateRawPassword("wrong password"); err == nil {
				t.Fatal("expected an incorrect password to be rejected")
			}
		})
	}
}

func TestSetPasswordUsesINoiStaticSalt(t *testing.T) {
	user := (&User{}).SetPassword("secret")
	if got, want := user.PwdHash, HashPwd(StaticHash("secret"), user.Salt); got != want {
		t.Fatalf("SetPassword used an unexpected static salt: got %q want %q", got, want)
	}
}
