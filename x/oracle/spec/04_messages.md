<!--
order: 4
-->

# Messages

## MsgAggregateExchangeRatePrevote

`Hash` is a hex string generated by the leading 20 bytes of the SHA256 hash (hex string) of a string of the format `{salt}:{exchange rate}{denom},...,{exchange rate}{denom}:{voter}`, the metadata of the actual `MsgAggregateExchangeRateVote` to follow in the next `VotePeriod`. You can use the `GetAggregateVoteHash()` function to help encode this hash. Note that since in the subsequent `MsgAggregateExchangeRateVote`, the salt will have to be revealed, the salt used must be regenerated for each prevote submission.

```go
// MsgAggregateExchangeRatePrevote - struct for aggregate prevoting on the ExchangeRateVote.
// The purpose of aggregate prevote is to hide vote exchange rates with hash
// which is formatted as hex string in SHA256("{salt}:{exchange rate}{denom},...,{exchange rate}{denom}:{voter}")
type MsgAggregateExchangeRatePrevote struct {
	Hash      AggregateVoteHash 
	Feeder    sdk.AccAddress    
	Validator sdk.ValAddress    
}
```

## MsgAggregateExchangeRateVote

The `MsgAggregateExchangeRateVote` contains the actual exchange rates vote. The `Salt` parameter must match the salt used to create the prevote, otherwise the voter cannot be rewarded.

```go
// MsgAggregateExchangeRateVote - struct for voting on various denominations' exchange rate against USD.
type MsgAggregateExchangeRateVote struct {
	Salt          string
	ExchangeRates string
	Feeder        sdk.AccAddress 
	Validator     sdk.ValAddress 
}
```


## MsgDelegateFeedConsent

Validators may also select to delegate voting rights to another key to prevent the block signing key from being kept online. To do so, they must submit a `MsgDelegateFeedConsent`, delegating their oracle voting rights to a `Delegate` that sign `MsgAggregateExchangeRatePrevote` and `MsgAggregateExchangeRateVote` on behalf of the validator.

> Delegate validators will likely require you to deposit some funds (in Black or Fury) which they can use to pay fees, sent in a separate MsgSend. This agreement is made off-chain and not enforced by the Blackfury protocol.

The `Operator` field contains the operator address of the validator (prefixed `mervaloper-`). The `Delegate` field is the account address (prefixed `black-`) of the delegate account that will be submitting exchange rate related votes and prevotes on behalf of the `Operator`.

```go
// MsgDelegateFeedConsent - struct for delegating oracle voting rights to another address.
type MsgDelegateFeedConsent struct {
	Operator sdk.ValAddress 
	Delegate sdk.AccAddress 
}
```