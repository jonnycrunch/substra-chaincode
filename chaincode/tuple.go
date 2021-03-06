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
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

// List of the possible tuple's status
const (
	StatusDoing    = "doing"
	StatusTodo     = "todo"
	StatusWaiting  = "waiting"
	StatusFailed   = "failed"
	StatusDone     = "done"
	StatusCanceled = "canceled"
	// The status aborted is still under discussion so the logic is already
	// implemented but it's value is the same as canceled for now.
	StatusAborted = "canceled"
)

// ------------------------------------------------
// Smart contracts related to multiple tuple types
// ------------------------------------------------

// queryModelDetails returns info about the testtuple and algo related to a traintuple
func queryModelDetails(db *LedgerDB, args []string) (outModelDetails outputModelDetails, err error) {
	inp := inputKey{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// get traintuple type
	traintupleType, err := db.GetAssetType(inp.Key)
	if err != nil {
		return
	}
	switch traintupleType {
	case TraintupleType:
		var out outputTraintuple
		out, err = getOutputTraintuple(db, inp.Key)
		if err != nil {
			return
		}
		outModelDetails.Traintuple = &out
	case CompositeTraintupleType:
		var out outputCompositeTraintuple
		out, err = getOutputCompositeTraintuple(db, inp.Key)
		if err != nil {
			return
		}
		outModelDetails.CompositeTraintuple = &out
	case AggregatetupleType:
		var out outputAggregatetuple
		out, err = getOutputAggregatetuple(db, inp.Key)
		if err != nil {
			return
		}
		outModelDetails.Aggregatetuple = &out
	}

	// get certified and non-certified testtuples related to traintuple
	testtupleKeys, err := db.GetIndexKeys("testtuple~traintuple~certified~key", []string{"testtuple", inp.Key})
	if err != nil {
		return
	}
	for _, testtupleKey := range testtupleKeys {
		// get testtuple and serialize it
		var outputTesttuple outputTesttuple
		outputTesttuple, err = getOutputTesttuple(db, testtupleKey)
		if err != nil {
			return
		}

		if outputTesttuple.Certified {
			outModelDetails.Testtuple = outputTesttuple
		} else {
			outModelDetails.NonCertifiedTesttuples = append(outModelDetails.NonCertifiedTesttuples, outputTesttuple)
		}
	}
	return
}

// queryModels returns all traintuples and associated testuples
func queryModels(db *LedgerDB, args []string) (outModels []outputModel, err error) {
	outModels = []outputModel{}

	if len(args) != 0 {
		err = errors.BadRequest("incorrect number of arguments, expecting nothing")
		return
	}

	// populate from regular traintuples
	traintupleKeys, err := db.GetIndexKeys("traintuple~algo~key", []string{"traintuple"})
	if err != nil {
		return
	}
	for _, traintupleKey := range traintupleKeys {
		var outputModel outputModel
		var out outputTraintuple

		out, err = getOutputTraintuple(db, traintupleKey)
		if err != nil {
			return
		}
		outputModel.Traintuple = &out
		outModels = append(outModels, outputModel)
	}

	// populate from composite traintuples
	compositeTraintupleKeys, err := db.GetIndexKeys("compositeTraintuple~algo~key", []string{"compositeTraintuple"})
	if err != nil {
		return
	}
	for _, compositeTraintupleKey := range compositeTraintupleKeys {
		var outputModel outputModel
		var out outputCompositeTraintuple

		out, err = getOutputCompositeTraintuple(db, compositeTraintupleKey)
		if err != nil {
			return
		}
		outputModel.CompositeTraintuple = &out
		outModels = append(outModels, outputModel)
	}

	// populate from composite traintuples
	aggregatetupleKeys, err := db.GetIndexKeys("aggregatetuple~algo~key", []string{"aggregatetuple"})
	if err != nil {
		return
	}
	for _, aggregatetupleKey := range aggregatetupleKeys {
		var outputModel outputModel
		var out outputAggregatetuple

		out, err = getOutputAggregatetuple(db, aggregatetupleKey)
		if err != nil {
			return
		}
		outputModel.Aggregatetuple = &out
		outModels = append(outModels, outputModel)
	}

	return
}

func queryModelPermissions(db *LedgerDB, args []string) (outputPermissions, error) {
	var out outputPermissions
	inp := inputKey{}
	err := AssetFromJSON(args, &inp)
	if err != nil {
		return out, err
	}
	modelHash := inp.Key
	keys, err := db.GetIndexKeys("tuple~modelHash~key", []string{"tuple", modelHash})
	if err != nil {
		return out, err
	}
	if len(keys) == 0 {
		return out, errors.NotFound("Could not find a model for hash %s", modelHash)
	}
	tupleKey := keys[0]
	tupleType, err := db.GetAssetType(tupleKey)
	if err != nil {
		return out, errors.Internal(err, "queryModelPermissions: could not retrieve model type with tupleKey %s", tupleKey)
	}

	// By default model is public processable
	modelPermissions := Permissions{}
	modelPermissions.Process.Public = true

	// get out-model and permissions from parent for head model in composite traintuple
	if tupleType == CompositeTraintupleType {
		tuple, err := db.GetCompositeTraintuple(tupleKey)
		if err != nil {
			return out, errors.Internal(err, "queryModelPermissions:")
		}
		// if `modelHash` refers to the head out-model, return the head out-model permissions
		// if `modelHash` refers to the trunk out-model, do nothing (default to "public processable")
		if tuple.OutHeadModel.OutModel.Hash == modelHash {
			modelPermissions = tuple.OutHeadModel.Permissions
		}
	}
	out.Fill(modelPermissions)
	return out, nil
}

// ----------------------------------------------------------
// Utils for smartcontracts related to  multiple tuple types
// ----------------------------------------------------------

func validateTupleOwner(db *LedgerDB, worker string) error {
	txCreator, err := GetTxCreator(db.cc)
	if err != nil {
		return err
	}
	if txCreator != worker {
		return errors.Forbidden("%s is not allowed to update tuple (%s)", txCreator, worker)
	}
	return nil
}

// check validity of traintuple update: consistent status and agent submitting the transaction
func checkUpdateTuple(db *LedgerDB, worker string, oldStatus string, newStatus string) error {
	if StatusAborted == newStatus {
		return nil
	}

	statusPossibilities := map[string]string{
		StatusWaiting: StatusTodo,
		StatusTodo:    StatusDoing,
		StatusDoing:   StatusDone}
	if statusPossibilities[oldStatus] != newStatus && newStatus != StatusFailed {
		return errors.BadRequest("cannot change status from %s to %s", oldStatus, newStatus)
	}
	return nil
}

// HashForKey to generate key for an asset
func HashForKey(objectType string, hashElements ...string) string {
	toHash := objectType
	sort.Strings(hashElements)
	for _, element := range hashElements {
		toHash += "," + element
	}
	sum := sha256.Sum256([]byte(toHash))
	return hex.EncodeToString(sum[:])
}

func determineStatusFromInModels(statuses []string) string {
	if stringInSlice(StatusFailed, statuses) {
		return StatusFailed
	}

	if stringInSlice(StatusAborted, statuses) {
		return StatusAborted
	}

	for _, s := range statuses {
		if s != StatusDone {
			return StatusWaiting
		}
	}
	return StatusTodo
}

func determineTupleStatus(db *LedgerDB, tupleStatus, computePlanID string) (string, error) {
	if tupleStatus != StatusWaiting || computePlanID == "" {
		return tupleStatus, nil
	}
	computePlan, err := db.GetComputePlan(computePlanID)
	if err != nil {
		return "", err
	}
	if stringInSlice(computePlan.State.Status, []string{StatusFailed, StatusCanceled}) {
		return StatusAborted, nil
	}
	return tupleStatus, nil
}

func createModelIndex(db *LedgerDB, modelHash, tupleKey string) error {
	return db.CreateIndex("tuple~modelHash~key", []string{"tuple", modelHash, tupleKey})
}
