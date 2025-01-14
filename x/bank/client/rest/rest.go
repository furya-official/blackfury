package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/furya-official/blackfury/x/bank/types"
)

type SetDenomMetadataProposalRequest struct {
	BaseReq     rest.BaseReq       `json:"base_req" yaml:"base_req"`
	Title       string             `json:"title" yaml:"title"`
	Description string             `json:"description" yaml:"description"`
	Deposit     sdk.Coins          `json:"deposit" yaml:"deposit"`
	Metadata    banktypes.Metadata `json:"metadata" yaml:"metadata"`
}

func SetDenomMetadataProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: banktypes.ModuleName,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req SetDenomMetadataProposalRequest

			if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
				return
			}

			req.BaseReq = req.BaseReq.Sanitize()
			if !req.BaseReq.ValidateBasic(w) {
				return
			}

			from, err := sdk.AccAddressFromBech32(req.BaseReq.From)
			if rest.CheckBadRequestError(w, err) {
				return
			}

			content := &types.SetDenomMetadataProposal{
				Title:       req.Title,
				Description: req.Description,
				Metadata:    req.Metadata,
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, from)
			if rest.CheckBadRequestError(w, err) {
				return
			}

			if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
				return
			}

			tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
		},
	}
}
