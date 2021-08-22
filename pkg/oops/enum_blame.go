package oops

type Blame uint8

const (
	BlameUnknown Blame = iota

	BlameClient
	BlameServer
	BlameDeveloper
	BlameThirdParty

	blameMAX
)

func (e Blame) String() string {
	code, ok := mapBlameToCode[e]
	if !ok {
		return "UNDEFINED"
	}
	return code
}

var mapBlameToCode = map[Blame]string{
	BlameUnknown:    "UNKNOWN",
	BlameClient:     "CLIENT",
	BlameServer:     "SERVER",
	BlameDeveloper:  "DEVELOPER",
	BlameThirdParty: "THIRD_PARTY",
}
