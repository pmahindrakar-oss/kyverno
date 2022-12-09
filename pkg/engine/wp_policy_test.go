package engine

import (
	"encoding/json"
	"testing"

	kyverno "github.com/kyverno/kyverno/api/kyverno/v1"
	enginecontext "github.com/kyverno/kyverno/pkg/engine/context"
	"github.com/kyverno/kyverno/pkg/engine/utils"
	"gotest.tools/assert"
)

var policyRaw = []byte(`{
    "apiVersion": "kyverno.io/v1",
    "kind": "ClusterPolicy",
    "metadata": {
        "name": "woven-planet-policy"
    },
    "spec": {
        "rules": [
            {
                "name": "general-purpose-cpu-pool",
                "match": {
                    "all": [
                        {
                            "resources": {
                                "kinds": [
                                    "Pod"
                                ],
                                "selector": {
                                    "matchLabels": {
                                        "platform": "flyte"
                                    }
                                }
                            }
                        }
                    ]
                },
                "exclude": {
                    "all": [
                        {
                            "resources": {
                                "selector": {
                                    "matchLabels": {
                                        "project": "trajectories",
                                        "task-name": "app-workflows-fs1-scene-tasks-detector-task"
                                    }
                                }
                            }
                        }
                    ]
                },
                "preconditions": {
                    "all": [
                        {
                            "key": "spark-kubernetes-driver",
                            "operator": "AnyNotIn",
                            "value": "{{ request.object.spec.containers[].name || ` + " `[]`" + ` }}"
                        },
                        {
                            "key": "spark-kubernetes-executor",
                            "operator": "AnyNotIn",
                            "value": "{{ request.object.spec.containers[].name || ` + " `[]`" + ` }}"
                        },
                        {
                            "key": "{{ request.object.spec.nodeSelector.\"l5.lyft.com/pool\" || '' }}",
                            "operator": "NotEqual",
                            "value": "eks-pdx-pool-gpu"
                        }
                    ]
                },
                "mutate": {
                    "patchesJson6902": "- op: add\n  path: \"/spec/affinity/nodeAffinity/requiredDuringSchedulingIgnoredDuringExecution/nodeSelectorTerms/-1/matchExpressions/-1\"\n  value: {\"key\": \"l5.lyft.com/pool\", \"values\": [\"eks-pdx-pool-cpu\"], \"operator\": \"In\"}\n- op: remove\n  path: \"/spec/nodeSelector\""
                }
            },
            {
                "name": "spark-driver-pod-resource-operator-pool",
                "match": {
                    "all": [
                        {
                            "resources": {
                                "kinds": [
                                    "Pod"
                                ],
                                "selector": {
                                    "matchLabels": {
                                        "platform": "flyte"
                                    }
                                }
                            }
                        }
                    ]
                },
                "preconditions": {
                    "all": [
                        {
                            "key": "spark-kubernetes-driver",
                            "operator": "AnyIn",
                            "value": "{{ request.object.spec.containers[].name || ` + " `[]`" + ` }}"
                        }
                    ]
                },
                "mutate": {
                    "patchesJson6902": "- op: remove\n  path: \"/spec/nodeSelector\"\n- op: add\n  path: \"/spec/tolerations/-\"\n  value: {\"key\":\"reserved\",\"operator\":\"Equal\",\"value\":\"operator\",\"effect\":\"NoSchedule\"}\n- op: add\n  path: \"/spec/affinity/nodeAffinity/requiredDuringSchedulingIgnoredDuringExecution/nodeSelectorTerms/-1/matchExpressions/-1\"\n  value: {\"key\": \"l5.lyft.com/pool\", \"values\": [\"eks-pdx-pool-operator\"], \"operator\": \"In\"}\n"
                }
            },
            {
                "name": "spark-executor-pod-resource-without-nodeselector",
                "match": {
                    "all": [
                        {
                            "resources": {
                                "kinds": [
                                    "Pod"
                                ],
                                "selector": {
                                    "matchLabels": {
                                        "platform": "flyte"
                                    }
                                }
                            }
                        }
                    ]
                },
                "exclude": {
                    "all": [
                        {
                            "resources": {
                                "selector": {
                                    "matchLabels": {
                                        "project": "trajectories",
                                        "task-name": "app-workflows-fs1-scene-tasks-detector-task"
                                    }
                                }
                            }
                        }
                    ]
                },
                "preconditions": {
                    "all": [
                        {
                            "key": "spark-kubernetes-executor",
                            "operator": "AnyIn",
                            "value": "{{ request.object.spec.containers[].name || ` + " `[]`" + ` }}"
                        },
                        {
                            "key": "{{ request.object.spec.nodeSelector.\"l5.lyft.com/pool\" || '' }}",
                            "operator": "NotEqual",
                            "value": "eks-pdx-pool-gpu"
                        }
                    ]
                },
                "mutate": {
                    "patchesJson6902": "- op: add\n  path: \"/spec/affinity/nodeAffinity/requiredDuringSchedulingIgnoredDuringExecution/nodeSelectorTerms/-1/matchExpressions/-1\"\n  value: {\"key\": \"l5.lyft.com/pool\", \"values\": [\"eks-pdx-pool-cpu\"], \"operator\": \"In\"}\n- op: remove\n  path: \"/spec/nodeSelector\""
                }
            },
            {
                "name": "spark-executor-pod-resource-with-nodeselector",
                "match": {
                    "all": [
                        {
                            "resources": {
                                "kinds": [
                                    "Pod"
                                ],
                                "selector": {
                                    "matchLabels": {
                                        "platform": "flyte"
                                    }
                                }
                            }
                        }
                    ]
                },
                "exclude": {
                    "all": [
                        {
                            "resources": {
                                "selector": {
                                    "matchLabels": {
                                        "project": "trajectories",
                                        "task-name": "app-workflows-fs1-scene-tasks-detector-task"
                                    }
                                }
                            }
                        }
                    ]
                },
                "preconditions": {
                    "all": [
                        {
                            "key": "spark-kubernetes-executor",
                            "operator": "AnyIn",
                            "value": "{{ request.object.spec.containers[].name || ` + " `[]`" + ` }}"
                        },
                        {
                            "key": "{{ request.object.spec.nodeSelector.\"l5.lyft.com/pool\" || '' }}",
                            "operator": "Equal",
                            "value": "eks-pdx-pool-gpu"
                        }
                    ]
                },
                "mutate": {
                    "patchesJson6902": "- op: add\n  path: \"/spec/containers/0/resources/limits/-\"\n  value: {\"nvidia.com/gpu\" : 1 }"
                }
            },
            {
                "name": "spark-executor-flyteondemandgpu-pool-with-nodeselector",
                "match": {
                    "all": [
                        {
                            "resources": {
                                "kinds": [
                                    "Pod"
                                ],
                                "selector": {
                                    "matchLabels": {
                                        "platform": "flyte",
                                        "project": "trajectories",
                                        "task-name": "app-workflows-fs1-scene-tasks-detector-task"
                                    }
                                }
                            }
                        }
                    ]
                },
                "preconditions": {
                    "all": [
                        {
                            "key": "spark-kubernetes-executor",
                            "operator": "AnyIn",
                            "value": "{{ request.object.spec.containers[].name || ` + " `[]`" + ` }}"
                        },
                        {
                            "key": "{{ request.object.spec.nodeSelector.\"l5.lyft.com/pool\" || '' }}",
                            "operator": "Equal",
                            "value": "eks-pdx-pool-gpu"
                        }
                    ]
                },
                "mutate": {
                    "patchesJson6902": "- op: add\n  path: \"/spec/containers/0/resources/limits/-\"\n  value: {\"nvidia.com/gpu\" : 1 }\n- op: add\n  path: \"/spec/tolerations/-\"\n  value: {\"key\":\"reserved\",\"operator\":\"Equal\",\"value\":\"flyteondemandgpu\",\"effect\":\"NoSchedule\"}\n- op: add\n  path: \"/spec/affinity/nodeAffinity/requiredDuringSchedulingIgnoredDuringExecution/nodeSelectorTerms/-1/matchExpressions/-1\"\n  value: {\"key\": \"l5.lyft.com/pool\", \"values\": [\"eks-pdx-pool-flyteondemandgpu\"], \"operator\": \"In\"}\n- op: remove\n  path: \"/spec/nodeSelector\""
                }
            },
            {
                "name": "non-spark-flyteondemandgpu-pool",
                "match": {
                    "all": [
                        {
                            "resources": {
                                "kinds": [
                                    "Pod"
                                ],
                                "selector": {
                                    "matchLabels": {
                                        "platform": "flyte",
                                        "project": "trajectories",
                                        "task-name": "app-workflows-fs1-scene-tasks-detector-task"
                                    }
                                }
                            }
                        }
                    ]
                },
                "preconditions": {
                    "all": [
                        {
                            "key": "{{ request.object.spec.nodeSelector.\"l5.lyft.com/pool\" || '' }}",
                            "operator": "Equal",
                            "value": "eks-pdx-pool-gpu"
                        },
                        {
                            "key": "spark-kubernetes-driver",
                            "operator": "AnyNotIn",
                            "value": "{{ request.object.spec.containers[].name || ` + " `[]`" + ` }}"
                        },
                        {
                            "key": "spark-kubernetes-executor",
                            "operator": "AnyNotIn",
                            "value": "{{ request.object.spec.containers[].name || ` + " `[]`" + ` }}"
                        }
                    ]
                },
                "mutate": {
                    "patchesJson6902": "- op: add\n  path: \"/spec/tolerations/-\"\n  value: {\"key\":\"reserved\",\"operator\":\"Equal\",\"value\":\"flyteondemandgpu\",\"effect\":\"NoSchedule\"}\n- op: add\n  path: \"/spec/affinity/nodeAffinity/requiredDuringSchedulingIgnoredDuringExecution/nodeSelectorTerms/-1/matchExpressions/-1\"\n  value: {\"key\": \"l5.lyft.com/pool\", \"values\": [\"eks-pdx-pool-flyteondemandgpu\"], \"operator\": \"In\"}\n- op: remove\n  path: \"/spec/nodeSelector\""
                }
            }
        ]
    }
}`)

