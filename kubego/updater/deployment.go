package updater

import (
	"bytes"
	"context"
	"fmt"
	"github.com/illuminatingKong/kongming-kit/kubego/util"
	"text/template"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var restartPatchTemplate = template.Must(template.New("restart-patch-template").Parse(`{
  "spec": {
    "template": {
      "metadata": {
        "annotations": {
          "updated-by-devops": "{{.Updated}}"
        },
        "labels": {
          "updated-by-devops": "{{.Updated}}"
        }
      }
    }
  }
}`))

func PatchDeployment(ns, name string, patchBytes []byte, cl client.Client) error {
	return patchObject(&appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      name,
		},
	}, patchBytes, cl)
}

// TODO: LOU（LOU is from zhangmen）: it is not the right way to restart a deployment, since it is a hack and it
// will generate a new revision which will pollute the revision history.
func RestartDeployment(ns, name string, cl client.Client) error {
	now := time.Now().Format(time.RFC3339Nano)
	payload := bytes.NewBufferString("")
	_ = restartPatchTemplate.Execute(payload, struct {
		Updated string
	}{now})

	if err := PatchDeployment(ns, name, payload.Bytes(), cl); err != nil {
		return fmt.Errorf("failed to restart %s/deploy/%s: %v", ns, name, err)
	}

	return nil
}

func DeleteDeployments(namespace string, selector labels.Selector, clientset *kubernetes.Clientset) error {
	deletePolicy := metav1.DeletePropagationForeground
	err := clientset.AppsV1().Deployments(namespace).DeleteCollection(
		context.TODO(),
		metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		},
		metav1.ListOptions{
			LabelSelector: selector.String(),
		},
	)

	return util.IgnoreNotFoundError(err)
}

func UpdateDeploymentImage(ns, name, container, image string, cl client.Client) error {
	patchBytes := []byte(fmt.Sprintf(`{"spec":{"template":{"spec":{"containers":[{"name":"%s","image":"%s"}]}}}}`, container, image))

	return PatchDeployment(ns, name, patchBytes, cl)
}

func ScaleDeployment(ns, name string, replicas int, cl client.Client) error {
	patchBytes := []byte(fmt.Sprintf(`{"spec":{"replicas": %d}}`, replicas))
	return PatchDeployment(ns, name, patchBytes, cl)
}

func CreateOrPatchDeployment(d *appsv1.Deployment, cl client.Client) error {
	return createOrPatchObject(d, cl)
}
