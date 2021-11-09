package accounts

type Account struct {
	CurrentBalance    Balance
	HistoricalBalance map[int]Balance
}

type Balance struct {
	Available float64
	Pledged   float64
	Locked    float64
	GasFees   float64
	BurnFees  float64
	Penalties float64
}
