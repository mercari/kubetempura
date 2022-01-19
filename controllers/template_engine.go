package controllers

import (
	"regexp"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func applyTemplate(obj unstructured.Unstructured, vars map[string]string, envVars []corev1.EnvVar) unstructured.Unstructured {
	o := obj.DeepCopy()
	o.Object = applyEnvVars(o.Object, envVars)
	o.Object = replacePlaceholdersRecursive(o.Object, vars).(map[string]interface{})
	return *o
}

func applyEnvVars(a map[string]interface{}, envVars []corev1.EnvVar) map[string]interface{} {
	var containers []interface{}
	switch a["kind"] {
	case "Deployment":
		containers = a["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})
	case "Job":
		containers = a["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})
	case "CronJob":
		containers = a["spec"].(map[string]interface{})["jobTemplate"].(map[string]interface{})["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})
	default:
		return a
	}

	for _, container := range containers {
		c := container.(map[string]interface{})
		_olds, ok := c["env"]
		if !ok {
			_olds = []interface{}{}
		}
		olds := _olds.([]interface{})
		for _, ne := range envVars {
			found := false
			for j, old := range olds {
				if ne.Name == old.(map[string]interface{})["name"] {
					olds[j] = envVarToMap(ne)
					found = true
					break
				}
			}
			if !found {
				olds = append(olds, envVarToMap(ne))
			}
		}
		if len(olds) != 0 {
			c["env"] = olds
		}
	}
	return a
}

func envVarToMap(envVar corev1.EnvVar) map[string]interface{} {
	m := map[string]interface{}{
		"name": envVar.Name,
	}
	if envVar.Value != "" {
		m["value"] = envVar.Value
	}

	// TODO: Support the "ValueFrom" syntax
	//if envVar.ValueFrom != "" {
	//	m["value"] = envVar.ValueFrom
	//}

	return m
}

func replacePlaceholdersRecursive(a interface{}, vars map[string]string) interface{} {
	switch aa := a.(type) {
	case string:
		return replacePlaceholders(aa, vars)
	case map[string]interface{}:
		for k, v := range aa {
			aa[k] = replacePlaceholdersRecursive(v, vars)
		}
	case []interface{}:
		for i, v := range aa {
			aa[i] = replacePlaceholdersRecursive(v, vars)
		}
	}
	return a
}

func replacePlaceholders(s string, vars map[string]string) string {
	for k, v := range vars {
		r := regexp.MustCompile(`\{\{\s*?` + k + `\s*?\}\}`)
		s = r.ReplaceAllString(s, v)
	}
	return s
}

func mergeResourceConfigs(rendered unstructured.Unstructured, existing unstructured.Unstructured) unstructured.Unstructured {
	n := *existing.DeepCopy()
	metadata, ok := n.Object["metadata"]
	if !ok {
		metadata = map[string]interface{}{}
		n.Object["metadata"] = metadata
	}
	m := metadata.(map[string]interface{})
	for k, v := range rendered.Object["metadata"].(map[string]interface{}) {
		m[k] = v
	}
	spec, ok := n.Object["spec"]
	if !ok {
		spec = map[string]interface{}{}
		n.Object["spec"] = spec
	}
	s := spec.(map[string]interface{})
	for k, v := range rendered.Object["spec"].(map[string]interface{}) {
		s[k] = v
	}
	return n
}
