package controllers

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestApplyTemplate(t *testing.T) {
	tests := []struct {
		name    string
		obj     unstructured.Unstructured
		vars    map[string]string
		envVars []corev1.EnvVar
		want    unstructured.Unstructured
	}{
		{
			name: "nothing to apply",
			obj: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			vars:    map[string]string{},
			envVars: []corev1.EnvVar{},
			want: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
		},
		{
			name: "apply vars",
			obj: unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"metadata": map[string]interface{}{
						"name": "foo-{{PR_NUMBER}}",
					},
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"name":  "echo",
										"image": "gcr.io/mercari/echo:{{CONTAINER_IMAGE_TAG}}",
									},
								},
							},
						},
					},
				},
			},
			vars: map[string]string{
				"PR_NUMBER":           "10",
				"CONTAINER_IMAGE_TAG": "123deadbeafdeadbeaf",
			},
			envVars: []corev1.EnvVar{},
			want: unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"metadata": map[string]interface{}{
						"name": "foo-10",
					},
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"name":  "echo",
										"image": "gcr.io/mercari/echo:123deadbeafdeadbeaf",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "add & override env vars",
			obj: unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"metadata": map[string]interface{}{
						"name": "foo-{{PR_NUMBER}}",
					},
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"name": "echo",
										"env": []interface{}{
											map[string]interface{}{
												"name":  "override",
												"value": "old value",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			vars: map[string]string{
				"PR_NUMBER": "10",
			},
			envVars: []corev1.EnvVar{
				{
					Name:  "new",
					Value: "new value",
				},
				{
					Name:  "override",
					Value: "override value",
				},
				{
					Name:  "override with placeholder",
					Value: "override value No. {{PR_NUMBER}}",
				},
			},
			want: unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"metadata": map[string]interface{}{
						"name": "foo-10",
					},
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"name": "echo",
										"env": []interface{}{
											map[string]interface{}{
												"name":  "override",
												"value": "override value",
											},
											map[string]interface{}{
												"name":  "new",
												"value": "new value",
											},
											map[string]interface{}{
												"name":  "override with placeholder",
												"value": "override value No. 10",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := applyTemplate(tt.obj, tt.vars, tt.envVars); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("applyTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeResourceConfigs(t *testing.T) {
	tests := []struct {
		name     string
		rendered unstructured.Unstructured
		exsiting unstructured.Unstructured
		want     unstructured.Unstructured
	}{
		{
			name: "empty",
			rendered: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{},
					"spec":     map[string]interface{}{},
				},
			},
			exsiting: unstructured.Unstructured{
				Object: map[string]interface{}{},
			},
			want: unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{},
					"spec":     map[string]interface{}{},
				},
			},
		},
		{
			name: "merge",
			rendered: unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Service",
					"metadata": map[string]interface{}{
						"name": "foo-11",
					},
					"spec": map[string]interface{}{
						"foo":  "bar-changed",
						"foo2": "bar2",
					},
				},
			},
			exsiting: unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Service",
					"metadata": map[string]interface{}{
						"name":              "foo-11",
						"creationTimestamp": "2021-12-16T01:40:49Z",
						"resourceVersion":   "1666331458",
						"uid":               "2cf24a3f-52ac-47c6-b30f-4d2755156549",
					},
					"spec": map[string]interface{}{
						"clusterIP": "10.0.255.0",
						"clusterIPs": []interface{}{
							"10.0.255.0",
						},
						"foo":  "bar",
						"foo2": "bar2",
					},
				},
			},
			want: unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Service",
					"metadata": map[string]interface{}{
						"name":              "foo-11",
						"creationTimestamp": "2021-12-16T01:40:49Z",
						"resourceVersion":   "1666331458",
						"uid":               "2cf24a3f-52ac-47c6-b30f-4d2755156549",
					},
					"spec": map[string]interface{}{
						"clusterIP": "10.0.255.0",
						"clusterIPs": []interface{}{
							"10.0.255.0",
						},
						"foo":  "bar-changed",
						"foo2": "bar2",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeResourceConfigs(tt.rendered, tt.exsiting); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("mergeResourceConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}
