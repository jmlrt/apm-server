// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package modelprocessor_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elastic/apm-server/model"
	"github.com/elastic/apm-server/model/modelprocessor"
)

func TestSetDefaultServiceEnvironment(t *testing.T) {
	nonEmptyMetadata := model.Metadata{Service: model.Service{Environment: "nonempty"}}
	defaultMetadata := model.Metadata{Service: model.Service{Environment: "default"}}

	processor := modelprocessor.SetDefaultServiceEnvironment{
		DefaultServiceEnvironment: "default",
	}
	testProcessBatchMetadata(t, &processor, nonEmptyMetadata, nonEmptyMetadata)
	testProcessBatchMetadata(t, &processor, model.Metadata{}, defaultMetadata)
}

func testProcessBatchMetadata(t *testing.T, processor model.BatchProcessor, in, out model.Metadata) {
	t.Helper()

	// Check that the model.Batch fields have not changed since this
	// test was last updated, to ensure we process all model types.
	var batchFields []string
	typ := reflect.TypeOf(model.Batch{})
	for i := 0; i < typ.NumField(); i++ {
		batchFields = append(batchFields, typ.Field(i).Name)
	}
	assert.ElementsMatch(t, []string{
		"Transactions",
		"Spans",
		"Metricsets",
		"Errors",
		"Profiles",
	}, batchFields)

	batch := &model.Batch{
		Transactions: []*model.Transaction{
			{Metadata: in},
		},
		Spans: []*model.Span{
			{Metadata: in},
		},
		Metricsets: []*model.Metricset{
			{Metadata: in},
		},
		Errors: []*model.Error{
			{Metadata: in},
		},
		Profiles: []*model.PprofProfile{
			{Metadata: in},
		},
	}
	err := processor.ProcessBatch(context.Background(), batch)
	require.NoError(t, err)

	expected := &model.Batch{
		Transactions: []*model.Transaction{
			{Metadata: out},
		},
		Spans: []*model.Span{
			{Metadata: out},
		},
		Metricsets: []*model.Metricset{
			{Metadata: out},
		},
		Errors: []*model.Error{
			{Metadata: out},
		},
		Profiles: []*model.PprofProfile{
			{Metadata: out},
		},
	}
	assert.Equal(t, expected, batch)
}