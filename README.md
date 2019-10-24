# ArgoCD Slack Notifier

Users of [ArgoCD](https://github.com/argoproj/argo-cd/) may wish to have automated slack messages sent when events occur within the CD pipeline. There is no built in mechanism in ArgoCD to do this - which would be out of scope - but there are _hooks_ which can be leveraged to have the same effect.

This repository builds a public docker image which can be used to send ArgoCD centric messages to a designated slack channel informing users of events that occur within the pipeline.

Argo supports four hook events: `PreSync`, `Sync`, `PostSync` and `SyncFail`. You can read more about the [circumstances under which they are triggered here](https://argoproj.github.io/argo-cd/user-guide/resource_hooks).

## Running

### Docker

The docker image is available at `andondev/argocd-slack-notifier`:

```bash
docker pull andondev/argocd-slack-notifier:latest
```

You can supply env variables from a file by running the following (assuming you have a `.env` with the required [environment variables](#environment-variables)):

```bash
docker run andondev/argocd-slack-notifier --env-file .env
```

You can also supply env variables indivudualy by running the following format:

```bash
docker run andondev/argocd-slack-notifier \
  -e NOTIFIER_ARGO_CD_EVENT_TYPE=PostSync \
  -e NOTIFIER_SLACK_URL=<SLACK_URL>
  [...]
```

### Locally (Go code)

The quickest way to test manually is to clone this repository then run the following:

```bash
cp .env.example .env
# now add your variables to .env
make run
```

## Deploying

The container is intended to be ran as a kuberenetes [Job](https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/).

It is also intended to be ran on each event trigger, so **do not give it a `metadata.name`** - instead, use `metdata.generateName` which will allow it to be triggered as many times as is required.

This also has an implication on how you should deploy it - instead of using `kubectl apply ...` **you must use `kubectl create ...`**

## Templating

This container is most useful when deployed to the same cluster that ArgoCD resides on. In this instance, it is useful to template your manifests. To get started quickly there is a [provided example of a helm template](/examples/helm/).

## Environment Variables

To keep the container as flexible as possible, several environment variables can be set so that you can customise the incoming messages:

Variable | Description | Default | Required
-------- | ----------- | ------- | --------
NOTIFIER_ARGO_CD_EVENT_TYPE | The event being observerd | none | `true`
NOTIFIER_ARGO_CD_BASE_URL | Your base URL for ArgoCD which is used to generate links to apps, projects and rollbacks | none | `true`
NOTIFIER_ARGO_CD_PROJECT | The name of the Argo project | none | `true`
NOTIFIER_ARGO_CD_APPLICATION | The name of the application | none | `true`
NOTIFIER_SLACK_URL | The incoming  wehbook URL generated for you Slack app | none | `true`
NOTIFIER_SLACK_CHANNEL | The slack channel to post to | none | `true`
NOTIFIER_SLACK_USERNAME | Your apps username | `Notifier` | `false`
NOTIFIER_SLACK_ICON_EMOJI | Your apps emoji | `:monkey:` | `false`
NOTIFIER_SLACK_TEXT | Any additional text to display in the message during a `PreSync` | `Deployment notification` | `false`
NOTIFIER_SLACK_FOOTER_TEXT | Any text to add to the footer (the current time is added automatically) | none | `false`
