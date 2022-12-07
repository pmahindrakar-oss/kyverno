package engine

import (
	"encoding/json"
	"reflect"
	"testing"

	kyverno "github.com/kyverno/kyverno/api/kyverno/v1"
	enginecontext "github.com/kyverno/kyverno/pkg/engine/context"
	"github.com/kyverno/kyverno/pkg/engine/utils"
	"gotest.tools/assert"
)

func TestAffinity(t *testing.T) {
	t.Run("non pod resource", func(t *testing.T) {
		policyRaw := []byte(`
	{
    "apiVersion": "kyverno.io/v1",
    "kind": "ClusterPolicy",
    "metadata": {
        "name": "spark-test-policy"
    },
    "spec": {
        "rules": [
            {
                "name": "toleration-and-affinity",
				  "match": {
					"resources": {
					  "kinds": [
						"Pod"
					  ]
					}
				  },
	            "preconditions": {
	                "any": [
	                    {
	                        "key": "spark-kubernetes-driver",
	                        "operator": "In",
	                        "value": "{{ request.object.spec.containers[].name ||` + " `[]`" + ` }}"
	                    }
	                ]
	            },
                "mutate": {
                    "patchesJson6902": "- op: add\n  path: \"/spec/tolerations/-\"\n  value: {\"key\":\"reserved\",\"operator\":\"Equal\",\"value\":\"operator\",\"effect\":\"NoSchedule\"}\n- op: add\n  path: \"/spec/affinity/nodeAffinity/requiredDuringSchedulingIgnoredDuringExecution/nodeSelectorTerms/-1/matchExpressions/-1\"\n  value: {\"key\": \"l5.lyft.com/pool\", \"values\": [\"eks-pdx-pool-operator\"], \"operator\": \"In\"}\n- op: remove\n  path: \"/spec/nodeSelector\""
                }
            }
        ]
    	}
	}`)
		resourceRaw := []byte(`{
				"apiVersion": "v1",
				"kind": "Deployment",
				"metadata": {
					"name": "foo",
					"namespace": "default"
				}
			}`)
		var policy kyverno.ClusterPolicy
		err := json.Unmarshal(policyRaw, &policy)
		if err != nil {
			t.Error(err)
		}
		resourceUnstructured, err := utils.ConvertToUnstructured(resourceRaw)
		assert.NilError(t, err)
		ctx := enginecontext.NewContext()
		err = enginecontext.AddResource(ctx, resourceRaw)
		if err != nil {
			t.Error(err)
		}
		value, err := ctx.Query("request.object.metadata.name")

		t.Log(value)
		if err != nil {
			t.Error(err)
		}
		policyContext := &PolicyContext{
			policy:      &policy,
			jsonContext: ctx,
			newResource: *resourceUnstructured}
		er := Mutate(policyContext)

		assert.Equal(t, len(er.PolicyResponse.Rules), 0)
	})

	t.Run("flyte spark driver pod resource operator pool", func(t *testing.T) {
		policyRaw := []byte(`
	{
    "apiVersion": "kyverno.io/v1",
    "kind": "ClusterPolicy",
    "metadata": {
        "name": "spark-test-policy"
    },
    "spec": {
        "rules": [
            {
                "name": "toleration-and-affinity",
				  "match": {
					"resources": {
					  "kinds": [
						"Pod"
					  ]
					}
				  },
	            "preconditions": {
	                "any": [
	                    {
	                        "key": "spark-kubernetes-driver",
	                        "operator": "In",
	                        "value": "{{ request.object.spec.containers[].name ||` + " `[]`" + ` }}"
	                    }
	                ]
	            },
                "mutate": {
                    "patchesJson6902": "- op: add\n  path: \"/spec/tolerations/-\"\n  value: {\"key\":\"reserved\",\"operator\":\"Equal\",\"value\":\"operator\",\"effect\":\"NoSchedule\"}\n- op: add\n  path: \"/spec/affinity/nodeAffinity/requiredDuringSchedulingIgnoredDuringExecution/nodeSelectorTerms/-1/matchExpressions/-1\"\n  value: {\"key\": \"l5.lyft.com/pool\", \"values\": [\"eks-pdx-pool-operator\"], \"operator\": \"In\"}\n- op: remove\n  path: \"/spec/nodeSelector\""
                }
            }
        ]
    	}
	}`)

		resourceRaw := []byte(`{
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
					]
				}
			}`)
		expectedTolerationPatch := []byte(`{"op":"add","path":"/spec/tolerations","value":[{"effect":"NoSchedule","key":"reserved","operator":"Equal","value":"operator"}]}`)
		expectedAffinityPatch := []byte(`{"op":"add","path":"/spec/affinity","value":{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"l5.lyft.com/pool","operator":"In","values":["eks-pdx-pool-operator"]}]}]}}}}`)
		var policy kyverno.ClusterPolicy
		err := json.Unmarshal(policyRaw, &policy)
		if err != nil {
			t.Error(err)
		}
		resourceUnstructured, err := utils.ConvertToUnstructured(resourceRaw)
		assert.NilError(t, err)
		ctx := enginecontext.NewContext()
		err = enginecontext.AddResource(ctx, resourceRaw)
		if err != nil {
			t.Error(err)
		}
		value, err := ctx.Query("request.object.metadata.name")

		t.Log(value)
		if err != nil {
			t.Error(err)
		}
		policyContext := &PolicyContext{
			policy:      &policy,
			jsonContext: ctx,
			newResource: *resourceUnstructured}
		er := Mutate(policyContext)
		t.Log(string(expectedTolerationPatch))

		assert.Equal(t, len(er.PolicyResponse.Rules), 1)
		assert.Equal(t, len(er.PolicyResponse.Rules[0].Patches), 2)
		t.Log(string(er.PolicyResponse.Rules[0].Patches[0]))
		if !reflect.DeepEqual(expectedTolerationPatch, er.PolicyResponse.Rules[0].Patches[0]) {
			t.Error("toleration patch  doesn't match")
		}
		if !reflect.DeepEqual(expectedAffinityPatch, er.PolicyResponse.Rules[0].Patches[1]) {
			t.Error("affinity patch doesn't match")
		}
	})
}
