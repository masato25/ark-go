package core

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"testing"
)

func TestListDelegates(t *testing.T) {
	arkapi := NewArkClient(nil)
	params := DelegateQueryParams{OrderBy: "rate:asc", Limit: 51}

	deleResponse, _, err := arkapi.ListDelegates(params)
	if deleResponse.Success {
		log.Println(t.Name(), "Success, returned ", deleResponse.TotalCount, "delegates, received:", len(deleResponse.Delegates))
		/*for _, element := range deleResponse.Delegates {
			log.Println(element.Username)
		}*/
	} else {
		t.Error(err.Error())
	}
}

func TestGetDelegateUsername(t *testing.T) {
	arkapi := NewArkClient(nil)
	deleUserName := "chris"

	if EnvironmentParams.Network.Type == DEVNET {
		deleUserName = "d_chris"
	}

	params := DelegateQueryParams{UserName: deleUserName}

	deleResponse, _, err := arkapi.GetDelegate(params)
	if deleResponse.Success {

		out, _ := json.Marshal(deleResponse.SingleDelegate)
		log.Println(t.Name(), "Success, returned", string(out))

	} else {
		t.Error(err.Error())
	}
}

func TestGetDelegatePubKey(t *testing.T) {
	arkapi := NewArkClient(nil)
	deleKey := "03e6397071866c994c519f114a9e7957d8e6f06abc2ca34dc9a96b82f7166c2bf9"

	if EnvironmentParams.Network.Type == DEVNET {
		deleKey = "02bcfa0951a92e7876db1fb71996a853b57f996972ed059a950d910f7d541706c9"
	}

	params := DelegateQueryParams{PublicKey: deleKey}

	deleResponse, _, err := arkapi.GetDelegate(params)
	if deleResponse.Success {

		out, _ := json.Marshal(deleResponse.SingleDelegate)
		log.Println(t.Name(), "Success, returned", string(out))

	} else {
		t.Error(err.Error())
	}
}

func TestGetDelegateVoters(t *testing.T) {
	arkapi := NewArkClient(nil)
	//params := DelegateQueryParams{PublicKey: "027acdf24b004a7b1e6be2adf746e3233ce034dbb7e83d4a900f367efc4abd0f21"}
	params := DelegateQueryParams{PublicKey: "02c7455bebeadde04728441e0f57f82f972155c088252bf7c1365eb0dc84fbf5de"}

	deleResponse, _, err := arkapi.GetDelegateVoters(params)
	if deleResponse.Success {

		//calculating vote weight
		balance := 0
		for _, element := range deleResponse.Accounts {
			intBalance, _ := strconv.Atoi(element.Balance)
			balance += intBalance
		}

		log.Println(t.Name(), "Success, returned", len(deleResponse.Accounts), "voters for delegate with weight", balance)

	} else {
		t.Error(err.Error())
	}
}

func TestGetDelegateVoteWeight(t *testing.T) {
	arkapi := NewArkClient(nil)
	//params := DelegateQueryParams{PublicKey: "027acdf24b004a7b1e6be2adf746e3233ce034dbb7e83d4a900f367efc4abd0f21"}
	params := DelegateQueryParams{PublicKey: "02c7455bebeadde04728441e0f57f82f972155c088252bf7c1365eb0dc84fbf5de"}

	voteWeight, _, _ := arkapi.GetDelegateVoteWeight(params)

	log.Println(t.Name(), "Success, returned delegate vote weight is", voteWeight)
}

func TestCalculcateVotersProfit(t *testing.T) {
	arkapi := NewArkClient(nil)
	//deleKey := "027acdf24b004a7b1e6be2adf746e3233ce034dbb7e83d4a900f367efc4abd0f21"
	deleKey := "02c7455bebeadde04728441e0f57f82f972155c088252bf7c1365eb0dc84fbf5de" //jar
	if EnvironmentParams.Network.Type == DEVNET {
		deleKey = "02bcfa0951a92e7876db1fb71996a853b57f996972ed059a950d910f7d541706c9"
	}

	params := DelegateQueryParams{PublicKey: deleKey}

	votersEarnings := arkapi.CalculateVotersProfit(params, 0.70, "", "", false, 0.0, true)

	log.Println(t.Name(), "Success", len(votersEarnings))
	//log.Println(t.Name(), "Success", votersEarnings)
	sumEarned := 0.0
	sumRatio := 0.0
	sumShareEarned := 0.0
	feeAmount := float64(len(votersEarnings)) * (float64(EnvironmentParams.Fees.Send) / SATOSHI)
	for _, element := range votersEarnings {
		log.Println(fmt.Sprintf("|%s|%15.8f|%15.8f|%15.8f|%15.8f|%4d|%25d|",
			element.Address,
			element.VoteWeight,
			element.VoteWeightShare,
			element.EarnedAmount100,
			element.EarnedAmountXX,
			element.VoteDuration,
			int(element.EarnedAmountXX*SATOSHI)))

		sumEarned += element.EarnedAmount100
		sumShareEarned += element.EarnedAmountXX
		sumRatio += element.VoteWeightShare
	}
	log.Println("Delegate wallet amount: ", sumEarned, "Ratio calc check sum: ", sumRatio, "Amount to voters: ", sumShareEarned, "Ratio shared: ", float64(sumShareEarned)/float64(sumEarned), "Lottery:", int64((sumEarned-sumShareEarned-feeAmount)*SATOSHI))
	log.Println(fmt.Sprintf("Payment fees: %2.2f", feeAmount))
}

func TestGetForgedData(t *testing.T) {
	arkapi := NewArkClient(nil)
	deleKey := "03e6397071866c994c519f114a9e7957d8e6f06abc2ca34dc9a96b82f7166c2bf9"
	if EnvironmentParams.Network.Type == DEVNET {
		deleKey = "02bcfa0951a92e7876db1fb71996a853b57f996972ed059a950d910f7d541706c9"
	}
	params := DelegateQueryParams{PublicKey: deleKey}

	resp, _, err := arkapi.GetForgedData(params)

	if resp.Success {
		log.Println(t.Name(), "Delegate forged:", resp.Forged, "fees:", resp.Fees, "rewards", resp.Rewards)
	} else {
		t.Error(err.Error())
	}

}

func TestGetVoteDuration(t *testing.T) {
	arkapi := NewArkClient(nil)
	//arkapi = arkapi.SetActiveConfiguration(DEVNET)
	//deleKey := "03d9ed6e7f29daf12ef925d4ce5753aade23c8cfd52a0427240fb30ad6ec232fed"

	duration := arkapi.GetVoteDuration("AdjCcEjtj9rAAVUWZofuMevEk69sfhRDvU")
	log.Println("Vote duration", duration)
}
