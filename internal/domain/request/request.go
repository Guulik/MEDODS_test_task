package request

type GetTokensRequest struct {
	UserGUID string `query:"user_guid"`
}
type RefreshTokensRequest struct {
	UserGUID string `query:"user_guid"`
}
