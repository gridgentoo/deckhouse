/*
Copyright 2021 Flant JSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package change_host_address

import (
	"testing"

	"github.com/deckhouse/deckhouse/go_lib/set"
	"github.com/stretchr/testify/assert"
)

func Test_changeHostAddress(t *testing.T) {

	tests := []struct {
		name     string
		podHosts []podHost
		assert   func(*testing.T, *mockPodClient)
	}{
		{
			name:     "does nothing when no pods available",
			podHosts: []podHost{},
			assert: func(t *testing.T, client *mockPodClient) {
				assert.Equal(t, client.deletedName.Size(), 0)
				assert.Equal(t, client.annotatedName.Size(), 0)
				assert.Equal(t, client.annotatedHost.Size(), 0)
			},
		},
		{
			name: "ignores pod without IP",
			podHosts: []podHost{{
				Name:      "a",
				IP:        "",
				InitialIP: "1.2.3.4",
			}},
			assert: func(t *testing.T, client *mockPodClient) {
				assert.Equal(t, 0, client.deletedName.Size())
				assert.Equal(t, 0, client.annotatedName.Size())
			},
		},
		{
			name: "annotates Pod with IP when Pod does not have the annotation",
			podHosts: []podHost{{
				Name:      "a",
				IP:        "1.2.3.4",
				InitialIP: "",
			}},
			assert: func(t *testing.T, client *mockPodClient) {
				assert.Equal(t, 0, client.deletedName.Size())

				assert.True(t, client.annotatedName.Has("a"), "expected name from the Pod")
				assert.True(t, client.annotatedHost.Has("1.2.3.4"), "expected IP from the Pod")
			},
		},
		{
			name: "deletes Pod when IP and initial IP do not equal",
			podHosts: []podHost{{
				Name:      "a",
				IP:        "1.2.3.4",
				InitialIP: "4.3.2.1",
			}},
			assert: func(t *testing.T, client *mockPodClient) {
				assert.Equal(t, 1, client.deletedName.Size())
				assert.True(t, client.deletedName.Has("a"))

				assert.Equal(t, 0, client.annotatedName.Size())
				assert.Equal(t, 0, client.annotatedHost.Size())
			},
		},
		{
			name: "processes multiple pods correclty",
			podHosts: []podHost{{
				Name:      "no-anno",
				IP:        "4.5.6.7",
			}, {
				Name:      "badhost",
				IP:        "2.3.4.5",
				InitialIP: "4.3.2.1",
			}, {
				Name:      "skipped",
			}, {
				Name:      "no-anno-2",
				IP:        "4.5.6.8",
			}},
			assert: func(t *testing.T, client *mockPodClient) {
				assert.Equal(t, 1, client.deletedName.Size())
				assert.True(t, client.deletedName.Has("badhost"))

				assert.Equal(t, 2, client.annotatedName.Size())

				assert.True(t, client.annotatedName.Has("no-anno"))
				assert.True(t, client.annotatedHost.Has("4.5.6.7"))

				assert.True(t, client.annotatedName.Has("no-anno-2"))
				assert.True(t, client.annotatedHost.Has("4.5.6.8"))
			},
		},
	}

	clientMock := newMockPodClient()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer clientMock.clean()

			changeHostAddress(clientMock, tt.podHosts)

			tt.assert(t, clientMock)
		})
	}
}

func newMockPodClient() *mockPodClient {
	return &mockPodClient{
		deletedName:   set.New(),
		annotatedName: set.New(),
		annotatedHost: set.New(),
	}
}

type mockPodClient struct {
	deletedName   set.Set
	annotatedName set.Set
	annotatedHost set.Set
}

func (pc *mockPodClient) Delete(name string) {
	pc.deletedName.Add(name)
}

func (pc *mockPodClient) AnnotateHost(name, host string) {
	pc.annotatedName.Add(name)
	pc.annotatedHost.Add(host)
}

func (pc *mockPodClient) clean() {
	pc.deletedName = set.New()
	pc.annotatedName = set.New()
	pc.annotatedHost = set.New()
}
