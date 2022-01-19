package github

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-playground/webhooks/v6/github"
	kubetempurav1 "github.com/mercari/kubetempura/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	path = "/webhooks"
)

var (
	log = ctrl.Log.WithName("github")
)

func Webhooks(c client.Client, githubWebHookSecret string) {
	hook, err := github.New(github.Options.Secret(githubWebHookSecret))
	if err != nil {
		log.Error(err, "Failed to initialize the GitHub library")
		return
	}

	log.Info("Github Webhooks started")

	serveMux := http.NewServeMux()
	serveMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		log.Info("request:", "request", r.RequestURI)
		payload, err := hook.Parse(r, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn't one of the ones asked to be parsed
			}
			log.Error(err, "Failed to parse the request from GitHub.")
			return
		}
		switch payload.(type) {
		case github.PullRequestPayload:
			pullRequest := payload.(github.PullRequestPayload)
			handlePREvent(pullRequest, c)
		}
		_, err = w.Write([]byte("OK"))
		if err != nil {
			log.Error(err, "Failed to return the response")
		}
	})
	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("I'm KubeTempura."))
		if err != nil {
			log.Error(err, "Failed to return the response")
		}
	})

	err = http.ListenAndServe(":3000", serveMux)
	if err != nil {
		log.Error(err, "Failed to listen and serve the github webhook.")
		return
	}
}

func handlePREvent(prp github.PullRequestPayload, c client.Client) {
	if !(prp.Action == "opened" ||
		prp.Action == "reopened" ||
		prp.Action == "synchronize" ||
		prp.Action == "closed") {
		return
	}
	reviewApps, err := getReviewApps(c)
	if err != nil {
		log.Error(err, "Failed to get ReviewApps")
		return
	}
	reviewApps = findReviewAppsByRepository(reviewApps, prp.Repository.HTMLURL)
	if len(reviewApps) == 0 {
		return
	}
	prNumber := strconv.FormatInt(prp.Number, 10)
	if prp.Action == "closed" {
		prClosed(reviewApps, prNumber, c)
		return
	}
	prUpdated(reviewApps, prNumber, prp.PullRequest.Head.Sha, c)
}

func getReviewApps(c client.Client) ([]kubetempurav1.ReviewApp, error) {
	var reviewApps = kubetempurav1.ReviewAppList{}
	err := c.List(context.Background(), &reviewApps)
	return reviewApps.Items, err
}

func findReviewAppsByRepository(reviewApps []kubetempurav1.ReviewApp, repository string) []kubetempurav1.ReviewApp {
	var ret []kubetempurav1.ReviewApp
	for _, reviewApp := range reviewApps {
		if reviewApp.Spec.GithubRepository == repository {
			ret = append(ret, reviewApp)
			continue
		}
	}
	return ret
}

func prClosed(reviewApps []kubetempurav1.ReviewApp, prNumber string, c client.Client) {
	for _, reviewApp := range reviewApps {
		pr := generatePRStruct(reviewApp, prNumber, "")
		err := c.Delete(context.Background(), &pr)
		if err != nil {
			log.Error(err, "Failed to delete the PR")
		}
	}
}

func prUpdated(reviewApps []kubetempurav1.ReviewApp, prNumber string, sha string, c client.Client) {
	for _, reviewApp := range reviewApps {
		log.Info("PR updated" + reviewApp.Name)
		pr := generatePRStruct(reviewApp, prNumber, sha)
		rendered := *pr.DeepCopy()
		_, err := ctrl.CreateOrUpdate(context.Background(), c, &pr, func() error {
			pr.Spec = rendered.Spec
			return ctrl.SetControllerReference(&reviewApp, &pr, c.Scheme())
		})
		if err != nil {
			log.Error(err, "Failed to create/update the PR")
		}
	}
}

func prName(reviewAppName string, prNumber string) string {
	return reviewAppName + "-pr" + prNumber
}

func generatePRStruct(reviewApp kubetempurav1.ReviewApp, prNumber string, sha string) kubetempurav1.PR {
	return kubetempurav1.PR{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "kubetempura.mercari.com/v1",
			Kind:       "PR",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      prName(reviewApp.Name, prNumber),
			Namespace: reviewApp.Namespace,
		},
		Spec: kubetempurav1.PRSpec{
			ParentReviewApp: reviewApp.Name,
			PRNumber:        prNumber,
			HeadCommitRef:   sha,
			EnvVars:         nil, // TODO: we can give environment variables with the future updates.
		},
	}
}