type RuleTest struct {
	Name               string
	Resource           []byte
	ExpectedPatches    [][]byte
	ExpectedPolicyName []string
}

var patchRulesTest = []RuleTest{
	RuleTest{
		Name: "flyte pod resource general purpose cpu pool",
		Resource: []byte(`{
				  "apiVersion": "v1",
				  "kind": "Pod",
				  "metadata": {
					"name": "foo",
					"namespace": "default",
					"labels" : {
					  "platform": "flyte"
					}
				  },

				  "spec": {
					"containers": [
					  {
						"image": "bar:v2",
						"name": "bar"
					  }
					]
				  }
		}`),
		ExpectedPatches: [][]byte{
			[]byte(`{"op":"add","path":"/spec/affinity","value":{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"l5.lyft.com/pool","operator":"In","values":["eks-pdx-pool-cpu"]}]}]}}}}`),
		},
		ExpectedPolicyName: []string{"general-purpose-cpu-pool"},
	},
	RuleTest{
		Name: "flyte spark driver pod resource operator pool",
		Resource: []byte(`{
				"apiVersion": "v1",
				"kind": "Pod",
				"metadata": {
					"name": "foo",
					"namespace": "default",
					"labels" : {
						"platform": "flyte"
					}
				},

				"spec": {
					"containers": [
						{
							"image": "sparkDriverContainer:v2",
							"name": "spark-kubernetes-driver"
						}
					],
					"nodeSelector": {
						"l5.lyft.com/pool" : "eks-pdx-pool-gpu"
					}
				}
		}`),
		ExpectedPatches: [][]byte{
			[]byte(`{"op":"add","path":"/spec/tolerations","value":[{"effect":"NoSchedule","key":"reserved","operator":"Equal","value":"operator"}]}`),
			[]byte(`{"op":"add","path":"/spec/affinity","value":{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"l5.lyft.com/pool","operator":"In","values":["eks-pdx-pool-operator"]}]}]}}}}`),
			[]byte(`{"op":"remove","path":"/spec/nodeSelector"}`),
		},
		ExpectedPolicyName: []string{"spark-driver-pod-resource-operator-pool"},
	},
	RuleTest{
		Name: "flyte spark executor pod general pool without node selector",
		Resource: []byte(`{
				"apiVersion": "v1",
				"kind": "Pod",
				"metadata": {
					"name": "foo",
					"namespace": "default",
					"labels" : {
						"platform": "flyte"
					}
				},

				"spec": {
					"containers": [
						{
							"image": "spark-kubernetes-executor:v2",
							"name": "spark-kubernetes-executor"
						}
					]
				}
			}`),
		ExpectedPatches: [][]byte{
			[]byte(`{"op":"add","path":"/spec/affinity","value":{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"l5.lyft.com/pool","operator":"In","values":["eks-pdx-pool-cpu"]}]}]}}}}`),
		},
		ExpectedPolicyName: []string{"spark-executor-pod-resource-without-nodeselector"},
	},
	RuleTest{
		Name: "flyte spark executor pod general pool with node selector",
		Resource: []byte(`{
				"apiVersion": "v1",
				"kind": "Pod",
				"metadata": {
					"name": "foo",
					"namespace": "default",
					"labels" : {
						"platform": "flyte"
					}
				},

				"spec": {
					"containers": [
						{
							"image": "spark-kubernetes-executor:v2",
							"name": "spark-kubernetes-executor"
						}
					],
					"nodeSelector": {
						"l5.lyft.com/pool" : "eks-pdx-pool-gpu"
					}
				}
			}`),
		ExpectedPatches: [][]byte{
			[]byte(`{"op":"add","path":"/spec/containers/0/resources","value":{"limits":[{"nvidia.com/gpu":1}]}}`),
		},
		ExpectedPolicyName: []string{"spark-executor-pod-resource-with-nodeselector"},
	},
	RuleTest{
		Name: "flyte spark executor pod flyteondemandgpu pool with node selector and project and task label",
		Resource: []byte(`{
				"apiVersion": "v1",
				"kind": "Pod",
				"metadata": {
					"name": "foo",
					"namespace": "default",
					"labels" : {
						"platform": "flyte",
						"project": "trajectories",
						"task-name": "app-workflows-fs1-scene-tasks-detector-task"
					}
				},

				"spec": {
					"containers": [
						{
							"image": "spark-kubernetes-executor:v2",
							"name": "spark-kubernetes-executor"
						}
					],
					"nodeSelector": {
						"l5.lyft.com/pool" : "eks-pdx-pool-gpu"
					}
				}
			}`),
		ExpectedPatches: [][]byte{
			[]byte(`{"op":"add","path":"/spec/containers/0/resources","value":{"limits":[{"nvidia.com/gpu":1}]}}`),
			[]byte(`{"op":"add","path":"/spec/tolerations","value":[{"effect":"NoSchedule","key":"reserved","operator":"Equal","value":"flyteondemandgpu"}]}`),
			[]byte(`{"op":"add","path":"/spec/affinity","value":{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"l5.lyft.com/pool","operator":"In","values":["eks-pdx-pool-flyteondemandgpu"]}]}]}}}}`),
			[]byte(`{"op":"remove","path":"/spec/nodeSelector"}`),
		},
		ExpectedPolicyName: []string{"spark-executor-flyteondemandgpu-pool-with-nodeselector"},
	},
	RuleTest{
		Name: "flyte non spark pod flyteondemandgpu pool with node selector and project and task label",
		Resource: []byte(`{
				"apiVersion": "v1",
				"kind": "Pod",
				"metadata": {
					"name": "foo",
					"namespace": "default",
					"labels" : {
						"platform": "flyte",
						"project": "trajectories",
						"task-name": "app-workflows-fs1-scene-tasks-detector-task"
					}
				},

				"spec": {
					"containers": [
						{
							"image": "normal-executor:v2",
							"name": "normal-executor"
						}
					],
					"nodeSelector": {
						"l5.lyft.com/pool" : "eks-pdx-pool-gpu"
					}
				}
			}`),
		ExpectedPatches: [][]byte{
			[]byte(`{"op":"add","path":"/spec/tolerations","value":[{"effect":"NoSchedule","key":"reserved","operator":"Equal","value":"flyteondemandgpu"}]}`),
			[]byte(`{"op":"add","path":"/spec/affinity","value":{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"l5.lyft.com/pool","operator":"In","values":["eks-pdx-pool-flyteondemandgpu"]}]}]}}}}`),
			[]byte(`{"op":"remove","path":"/spec/nodeSelector"}`),
		},
		ExpectedPolicyName: []string{"non-spark-flyteondemandgpu-pool"},
	},
}

