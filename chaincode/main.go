// Copyright 2018 Owkin, inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"chaincode/errors"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// SubstraChaincode is a Receiver for Chaincode shim functions
type SubstraChaincode struct {
}

// Create a global logger for the chaincode. Its default level is Info
var logger = shim.NewLogger("substra-chaincode")

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
// TODO!!!!
func (t *SubstraChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal
	args := stub.GetStringArgs()
	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting nothing...")
	}
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode.
func (t *SubstraChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	start := time.Now()
	// Log all input for potential debug later on.
	logger.Infof("Args received by the chaincode: %#v", stub.GetStringArgs())

	// Seed with a timestamp from the channel header so the chaincode's output
	// stay determinist for each transaction. It's necessary because endorsers
	// will compare their own output to the proposal.
	timestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return formatErrorResponse(err)
	}
	seedTime := time.Unix(timestamp.GetSeconds(), int64(timestamp.GetNanos()))
	rand.Seed(seedTime.UnixNano())

	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	db := NewLedgerDB(stub)

	var result interface{}
	switch fn {
	case "createComputePlan":
		result, err = createComputePlan(db, args)
	case "createTesttuple":
		result, err = createTesttuple(db, args)
	case "createTraintuple":
		result, err = createTraintuple(db, args)
	case "createCompositeTraintuple":
		result, err = createCompositeTraintuple(db, args)
	case "createAggregatetuple":
		result, err = createAggregatetuple(db, args)
	case "cancelComputePlan":
		result, err = cancelComputePlan(db, args)
	case "logFailTest":
		result, err = logFailTest(db, args)
	case "logFailTrain":
		result, err = logFailTrain(db, args)
	case "logFailCompositeTrain":
		result, err = logFailCompositeTrain(db, args)
	case "logFailAggregate":
		result, err = logFailAggregate(db, args)
	case "logStartTest":
		result, err = logStartTest(db, args)
	case "logStartTrain":
		result, err = logStartTrain(db, args)
	case "logStartCompositeTrain":
		result, err = logStartCompositeTrain(db, args)
	case "logStartAggregate":
		result, err = logStartAggregate(db, args)
	case "logSuccessTest":
		result, err = logSuccessTest(db, args)
	case "logSuccessTrain":
		result, err = logSuccessTrain(db, args)
	case "logSuccessCompositeTrain":
		result, err = logSuccessCompositeTrain(db, args)
	case "logSuccessAggregate":
		result, err = logSuccessAggregate(db, args)
	case "queryAlgo":
		result, err = queryAlgo(db, args)
	case "queryAlgos":
		result, err = queryAlgos(db, args)
	case "queryCompositeAlgo":
		result, err = queryCompositeAlgo(db, args)
	case "queryCompositeAlgos":
		result, err = queryCompositeAlgos(db, args)
	case "queryAggregateAlgo":
		result, err = queryAggregateAlgo(db, args)
	case "queryAggregateAlgos":
		result, err = queryAggregateAlgos(db, args)
	case "queryDataManager":
		result, err = queryDataManager(db, args)
	case "queryDataManagers":
		result, err = queryDataManagers(db, args)
	case "queryDataSamples":
		result, err = queryDataSamples(db, args)
	case "queryDataset":
		result, err = queryDataset(db, args)
	case "queryFilter":
		result, err = queryFilter(db, args)
	case "queryModelDetails":
		result, err = queryModelDetails(db, args)
	case "queryModelPermissions":
		result, err = queryModelPermissions(db, args)
	case "queryModels":
		result, err = queryModels(db, args)
	case "queryObjective":
		result, err = queryObjective(db, args)
	case "queryObjectiveLeaderboard":
		result, err = queryObjectiveLeaderboard(db, args)
	case "queryObjectives":
		result, err = queryObjectives(db, args)
	case "queryTesttuple":
		result, err = queryTesttuple(db, args)
	case "queryTesttuples":
		result, err = queryTesttuples(db, args)
	case "queryTraintuple":
		result, err = queryTraintuple(db, args)
	case "queryCompositeTraintuple":
		result, err = queryCompositeTraintuple(db, args)
	case "queryAggregatetuple":
		result, err = queryAggregatetuple(db, args)
	case "queryTraintuples":
		result, err = queryTraintuples(db, args)
	case "queryCompositeTraintuples":
		result, err = queryCompositeTraintuples(db, args)
	case "queryAggregatetuples":
		result, err = queryAggregatetuples(db, args)
	case "queryComputePlan":
		result, err = queryComputePlan(db, args)
	case "queryComputePlans":
		result, err = queryComputePlans(db, args)
	case "registerAlgo":
		result, err = registerAlgo(db, args)
	case "registerCompositeAlgo":
		result, err = registerCompositeAlgo(db, args)
	case "registerAggregateAlgo":
		result, err = registerAggregateAlgo(db, args)
	case "registerDataManager":
		result, err = registerDataManager(db, args)
	case "registerDataSample":
		result, err = registerDataSample(db, args)
	case "registerObjective":
		result, err = registerObjective(db, args)
	case "updateComputePlan":
		result, err = updateComputePlan(db, args)
	case "updateDataManager":
		result, err = updateDataManager(db, args)
	case "updateDataSample":
		result, err = updateDataSample(db, args)
	case "registerNode":
		result, err = registerNode(db, args)
	case "queryNodes":
		result, err = queryNodes(db, args)
	default:
		err = errors.BadRequest("function \"%s\" not implemented", fn)
	}
	// Invoke duration
	duration := int(time.Since(start).Nanoseconds()) / 1e6
	logger.Infof("Response from chaincode (in %dms): %#v, error: %s", duration, result, err)
	// Return the result as success payload
	if err != nil {
		return formatErrorResponse(err)
	}
	// Send event if there is any. It's done in one batch since we can only send
	// one event per call
	err = db.SendEvent()
	if err != nil {
		return formatErrorResponse(errors.Internal("could not send event: %s", err.Error()))
	}
	// Marshal to json the smartcontract result
	resp, err := json.Marshal(result)
	if err != nil {
		return formatErrorResponse(errors.Internal("could not format response: %s", err.Error()))
	}
	return shim.Success(resp)
}

func formatErrorResponse(err error) peer.Response {
	e := errors.Wrap(err)
	status := e.HTTPStatusCode()

	errStruct := map[string]interface{}{
		"error": e.Error(),
		// Serialize status in the message until fabric-sdk-py allows subtrabac to
		// access the status
		"status": status,
	}
	for k, v := range e.GetContext() {
		errStruct[k] = v
	}

	payload, _ := json.Marshal(errStruct)
	return peer.Response{
		Message: string(payload),
		Payload: payload,
		Status:  int32(status),
	}
}

// main function starts up the chaincode in the container during instantiate
func main() {
	// TODO use the same level as the shim or an env variable
	logger.SetLevel(shim.LogDebug)
	if err := shim.Start(new(SubstraChaincode)); err != nil {
		fmt.Printf("Error starting SubstraChaincode chaincode: %s", err)
	}
}
