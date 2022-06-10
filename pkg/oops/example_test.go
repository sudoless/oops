package oops_test

import (
	"fmt"
	"go.sdls.io/oops/pkg/oops"
)

func ExampleDefine() {
	var (
		ErrAuthMissing        = oops.Define().Type("auth").Code("auth_missing").StatusCode(401)
		ErrAuthBadCredentials = oops.Define().Type("auth").Code("auth_bad").StatusCode(401)
		ErrAuthExpired        = oops.Define().Type("auth").Code("auth_expired").StatusCode(401)
	)

	fmt.Println(ErrAuthMissing.Yeet("please provide the 'Authorization' header").String())
	fmt.Println(ErrAuthBadCredentials.Yeet("invalid username and/or password").String())
	fmt.Println(ErrAuthExpired.Yeet("token expired at '%s'", "2022-06-14").String())

	// Output:
	// [auth] auth_missing : please provide the 'Authorization' header
	// [auth] auth_bad : invalid username and/or password
	// [auth] auth_expired : token expired at '2022-06-14'
}

func ExampleErrorDefined_Group() {
	var (
		errAuthGroup = oops.Define().Type("auth").StatusCode(401).Group().PrefixCode("auth_")

		ErrAuthMissing        = errAuthGroup.Code("missing")
		ErrAuthBadCredentials = errAuthGroup.Code("bad")
		ErrAuthExpired        = errAuthGroup.Code("expired")
	)

	fmt.Println(ErrAuthMissing.Yeet("please provide the 'Authorization' header").String())
	fmt.Println(ErrAuthBadCredentials.Yeet("invalid username and/or password").String())
	fmt.Println(ErrAuthExpired.Yeet("token expired at '%s'", "2022-06-14").String())

	// Output:
	// [auth] auth_missing : please provide the 'Authorization' header
	// [auth] auth_bad : invalid username and/or password
	// [auth] auth_expired : token expired at '2022-06-14'
}
