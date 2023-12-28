package btobet

// Error bto bet Error
type Error struct {
	Description string
	ErrorNo     int
}

// RegistrationResponse ...
type RegistrationResponse struct {
	Errors       []Error
	IsSuccessful bool
}

// LoginRequest web login request
type LoginRequest struct {
	Username  string `json:"login,omitempty" binding:"required"`
	Password  string `json:"password,omitempty" binding:"required"`
	IPaddress string `json:"ipAddress,omitempty" binding:"required"`
}

// LoginResponse from bto bet
type LoginResponse struct {
	Errors       []Error
	IsSuccessful bool
}

// CustomerDetails response from bto bet
type CustomerDetails struct {
	Errors       []Error
	IsSuccessful bool
	Customer     Customer
	Metadata     map[string]any
}

// Customer customer data
type Customer struct {
	Balance Balance
	Account Account
}

// Balance customer account balance
type Balance struct {
	Bonus       float32
	CurrencyISO string
	Real        float32
}

// Account the customer account details
type Account struct {
	InternalID string
}

// BetSlipItem is a slip within a users jackpot bet
type BetSlipItem struct {
	EventCode        string
	OutcomeShortCode string
}

// BetSlipRequest used to submit bets to clients
type BetSlipRequest struct {
	Mobile       string
	Stake        string
	BetSlipItems []BetSlipItem
}

// BetSlipResponse response from service provider
type BetSlipResponse struct {
	BetSlipID       int
	ResponseCode    int
	Timestamp       string
	ResponseMessage string
	Status          string
}

// BetStatusResponse ...
type BetStatusResponse struct {
	BetSlipID    int            `json:"BetSlipID,omitempty"`
	PotentialWin map[string]any `json:"PotentialWin,omitempty"`
	Status       string         `json:"Status,omitempty"`
}

// MarketResponse ...
type MarketResponse struct {
	EventCode       int      `json:"EventCode,omitempty"`
	ResponseCode    int      `json:"ResponseCode,omitempty"`
	ResponseMessage string   `json:"ResponseMessage,omitempty"`
	Markets         []Market `json:"Markets,omitempty"`
}

// Market ...
type Market struct {
	Name      string    `json:"Name,omitempty"`
	ShortCode string    `json:"ShortCode,omitempty"`
	Outcomes  []Outcome `json:"Outcomes,omitempty"`
}

// Outcome ...
type Outcome struct {
	Outcome   string  `json:"Outcome,omitempty"`
	Odd       float32 `json:"Odd,omitempty"`
	ShortCode string  `json:"ShortCode,omitempty"`
}
