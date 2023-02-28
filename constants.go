package btobet

const (
	registerCustomerURL   = "https://api-bo-stm04.btobet.games/Services/GamingPortalService.svc/CustomerCreateAccountV4"
	loginURL              = "https://api-bo-stm04.btobet.games/Services/GamingPortalService.svc/CustomerLoginAccountV2"
	getCustomerDetailsURL = "https://api-bo-stm04.btobet.games/Services/GamingPortalService.svc/CustomerGetDetailsAlternativeV2"
	addPaymentAccountURL  = "https://api-bo-stm04.btobet.games/Services/GamingPortalService.svc/CustomerAddPaymentAccountsV2"
	getCustomerBonusesURL = "https://api-bo-stm04.btobet.games/Services/GamingPortalService.svc/BonusGetCurrentBonusesForCustomer"
	withdrawURL           = "https://payment-bo-stm04.btobet.games/Services/PosPayment.svc/AgentWithdrawalProcessV3"
	placeBetURL           = "https://sports-stm04-core.btobet.games/rest/smsbetting/Place"
	checkSlipURL          = "https://sports-stm04-core.btobet.games/rest/smsbetting/CheckBetSlip?mobile=%s&betslipid=%s"
	getMarketsURL         = "https://sports-stm04-core.btobet.games/rest/smsbetting/GetMarkets?mobile=%s&eventCode=%s&marketCode=%s"
)
