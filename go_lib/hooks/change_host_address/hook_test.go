package change_host_address

import (
	"github.com/deckhouse/deckhouse/go_lib/set"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_changeHostAddress(t *testing.T) {
	mock := newMockPodClient()

	type args struct {
		podClient podClient
		podHosts  []podHost
	}
	tests := []struct {
		name   string
		args   args
		assert func(*testing.T, *mockPodClient)
	}{
		{
			name: "does nothing when no pods available",
			args: args{
				podClient: mock,
				podHosts:  []podHost{},
			},
			assert: func(t *testing.T, client *mockPodClient) {
				assert.Equal(t, client.deletedName.Size(), 0)
				assert.Equal(t, client.annotatedName.Size(), 0)
				assert.Equal(t, client.annotatedHost.Size(), 0)
			},
		},
		{
			name: "ignores pod without IP",
			args: args{
				podClient: mock,
				podHosts: []podHost{{
					Name:      "a",
					IP:        "",
					InitialIP: "1.2.3.4",
				}},
			},
			assert: func(t *testing.T, client *mockPodClient) {
				assert.Equal(t, 0, client.deletedName.Size())
				assert.Equal(t, 0, client.annotatedName.Size())
				assert.Equal(t, 0, client.annotatedHost.Size())
			},
		},
		{
			name: "annotates Pod with IP when Pod does not have the annotation",
			args: args{
				podClient: mock,
				podHosts: []podHost{{
					Name:      "a",
					IP:        "1.2.3.4",
					InitialIP: "",
				}},
			},
			assert: func(t *testing.T, client *mockPodClient) {
				assert.Equal(t, 0, client.deletedName.Size())
				assert.True(t, client.annotatedHost.Has("1.2.3.4"), "expected IP from the Pod")
				assert.True(t, client.annotatedName.Has("a"), "expected name from the Pod")
			},
		},
		{
			name: "deletes Pod when IP and initial IP do not equal",
			args: args{
				podClient: mock,
				podHosts: []podHost{{
					Name:      "a",
					IP:        "1.2.3.4",
					InitialIP: "4.3.2.1",
				}},
			},
			assert: func(t *testing.T, client *mockPodClient) {
				assert.Equal(t, 1, client.deletedName.Size())
				assert.True(t, client.deletedName.Has("a"))
				assert.Equal(t, 0, client.annotatedName.Size())
				assert.Equal(t, 0, client.annotatedHost.Size())
			},
		},
		{
			name: "process multiple pods correclty",
			args: args{
				podClient: mock,
				podHosts: []podHost{{
					Name:      "annotated",
					IP:        "1.2.3.4",
					InitialIP: "",
				}, {
					Name:      "deleted",
					IP:        "2.3.4.5",
					InitialIP: "4.3.2.1",
				}, {
					Name:      "skipped",
					IP:        "",
					InitialIP: "",
				}},
			},
			assert: func(t *testing.T, client *mockPodClient) {
				assert.Equal(t, 1, client.deletedName.Size())
				assert.True(t, client.deletedName.Has("deleted"))

				assert.Equal(t, 1, client.annotatedName.Size())
				assert.Equal(t, 1, client.annotatedHost.Size())
				assert.True(t, client.annotatedName.Has("annotated"))
				assert.True(t, client.annotatedHost.Has("1.2.3.4"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.args.podClient.(*mockPodClient)
			defer mock.clean()

			changeHostAddress(tt.args.podClient, tt.args.podHosts)
			tt.assert(t, mock)
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