func TestAffinity(t *testing.T) {

	t.Run("flyte woven planet policy", func(t *testing.T) {
		var policy kyverno.ClusterPolicy
		err := json.Unmarshal(policyRaw, &policy)
		if err != nil {
			t.Error(err)
		}
		for _, test := range patchRulesTest {
			resourceUnstructured, err := utils.ConvertToUnstructured(test.Resource)
			assert.NilError(t, err)
			ctx := enginecontext.NewContext()
			err = enginecontext.AddResource(ctx, test.Resource)
			if err != nil {
				t.Error(err)
			}
			policyContext := &PolicyContext{
				policy:      &policy,
				jsonContext: ctx,
				newResource: *resourceUnstructured}
			resp := Mutate(policyContext)
			patches := [][]byte{}
			policies := []string{}
			if resp != nil {
				for _, r := range resp.PolicyResponse.Rules {
					if len(r.Patches) > 0 {
						policies = append(policies, r.Name)
						for _, p := range r.Patches {
							patches = append(patches, p)
						}
					}
				}
			}
			assert.Equal(t, len(patches), len(test.ExpectedPatches))
			for index, expectedPatch := range test.ExpectedPatches {
				assert.Equal(t, string(expectedPatch), string(patches[index]))
			}
			for i, _ := range policies {
				assert.Equal(t, policies[i], test.ExpectedPolicyName[i])
			}
		}
	})
}
