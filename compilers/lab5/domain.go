package main

type DomainTag int

var (
	WS DomainTag = 22
	K1 DomainTag = 1
	K2 DomainTag = 2
	K3 DomainTag = 3
	K4 DomainTag = 4
	K5 DomainTag = 5
	K6 DomainTag = 6
	K7 DomainTag = 7
	K8 DomainTag = 8
	KW DomainTag = 10
	O1 DomainTag = 11
	OP DomainTag = 13
	ID DomainTag = 14
	C1 DomainTag = 15
	C2 DomainTag = 16
	C3 DomainTag = 17
	C4 DomainTag = 18
	CM DomainTag = 19
	ST DomainTag = 20
	NU DomainTag = 21

	UK DomainTag = -1
)

var (
	Spaces    = 0
	Digit     = 1
	Letter    = 2
	BackSlash = 3
	Slash     = 4
	LetterU   = 5
	LetterP   = 6
	LetterD   = 7
	LetterA   = 8
	LetterT   = 9
	LetterE   = 10
	LetterW   = 11
	LetterH   = 12
	LetterR   = 13
	NotEqual  = 14
	NewLine   = 15
	Equals    = 16
	Others    = 17

	Eof = -1
)

var transes = map[DomainTag][18]DomainTag{
	//   SP  DI  LE   /  \    u   p   d   a   t   e   w   h   r  !   \n   =
	ST: {WS, NU, ID, C1, UK, K1, ID, ID, ID, ID, ID, K6, ID, ID, O1, WS, O1, UK},
	WS: {WS, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, WS, UK, UK},
	K1: {UK, ID, ID, UK, UK, ID, K2, ID, ID, ID, ID, ID, ID, ID, UK, UK, UK, UK},
	K2: {UK, ID, ID, UK, UK, ID, ID, K3, ID, ID, ID, ID, ID, ID, UK, UK, UK, UK},
	K3: {UK, ID, ID, UK, UK, ID, ID, ID, K4, ID, ID, ID, ID, ID, UK, UK, UK, UK},
	K4: {UK, ID, ID, UK, UK, ID, ID, ID, ID, K5, ID, ID, ID, ID, UK, UK, UK, UK},
	K5: {UK, ID, ID, UK, UK, ID, ID, ID, ID, ID, KW, ID, ID, ID, UK, UK, UK, UK},
	K6: {UK, ID, ID, UK, UK, ID, ID, ID, ID, ID, ID, ID, K7, ID, UK, UK, UK, UK},
	K7: {UK, ID, ID, UK, UK, ID, ID, ID, ID, ID, K8, ID, ID, ID, UK, UK, UK, UK},
	K8: {UK, ID, ID, UK, UK, ID, ID, ID, ID, ID, ID, ID, ID, K5, UK, UK, UK, UK},
	O1: {UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, OP, UK},
	OP: {UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK},
	ID: {UK, ID, ID, UK, UK, ID, ID, ID, ID, ID, ID, ID, ID, ID, UK, UK, UK, UK},
	C1: {UK, UK, UK, UK, C2, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK},
	C2: {C2, C2, C2, C2, C3, C2, C2, C2, C2, C2, C2, C2, C2, C2, C2, CM, C2, C2},
	C3: {C2, C2, C2, C2, C2, C2, C2, C2, C2, C2, C2, C2, C2, C2, C2, C4, C2, C2},
	C4: {C2, C2, C2, C2, C2, C2, C2, C2, C2, C2, C2, C2, C2, C2, C2, CM, C2, C2},
	CM: {UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK},
	KW: {UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK},
	NU: {UK, NU, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK, UK},
}
